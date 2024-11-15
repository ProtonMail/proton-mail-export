package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ProtonMail/export-tool/internal"
	"github.com/ProtonMail/export-tool/internal/apiclient"
	"github.com/ProtonMail/export-tool/internal/mail"
	"github.com/ProtonMail/export-tool/internal/reporter"
	"github.com/ProtonMail/export-tool/internal/sentry"
	"github.com/ProtonMail/export-tool/internal/session"
	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
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
		Name:    "dir",
		Aliases: []string{"f"},
		EnvVars: []string{"ET_DIR"},
	}
)

func Run() {
	folder, err := getDefaultOperationFolder()
	if err != nil {
		fatal(err)
	}

	if err := initApp(filepath.Join(folder, "logs"), func() {}); err != nil {
		fatal(err)
	}

	defer closeApp()

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
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("\nFatal error: %v\n", err)
	logrus.WithError(err).Fatal("Fatal error")
}

func run(ctx *cli.Context) error {
	panicHandler := sentry.NewPanicHandler(func() {})
	defer async.HandlePanic(panicHandler)

	printHeader()
	checkForNewVersion()

	fmt.Printf("\nSession log: %v\n\n", filepath.FromSlash(state.logPath))

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
	fmt.Printf("This program is licensed under the GNU General Public License v3\n")
	fmt.Printf("Get support at https://proton.me/support/proton-mail-export-tool\n\n")
}

func checkForNewVersion() {
	fmt.Print("Checking for new version... ")
	if internal.HasNewVersion() {
		fmt.Println("A new version is available at: https://proton.me/support/proton-mail-export-tool")
	} else {
		fmt.Println("Your version is up to date")
	}
}

func getAPIURL() string {
	url := os.Getenv("ET_API_URL")
	if len(url) == 0 {
		url = internal.ETDefaultAPIURL
	}

	return url
}

func newSession(panicHandler async.PanicHandler) (*session.Session, error) {
	sessionCb := CliCallback{}
	builder, err := apiclient.NewProtonAPIClientBuilder(getAPIURL(), panicHandler, sessionCb)
	if err != nil {
		return nil, err
	}

	clientBuilder := apiclient.NewAutoRetryClientBuilder(
		builder,
		&apiclient.SleepRetryStrategyBuilder{},
	)

	return session.NewSession(clientBuilder, sessionCb, panicHandler, reporter.NullReporter{}, false), nil
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
				printError(err)
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
				printError(err)
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
				printError(err)
				if err := creds.nextAttempt(); err != nil {
					return err
				}
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
	fmt.Printf("Starting backup - Path=\"%v\"\n", filepath.FromSlash(exportTask.GetExportPath()))
	err := exportTask.Run(ctx, newCliReporter())
	if err == nil {
		fmt.Println("Backup finished")
	}

	return err
}

func runRestore(ctx context.Context, backupPath string, session *session.Session) error {
	restoreTask, err := mail.NewRestoreTask(ctx, backupPath, session)
	if err != nil {
		return err
	}

	fmt.Println("Starting restore")
	err = restoreTask.Run(newCliReporter())
	if err == nil {
		fmt.Println("Restore finished")
	}
	printRestoreTaskSummary(restoreTask)
	return err
}

func printRestoreTaskSummary(task *mail.RestoreTask) {
	fmt.Printf("Importable emails: %v\n", task.GetImportableCount())
	fmt.Printf("Successful imports: %v\n", task.GetImportedCount())
	fmt.Printf("Failed imports: %v\n", task.GetFailedCount())
	fmt.Printf("Skipped imports: %v\n", task.GetSkippedCount())
}

func initApp(defaultOperationPath string, onRecover func()) error {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if err := sentry.InitSentry(); err != nil {
		return err
	}

	if state.file != nil {
		return errors.New("application has already been initialized")
	}

	if err := os.MkdirAll(defaultOperationPath, 0o700); err != nil {
		return err
	}

	state.logPath = filepath.Join(defaultOperationPath, internal.NewLogFileName())
	file, err := os.OpenFile(state.logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		ForceQuote:       true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		TimestampFormat:  "2006-01-02 15:04:05.000",
	})
	internal.LogPrelude()

	if onRecover != nil {
		state.onRecover = onRecover
	} else {
		state.onRecover = func() {
			os.Exit(-200)
		}
	}

	state.reporter = sentry.NewReporter()

	return nil
}

func closeApp() {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if state.file != nil {
		logrus.SetOutput(os.Stdout)
		if err := state.file.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close log file")
		} else {
			state.file = nil
		}
	}
}

func printError(err error) {
	var apiError *proton.APIError
	if errors.As(err, &apiError) {
		fmt.Println(apiError.Message)
		return
	}

	fmt.Println(err)
}

type globalState struct {
	mutex     sync.Mutex
	file      *os.File
	logPath   string
	onRecover func()
	reporter  reporter.Reporter
}

//nolint:gochecknoglobals
var state globalState
