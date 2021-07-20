package handler

import (
	"encoding/json"
	"fmt"
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/helper"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

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

	errValidate := validatorCustom(input)
	if errValidate != nil {
		response := helper.ApiResponse("Missing required input for campaign creation", http.StatusForbidden, "failed", errValidate.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

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

func (h *campaignHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	var response helper.Response
	var isUploaded helper.UploadStatus
	var input campaign.CreateCampaignImageInput

	campaignId, _ := strconv.Atoi(r.FormValue("campaign_id"))
	input.CampaignId = campaignId
	input.IsPrimary, _ = strconv.Atoi(r.FormValue("is_primary"))

	errValidate := validatorCustom(input)
	if errValidate != nil {
		response = helper.ApiResponse("Missing required input for campaign image", http.StatusForbidden, "failed", errValidate.Error())
		json.NewEncoder(w).Encode(response)
		return
	}

	userId, _ := strconv.Atoi(w.Header().Get("currentUser"))
	campaignOwner, err := h.service.GetCampaignById(campaignId)
	if err != nil {
		isUploaded.IsUploaded = false
		response = helper.ApiResponse("Invalid campaign ID", http.StatusForbidden, "failed", nil)
		log.Info("Can't process file")
	}

	if userId != campaignOwner.UserId {
		isUploaded.IsUploaded = false
		response = helper.ApiResponse("Access unauthorized", http.StatusForbidden, "failed", nil)
		log.Info("Not authorized")
		json.NewEncoder(w).Encode(response)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("campaign-image")
	fileExtension := getFileExtension(handler)
	

	if err != nil {
		isUploaded.IsUploaded = false
		response = helper.ApiResponse("File uploading error", http.StatusForbidden, "failed", err.Error())
		log.Info("Can't process file")
	}
	defer file.Close()

	path := "images/"
	tempFile, err := ioutil.TempFile(path, fmt.Sprintf("%d-campaign-*.%s", campaignId, fileExtension))
	fileName := tempFile.Name()

	if err != nil {
        log.Info(err)
    }
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
        log.Info(err)
    }

	_, err = h.service.SaveCampaignImage(input, fileName)

	if err != nil {
        log.Info(err)
		response = helper.ApiResponse("File upload failed", http.StatusBadRequest, "failed", err.Error())
		json.NewEncoder(w).Encode(response)
		return
    } else {
		isUploaded.IsUploaded = true
		response = helper.ApiResponse("File successfully uploaded", http.StatusOK, "success", isUploaded)
	}
	tempFile.Write(fileBytes)

	json.NewEncoder(w).Encode(response)
}

func getFileExtension(fileHandler *multipart.FileHeader) string {
	filename := fileHandler.Filename
	getFileExtension := strings.Split(filename, ".")
	fileExtension := getFileExtension[len(getFileExtension) - 1]
	return fileExtension
}