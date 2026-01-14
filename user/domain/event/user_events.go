package event

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// UserCreatedEvent 用户创建事件
type UserCreatedEvent struct {
	UserID     uuid.UUID
	Username   string
	Email      string
	occurredAt time.Time
}

func NewUserCreatedEvent(userID uuid.UUID, username, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		UserID:     userID,
		Username:   username,
		Email:      email,
		occurredAt: time.Now(),
	}
}

func (e *UserCreatedEvent) EventName() string {
	return "user.created"
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// UserActivatedEvent 用户激活事件
type UserActivatedEvent struct {
	UserID     uuid.UUID
	occurredAt time.Time
}

func NewUserActivatedEvent(userID uuid.UUID) *UserActivatedEvent {
	return &UserActivatedEvent{
		UserID:     userID,
		occurredAt: time.Now(),
	}
}

func (e *UserActivatedEvent) EventName() string {
	return "user.activated"
}

func (e *UserActivatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// UserPasswordChangedEvent 用户密码修改事件
type UserPasswordChangedEvent struct {
	UserID     uuid.UUID
	occurredAt time.Time
}

func NewUserPasswordChangedEvent(userID uuid.UUID) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		UserID:     userID,
		occurredAt: time.Now(),
	}
}

func (e *UserPasswordChangedEvent) EventName() string {
	return "user.password_changed"
}

func (e *UserPasswordChangedEvent) OccurredAt() time.Time {
	return e.occurredAt
}
