package services

import (
	"context"
	"fmt"
	"subscribe_project/internal/models"
	"subscribe_project/internal/repository"
	"subscribe_project/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubscription(ctx context.Context, id string) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, req models.UpdateSubscriptionRequest) error
	DeleteSubscription(ctx context.Context, id string) error
	ListSubscriptions(ctx context.Context, page, limit int) ([]models.Subscription, error)
	GetSummary(ctx context.Context, req models.SummaryRequest) (*models.SubscriptionSummary, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	logger.Log.WithField("component", "subscription_service").Info("Creating new subscription service")
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req models.CreateSubscriptionRequest) (*models.Subscription, error) {
	logger.Log.WithFields(logrus.Fields{
		"method":       "CreateSubscription",
		"service_name": req.ServiceName,
		"user_id":      req.UserID,
		"price":        req.Price,
		"start_date":   req.StartDate,
	}).Info("Creating new subscription")

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": req.UserID,
			"method":  "CreateSubscription",
		}).Error("Invalid user_id format")
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":      err.Error(),
			"start_date": req.StartDate,
			"method":     "CreateSubscription",
		}).Error("Invalid start_date format")
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		ed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"error":    err.Error(),
				"end_date": *req.EndDate,
				"method":   "CreateSubscription",
			}).Error("Invalid end_date format")
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = &ed
		logger.Log.WithField("end_date", ed.Format("2006-01-02")).Debug("Parsed end date")
	}

	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	logger.Log.WithFields(logrus.Fields{
		"subscription_id": subscription.ID.String(),
		"start_date":      startDate.Format("2006-01-02"),
		"has_end_date":    endDate != nil,
		"method":          "CreateSubscription",
	}).Debug("Subscription object created")

	err = s.repo.Create(ctx, subscription)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":           err.Error(),
			"subscription_id": subscription.ID.String(),
			"method":          "CreateSubscription",
		}).Error("Failed to create subscription in repository")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"subscription_id": subscription.ID.String(),
		"service_name":    subscription.ServiceName,
		"price":           subscription.Price,
		"user_id":         subscription.UserID.String(),
		"method":          "CreateSubscription",
	}).Info("Subscription created successfully")

	return subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id string) (*models.Subscription, error) {
	logger.Log.WithFields(logrus.Fields{
		"method": "GetSubscription",
		"id":     id,
	}).Info("Getting subscription")

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "GetSubscription",
		}).Error("Invalid subscription id format")
		return nil, fmt.Errorf("invalid subscription id: %w", err)
	}

	subscription, err := s.repo.GetByID(ctx, subscriptionID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "GetSubscription",
		}).Error("Failed to get subscription from repository")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"id":           subscription.ID.String(),
		"service_name": subscription.ServiceName,
		"method":       "GetSubscription",
	}).Debug("Subscription retrieved successfully")

	return subscription, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id string, req models.UpdateSubscriptionRequest) error {
	logger.Log.WithFields(logrus.Fields{
		"method": "UpdateSubscription",
		"id":     id,
		"fields_to_update": map[string]interface{}{
			"service_name": req.ServiceName != nil,
			"price":        req.Price != nil,
			"end_date":     req.EndDate != nil,
		},
	}).Info("Updating subscription")

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "UpdateSubscription",
		}).Error("Invalid subscription id format")
		return fmt.Errorf("invalid subscription id: %w", err)
	}

	logger.Log.WithFields(logrus.Fields{
		"id":     id,
		"method": "UpdateSubscription",
	}).Debug("Attempting to update subscription in repository")

	err = s.repo.Update(ctx, subscriptionID, &req)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "UpdateSubscription",
		}).Error("Failed to update subscription in repository")
		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"id":     id,
		"method": "UpdateSubscription",
	}).Info("Subscription updated successfully")

	return nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	logger.Log.WithFields(logrus.Fields{
		"method": "DeleteSubscription",
		"id":     id,
	}).Info("Deleting subscription")

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "DeleteSubscription",
		}).Error("Invalid subscription id format")
		return fmt.Errorf("invalid subscription id: %w", err)
	}

	logger.Log.WithFields(logrus.Fields{
		"id":     id,
		"method": "DeleteSubscription",
	}).Debug("Attempting to delete subscription from repository")

	err = s.repo.Delete(ctx, subscriptionID)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"id":     id,
			"method": "DeleteSubscription",
		}).Error("Failed to delete subscription from repository")
		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"id":     id,
		"method": "DeleteSubscription",
	}).Info("Subscription deleted successfully")

	return nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, page, limit int) ([]models.Subscription, error) {
	logger.Log.WithFields(logrus.Fields{
		"method": "ListSubscriptions",
		"page":   page,
		"limit":  limit,
	}).Info("Listing subscriptions")

	if limit <= 0 {
		limit = 10
		logger.Log.WithField("new_limit", limit).Debug("Limit adjusted to default")
	}
	if page <= 0 {
		page = 1
		logger.Log.WithField("new_page", page).Debug("Page adjusted to default")
	}
	offset := (page - 1) * limit

	logger.Log.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
		"method": "ListSubscriptions",
	}).Debug("Fetching subscriptions from repository")

	subscriptions, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"page":   page,
			"limit":  limit,
			"method": "ListSubscriptions",
		}).Error("Failed to list subscriptions from repository")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"count":        len(subscriptions),
		"method":       "ListSubscriptions",
		"current_page": page,
		"per_page":     limit,
	}).Info("Subscriptions listed successfully")

	return subscriptions, nil
}

func (s *subscriptionService) GetSummary(ctx context.Context, req models.SummaryRequest) (*models.SubscriptionSummary, error) {
	logger.Log.WithFields(logrus.Fields{
		"method":     "GetSummary",
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"filters": map[string]interface{}{
			"has_user_id":      req.UserID != nil,
			"has_service_name": req.ServiceName != nil,
		},
	}).Info("Getting subscription summary")

	totalCost, err := s.repo.GetSummary(ctx, req)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":        err.Error(),
			"start_date":   req.StartDate,
			"end_date":     req.EndDate,
			"method":       "GetSummary",
			"user_id":      req.UserID,
			"service_name": req.ServiceName,
		}).Error("Failed to get subscription summary from repository")
		return nil, err
	}

	summary := &models.SubscriptionSummary{
		TotalCost: totalCost,
	}

	logger.Log.WithFields(logrus.Fields{
		"total_cost": totalCost,
		"method":     "GetSummary",
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
	}).Info("Subscription summary calculated successfully")

	return summary, nil
}
