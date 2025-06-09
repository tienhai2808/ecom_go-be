package smtp

type SMTPService interface {
	SendEmail(to, subject, htmlBody string) error
}