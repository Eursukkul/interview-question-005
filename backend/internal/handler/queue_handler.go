package handler

import (
	"errors"

	"example.com/interview-question-005/backend/internal/model"
	"example.com/interview-question-005/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type QueueHandler struct {
	service *service.QueueService
}

func NewQueueHandler(service *service.QueueService) *QueueHandler {
	return &QueueHandler{service: service}
}

func (h *QueueHandler) Current(c *fiber.Ctx) error {
	state, err := h.service.Current(c.Context())
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "failed to get current queue")
	}
	return c.Status(fiber.StatusOK).JSON(model.QueueResponse{
		CurrentQueue: state.CurrentQueue,
		IssuedAt:     state.IssuedAt,
	})
}

func (h *QueueHandler) Next(c *fiber.Ctx) error {
	state, err := h.service.Next(c.Context())
	if err != nil {
		if errors.Is(err, service.ErrQueueLimitReached) {
			return writeError(c, fiber.StatusConflict, service.ErrQueueLimitReached.Error())
		}
		if errors.Is(err, service.ErrInvalidQueueState) {
			return writeError(c, fiber.StatusUnprocessableEntity, err.Error())
		}
		return writeError(c, fiber.StatusInternalServerError, "failed to generate queue")
	}
	return c.Status(fiber.StatusOK).JSON(model.QueueResponse{
		QueueNumber: state.CurrentQueue,
		IssuedAt:    state.IssuedAt,
	})
}

func (h *QueueHandler) Reset(c *fiber.Ctx) error {
	state, err := h.service.Reset(c.Context())
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "failed to reset queue")
	}
	return c.Status(fiber.StatusOK).JSON(model.QueueResponse{
		QueueNumber: state.CurrentQueue,
		IssuedAt:    state.IssuedAt,
	})
}

func Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func writeError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(model.ErrorResponse{Error: message})
}
