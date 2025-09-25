package handlers

import (
	"net/http"
	"strings"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/usecase"
)

type ForecastHandler struct {
	uc usecase.ForecastUsecase
}

func NewForecastHandler(uc usecase.ForecastUsecase) *ForecastHandler {
	return &ForecastHandler{uc: uc}
}

func (h *ForecastHandler) Get(w http.ResponseWriter, r *http.Request) {
	dateStr := strings.TrimSpace(r.URL.Query().Get("date"))
	if dateStr == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[dto.ForecastData]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "missing 'date' (YYYY-MM-DD)"},
		})
		return
	}
	d, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[dto.ForecastData]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "invalid date format (want YYYY-MM-DD)"},
		})
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))

	data, err := h.uc.Get(r.Context(), d, q)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[dto.ForecastData]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse[dto.ForecastData]{
		Data:  &data,
		Meta:  nil,
		Error: nil,
	})
}
