package handler

import (
	"encoding/json"
	"golang-crowdfunding/helper"
	"golang-crowdfunding/transaction"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(s transaction.Service) *transactionHandler {
	return &transactionHandler{s}
}

func (h *transactionHandler) GetCampaignTransactions(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	campaignId, _ := strconv.Atoi(vars["campaign_id"])
	userId, _ := strconv.Atoi(w.Header().Get("currentUser"))
	transactions, err := h.service.GetTransactionsByCampaignId(campaignId, userId)

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		response := helper.ApiResponse("Invalid parameters", http.StatusUnprocessableEntity, "failed", err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	

	response := helper.ApiResponse("Transactions data", http.StatusOK, "success", transaction.FormatTransactionArr(transactions))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}