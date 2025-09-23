package smtp

import (
	"backend/config"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type smtpServiceImpl struct {
	auth     smtp.Auth
	template *template.Template
	cfg      *config.Config
}

func NewSMTPService(cfg *config.Config) SMTPService {
	auth := smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Host)

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Printf("Lỗi load template email: %v", err)
	}

	return &smtpServiceImpl{
		auth,
		tmpl,
		cfg,
	}
}

func (s *smtpServiceImpl) SendEmail(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.cfg.SMTP.Host, s.cfg.SMTP.Port)

	data := struct {
		Subject string
		Body    template.HTML
		AppName string
	}{
		Subject: subject,
		Body:    template.HTML(htmlBody),
		AppName: s.cfg.App.Name,
	}

	var buf bytes.Buffer
	if err := s.template.Execute(&buf, data); err != nil {
		return fmt.Errorf("lỗi render template email: %v", err)
	}

	msg := buildHTMLMessage(s.cfg.SMTP.User, to, subject, buf.String())

	if err := smtp.SendMail(addr, s.auth, s.cfg.SMTP.User, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("lỗi gửi email tới %s: %v", to, err)
	}

	return nil
}

func buildHTMLMessage(from, to, subject, body string) string {
	return fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset="UTF-8"

%s
`, from, to, subject, body)
}
