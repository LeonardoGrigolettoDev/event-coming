package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct{ db *pgxpool.Pool }

func NewEventRepository(db *pgxpool.Pool) *EventRepository { return &EventRepository{db: db} }

func (r *EventRepository) Create(ctx context.Context, e *models.Event) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now
	var payload []byte
	if e.Payload != nil {
		b, _ := json.Marshal(e.Payload)
		payload = b
	}
	_, err := r.db.Exec(ctx, `INSERT INTO events (id,scheduler_id,contact_id,event_type,status,payload,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, e.ID, e.SchedulerID, e.ContactID, e.EventType, e.Status, payload, e.CreatedAt, e.UpdatedAt)
	return err
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	row := r.db.QueryRow(ctx, `SELECT id,scheduler_id,contact_id,event_type,status,payload,created_at,updated_at FROM events WHERE id=$1`, id)
	var e models.Event
	var payload []byte
	if err := row.Scan(&e.ID, &e.SchedulerID, &e.ContactID, &e.EventType, &e.Status, &payload, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, err
	}
	if payload != nil {
		var m map[string]interface{}
		_ = json.Unmarshal(payload, &m)
		e.Payload = m
	}
	return &e, nil
}

func (r *EventRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Event, error) {
	where, args := buildWhere(filters, 1)
	q := `SELECT id,scheduler_id,contact_id,event_type,status,payload,created_at,updated_at FROM events ` + where + ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Event
	for rows.Next() {
		var e models.Event
		var payload []byte
		if err := rows.Scan(&e.ID, &e.SchedulerID, &e.ContactID, &e.EventType, &e.Status, &payload, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		if payload != nil {
			var m map[string]interface{}
			_ = json.Unmarshal(payload, &m)
			e.Payload = m
		}
		out = append(out, e)
	}
	return out, nil
}

func (r *EventRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE events SET status=$1,updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id)
	return err
}
