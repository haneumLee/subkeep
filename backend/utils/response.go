package utils

import (
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginationMeta holds pagination metadata.
type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Success sends a 200 OK response with data.
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMessage sends a 200 OK response with data and a message.
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response with data.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Success: true,
		Message: "리소스가 생성되었습니다",
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNoContent)
}

// Paginated sends a 200 OK response with paginated data.
func Paginated(c *fiber.Ctx, data interface{}, page, perPage int, totalItems int64) error {
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))

	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			TotalItems: totalItems,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	})
}

// Error sends an error response using an AppError.
func Error(c *fiber.Ctx, err *AppError) error {
	return c.Status(err.Code).JSON(Response{
		Success: false,
		Error: fiber.Map{
			"code":    err.Code,
			"message": err.Message,
			"detail":  err.Detail,
		},
	})
}

// ErrorWithStatus sends an error response with a custom HTTP status code.
func ErrorWithStatus(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: fiber.Map{
			"code":    status,
			"message": message,
		},
	})
}
