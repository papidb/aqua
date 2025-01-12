package notification

import (
	context "context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	mu            sync.Mutex
	notifications map[string][]*Notification // UserID -> Notifications
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: make(map[string][]*Notification),
	}
}

func (s *NotificationService) AddNotification(userID, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	notification := &Notification{
		Id:        uuid.New().String(),
		UserId:    userID,
		Message:   message,
		CreatedAt: time.Now().String(),
	}
	s.notifications[userID] = append(s.notifications[userID], notification)
}

func (s *NotificationService) GetNotifications(userID string) []*Notification {
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

// NotificationServiceImpl extends NotificationService and implements gRPC methods.
type NotificationServiceImpl struct {
	*NotificationService
	UnimplementedNotificationServiceServer
}

// AddNotification (gRPC) adds a notification for a user.
func (s *NotificationServiceImpl) AddNotification(ctx context.Context, req *AddNotificationRequest) (*Empty, error) {
	s.NotificationService.AddNotification(req.UserId, req.Message)
	log.Printf("Added notification for user %s: %s", req.UserId, req.Message)
	return &Empty{}, nil
}

// GetAllNotifications (gRPC) retrieves all notifications for a user.
func (s *NotificationServiceImpl) GetAllNotifications(ctx context.Context, req *GetAllNotificationsRequest) (*NotificationList, error) {
	notifications := s.NotificationService.GetNotifications(req.UserId)
	return &NotificationList{Notifications: notifications}, nil
}

// ClearNotification (gRPC) clears a specific notification by ID.
func (s *NotificationServiceImpl) ClearNotification(ctx context.Context, req *ClearNotificationRequest) (*Empty, error) {
	s.NotificationService.ClearNotification(req.UserId, req.NotificationId)
	return &Empty{}, nil
}

// ClearAllNotifications (gRPC) clears all notifications for a user.
func (s *NotificationServiceImpl) ClearAllNotifications(ctx context.Context, req *ClearAllNotificationsRequest) (*Empty, error) {
	s.NotificationService.ClearAllNotifications(req.UserId)
	log.Printf("Cleared all notifications for user %s", req.UserId)
	return &Empty{}, nil
}
