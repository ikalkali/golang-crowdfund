package handler

import (
	"encoding/json"
	"fmt"
	"golang-crowdfunding/helper"
	"golang-crowdfunding/transaction"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

func (h *transactionHandler) GetTransactionsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(w.Header().Get("currentUser"))
	transactions, err := h.service.GetTransactionsByUserId(userId)
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		response := helper.ApiResponse("Invalid parameters", http.StatusUnprocessableEntity, "failed", err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := helper.ApiResponse("Transactions data", http.StatusOK, "success", transaction.FormatUserTransactionArr(transactions))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *transactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input transaction.CreateTransactionInput
	fmt.Println("HANDLER CALLED")
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Invalid request body")
		return
	}
	json.Unmarshal(reqBody, &input)

	userId := w.Header().Get("currentUser")
	input.UserId, _ = strconv.Atoi(userId)

	errValidate := validatorCustom(input)
	if errValidate != nil {
		response := helper.ApiResponse("Missing required input for campaign creation", http.StatusForbidden, "failed", errValidate.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Invalid request body")
		return
	}

	response := helper.ApiResponse("User created successfully", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	json.NewEncoder(w).Encode(response)
}

func (h *transactionHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	var input transaction.TransactionNotificationInput
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Invalid request body")
		return
	}
	json.Unmarshal(reqBody, &input)

	err = h.service.ProcessPayment(input)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Can't process payment")
		return
	}

	json.NewEncoder(w).Encode(input)
}