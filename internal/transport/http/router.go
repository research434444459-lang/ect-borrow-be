// internal/transport/http/router.go
package transport

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/repository"
	"ect-borrow-be/internal/transport/http/handlers"
	"ect-borrow-be/internal/usecase"

	"google.golang.org/api/sheets/v4"
)

// ให้ repo ที่มี Google Sheets client สามารถถูก assert ได้
type sheetsProvider interface {
	Service() *sheets.Service
}

func NewRouter(repo repository.Repository, cfg config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(corsMiddleware(cfg))

	// ===== routes เดิม =====
	devices := handlers.NewDeviceHandler(repo, cfg)
	r.Get("/api/devices", devices.GetDevices)

	requests := handlers.NewRequestHandler(repo)
	r.Get("/api/requests", requests.GetRequests)

	reqDetail := handlers.NewRequestDetailHandler(repo)
	// 2025-09-10
	r.Get("/api/requests/{studentId}/{date}", reqDetail.GetRequestDetail)
	// 2025/09/10
	r.Get("/api/requests/{studentId}/*", reqDetail.GetRequestDetail)

	adminToday := handlers.NewAdminTodayHandler(repo)
	r.Get("/api/admin/today", adminToday.GetAdminToday)

	adminReqList := handlers.NewAdminRequestsHandler(repo)
	r.Get("/api/admin/requests", adminReqList.GetAdminRequests)

	adminReqDetail := handlers.NewAdminRequestDetailHandler(repo)
	r.Get("/api/admin/requests/{requestId}", adminReqDetail.Get)
	r.Patch("/api/admin/requests/{requestId}", adminReqDetail.Patch)

	// ===== ใหม่: Forecast =====
	if sp, ok := any(repo).(sheetsProvider); ok && sp.Service() != nil {
		svc := sp.Service()
		frepo := repository.NewForecastRepository(svc, cfg)
		fuc := usecase.NewForecastUsecase(frepo)
		fh := handlers.NewForecastHandler(fuc)
		r.Get("/api/forecast", fh.Get)
	} else {
		// ลงทะเบียน route ไว้เสมอเพื่อเลี่ยง 404 และส่ง error JSON ที่อ่านออก
		r.Get("/api/forecast", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"data":null,"meta":null,"error":{"code":"FORECAST_DISABLED","message":"Sheets service not available or repo.Service() missing"}}`))
		})
	}

	// health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return r
	
}
func corsMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	// รองรับหลาย origin ได้: แยกด้วยจุลภาคใน ENV ถ้าต้องการ
	// เช่น FRONTEND_ORIGIN="http://localhost:4200,https://app.example.com"
	allowed := map[string]struct{}{}
	for _, o := range strings.Split(cfg.FrontendOrigin, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = struct{}{}
		}
	}
	// กันพลาด: เผื่อค่า default dev
	if len(allowed) == 0 {
		allowed["http://localhost:4200"] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; ok {
				// ถ้าจะใช้ credential (cookies/authorization) ต้อง “สะท้อน origin ตรงๆ”
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin") // ให้ cache-aware
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			}

			// preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
