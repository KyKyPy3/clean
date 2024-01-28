package email

import (
	"github.com/KyKyPy3/clean/pkg/email"
	"github.com/KyKyPy3/clean/pkg/logger"
)

type Client struct {
	emailClient *email.Client
	logger      logger.Logger
}

func New(logger logger.Logger) *Client {
	emailClient := email.New(logger)

	return &Client{
		emailClient: emailClient,
	}
}
