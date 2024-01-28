package email

import (
	"context"

	"github.com/KyKyPy3/clean/pkg/logger"
)

type Options struct {
	ServerAddress string
}

type MailOptions struct {
	sender     string
	recipients string
	subject    string
	body       string
}

type Client struct {
	logger logger.Logger
}

func New(logger logger.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) SendEmail(_ context.Context, options MailOptions) error {
	c.logger.Debugf(
		"Sending email with options: sender: %s, recipients: %s, subject: %s, body: %s",
		options.sender,
		options.recipients,
		options.subject,
		options.body,
	)

	return nil
}
