package service

import (
	"errors"
	"fmt"

	"hinoob.net/learn-go/internal/model"
	"hinoob.net/learn-go/internal/pkg/websocket"
	"hinoob.net/learn-go/internal/repository"
)

// LiveService handles the business logic for live streaming.
type LiveService struct {
	Hub *websocket.Hub
}

func NewLiveService(hub *websocket.Hub) *LiveService {
	return &LiveService{Hub: hub}
}

// StartStreamForCourse starts a live stream for a specific course.
// It creates a chat room and notifies all students in the class.
func (s *LiveService) StartStreamForCourse(courseID uint, teacherID uint) error {
	// 1. Validate that the teacher is actually teaching this course
	course, err := repository.GetCourseByID(courseID) // This function needs to be created
	if err != nil {
		return errors.New("course not found")
	}
	if course.TeacherID != teacherID {
		return errors.New("you are not the teacher for this course")
	}

	// 2. Create the live room in the hub
	if _, ok := s.Hub.LiveRooms[courseID]; !ok {
		s.Hub.LiveRooms[courseID] = make(map[*websocket.Client]bool)
	}

	// 3. Notify students in the class
	class, err := repository.GetClassByID(course.ClassID)
	if err != nil {
		return errors.New("class not found for the course")
	}

	notification := &model.Message{
		// This is a "system" message, so we don't have a specific recipient.
		// The logic here is simplified. A real implementation would be more robust.
		Content:     fmt.Sprintf(`{"type": "live_start", "course_id": %d, "course_name": "%s"}`, courseID, course.SubjectName),
		MessageType: model.MsgTypeText,
	}

	for _, student := range class.Students {
		if client, ok := s.Hub.Clients[student.ID]; ok {
			// In a real system, you'd send a structured message.
			// Here, we send the raw content for simplicity.
			client.Send <- []byte(notification.Content)
		}
	}

	return nil
}

// EndStreamForCourse ends a live stream for a course and cleans up the chat room.
func (s *LiveService) EndStreamForCourse(courseID uint, teacherID uint) error {
	// Optional: Validate teacher ownership again

	// Delete the room from the hub, which stops all broadcasts and allows GC
	if _, ok := s.Hub.LiveRooms[courseID]; ok {
		delete(s.Hub.LiveRooms, courseID)
	}

	// Notify students that the stream has ended (optional)

	return nil
}
