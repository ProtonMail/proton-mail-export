package app

import (
	"errors"

	"github.com/urfave/cli/v2"
)

type credentials struct {
	username     string
	password     []byte
	totp         string
	mboxPassword []byte
	attemptCount int
}

func newCredentialsFromCLI(ctx *cli.Context) *credentials {
	return &credentials{
		username:     ctx.String(flagUsername.Name),
		password:     []byte(ctx.String(flagPassword.Name)),
		totp:         ctx.String(flagTOTP.Name),
		mboxPassword: []byte(ctx.String(flagMBoxPassword.Name)),
	}
}

func (c *credentials) nextAttempt() error {
	if c.attemptCount++; c.attemptCount >= 5 {
		return errors.New("failed to login: too many attempts")
	}

	c.username = ""
	c.password = nil
	c.totp = ""
	c.mboxPassword = nil

	return nil
}
