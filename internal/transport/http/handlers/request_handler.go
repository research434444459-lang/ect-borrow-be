package handlers

import (
	"net/http"
	"strings"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"
)

type RequestHandler struct {
	uc usecase.RequestUsecase
}

func NewRequestHandler(repo repository.Repository) *RequestHandler {
	return &RequestHandler{uc: usecase.NewRequestUsecase(repo)}
}

func (h *RequestHandler) GetRequests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	student := strings.TrimSpace(q.Get("student"))
	date := strings.TrimSpace(q.Get("date"))

	items, total, err := h.uc.SearchRequests(r.Context(), student, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	meta := dto.ListMeta{Total: total}
	writeJSON(w, http.StatusOK, dto.APIResponse[[]dto.RequestItem]{
		Data:  &items,
		Meta:  meta,
		Error: nil,
	})
}
