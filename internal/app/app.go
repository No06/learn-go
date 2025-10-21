package app

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	apihandlers "learn-go/internal/api/http"
	"learn-go/internal/api/ws"
	"learn-go/internal/config"
	"learn-go/internal/domain"
	gormrepo "learn-go/internal/repository/gorm"
	"learn-go/internal/service"
	"learn-go/pkg/logger"
	"learn-go/pkg/middleware"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Application wires up services and transports.
type Application struct {
	cfg    config.AppConfig
	engine *gin.Engine
	db     *gorm.DB
	log    *logger.Logger
}

// New creates the application, preparing dependencies.
func New() (*Application, error) {
	cfg := config.Load()

	log := logger.New()

	gin.SetMode(gin.ReleaseMode)
	if cfg.Environment == "local" || cfg.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	}

	dialector, err := resolveDialector(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	accountRepo := gormrepo.NewAccountStore(db)
	teacherRepo := gormrepo.NewTeacherStore(db)
	studentRepo := gormrepo.NewStudentStore(db)
	departmentRepo := gormrepo.NewDepartmentStore(db)
	classRepo := gormrepo.NewClassStore(db)
	teacherStudentRepo := gormrepo.NewTeacherStudentStore(db)
	assignmentRepo := gormrepo.NewAssignmentStore(db)
	submissionRepo := gormrepo.NewSubmissionStore(db)
	submissionCommentRepo := gormrepo.NewSubmissionCommentStore(db)
	noteRepo := gormrepo.NewNoteStore(db)
	noteCommentRepo := gormrepo.NewNoteCommentStore(db)
	conversationRepo := gormrepo.NewConversationStore(db)
	messageRepo := gormrepo.NewMessageStore(db)
	receiptRepo := gormrepo.NewMessageReceiptStore(db)

	authService := service.NewAuthService(accountRepo, cfg)
	adminService := service.NewAdminService(accountRepo, teacherRepo, studentRepo, departmentRepo, classRepo, teacherStudentRepo)
	assignmentService := service.NewAssignmentService(assignmentRepo, submissionRepo, submissionCommentRepo)
	conversationService := service.NewConversationService(conversationRepo, messageRepo, receiptRepo, accountRepo)
	noteService := service.NewNoteService(noteRepo, accountRepo)
	noteCommentService := service.NewNoteCommentService(noteRepo, noteCommentRepo, accountRepo)
	wsHub := ws.NewHub()

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	handler := apihandlers.NewHandler(authService, adminService, assignmentService, conversationService, noteService, noteCommentService, wsHub)

	adminGuard := middleware.JWTAuth(middleware.AuthConfig{Secret: cfg.JWTSecret, AllowedRoles: []string{string(domain.RoleAdmin)}})
	teacherGuard := middleware.JWTAuth(middleware.AuthConfig{Secret: cfg.JWTSecret, AllowedRoles: []string{string(domain.RoleTeacher), string(domain.RoleAdmin)}})
	studentGuard := middleware.JWTAuth(middleware.AuthConfig{Secret: cfg.JWTSecret, AllowedRoles: []string{string(domain.RoleStudent), string(domain.RoleTeacher), string(domain.RoleAdmin)}})

	handler.RegisterRoutes(engine, adminGuard, teacherGuard, studentGuard)

	return &Application{cfg: cfg, engine: engine, db: db, log: log}, nil
}

// Run starts the HTTP server.
func (a *Application) Run() error {
	address := fmt.Sprintf(":%s", a.cfg.HTTPPort)
	a.log.Printf("starting http server on %s", address)
	return a.engine.Run(address)
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.School{},
		&domain.Account{},
		&domain.Teacher{},
		&domain.Student{},
		&domain.TeacherStudentLink{},
		&domain.Department{},
		&domain.Class{},
		&domain.Course{},
		&domain.CourseSlot{},
		&domain.CourseSession{},
		&domain.Assignment{},
		&domain.AssignmentQuestion{},
		&domain.AssignmentSubmission{},
		&domain.SubmissionItem{},
		&domain.SubmissionComment{},
		&domain.Conversation{},
		&domain.ConversationMember{},
		&domain.Message{},
		&domain.MessageReceipt{},
		&domain.Note{},
		&domain.NoteComment{},
	)
}

func resolveDialector(cfg config.AppConfig) (gorm.Dialector, error) {
	driver := strings.ToLower(cfg.DatabaseDriver)
	switch driver {
	case "postgres", "postgresql":
		return postgres.Open(cfg.DatabaseDSN), nil
	case "sqlite", "sqlite3", "":
		return sqlite.Open(cfg.DatabaseDSN), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DatabaseDriver)
	}
}
