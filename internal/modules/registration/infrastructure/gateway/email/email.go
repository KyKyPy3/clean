package email

import (
	"context"

	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/pkg/email"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type Client struct {
	emailClient *email.Client
	logger      logger.Logger
}

func New(emailClient *email.Client, logger logger.Logger) *Client {
	return &Client{
		emailClient: emailClient,
		logger:      logger,
	}
}

func (c *Client) Send(ctx context.Context, destination common.Email, subject, body string) error {
	c.logger.Debugf("Send email")

	return nil
}
