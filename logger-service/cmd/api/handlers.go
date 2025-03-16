package main

import (
	"log"
	"logger-service/data"
	"net/http"

	"github.com/tsawler/toolbox"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, req *http.Request) {
	var requestPayload JsonPayload

	log.Print("Reading JSON")
	err := tools.ReadJSON(w, req, &requestPayload)
	log.Printf("request payload: %v", requestPayload)

	if err != nil {
		tools.ErrorJSON(w, err, http.StatusBadRequest)
		log.Printf("Unable to read req body. Err : %v", err)
		return
	}

	event := data.LogEntry{Name: requestPayload.Name, Data: requestPayload.Data}
	log.Printf("created new entry: %v", event)
	err = app.Models.LogEntry.Insert(event)

	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = tools.WriteJSON(w, http.StatusCreated, toolbox.JSONResponse{Error: false, Message: "logged"})

	if err != nil {
		tools.ErrorJSON(w, err, http.StatusInternalServerError)
	}

}
