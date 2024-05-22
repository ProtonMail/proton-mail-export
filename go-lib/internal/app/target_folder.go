package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

func getTargetFolder(ctx *cli.Context, operation Operation, username string) (string, error) {
	argsPath := ctx.String(flagFolder.Name)
	if len(argsPath) == 0 {
		return readTargetFolderFromCLI(operation, username)
	}
	return validateTargetFolder(operation, argsPath)
}

func readTargetFolderFromCLI(operation Operation, username string) (string, error) {
	defaultDir := filepath.Join(getDownloadDir(), username)
	useDefault, err := readYesNo(fmt.Sprintf("Use default folder '%s' for %s? (Y/N): ", defaultDir, operationToString(operation)), retryCount)
	if err != nil {
		return "", err
	}

	if useDefault {
		return validateTargetFolder(operation, defaultDir)
	}

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < retryCount; i++ {
		fmt.Printf("Enter the path of the target folder: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Printf("Error: please provide a path\n")
		}

		path, err := validateTargetFolder(operation, input)
		if err == nil {
			return path, nil
		}

		fmt.Printf("Error: %v\n", err)
	}

	return "", nil
}

func validateTargetFolder(operation Operation, path string) (string, error) {
	if (runtime.GOOS != "windows") && (strings.HasPrefix(path, "~")) {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1) // we do not support named home, such as `~john/test`
	}

	fullPath, err := filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		return "", err
	}

	if operation == operationBackup {
		if err = os.MkdirAll(fullPath, 0o700); err != nil {
			return "", err
		}
	}

	if operation == operationRestore {
		stat, err := os.Stat(fullPath)
		if err != nil {
			return "", err
		}
		if !stat.IsDir() {
			return "", errors.New("target folder is not a directory")
		}
	}

	return fullPath, nil
}
