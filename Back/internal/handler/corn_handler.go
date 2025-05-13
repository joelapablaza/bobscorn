package handler

import (
	"log"

	serviceinterfaces "bobscorn/internal/service/interfaces"

	"github.com/gofiber/fiber/v2"
)

type CornHandler struct {
	service serviceinterfaces.CornService
}

func NewCornHandler(service serviceinterfaces.CornService) *CornHandler {
	return &CornHandler{
		service: service,
	}
}

func (h *CornHandler) BuyCorn(c *fiber.Ctx) error {
	clientIP := c.IP()

	allowed, err := h.service.CanBuyCorn(c.Context(), clientIP)
	if err != nil {
		log.Printf("Error processing corn purchase for IP %s: %v", clientIP, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An internal error occurred."})
	}

	if !allowed {
		log.Printf("Rate limit exceeded for IP %s", clientIP)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "Too Many Requests. Please wait a minute."})
	}

	log.Printf("Corn sold to IP %s", clientIP)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "🌽 Corn successfully purchased!"})
}
