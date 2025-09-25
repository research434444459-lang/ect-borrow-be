package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"
	
)

type AdminRequestsHandler struct {
	uc usecase.AdminRequestsUsecase
}

func NewAdminRequestsHandler(repo repository.Repository) *AdminRequestsHandler {
	return &AdminRequestsHandler{uc: usecase.NewAdminRequestsUsecase(repo)}
}

// GET /api/admin/requests?q=...&student=...&date=YYYY-MM-DD&status=...&page=1&pageSize=50
func (h *AdminRequestsHandler) GetAdminRequests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	qstr := strings.TrimSpace(q.Get("q"))
	student := strings.TrimSpace(q.Get("student"))
	date := strings.TrimSpace(q.Get("date"))
	status := strings.TrimSpace(q.Get("status"))

	parseInt := func(s string, def int) int {
		if s == "" {
			return def
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return def
		}
		return n
	}
	page := parseInt(q.Get("page"), 1)
	pageSize := parseInt(q.Get("pageSize"), 50)

	items, total, err := h.uc.SearchAdminRequests(r.Context(), qstr, student, date, status, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	data := dto.AdminRequestListData{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	meta := dto.AdminRequestListMeta{
		Timezone:   "Asia/Bangkok",
		ServerTime: time.Now().In(loc).Format(time.RFC3339),
	}

	writeJSON(w, http.StatusOK, dto.APIResponse[dto.AdminRequestListData]{
		Data:  &data,
		Meta:  meta,
		Error: nil,
	})
}
