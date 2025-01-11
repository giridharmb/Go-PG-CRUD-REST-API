package main

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MetadataHandler struct {
	service *MetadataService
}

func NewMetadataHandler(service *MetadataService) *MetadataHandler {
	return &MetadataHandler{service: service}
}

// Response structures for consistent API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *MetadataHandler) Create(c *fiber.Ctx) error {
	var requestBody struct {
		MyKey   string         `json:"my_key"`
		MyValue map[string]any `json:"my_value"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: "Request body must be valid JSON with my_key and my_value fields",
		})
	}

	if requestBody.MyKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required field",
			Details: "my_key cannot be empty",
		})
	}

	if requestBody.MyValue == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required field",
			Details: "my_value cannot be empty",
		})
	}

	// Convert MyValue to JSON
	jsonValue, err := json.Marshal(requestBody.MyValue)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid JSON value",
			Details: "my_value must be a valid JSON object",
		})
	}

	entry := MetadataEntry{
		MyKey:   requestBody.MyKey,
		MyValue: jsonValue,
	}

	if err := h.service.Create(&entry); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error:   "Key already exists",
				Details: "Use PUT /api/metadata for updating existing entries",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to create entry",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{
		Message: "Entry created successfully",
		Data:    entry,
	})
}

func (h *MetadataHandler) Get(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Missing key parameter",
		})
	}

	entry, err := h.service.Get(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Entry not found",
				Details: "No metadata entry exists with the specified key",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to retrieve entry",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Data: entry,
	})
}

func (h *MetadataHandler) Update(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Missing key parameter",
		})
	}

	var requestBody struct {
		MyValue map[string]any `json:"my_value"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: "Request body must be valid JSON with my_value field",
		})
	}

	if requestBody.MyValue == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required field",
			Details: "my_value cannot be empty",
		})
	}

	// Convert MyValue to JSON
	jsonValue, err := json.Marshal(requestBody.MyValue)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid JSON value",
			Details: "my_value must be a valid JSON object",
		})
	}

	entry := MetadataEntry{
		MyKey:   key,
		MyValue: jsonValue,
	}

	if err := h.service.Update(&entry); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Entry not found",
				Details: "No metadata entry exists with the specified key",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to update entry",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "Entry updated successfully",
		Data:    entry,
	})
}

func (h *MetadataHandler) PatchUpdate(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Missing key parameter",
		})
	}

	var partialValue map[string]any
	if err := c.BodyParser(&partialValue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: "Request body must be a valid JSON object",
		})
	}

	if len(partialValue) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Empty update",
			Details: "No fields provided for update",
		})
	}

	if err := h.service.PatchUpdate(key, partialValue); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Entry not found",
				Details: "No metadata entry exists with the specified key",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to update entry",
			Details: err.Error(),
		})
	}

	// Fetch the updated entry to return in response
	updatedEntry, err := h.service.Get(key)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(SuccessResponse{
			Message: "Entry updated successfully",
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "Entry updated successfully",
		Data:    updatedEntry,
	})
}

func (h *MetadataHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Missing key parameter",
		})
	}

	if err := h.service.Delete(key); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error:   "Entry not found",
				Details: "The specified key does not exist or was already deleted",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete entry",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "Entry deleted successfully",
	})
}

func (h *MetadataHandler) DeleteAll(c *fiber.Ctx) error {
	if err := h.service.DeleteAll(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to delete all entries",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "All entries deleted successfully",
	})
}

func (h *MetadataHandler) Upsert(c *fiber.Ctx) error {
	var requestBody struct {
		MyKey   string         `json:"my_key"`
		MyValue map[string]any `json:"my_value"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request format",
			Details: "Request body must be valid JSON with my_key and my_value fields",
		})
	}

	if requestBody.MyKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required field",
			Details: "my_key cannot be empty",
		})
	}

	if requestBody.MyValue == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required field",
			Details: "my_value cannot be empty",
		})
	}

	// Convert MyValue to JSON
	jsonValue, err := json.Marshal(requestBody.MyValue)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid JSON value",
			Details: "my_value must be a valid JSON object",
		})
	}

	entry := MetadataEntry{
		MyKey:   requestBody.MyKey,
		MyValue: jsonValue,
	}

	if err := h.service.Upsert(&entry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to upsert entry",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Message: "Entry upserted successfully",
		Data:    entry,
	})
}
