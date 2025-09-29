package main

import (
	"ect-borrow-be/internal/config"
	"ect-borrow-be/internal/repository"
	transport "ect-borrow-be/internal/transport/http"
	"log"
	"net/http"
	"strings"
)

// ---- CORS middleware (วางนอก func main) ----
func allowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	// ✅ อนุญาตโดเมน FE ของคุณ (ปรับให้ตรงของจริง)
	if origin == "https://776c1340.ect-borrow.pages.dev" {
		return true
	}
	// ถ้ามี custom domain ก็เพิ่มได้
	if origin == "https://ect-inventory.com" || origin == "https://www.ect-inventory.com" {
		return true
	}
	// ✅ อนุญาต preview ของ Cloudflare Pages (โดเมนแบบสุ่ม)
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
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			// ถ้าไม่มีคุกกี้ ไม่ต้องใส่ Allow-Credentials
			// w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
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

	r := transport.NewRouter(repo, cfg) // r ต้องเป็น http.Handler

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)

	// ⬇️ ต้องเป็นแบบนี้ (ครอบด้วย withCORS)
	if err := http.ListenAndServe(addr, withCORS(r)); err != nil {
		log.Fatal(err)
	}
}
