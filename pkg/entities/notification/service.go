package notification

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	mu            sync.Mutex
	notifications map[string][]Notification // UserID -> Notifications
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: make(map[string][]Notification),
	}
}

func (s *NotificationService) AddNotification(userID, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	notification := Notification{
		Id:        uuid.New().String(),
		UserId:    userID,
		Message:   message,
		CreatedAt: time.Now().String(),
	}
	s.notifications[userID] = append(s.notifications[userID], notification)
}

func (s *NotificationService) GetNotifications(userID string) []Notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.notifications[userID]
}

func (s *NotificationService) ClearNotification(userID, notificationID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	notifications := s.notifications[userID]
	for i, n := range notifications {
		if n.Id == notificationID {
			s.notifications[userID] = append(notifications[:i], notifications[i+1:]...)
			break
		}
	}
}

func (s *NotificationService) ClearAllNotifications(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.notifications, userID)
}
