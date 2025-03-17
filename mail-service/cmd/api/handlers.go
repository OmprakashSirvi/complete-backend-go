package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/tsawler/toolbox"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	logger := appLogger.WithFields(logrus.Fields{
		"url": r.URL.Path,
		"method": r.Method,
	})
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	
	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		logger.Errorf("error while reading request payload to mailMessage struct: %v", err)
		tools.ErrorJSON(w, err)
		return
	}
	
	logger.Tracef("constructed mail message: %v", requestPayload)
	msg := Message {
		From: requestPayload.From,
		To: requestPayload.To,
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg, logger)
	if err != nil {
		logger.Errorf("error while sending SMTP message: %v", err)
		tools.ErrorJSON(w, err)
		return
	}

	payload := toolbox.JSONResponse {
		Error: false,
		Message: "sent to " + requestPayload.To,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
