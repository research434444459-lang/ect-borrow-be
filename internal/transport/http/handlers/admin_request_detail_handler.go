package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"errors" // ⬅️ เพิ่ม
    "google.golang.org/api/googleapi" // ⬅️ เพิ่ม

	"ect-borrow-be/internal/dto"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type AdminRequestDetailHandler struct {
	uc usecase.AdminRequestDetailUsecase
}

func NewAdminRequestDetailHandler(repo repository.Repository) *AdminRequestDetailHandler {
	return &AdminRequestDetailHandler{uc: usecase.NewAdminRequestDetailUsecase(repo)}
}

// GET /api/admin/requests/{requestId}
func (h *AdminRequestDetailHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "requestId"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "requestId is required"},
		})
		return
	}

	detail, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if err == usecase.ErrNotFound {
			writeJSON(w, http.StatusNotFound, dto.APIResponse[any]{
				Data:  nil,
				Meta:  nil,
				Error: &dto.APIError{Code: "NOT_FOUND", Message: "request not found"},
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	meta := dto.AdminRequestDetailMeta{
		Timezone:   "Asia/Bangkok",
		ServerTime: time.Now().In(loc).Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, dto.APIResponse[dto.AdminRequestDetail]{
		Data:  &detail,
		Meta:  meta,
		Error: nil,
	})
}

// PATCH /api/admin/requests/{requestId}
func (h *AdminRequestDetailHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "requestId"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "requestId is required"},
		})
		return
	}

	var payload dto.AdminRequestUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.APIResponse[any]{
			Data:  nil,
			Meta:  nil,
			Error: &dto.APIError{Code: "BAD_REQUEST", Message: "invalid JSON body"},
		})
		return
	}

	detail, err := h.uc.UpdateByID(r.Context(), id, payload)
if err != nil {
    if err == usecase.ErrNotFound {
        writeJSON(w, http.StatusNotFound, dto.APIResponse[any]{
            Data:  nil,
            Meta:  nil,
            Error: &dto.APIError{Code: "NOT_FOUND", Message: "request not found"},
        })
        return
    }

    // ✅ ถ้าเป็น error จาก Google Sheets → ส่ง code และ message จริงกลับไป
    var gerr *googleapi.Error
    if errors.As(err, &gerr) && gerr != nil {
        writeJSON(w, gerr.Code, dto.APIResponse[any]{
            Data:  nil,
            Meta:  nil,
            Error: &dto.APIError{
                Code:    "GOOGLE_SHEETS_ERROR",
                Message: gerr.Message, // เช่น "The caller does not have permission" / "Insufficient authentication scopes" / ฯลฯ
            },
        })
        return
    }

    // อื่น ๆ ส่งข้อความจริงจาก server ออกไปเพื่อดีบัก
    writeJSON(w, http.StatusInternalServerError, dto.APIResponse[any]{
        Data:  nil,
        Meta:  nil,
        Error: &dto.APIError{Code: "INTERNAL_ERROR", Message: err.Error()},
    })
    return
}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	meta := dto.AdminRequestDetailMeta{
		Timezone:   "Asia/Bangkok",
		ServerTime: time.Now().In(loc).Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, dto.APIResponse[dto.AdminRequestDetail]{
		Data:  &detail,
		Meta:  meta,
		Error: nil,
	})
}
