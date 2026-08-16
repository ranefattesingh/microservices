package handler

import (
	"net/http"
	"strconv"

	"github.com/ranefattesingh/microservices/user/handler/dto"
	"github.com/ranefattesingh/microservices/user/json"
	"go.uber.org/zap"
)

// GetAllUsers returns all users in a paged-style payload with total count.
// @Summary List users
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} json.Error
// @Router /v1/users/ [get]
func (h *userHandle) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := h.service.GetAllUsers(r.Context(), page, limit)
	if err != nil {
		h.logger.Error("get all users: fail", zap.Error(err))
		handleError(w, err)

		return
	}

	result := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, dto.FromServiceModel(user))
	}

	json.Respond(w).JSON(map[string]any{
		"page":        page,
		"limit":       limit,
		"total_count": total,
		"count":       len(result),
		"users":       result,
	})
}
