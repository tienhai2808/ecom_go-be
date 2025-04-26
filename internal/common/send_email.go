package common

import (
	"bytes"
	"backend/internal/config"
	"fmt"
	"html/template"
	"net/smtp"
)

type EmailSender interface {
	SendEmail(to, subject, htmlBody string) error
}

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

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>{{.Subject}}</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f4f4f4;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 8px;">
        <h2 style="color: #333;">{{.AppName}}</h2>
        <h3>{{.Subject}}</h3>
        <p>{{.Body}}</p>
        <p style="color: #777;">Email này được gửi từ {{.AppName}}. Vui lòng không trả lời trực tiếp.</p>
    </div>
</body>
</html>
`
