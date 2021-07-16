package handler

import (
	"encoding/json"
	"golang-crowdfunding/helper"
	"golang-crowdfunding/user"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)


type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input user.RegisterUserInput
	
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Can't react database, please try again later")
		return
	}
	json.Unmarshal(reqBody, &input)

	errValidate := validatorCustom(input)
	if errValidate != nil {
		response := helper.ApiResponse("Invalid input", http.StatusUnprocessableEntity, "failed", errValidate)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Can't react database, please try again later")
		return
	}

	formatter := user.FormatUser(newUser, "halo")
	
	response := helper.ApiResponse("Account has been registerd", http.StatusOK, "success", formatter)
	
	
	json.NewEncoder(w).Encode(response)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input user.LoginInput

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Can't react database, please try again later")
		return
	}
	json.Unmarshal(reqBody, &input)

	loggedUser, err := h.userService.Login(input)
	if err != nil {
		response := helper.ApiResponse("Invalid email or password", http.StatusForbidden, "failed", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	formatter := user.FormatUser(loggedUser, "halo")
	response := helper.ApiResponse("Berhasil login", http.StatusOK, "success", formatter)
	json.NewEncoder(w).Encode(response)

}

func validatorCustom(input interface{}) error {
	var validate *validator.Validate = validator.New()
	errValidate := validate.Struct(input)
	return errValidate
}