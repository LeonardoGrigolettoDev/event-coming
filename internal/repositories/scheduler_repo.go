package repositories

import (
	"context"
	"time"

	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulerRepository struct{ db *pgxpool.Pool }

func NewSchedulerRepository(db *pgxpool.Pool) *SchedulerRepository {
	return &SchedulerRepository{db: db}
}

func (r *SchedulerRepository) Create(ctx context.Context, s *models.Scheduler) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := r.db.Exec(ctx, `INSERT INTO schedulers (id,name,description,schedule_type,cron_expression,start_date,end_date,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, s.ID, s.Name, s.Description, s.ScheduleType, s.CronExpression, s.StartDate, s.EndDate, s.CreatedAt, s.UpdatedAt)
	return err
}

func (r *SchedulerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Scheduler, error) {
	row := r.db.QueryRow(ctx, `SELECT id,name,description,schedule_type,cron_expression,start_date,end_date,created_at,updated_at FROM schedulers WHERE id=$1`, id)
	var s models.Scheduler
	if err := row.Scan(&s.ID, &s.Name, &s.Description, &s.ScheduleType, &s.CronExpression, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SchedulerRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Scheduler, error) {
	where, args := buildWhere(filters, 1)
	q := `SELECT id,name,description,schedule_type,cron_expression,start_date,end_date,created_at,updated_at FROM schedulers ` + where + ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]models.Scheduler, 0)
	for rows.Next() {
		var s models.Scheduler
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.ScheduleType, &s.CronExpression, &s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
