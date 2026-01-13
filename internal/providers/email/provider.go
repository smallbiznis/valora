package email

import "context"

type Provider interface {
	Send(ctx context.Context, to []string, subject string, htmlBody string) error
	SendTemplate(ctx context.Context, to []string, templateName string, data interface{}) error
}

type NoOpProvider struct{}

func (p *NoOpProvider) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	return nil
}

func (p *NoOpProvider) SendTemplate(ctx context.Context, to []string, templateName string, data interface{}) error {
	return nil
}
