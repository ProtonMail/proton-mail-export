package app

/*
#include <stdlib.h>
void etExportMailCallbackOnProgress() {}
void etCallOnRecover() {}
void etSessionCallbackOnNetworkLost() {}
void etSessionCallbackOnNetworkRestored() {}
*/
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/ProtonMail/export-tool/internal/sentry"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/gluon/async"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	flagUsername = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "username",
		Aliases: []string{"u"},
		EnvVars: []string{"ET_USER_EMAIL"},
	} //nolint:gochecknoglobals
	flagPassword = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "password",
		Aliases: []string{"p"},
		EnvVars: []string{"ET_USER_PASSWORD"},
	}
	flagMBoxPassword = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "mbox-password",
		Aliases: []string{"m"},
		EnvVars: []string{"ET_USER_MAILBOX_PASSWORD"},
	}
	flagTOTP = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "totp",
		Aliases: []string{"t"},
		EnvVars: []string{"ET_TOTP_CODE"},
	}
)

func Run() {
	app := &cli.App{
		Name:   "proton-mail-export-cli",
		Action: run,
		Flags: []cli.Flag{
			flagUsername,
			flagPassword,
			flagMBoxPassword,
			flagTOTP,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Fatal error")
	}
}

func run(ctx *cli.Context) error {
	printHeader()
	checkForNewVersion()

	panicHandler := sentry.NewPanicHandler(func() {})
	defer async.HandlePanic(panicHandler)

	s, err := newSession(panicHandler)
	if err != nil {
		return err
	}

	fmt.Printf("DownloadDir: %v\n", getDownloadDir())

	if err = login(ctx, s); err != nil {
		return err
	}

	logrus.Info("Login was successful")
	return nil
}

func printHeader() {
	fmt.Printf("Proton Mail Export Tool (%v) (c) Proton AG, Switzerland\n", internal.ETVersionString)
	fmt.Println("This program is licensed under the GNU General Public License v3")
	fmt.Println("Get support at https://proton.me/support/proton-mail-export-tool")
}

func checkForNewVersion() {
	if internal.HasNewVersion() {
		fmt.Println("A new version is available at: https://proton.me/support/proton-mail-export-tool")
	}
}

func newSession(panicHandler async.PanicHandler) (*session.Session, error) {
	sessionCb := CliCallback{}
	builder, err := apiclient.NewProtonAPIClientBuilder(internal.ETDefaultAPIURL, panicHandler, sessionCb)
	if err != nil {
		logrus.WithError(err).Fatal("Fatal error")
	}

	clientBuilder := apiclient.NewAutoRetryClientBuilder(
		builder,
		&apiclient.SleepRetryStrategyBuilder{},
	)

	return session.NewSession(clientBuilder, sessionCb, panicHandler, reporter.NullReporter{}), nil
}

type CliCallback struct{}

func (n CliCallback) OnNetworkRestored() {
	fmt.Println("Network restored")
}

func (n CliCallback) OnNetworkLost() {
	fmt.Println("Network lost")
}

func login(ctx *cli.Context, s *session.Session) error {
	creds := newCredentialsFromCli(ctx)
	var err error
	for {
		switch s.LoginState() {
		case session.LoginStateLoggedOut:
			if len(creds.username) == 0 {
				if creds.username, err = readLine("Enter your username: "); err != nil {
					return err
				}
			}
			if len(creds.password) == 0 {
				if creds.password, err = readPassword("Enter your password: "); err != nil {
					return err
				}
			}
			if err := s.Login(ctx.Context, creds.username, creds.password); err != nil {
				if err := creds.nextAttempt(); err != nil {
					return err
				}
			}
		case session.LoginStateAwaitingTOTP:
			if len(creds.totp) == 0 {
				if creds.totp, err = readLine("Enter the code from your authenticator app: "); err != nil {
					return err
				}
			}
			if err := s.SubmitTOTP(ctx.Context, creds.totp); err != nil {
				if err := creds.nextAttempt(); err != nil {
					return err
				}
			}
		case session.LoginStateAwaitingMailboxPassword:
			if len(creds.mboxPassword) == 0 {
				if creds.mboxPassword, err = readPassword("Enter you mailbox password: "); err != nil {
					return err
				}
			}

			if err := s.SubmitMailboxPassword(
				apiclient.NewProtonMailboxPasswordValidator(s.GetUser(), s.GetUserSalts()),
				creds.mboxPassword,
			); err != nil {
				return err
			}
		case session.LoginStateAwaitingHV:
			url, err := s.GetHVSolveURL()
			if err != nil {
				return err
			}

			fmt.Printf("Human Verification requested. Please open the URL below in a  browser and "+
				" press ENTER when the challenge has been completed.\n\n%s\n\n", url)
			waitForReturn()

			if err := s.MarkHVSolved(ctx.Context); err != nil {
				return err
			}
		case session.LoginStateLoggedIn:
			return nil
		default:
			return fmt.Errorf("unknown login state: %v", s.LoginState())
		}
	}
}

type credentials struct {
	username     string
	password     []byte
	totp         string
	mboxPassword []byte
	attemptCount int
}

func newCredentialsFromCli(ctx *cli.Context) *credentials {
	return &credentials{
		username:     ctx.String(flagUsername.Name),
		password:     []byte(ctx.String(flagPassword.Name)),
		totp:         ctx.String(flagTOTP.Name),
		mboxPassword: []byte(ctx.String(flagMBoxPassword.Name)),
	}
}

func (c *credentials) nextAttempt() error {
	if c.attemptCount++; c.attemptCount >= 5 {
		return errors.New("failed to login: max attempts reached")
	}

	c.username = ""
	c.password = nil
	c.totp = ""
	c.mboxPassword = nil

	return nil
}

func readLine(prompt string) (string, error) {
	if len(prompt) > 0 {
		fmt.Print(prompt)
	}

	result, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func readPassword(prompt string) ([]byte, error) {
	if len(prompt) > 0 {
		fmt.Print(prompt)
	}

	result, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}

	fmt.Println()

	return result, nil
}

func waitForReturn() {
	_, _ = bufio.NewReader(os.Stdin).ReadSlice('\n')
}

func defaultDownloadFolder() string {
	return getDownloadDir()
}
