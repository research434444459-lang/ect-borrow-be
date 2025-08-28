package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"
)

type DeviceHandler struct {
	uc usecase.DeviceUsecase
}

func NewDeviceHandler(repo repository.Repository, cfg config.Config) *DeviceHandler {
	return &DeviceHandler{uc: usecase.NewDeviceUsecase(repo, cfg)}
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	category := strings.TrimSpace(q.Get("category"))
	search := strings.TrimSpace(q.Get("q"))

	if category == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "category is required"},
		})
		return
	}

	data, err := h.uc.SearchDevices(r.Context(), category, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	writeJSON(w, http.StatusOK, dto.APIResponse[dto.DevicesData]{
		Data:  &data,
		Meta:  nil, // keep top-level meta null per spec
		Error: nil,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
