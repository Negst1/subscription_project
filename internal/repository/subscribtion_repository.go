package repository

import (
	"context"
	"fmt"
	"strings"
	"subscribe_project/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, update *models.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	GetSummary(ctx context.Context, req models.SummaryRequest) (int, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, service_name, price, user_id, 
			start_date, end_date, created_at, updated_at
		)
		VALUES (
			:id, :service_name, :price, :user_id, 
			:start_date, :end_date, :created_at, :updated_at
		)`

	sub.ID = uuid.New()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, sub)
	return err
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1`
	err := r.db.GetContext(ctx, &sub, query, id)
	return &sub, err
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, update *models.UpdateSubscriptionRequest) error {
	query := "UPDATE subscriptions SET updated_at = $1"
	args := []interface{}{time.Now()}
	argIndex := 2

	if update.ServiceName != nil {
		query += fmt.Sprintf(", service_name = $%d", argIndex)
		args = append(args, *update.ServiceName)
		argIndex++
	}

	if update.Price != nil {
		query += fmt.Sprintf(", price = $%d", argIndex)
		args = append(args, *update.Price)
		argIndex++
	}

	if update.EndDate != nil {
		if *update.EndDate == "" {
			query += fmt.Sprintf(", end_date = $%d", argIndex)
			args = append(args, nil)
		} else {
			endDate, _ := time.Parse("01-2006", *update.EndDate)
			query += fmt.Sprintf(", end_date = $%d", argIndex)
			args = append(args, endDate)
		}
		argIndex++
	}

	query += " WHERE id = $" + fmt.Sprint(argIndex)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *subscriptionRepo) List(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	query := `SELECT * FROM subscriptions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &subscriptions, query, limit, offset)
	return subscriptions, err
}

func (r *subscriptionRepo) GetSummary(ctx context.Context, req models.SummaryRequest) (int, error) {
	startDate, _ := time.Parse("01-2006", req.StartDate)
	endDate, _ := time.Parse("01-2006", req.EndDate)

	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
	          WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`
	args := []interface{}{endDate, startDate}

	conditions := []string{}

	if req.UserID != nil {
		userID, _ := uuid.Parse(*req.UserID)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, userID)
	}

	if req.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(args)+1))
		args = append(args, *req.ServiceName)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var totalCost int
	err := r.db.GetContext(ctx, &totalCost, query, args...)
	return totalCost, err
}
