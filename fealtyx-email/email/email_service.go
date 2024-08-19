package email

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"io"
	"os"
)

type EmailData struct {
	Content string
}

type SESConfig struct {
	Sender          string
	Receiver        string
	Subject         string
	TemplatePath    string
	Region          string
	AccessKeyId     string
	SecretAccessKey string
}

type EmailService struct {
	config SESConfig
}

func NewEmailService(config SESConfig) *EmailService {
	return &EmailService{config: config}
}

func (s *EmailService) GenerateEmailTemplate(templatePath string, data *EmailData) (string, error) {
	var templateBuffer bytes.Buffer

	file, err := os.Open(templatePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	htmlData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	htmlTemplate := template.Must(template.New("email").Parse(string(htmlData)))

	if data != nil {
		err = htmlTemplate.Execute(&templateBuffer, data)
	} else {
		err = htmlTemplate.Execute(&templateBuffer, nil)
	}
	if err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}

func (s *EmailService) SendEmail(data *EmailData) error {
	html, err := s.GenerateEmailTemplate(s.config.TemplatePath, data)
	if err != nil {
		return fmt.Errorf("failed to generate email template: %w", err)
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.config.Region),
		Credentials: credentials.NewStaticCredentials(s.config.AccessKeyId, s.config.SecretAccessKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}
	service := ses.New(sess)
	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(s.config.Receiver),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(s.config.Subject),
			},
		},
		Source: aws.String(s.config.Sender),
	}
	_, err = service.SendEmail(emailInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return fmt.Errorf("aws error: %w", aerr)
		}
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
