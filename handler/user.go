package handler

import (
	"encoding/json"
	"golang-crowdfunding/user"
	"io/ioutil"
	"net/http"
)


type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input user.RegisterUserInput
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &input)

	w.Header().Set("Content-Type", "application/json")
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(newUser)
}