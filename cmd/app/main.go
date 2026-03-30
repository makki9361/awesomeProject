package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"awesomeProject/internal/config"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middleware"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load config")
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	ruleRepo := repository.NewRuleRepository(db)

	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	ruleService := service.NewRuleService(ruleRepo)

	userHandler := handlers.NewUserHandler(userService, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryService, logger)
	ruleHandler := handlers.NewRuleHandler(ruleService, logger)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/login", userHandler.Login).Methods("GET", "POST")

	protected := apiRouter.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(userRepo, logger))
	protected.Use(middleware.LoggingMiddleware(logger))

	protected.HandleFunc("/users", middleware.RequireAdmin(userHandler.ListUsers)).Methods("GET")
	protected.HandleFunc("/users", middleware.RequireAdmin(userHandler.CreateUser)).Methods("POST")
	protected.HandleFunc("/users/{id:[0-9]+}", middleware.RequireAdmin(userHandler.GetUser)).Methods("GET")
	protected.HandleFunc("/users/{id:[0-9]+}", middleware.RequireAdmin(userHandler.UpdateUser)).Methods("PUT")
	protected.HandleFunc("/users/{id:[0-9]+}", middleware.RequireAdmin(userHandler.DeleteUser)).Methods("DELETE")
	protected.HandleFunc("/users/me", userHandler.GetMe).Methods("GET")

	protected.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	protected.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.GetCategory).Methods("GET")
	protected.HandleFunc("/categories", middleware.RequireAdmin(categoryHandler.CreateCategory)).Methods("POST")
	protected.HandleFunc("/categories/{id:[0-9]+}", middleware.RequireAdmin(categoryHandler.UpdateCategory)).Methods("PUT")
	protected.HandleFunc("/categories/{id:[0-9]+}", middleware.RequireAdmin(categoryHandler.DeleteCategory)).Methods("DELETE")

	protected.HandleFunc("/rules", ruleHandler.ListRules).Methods("GET")
	protected.HandleFunc("/rules/{id:[0-9]+}", ruleHandler.GetRule).Methods("GET")
	protected.HandleFunc("/rules", middleware.RequireAdmin(ruleHandler.CreateRule)).Methods("POST")
	protected.HandleFunc("/rules", middleware.RequireAdmin(ruleHandler.CreateRule)).Methods("POST")
	protected.HandleFunc("/rules/{id:[0-9]+}", middleware.RequireAdmin(ruleHandler.UpdateRule)).Methods("PUT")
	protected.HandleFunc("/rules/{id:[0-9]+}", middleware.RequireAdmin(ruleHandler.DeleteRule)).Methods("DELETE")
	protected.HandleFunc("/rules/{id:[0-9]+}/publish", middleware.RequireAdmin(ruleHandler.PublishRule)).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	workDir, _ := os.Getwd()
	frontendDir := filepath.Join(workDir, "frontend")

	logger.WithField("frontend_dir", frontendDir).Info("Looking for frontend files")

	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		logger.WithField("path", frontendDir).Warn("Frontend directory not found")
	} else {
		logger.WithField("path", frontendDir).Info("Serving frontend from")
		fs := http.FileServer(http.Dir(frontendDir))
		router.PathPrefix("/").Handler(fs)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.WithField("addr", serverAddr).Info("Starting server")

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logger.WithError(err).Fatal("Server failed")
	}
}
