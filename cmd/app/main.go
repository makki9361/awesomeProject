package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/middleware"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/service"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	router.Use(middleware.LoggingMiddleware(logger))

	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", userHandler.ListUsers).Methods("GET")

	router.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.GetCategory).Methods("GET")
	router.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id:[0-9]+}", categoryHandler.DeleteCategory).Methods("DELETE")
	router.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")

	router.HandleFunc("/rules", ruleHandler.CreateRule).Methods("POST")
	router.HandleFunc("/rules/{id:[0-9]+}", ruleHandler.GetRule).Methods("GET")
	router.HandleFunc("/rules/{id:[0-9]+}", ruleHandler.UpdateRule).Methods("PUT")
	router.HandleFunc("/rules/{id:[0-9]+}", ruleHandler.DeleteRule).Methods("DELETE")
	router.HandleFunc("/rules", ruleHandler.ListRules).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.WithField("addr", serverAddr).Info("Starting server")

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logger.WithError(err).Fatal("Server failed")
	}
}
