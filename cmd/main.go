package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"golang-rbac/internal/config"
	_handler "golang-rbac/internal/delivery/http"
	"golang-rbac/internal/delivery/http/middleware"
	_repo "golang-rbac/internal/repository"
	_usecase "golang-rbac/internal/usecase"
)

func main() {
	// 0. Load konfigurasi dari .env / environment variable
	cfg := config.Load()

	// 1. Inisialisasi Database MySQL
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("gagal konek ke database: %v", err)
	}

	// 2. Dependency Injection
	userRepo := _repo.NewMysqlUserRepository(db)
	userUsecase := _usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	userHandler := _handler.NewUserHandler(userUsecase)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// 3. Routing & RBAC Definition
	mux := http.NewServeMux()

	// Public Endpoints
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/register", userHandler.Register)

	// Protected Endpoints dengan Tingkat Role Berbeda
	adminPage := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dashboard Admin: Rahasia Negara"))
	})
	editorPage := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halaman Artikel: Admin dan Editor bisa masuk"))
	})
	viewerPage := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halaman Publik: Semua role bisa lihat"))
	})

	// Terapkan Middleware Role
	mux.Handle("/admin", authMiddleware.RoleMiddleware("Admin")(adminPage))
	mux.Handle("/edit-article", authMiddleware.RoleMiddleware("Admin", "Editor")(editorPage))
	mux.Handle("/view", authMiddleware.RoleMiddleware("Admin", "Editor", "Viewer")(viewerPage))

	// Endpoint baru: profil user yang sedang login, contoh pemakaian JWT claims dari context
	mux.Handle("/me", authMiddleware.RoleMiddleware("Admin", "Editor", "Viewer")(http.HandlerFunc(userHandler.Profile)))

	// Run Server
	fmt.Printf("Server berjalan di port :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
