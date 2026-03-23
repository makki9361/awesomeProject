package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"awesomeProject/internal/models"
	"awesomeProject/internal/service"
)

type RuleHandler struct {
	service *service.RuleService
	logger  *logrus.Logger
}

func NewRuleHandler(service *service.RuleService, logger *logrus.Logger) *RuleHandler {
	return &RuleHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RuleHandler) CreateRule(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRuleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rule, err := h.service.CreateRule(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create rule")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) GetRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	rule, err := h.service.GetRule(id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get rule")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rule == nil {
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

func (h *RuleHandler) UpdateRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRule(id, &req); err != nil {
		h.logger.WithError(err).Error("Failed to update rule")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rule updated successfully"})
}

func (h *RuleHandler) DeleteRule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid rule ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteRule(id); err != nil {
		h.logger.WithError(err).Error("Failed to delete rule")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rule deleted successfully"})
}

func (h *RuleHandler) ListRules(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")
	status := r.URL.Query().Get("status")
	createdByStr := r.URL.Query().Get("created_by")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	var categoryID *int
	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			categoryID = &id
		}
	}

	var createdBy *int
	if createdByStr != "" {
		id, err := strconv.Atoi(createdByStr)
		if err == nil {
			createdBy = &id
		}
	}

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	rules, err := h.service.ListRules(categoryID, &status, createdBy, page, pageSize)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list rules")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rules == nil {
		rules = []models.RuleWithCategory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}
