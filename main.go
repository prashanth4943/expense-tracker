package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prashanth4943/expense-tracker/internal/db"
	"github.com/prashanth4943/expense-tracker/internal/handlers"
	"github.com/prashanth4943/expense-tracker/internal/middleware"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// --- Config from environment (Railway / Render friendly) ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "expenses.db"
	}

	store, err := db.New(dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	log.Printf("database ready: %s", dsn)

	h := handlers.NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("GET /expenses", h.ListExpenses)
	mux.HandleFunc("GET /categories", h.ListCategories)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve embedded Vue frontend for all other routes (SPA fallback)
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("embed frontend: %v", err)
	}
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/", spaHandler(fileServer, distFS))

	// --- Server ---
	handler := middleware.Logger(middleware.CORS(mux))

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// spaHandler serves static files and falls back to index.html for SPA routing.
func spaHandler(fileServer http.Handler, distFS fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if path == "" {
			path = "index.html"
		}

		_, err := distFS.Open(path)
		if err != nil {
			// index, err := frontendFS.ReadFile("frontend/dist/index.html")
			index, err := fs.ReadFile(distFS, "index.html")
			if err != nil {
				http.Error(w, "frontend not built", http.StatusInternalServerError)
				return
			}
			// w.Header().Set("Content-Type", "text/html")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(index)
			return
		}

		fileServer.ServeHTTP(w, r)
	})
}
