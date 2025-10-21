package repository

import (
	"context"
	"time"

	"learn-go/internal/domain"
)

// AccountRepository defines persistence for accounts.
type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByIdentifier(ctx context.Context, schoolID, identifier string) (*domain.Account, error)
	FindByID(ctx context.Context, id string) (*domain.Account, error)
	ListByRole(ctx context.Context, schoolID string, role domain.Role, page, size int) ([]domain.Account, int64, error)
}

// TeacherRepository handles teacher profile persistence.
type TeacherRepository interface {
	Create(ctx context.Context, teacher *domain.Teacher) error
	GetByNumber(ctx context.Context, schoolID, number string) (*domain.Teacher, error)
	GetByID(ctx context.Context, id string) (*domain.Teacher, error)
}

// StudentRepository handles student profile persistence.
type StudentRepository interface {
	Create(ctx context.Context, student *domain.Student) error
	GetByNumber(ctx context.Context, schoolID, number string) (*domain.Student, error)
	GetByID(ctx context.Context, id string) (*domain.Student, error)
}

// TeacherStudentRepository manages relationships between teachers and students.
type TeacherStudentRepository interface {
	BindTeachers(ctx context.Context, studentID string, teacherIDs []string) error
}

// DepartmentRepository handles departments.
type DepartmentRepository interface {
	Create(ctx context.Context, department *domain.Department) error
	List(ctx context.Context, schoolID string) ([]domain.Department, error)
}

// ClassRepository handles classes.
type ClassRepository interface {
	Create(ctx context.Context, class *domain.Class) error
	ListByDepartment(ctx context.Context, schoolID, departmentID string) ([]domain.Class, error)
	GetByID(ctx context.Context, id string) (*domain.Class, error)
}

// AssignmentRepository handles assignments.
type AssignmentRepository interface {
	Create(ctx context.Context, assignment *domain.Assignment, questions []domain.AssignmentQuestion) error
	Get(ctx context.Context, id string) (*domain.Assignment, []domain.AssignmentQuestion, error)
}

// SubmissionRepository handles student submissions.
type SubmissionRepository interface {
	CreateOrUpdate(ctx context.Context, submission *domain.AssignmentSubmission, items []domain.SubmissionItem) error
	ListByAssignment(ctx context.Context, assignmentID string) ([]domain.AssignmentSubmission, error)
	ListItemsBySubmissionIDs(ctx context.Context, submissionIDs []string) ([]domain.SubmissionItem, error)
	GetByAssignmentAndStudent(ctx context.Context, assignmentID, studentID string) (*domain.AssignmentSubmission, []domain.SubmissionItem, error)
	GetByID(ctx context.Context, submissionID string) (*domain.AssignmentSubmission, []domain.SubmissionItem, error)
	UpdateGrades(ctx context.Context, submission *domain.AssignmentSubmission, items []domain.SubmissionItem) error
}

// SubmissionCommentRepository handles submission review comments.
type SubmissionCommentRepository interface {
	Create(ctx context.Context, comment *domain.SubmissionComment) error
	ListBySubmission(ctx context.Context, submissionID string) ([]domain.SubmissionComment, error)
}

// NoteRepository handles note persistence.
type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	Update(ctx context.Context, note *domain.Note) error
	FindByID(ctx context.Context, id string) (*domain.Note, error)
	ListByOwner(ctx context.Context, ownerID string, includeDeleted bool, status string) ([]domain.Note, error)
	ListPublishedBySchool(ctx context.Context, schoolID string) ([]domain.Note, error)
	SoftDelete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
}

// NoteCommentRepository handles note comment persistence.
type NoteCommentRepository interface {
	Create(ctx context.Context, comment *domain.NoteComment) error
	ListByNote(ctx context.Context, noteID string) ([]domain.NoteComment, error)
}

// ConversationRepository handles conversation persistence and membership.
type ConversationRepository interface {
	Create(ctx context.Context, conversation *domain.Conversation, members []domain.ConversationMember) error
	GetByID(ctx context.Context, id string) (*domain.Conversation, error)
	ListByAccount(ctx context.Context, accountID string) ([]domain.Conversation, error)
	GetMembers(ctx context.Context, conversationID string) ([]domain.ConversationMember, error)
	IsMember(ctx context.Context, conversationID, accountID string) (bool, error)
	FindDirectBetween(ctx context.Context, schoolID string, participantIDs [2]string) (*domain.Conversation, error)
	UpdateTimestamp(ctx context.Context, conversationID string, ts time.Time) error
}

// MessageRepository handles chat messages.
type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	ListByConversation(ctx context.Context, conversationID string, limit int, beforeID string) ([]domain.Message, error)
	GetLastByConversation(ctx context.Context, conversationID string) (*domain.Message, error)
	GetByID(ctx context.Context, id string) (*domain.Message, error)
}

// MessageReceiptRepository records read state for messages.
type MessageReceiptRepository interface {
	CreateBatch(ctx context.Context, receipts []domain.MessageReceipt) error
	CountUnread(ctx context.Context, accountID, conversationID string) (int64, error)
	MarkReadUpTo(ctx context.Context, accountID, conversationID string, ts time.Time) error
}
