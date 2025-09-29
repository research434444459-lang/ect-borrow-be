package main

import (
	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/repository"
	transport "ect-borrow-be/internal/transport/http"
	"log"
	"net/http"
	"strings"
)

// ---- CORS middleware ----
func allowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	// ✅ อนุญาตโดเมน FE ของคุณ (ปรับได้ตามจริง)
	if origin == "https://776c1340.ect-borrow.pages.dev" {
		return true
	}
	// ตัวอย่าง custom domain
	if origin == "https://ect-inventory.com" || origin == "https://www.ect-inventory.com" {
		return true
	}
	// ✅ อนุญาต Cloudflare Pages preview ทั้งหมด
	if strings.HasSuffix(origin, ".pages.dev") {
		return true
	}
	return false
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			// ⭐ เพิ่ม PATCH และจัดชุด header ที่ต้องใช้กับ Bearer token
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")
			// ถ้าใช้ cookie auth ค่อยเปิดบรรทัดนี้
			// w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// ตอบ preflight ให้จบที่นี่
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent) // 204
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()

	repo, err := repository.NewSheetsRepository(cfg)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	r := transport.NewRouter(repo, cfg) // r เป็น http.Handler

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)

	// ⬇️ ครอบ router ด้วย withCORS เสมอ
	if err := http.ListenAndServe(addr, withCORS(r)); err != nil {
		log.Fatal(err)
	}
}
