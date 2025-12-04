package repositories

import (
	"context"
	"time"

	"event-coming/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(ctx context.Context, c *models.Contact) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.Exec(ctx, `INSERT INTO contacts (id,name,phone_number,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)`, c.ID, c.Name, c.PhoneNumber, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *ContactRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	row := r.db.QueryRow(ctx, `SELECT id,name,phone_number,created_at,updated_at FROM contacts WHERE id=$1`, id)
	var c models.Contact
	if err := row.Scan(&c.ID, &c.Name, &c.PhoneNumber, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContactRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Contact, error) {
	where, args := buildWhere(filters, 1)
	q := `SELECT id,name,phone_number,created_at,updated_at FROM contacts ` + where + ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Contact
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.Name, &c.PhoneNumber, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}
