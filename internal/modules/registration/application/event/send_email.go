package event

import (
	"context"
	"fmt"

	"github.com/KyKyPy3/clean/internal/application/core"
	"github.com/KyKyPy3/clean/internal/domain/common"
	"github.com/KyKyPy3/clean/internal/modules/registration/application/ports"
	"github.com/KyKyPy3/clean/pkg/logger"
)

const SendEmailKind = "SendEmail"

type SendEmailCommand struct {
	ID    string
	Email common.Email
}

func (c SendEmailCommand) Type() core.CommandType {
	return SendEmailKind
}

var _ core.Command = (*SendEmailCommand)(nil)

type SendEmail struct {
	sender ports.EmailSender
	logger logger.Logger
}

func NewSendEmail(logger logger.Logger, sender ports.EmailSender) *SendEmail {
	return &SendEmail{
		sender: sender,
		logger: logger,
	}
}

func (s *SendEmail) Handle(ctx context.Context, command core.Command) error {
	sendCommand, ok := command.(SendEmailCommand)
	if !ok {
		return fmt.Errorf("command type %s: %w", command.Type(), core.ErrUnexpectedCommand)
	}

	return s.sender.Send(ctx, sendCommand.Email, "Registration email", "Test body")
}
