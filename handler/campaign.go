package handler

import (
	"encoding/json"
	"fmt"
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/helper"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)


type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.FormValue("userId"))

	w.Header().Set("Content-Type", "application/json")


	campaigns, err := h.service.GetCampaigns(userId)
	if err != nil {
		response := helper.ApiResponse("Invalid parameters", http.StatusBadRequest, "failed", nil )
		log.Info("Can't get campaigns")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := helper.ApiResponse("Campaigns data", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	json.NewEncoder(w).Encode(response)
}

func(h *campaignHandler) GetCampaignById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId, _ := strconv.Atoi(vars["campaign_id"])
	campaignDetail, err := h.service.GetCampaignById(campaignId)
	if err != nil {
		response := helper.ApiResponse("Invalid parameters", http.StatusBadRequest, "failed", nil )
		log.Info("Can't get campaign from the given ID")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := helper.ApiResponse("Campaigns data", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	json.NewEncoder(w).Encode(response)

}

func (h *campaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var input campaign.CreateCampaignInput
	log.Info("Create campaign called")
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

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Invalid request body")
		return
	}

	formatter := campaign.FormatCampaign(newCampaign)
	response := helper.ApiResponse("User created successfully", http.StatusOK, "success", formatter)
	json.NewEncoder(w).Encode(response)
}

func (h *campaignHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	campaignId, _ := strconv.Atoi(vars["campaign_id"])

	userId, _ := strconv.Atoi(w.Header().Get("currentUser"))
	getCampaign, err := h.service.GetCampaignById(campaignId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("No campaign with that id found")
		return
	}

	fmt.Println(getCampaign)
	if getCampaign.UserId != userId {
		response := helper.ApiResponse("You are not authorized to perform this action", http.StatusForbidden, "failed", nil)
		json.NewEncoder(w).Encode(response)
		return
	}

	var input campaign.CreateCampaignInput
	log.Info("Update campaign called")
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Invalid request body")
		return
	}
	json.Unmarshal(reqBody, &input)

	updatedCampaign, err := h.service.Update(campaignId, input)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Can't save changes, please try again later")
		return
	}

	formatter := campaign.FormatCampaign(updatedCampaign)
	response := helper.ApiResponse("Campaign updated", http.StatusOK, "success", formatter)
	json.NewEncoder(w).Encode(response)
}