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

type client struct {
	logger logger.Logger
}

func New(logger logger.Logger) *client {
	return &client{
		logger: logger,
	}
}

func (c *client) SendEmail(_ context.Context, options MailOptions) error {
	c.logger.Debugf("Sending email with options: %+v", options)

	return nil
}
