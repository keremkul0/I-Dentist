package signUpClinic

import (
	"context"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"

	"github.com/gofiber/fiber/v2"
)

// SignUpClinicService defines the interface for clinic signup service
type SignUpClinicService interface {
	SignUpClinic(ctx context.Context, clinic clinic.Clinic, userCacheKey string) (clinic.Clinic, user.UserGetModel, error)
}

// SignUpClinicController handles clinic signup HTTP requests
type SignUpClinicController struct {
	service SignUpClinicService
}

// NewSignUpClinicController creates a new instance of SignUpClinicController
func NewSignUpClinicController(service SignUpClinicService) *SignUpClinicController {
	return &SignUpClinicController{service: service}
}

// signUpClinicRequest represents the incoming JSON request structure
type signUpClinicRequest struct {
	Clinic *clinic.Clinic `json:"clinic"`
	ID     *string        `json:"id"`
}

// signUpClinicResponse represents the outgoing JSON response structure
type signUpClinicResponse struct {
	Clinic clinic.Clinic     `json:"clinic"`
	User   user.UserGetModel `json:"user"`
}

// SignUpClinic handles the HTTP request for signing up a clinic
func (h *SignUpClinicController) SignUpClinic(c *fiber.Ctx) error {
	ctx := c.Context()
	var req signUpClinicRequest

	// Decode the incoming JSON request
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	// Validate required fields
	if req.Clinic == nil || req.ID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Clinic ID information are required",
		})
	}

	// Call the service layer to process the signup
	updatedClinic, userGetModel, err := h.service.SignUpClinic(ctx, *req.Clinic, *req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Clinic signup failed:",
		})
	}

	// Prepare the response
	response := signUpClinicResponse{
		Clinic: updatedClinic,
		User:   userGetModel,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
