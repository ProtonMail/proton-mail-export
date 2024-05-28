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
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/ProtonMail/export-tool/internal/sentry"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/gluon/async"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const retryCount = 5

var (
	flagUsername = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "username",
		Aliases: []string{"u"},
		EnvVars: []string{"ET_USER_EMAIL"},
	}
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
	flagOperation = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "operation",
		Aliases: []string{"o"},
		EnvVars: []string{"ET_OPERATION"},
	}
	flagFolder = &cli.StringFlag{ //nolint:gochecknoglobals
		Name:    "folder",
		Aliases: []string{"f"},
		EnvVars: []string{"ET_FOLDER"},
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
			flagOperation,
			flagFolder,
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

	session, err := newSession(panicHandler)
	if err != nil {
		return err
	}

	operation, err := getOperation(ctx)
	if err != nil {
		return err
	}

	if err = login(ctx, session); err != nil {
		return err
	}

	dir, err := getTargetFolder(ctx, operation, session.GetUser().Email)
	if err != nil {
		return err
	}

	if operation == operationBackup {
		return runBackup(ctx.Context, dir, session)
	}

	if operation == operationRestore {
		return runRestore(ctx.Context, dir, session)
	}

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
	creds := newCredentialsFromCLI(ctx)
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

func runBackup(ctx context.Context, exportPath string, session *session.Session) error {
	exportTask := mail.NewExportTask(ctx, exportPath, session)

	return exportTask.Run(ctx, newCliReporter())
}

func runRestore(ctx context.Context, backupPath string, session *session.Session) error {
	restoreTask, err := mail.NewRestoreTask(ctx, backupPath, session)
	if err != nil {
		return err
	}

	return restoreTask.Run(newCliReporter())
}
