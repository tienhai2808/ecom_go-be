package types

import "html/template"

type EmailTemplateData struct {
	Subject string
	Body    template.HTML
	AppName string
}
