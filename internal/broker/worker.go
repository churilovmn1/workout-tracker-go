package broker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/churilovmn1/workout-tracker/internal/models"
)

// consumeTimeout bounds each blocking BRPOP so the consume loop can observe
// context cancellation promptly.
const consumeTimeout = 5 * time.Second

// reminderTick is how often the reminder scheduler scans for upcoming sessions.
const reminderTick = time.Minute

// Notifier delivers a plain-text message to a Telegram chat. *bot.Bot satisfies
// this. It is optional — a nil Notifier disables outbound Telegram messages.
type Notifier interface {
	Notify(chatID int64, text string) error
}

type userGetter interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type workoutGetter interface {
	GetByID(ctx context.Context, id int) (*models.Workout, error)
}

type scheduleLister interface {
	ListUpcoming(ctx context.Context, from, to time.Time) ([]models.ScheduleEntry, error)
}

// Worker drains the broker queue and runs the reminder scheduler.
type Worker struct {
	broker    *RedisBroker
	notifier  Notifier
	users     userGetter
	workouts  workoutGetter
	schedules scheduleLister

	mu       sync.Mutex
	reminded map[int]bool // schedule entry IDs already reminded this process lifetime
}

// NewWorker wires a worker. notifier may be nil (Telegram disabled).
func NewWorker(b *RedisBroker, notifier Notifier, users userGetter, workouts workoutGetter, schedules scheduleLister) *Worker {
	return &Worker{
		broker:    b,
		notifier:  notifier,
		users:     users,
		workouts:  workouts,
		schedules: schedules,
		reminded:  make(map[int]bool),
	}
}

// Run starts the consume loop and the reminder scheduler. It blocks until ctx
// is cancelled, so callers typically launch it in a goroutine.
func (w *Worker) Run(ctx context.Context) {
	log.Println("broker worker started")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); w.consumeLoop(ctx) }()
	go func() { defer wg.Done(); w.reminderLoop(ctx) }()
	wg.Wait()
	log.Println("broker worker stopped")
}

// consumeLoop drains events and dispatches them.
func (w *Worker) consumeLoop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		event, err := w.broker.Consume(ctx, consumeTimeout)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("broker consume error: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if event == nil {
			continue // timeout, no event
		}
		w.handle(ctx, event)
	}
}

// handle dispatches a single event to its handler.
func (w *Worker) handle(ctx context.Context, e *Event) {
	switch e.Type {
	case EventWorkoutCreated:
		log.Printf("event %s: workout #%d by user #%d", e.Type, e.Payload.WorkoutID, e.Payload.UserID)
	case EventWorkoutCommented:
		w.handleWorkoutCommented(ctx, e)
	case EventScheduleReminder:
		w.handleScheduleReminder(ctx, e)
	default:
		log.Printf("broker: unknown event type %q", e.Type)
	}
}

// handleWorkoutCommented notifies the workout owner that their trainer commented.
func (w *Worker) handleWorkoutCommented(ctx context.Context, e *Event) {
	workout, err := w.workouts.GetByID(ctx, e.Payload.WorkoutID)
	if err != nil {
		log.Printf("broker: load workout #%d: %v", e.Payload.WorkoutID, err)
		return
	}
	user, err := w.users.GetByID(ctx, workout.UserID)
	if err != nil {
		log.Printf("broker: load user #%d: %v", workout.UserID, err)
		return
	}
	w.notify(user, "💬 Тренер оставил комментарий к тренировке «"+workout.Title+"»:\n"+e.Payload.Comment)
}

// handleScheduleReminder reminds a client about an upcoming session.
func (w *Worker) handleScheduleReminder(ctx context.Context, e *Event) {
	user, err := w.users.GetByID(ctx, e.Payload.ClientID)
	if err != nil {
		log.Printf("broker: load client #%d: %v", e.Payload.ClientID, err)
		return
	}
	w.notify(user, "⏰ Напоминание: через час тренировка «"+e.Payload.Title+"» в "+e.Payload.ScheduledAt)
}

// notify sends a Telegram message if the user has a chat id and a notifier exists.
func (w *Worker) notify(user *models.User, text string) {
	if w.notifier == nil {
		log.Printf("broker: notifier disabled, dropping message to user #%d", user.ID)
		return
	}
	if user.TelegramChatID == nil {
		log.Printf("broker: user #%d has no telegram chat id, skipping", user.ID)
		return
	}
	if err := w.notifier.Notify(*user.TelegramChatID, text); err != nil {
		log.Printf("broker: notify user #%d: %v", user.ID, err)
	}
}

// reminderLoop scans every minute for sessions starting in ~1 hour and publishes
// a schedule.reminder event for each (deduplicated per process lifetime).
func (w *Worker) reminderLoop(ctx context.Context) {
	ticker := time.NewTicker(reminderTick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.scanReminders(ctx)
		}
	}
}

func (w *Worker) scanReminders(ctx context.Context) {
	// Window: sessions starting between 59 and 60 minutes from now. With a
	// 1-minute tick each session falls in the window exactly once.
	from := time.Now().Add(59 * time.Minute)
	to := time.Now().Add(60 * time.Minute)

	entries, err := w.schedules.ListUpcoming(ctx, from, to)
	if err != nil {
		log.Printf("broker: scan reminders: %v", err)
		return
	}

	for _, e := range entries {
		w.mu.Lock()
		already := w.reminded[e.ID]
		if !already {
			w.reminded[e.ID] = true
		}
		w.mu.Unlock()
		if already {
			continue
		}

		err := w.broker.Publish(ctx, NewEvent(EventScheduleReminder, Payload{
			ScheduleID:  e.ID,
			ClientID:    e.ClientID,
			Title:       e.Title,
			ScheduledAt: e.ScheduledAt.Format("15:04 02.01.2006"),
		}))
		if err != nil {
			log.Printf("broker: publish reminder for schedule #%d: %v", e.ID, err)
		}
	}
}
