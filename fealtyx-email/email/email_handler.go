package email

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	emailService *EmailService
}

func NewHandler(emailService *EmailService) *Handler {
	return &Handler{
		emailService: emailService,
	}
}

func (h *Handler) Init(r *mux.Router) {
	r.HandleFunc("/send-email", h.SendEmail).Methods("POST")
}

func (h *Handler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	data := &EmailData{Content: req.Content}
	err := h.emailService.SendEmail(data)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	resp := SendEmailResponse{Message: "Email sent successfully!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
}
