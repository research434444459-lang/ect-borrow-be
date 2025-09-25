package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"
)

type RequestDetailHandler struct {
	uc usecase.RequestDetailUsecase
}

func NewRequestDetailHandler(repo repository.Repository) *RequestDetailHandler {
	return &RequestDetailHandler{uc: usecase.NewRequestDetailUsecase(repo)}
}

// รองรับทั้ง /api/requests/{studentId}/{date} และ /api/requests/{studentId}/* (เช่น 2025/09/12)
func (h *RequestDetailHandler) GetRequestDetail(w http.ResponseWriter, r *http.Request) {
	studentId := chi.URLParam(r, "studentId")
	date := chi.URLParam(r, "date")
	if date == "" {
		date = chi.URLParam(r, "*") // จะได้ "2025/09/12"
	}
	date = strings.TrimPrefix(date, "/")

	if studentId == "" || date == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data: nil, Meta: nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "studentId and date are required"},
		})
		return
	}

	detail, err := h.uc.GetRequestDetail(r.Context(), studentId, date)
	if err != nil {
		if err == usecase.ErrNotFound {
			writeJSON(w, http.StatusNotFound, dto.APIResponse[any]{
				Data: nil, Meta: nil,
				Error: &dto.APIError{Code: "NOT_FOUND", Message: "request not found"},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data: nil, Meta: nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	meta := dto.ListMeta{UpdatedAt: time.Now().Format(time.RFC3339)}
	writeJSON(w, http.StatusOK, dto.APIResponse[dto.RequestDetail]{Data: &detail, Meta: meta, Error: nil})
}
