package repositories

import (
	"context"
	"encoding/json"
	"time"

	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConsolidatedRepository struct{ db *pgxpool.Pool }

func NewConsolidatedRepository(db *pgxpool.Pool) *ConsolidatedRepository {
	return &ConsolidatedRepository{db: db}
}

func (r *ConsolidatedRepository) Create(ctx context.Context, c *models.Consolidated) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.UpdatedAt = time.Now().UTC()
	var timeline []byte
	var payload []byte
	if c.Timeline != nil {
		b, _ := json.Marshal(c.Timeline)
		timeline = b
	}
	if c.PayloadFinal != nil {
		b, _ := json.Marshal(c.PayloadFinal)
		payload = b
	}
	_, err := r.db.Exec(ctx, `INSERT INTO consolidated (id,scheduler_id,contact_id,final_status,timeline,payload_final,started_at,finished_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, c.ID, c.SchedulerID, c.ContactID, c.FinalStatus, timeline, payload, c.StartedAt, c.FinishedAt, c.UpdatedAt)
	return err
}

func (r *ConsolidatedRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Consolidated, error) {
	row := r.db.QueryRow(ctx, `SELECT id,scheduler_id,contact_id,final_status,timeline,payload_final,started_at,finished_at,updated_at FROM consolidated WHERE id=$1`, id)
	var c models.Consolidated
	var timeline []byte
	var payload []byte
	if err := row.Scan(&c.ID, &c.SchedulerID, &c.ContactID, &c.FinalStatus, &timeline, &payload, &c.StartedAt, &c.FinishedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	if timeline != nil {
		var m map[string]interface{}
		_ = json.Unmarshal(timeline, &m)
		c.Timeline = m
	}
	if payload != nil {
		var m map[string]interface{}
		_ = json.Unmarshal(payload, &m)
		c.PayloadFinal = m
	}
	return &c, nil
}

func (r *ConsolidatedRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Consolidated, error) {
	where, args := buildWhere(filters, 1)
	q := `SELECT id,scheduler_id,contact_id,final_status,timeline,payload_final,started_at,finished_at,updated_at FROM consolidated ` + where + ` ORDER BY updated_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Consolidated
	for rows.Next() {
		var c models.Consolidated
		var timeline []byte
		var payload []byte
		if err := rows.Scan(&c.ID, &c.SchedulerID, &c.ContactID, &c.FinalStatus, &timeline, &payload, &c.StartedAt, &c.FinishedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if timeline != nil {
			var m map[string]interface{}
			_ = json.Unmarshal(timeline, &m)
			c.Timeline = m
		}
		if payload != nil {
			var m map[string]interface{}
			_ = json.Unmarshal(payload, &m)
			c.PayloadFinal = m
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *ConsolidatedRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE consolidated SET final_status=$1,updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id)
	return err
}
