package pushbullet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"net/http"
	"time"
)

const (
	tokEnv = "NOTI_PUSHBULLET_TOK"

	// API is the Pushbullet API endpoint.
	API = "https://api.pushbullet.com/v2/pushes"
)

var (
	errNoTok = ConfigErrror{Env: tokEnv, Reason: "missing"}
)

type apiRequest struct {
	Body  string `json:"body"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type ApiResponse struct {
	Active                  bool    `json:"active"`
	Iden                    string  `json:"iden"`
	Created                 float64 `json:"created"`
	Modified                float64 `json:"modified"`
	Type                    string  `json:"type"`
	Dismissed               bool    `json:"dismissed"`
	Direction               string  `json:"direction"`
	SenderIden              string  `json:"sender_iden"`
	SenderEmail             string  `json:"sender_email"`
	SenderEmailNormalized   string  `json:"sender_email_normalized"`
	SenderName              string  `json:"sender_name"`
	ReceiverIden            string  `json:"receiver_iden"`
	ReceiverEmail           string  `json:"receiver_email"`
	ReceiverEmailNormalized string  `json:"receiver_email_normalized"`
	Title                   string  `json:"title"`
	Body                    string  `json:"body"`
	Error                   struct {
		Code    string `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
		Cat     string `json:"cat"`
	} `json:"error"`
	ErrorCode string `json:"error_code"`
}

type configuration struct {
	accessToken string
}

// Notify sends a push request to the Pushbullet API.
func Notify(n Params) error {
	if n.Token == "" {
		return errors.New("Token not found")
	}
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(apiRequest{
		Body:  n.Message,
		Title: n.Title,
		Type:  "note",
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", n.API, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Access-Token", n.Token)
	req.Header.Set("Content-Type", "application/json")

	webClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := webClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return fmt.Errorf("decoding response: %s", err)
	}

	if r.ErrorCode != "" {
		return APIError{Site: "pushbullet", Msg: r.Error.Message}
	}

	return nil
}
