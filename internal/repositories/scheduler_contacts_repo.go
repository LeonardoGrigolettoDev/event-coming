package repositories

import (
	"context"
	"time"

	"event-coming/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulerContactsRepository struct{ db *pgxpool.Pool }

func NewSchedulerContactsRepository(db *pgxpool.Pool) *SchedulerContactsRepository {
	return &SchedulerContactsRepository{db: db}
}

func (r *SchedulerContactsRepository) Create(ctx context.Context, sc *models.SchedulerContact) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	sc.CreatedAt = time.Now().UTC()
	_, err := r.db.Exec(ctx, `INSERT INTO scheduler_contacts (id,scheduler_id,contact_id,created_at) VALUES ($1,$2,$3,$4)`, sc.ID, sc.SchedulerID, sc.ContactID, sc.CreatedAt)
	return err
}

func (r *SchedulerContactsRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.SchedulerContact, error) {
	where, args := buildWhere(filters, 1)
	q := `SELECT id,scheduler_id,contact_id,created_at FROM scheduler_contacts ` + where + ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.SchedulerContact
	for rows.Next() {
		var sc models.SchedulerContact
		if err := rows.Scan(&sc.ID, &sc.SchedulerID, &sc.ContactID, &sc.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, nil
}
