package controllers

import (
	"bytes"
	"encoding/json"
	models "go-gorilla-mongo/cmd/api/models/schema"
	"go-gorilla-mongo/cmd/api/utils"

	// models "go-gorilla-mongo/cmd/api/models/schema"
	"net/http"
)

type Request struct {
	Data json.RawMessage
}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Register(w http.ResponseWriter, r *http.Request) {
	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	var user models.User
	err = json.NewDecoder(bytes.NewReader(request.Data)).Decode(&user)
}
