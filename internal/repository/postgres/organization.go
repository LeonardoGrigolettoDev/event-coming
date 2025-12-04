package postgres

import (
	"context"
	"fmt"

	"event-coming/internal/domain"
	"event-coming/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type organizationRepository struct {
	db *pgxpool.Pool
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *pgxpool.Pool) repository.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, type, subscription_plan, max_events, max_participants, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		org.ID, org.Name, org.Type, org.SubscriptionPlan,
		org.MaxEvents, org.MaxParticipants, org.Active,
	).Scan(&org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `
		SELECT id, name, type, subscription_plan, max_events, max_participants, active, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org domain.Organization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Type, &org.SubscriptionPlan,
		&org.MaxEvents, &org.MaxParticipants, &org.Active,
		&org.CreatedAt, &org.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, id uuid.UUID, input *domain.UpdateOrganizationInput) error {
	query := `
		UPDATE organizations
		SET 
			name = COALESCE($2, name),
			subscription_plan = COALESCE($3, subscription_plan),
			max_events = COALESCE($4, max_events),
			max_participants = COALESCE($5, max_participants),
			active = COALESCE($6, active)
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id, input.Name, input.SubscriptionPlan, input.MaxEvents, input.MaxParticipants, input.Active)
	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

func (r *organizationRepository) List(ctx context.Context, page, perPage int) ([]*domain.Organization, int64, error) {
	offset := (page - 1) * perPage

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM organizations`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, name, type, subscription_plan, max_events, max_participants, active, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		var org domain.Organization
		err := rows.Scan(
			&org.ID, &org.Name, &org.Type, &org.SubscriptionPlan,
			&org.MaxEvents, &org.MaxParticipants, &org.Active,
			&org.CreatedAt, &org.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization: %w", err)
		}
		orgs = append(orgs, &org)
	}

	return orgs, total, nil
}
