package handlers

import (
	"net/http"
	"strings"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"
)

type AdminTodayHandler struct {
	uc usecase.AdminTodayUsecase
}

func NewAdminTodayHandler(repo repository.Repository) *AdminTodayHandler {
	return &AdminTodayHandler{uc: usecase.NewAdminTodayUsecase(repo)}
}

// GET /api/admin/today?date=YYYY-MM-DD
func (h *AdminTodayHandler) GetAdminToday(w http.ResponseWriter, r *http.Request) {
	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if date == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data: nil, Meta: nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "date is required (YYYY-MM-DD)"},
		})
		return
	}

	data, meta, err := h.uc.GetAdminToday(r.Context(), date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data: nil, Meta: nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse[dto.AdminTodayData]{
		Data: &data, Meta: meta, Error: nil,
	})
}
