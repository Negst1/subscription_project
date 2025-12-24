package handlers

import (
	"strconv"

	"subscribe_project/internal/models"
	"subscribe_project/internal/services"
	"subscribe_project/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
	service services.SubscriptionService
}

func NewSubscriptionHandler(service services.SubscriptionService) *SubscriptionHandler {
	logger.Log.WithField("component", "subscription_handler").Info("Creating new subscription handler")
	return &SubscriptionHandler{service: service}
}

// CreateSubscription создает новую подписку
// @Summary Создать подписку
// @Description Создает новую подписку на сервис
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.CreateSubscriptionRequest true "Данные для создания подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	logger.Log.WithFields(logrus.Fields{
		"handler": "CreateSubscription",
		"method":  c.Method(),
		"path":    c.Path(),
		"ip":      c.IP(),
	}).Info("Received request to create subscription")

	var req models.CreateSubscriptionRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "CreateSubscription",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":      "CreateSubscription",
		"service_name": req.ServiceName,
		"user_id":      req.UserID,
		"price":        req.Price,
	}).Debug("Request body parsed successfully")

	subscription, err := h.service.CreateSubscription(c.Context(), req)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":        err.Error(),
			"handler":      "CreateSubscription",
			"service_name": req.ServiceName,
			"user_id":      req.UserID,
		}).Error("Service failed to create subscription")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":         "CreateSubscription",
		"subscription_id": subscription.ID.String(),
		"status_code":     fiber.StatusCreated,
	}).Info("Subscription created successfully, sending response")

	return c.Status(fiber.StatusCreated).JSON(subscription)
}

// GetSubscription получает подписку по ID
// @Summary Получить подписку
// @Description Возвращает информацию о подписке по её ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.WithFields(logrus.Fields{
		"handler": "GetSubscription",
		"method":  c.Method(),
		"path":    c.Path(),
		"id":      id,
		"ip":      c.IP(),
	}).Info("Received request to get subscription")

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "GetSubscription",
			"id":      id,
		}).Error("Invalid subscription ID format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid subscription ID",
		})
	}

	subscription, err := h.service.GetSubscription(c.Context(), id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "GetSubscription",
			"id":      id,
		}).Warn("Subscription not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Subscription not found",
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":         "GetSubscription",
		"subscription_id": id,
		"service_name":    subscription.ServiceName,
		"status_code":     fiber.StatusOK,
	}).Info("Subscription retrieved successfully, sending response")

	return c.JSON(subscription)
}

// UpdateSubscription обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет информацию о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param request body models.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.WithFields(logrus.Fields{
		"handler": "UpdateSubscription",
		"method":  c.Method(),
		"path":    c.Path(),
		"id":      id,
		"ip":      c.IP(),
	}).Info("Received request to update subscription")

	var req models.UpdateSubscriptionRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "UpdateSubscription",
			"id":      id,
		}).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":          "UpdateSubscription",
		"id":               id,
		"has_service_name": req.ServiceName != nil,
		"has_price":        req.Price != nil,
		"has_end_date":     req.EndDate != nil,
	}).Debug("Request body parsed successfully")

	if err := h.service.UpdateSubscription(c.Context(), id, req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "UpdateSubscription",
			"id":      id,
		}).Error("Service failed to update subscription")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":     "UpdateSubscription",
		"id":          id,
		"status_code": fiber.StatusOK,
	}).Info("Subscription updated successfully, sending response")

	return c.JSON(fiber.Map{"message": "Subscription updated successfully"})
}

// DeleteSubscription удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	id := c.Params("id")

	logger.Log.WithFields(logrus.Fields{
		"handler": "DeleteSubscription",
		"method":  c.Method(),
		"path":    c.Path(),
		"id":      id,
		"ip":      c.IP(),
	}).Info("Received request to delete subscription")

	if err := h.service.DeleteSubscription(c.Context(), id); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "DeleteSubscription",
			"id":      id,
		}).Error("Service failed to delete subscription")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":     "DeleteSubscription",
		"id":          id,
		"status_code": fiber.StatusOK,
	}).Info("Subscription deleted successfully, sending response")

	return c.JSON(fiber.Map{"message": "Subscription deleted successfully"})
}

// ListSubscriptions получает список подписок
// @Summary Список подписок
// @Description Возвращает список подписок с пагинацией
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *fiber.Ctx) error {
	logger.Log.WithFields(logrus.Fields{
		"handler": "ListSubscriptions",
		"method":  c.Method(),
		"path":    c.Path(),
		"ip":      c.IP(),
		"query":   c.OriginalURL(),
	}).Info("Received request to list subscriptions")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	logger.Log.WithFields(logrus.Fields{
		"handler": "ListSubscriptions",
		"page":    page,
		"limit":   limit,
	}).Debug("Query parameters parsed")

	subscriptions, err := h.service.ListSubscriptions(c.Context(), page, limit)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "ListSubscriptions",
			"page":    page,
			"limit":   limit,
		}).Error("Service failed to list subscriptions")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":     "ListSubscriptions",
		"count":       len(subscriptions),
		"status_code": fiber.StatusOK,
		"page":        page,
		"limit":       limit,
	}).Info("Subscriptions listed successfully, sending response")

	return c.JSON(subscriptions)
}

// GetSummary получает сводку по подпискам
// @Summary Сводка по подпискам
// @Description Возвращает общую стоимость подписок за период
// @Tags summary
// @Accept json
// @Produce json
// @Param request body models.SummaryRequest true "Параметры фильтрации"
// @Success 200 {object} models.SubscriptionSummary
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/summary [post]
func (h *SubscriptionHandler) GetSummary(c *fiber.Ctx) error {
	logger.Log.WithFields(logrus.Fields{
		"handler": "GetSummary",
		"method":  c.Method(),
		"path":    c.Path(),
		"ip":      c.IP(),
	}).Info("Received request to get summary")

	var req models.SummaryRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "GetSummary",
		}).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":          "GetSummary",
		"start_date":       req.StartDate,
		"end_date":         req.EndDate,
		"has_user_id":      req.UserID != nil,
		"has_service_name": req.ServiceName != nil,
	}).Debug("Request body parsed successfully")

	summary, err := h.service.GetSummary(c.Context(), req)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":      err.Error(),
			"handler":    "GetSummary",
			"start_date": req.StartDate,
			"end_date":   req.EndDate,
		}).Error("Service failed to calculate summary")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Log.WithFields(logrus.Fields{
		"handler":     "GetSummary",
		"total_cost":  summary.TotalCost,
		"status_code": fiber.StatusOK,
		"start_date":  req.StartDate,
		"end_date":    req.EndDate,
	}).Info("Summary calculated successfully, sending response")

	return c.JSON(summary)
}
