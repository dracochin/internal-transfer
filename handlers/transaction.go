package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"internal-transfer/models"
	"internal-transfer/models/mapper"
	"internal-transfer/repository"
	"internal-transfer/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: repository.NewRepository(db),
	}
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	createAccountServiceInput, err := mapper.CreateAccountRequestToCreateAccountInput(acc)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
	}
	err = h.repo.CreateAccount(*createAccountServiceInput)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid account ID")
		return
	}

	acc, err := h.repo.GetAccount(id)
	if err != nil {
		utils.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSONResponse(w, acc)
}

func (h *Handler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var t models.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if t.IdempotencyKey == "" {
		utils.JSONError(w, http.StatusBadRequest, "idempotency_key is required")
		return
	}

	tx, err := h.repo.DB.Begin()
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	exists, err := h.repo.TransactionExists(tx, t.IdempotencyKey)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if exists {
		// Idempotent: transaction already processed
		w.WriteHeader(http.StatusOK)
		return
	}
	input, err := mapper.CreateTransactionRequestToCreateTransactionInput(t)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.repo.CreateTransactionWithEntries(tx, *input)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
