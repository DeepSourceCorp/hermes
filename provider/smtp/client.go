package smtp

import (
	"net/smtp"
	"strings"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

type Client struct {
	HTTPClient provider.IHTTPClient
}

type SendMessageRequest struct {
	FromEmail    string   `json:"from_email"`
	FromPassword string   `json:"from_password"`
	ToEmail      []string `json:"to_email"`
	SMTPHost     string   `json:"smtp_host"`
	SMTPPort     string   `json:"smtp_port"`
	Subject      string   `json:"subject"`
	Message      string   `json:"message"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (*Client) SendMessage(request *SendMessageRequest) (interface{}, domain.IError) {
	from := request.FromEmail
	password := request.FromPassword
	to := request.ToEmail
	smtpHost := request.SMTPHost
	smtpPort := request.SMTPPort
	subject := request.Subject
	message := request.Message

	auth := smtp.PlainAuth("", from, password, smtpHost)

	toHeader := strings.Join(to, ",")

	smtpMessage := ""
	smtpMessage += "From: " + from + "\n"
	smtpMessage += "To: " + toHeader + "\n"
	smtpMessage += "Subject: " + subject + "\n\n"
	smtpMessage += message

	resp := new(SendMessageResponse)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(smtpMessage))
	if err != nil {
		return resp, errFailedSendTemporary("something went wrong while sending mail through SMTP")
	}

	resp.Ok = true
	return resp, nil
}
