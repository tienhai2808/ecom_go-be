package smtp

import (
	"backend/internal/config"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type SMTPSender struct {
	auth     smtp.Auth
	host     string
	port     string
	from     string
	appName  string
	template *template.Template
}

func NewSMTPSender(cfg *config.AppConfig) EmailSender {
	auth := smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Host)

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		fmt.Printf("🚨 Lỗi load template email: %v\n", err)
	}

	return &SMTPSender{
		auth:     auth,
		host:     cfg.SMTP.Host,
		port:     cfg.SMTP.Port,
		from:     cfg.SMTP.User,
		appName:  cfg.App.Name,
		template: tmpl,
	}
}

func (s *SMTPSender) SendEmail(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	data := struct {
		Subject string
		Body    template.HTML
		AppName string
	}{
		Subject: subject,
		Body:    template.HTML(htmlBody),
		AppName: s.appName,
	}

	var buf bytes.Buffer
	if err := s.template.Execute(&buf, data); err != nil {
		return fmt.Errorf("🚨 Lỗi render template email: %v", err)
	}

	msg := buildHTMLMessage(s.from, to, subject, buf.String())

	if err := smtp.SendMail(addr, s.auth, s.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("🚨 Lỗi gửi email tới %s: %v", to, err)
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
