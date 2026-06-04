package broker

import "time"

// Event type identifiers published onto the broker queue.
const (
	// EventWorkoutCreated is emitted when a client logs a new workout.
	EventWorkoutCreated = "workout.created"
	// EventWorkoutCommented is emitted when a trainer leaves a comment on a workout.
	EventWorkoutCommented = "workout.commented"
	// EventScheduleReminder is emitted ~1 hour before a scheduled session starts.
	EventScheduleReminder = "schedule.reminder"
)

// Event is a single message placed on the broker queue. Only the Payload fields
// relevant to Type are populated.
type Event struct {
	Type      string    `json:"type"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// Payload carries the event-specific data.
type Payload struct {
	// workout.created / workout.commented
	WorkoutID int    `json:"workout_id,omitempty"`
	UserID    int    `json:"user_id,omitempty"` // workout owner / client
	Title     string `json:"title,omitempty"`
	Comment   string `json:"comment,omitempty"`

	// schedule.reminder
	ScheduleID  int    `json:"schedule_id,omitempty"`
	ClientID    int    `json:"client_id,omitempty"`
	ScheduledAt string `json:"scheduled_at,omitempty"` // "15:04 02.01.2006"
}

// NewEvent builds an event with the current timestamp.
func NewEvent(eventType string, payload Payload) Event {
	return Event{Type: eventType, Payload: payload, CreatedAt: time.Now()}
}
