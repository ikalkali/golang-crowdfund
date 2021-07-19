package handler

import (
	"encoding/json"
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/helper"
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