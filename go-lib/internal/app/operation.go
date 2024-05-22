package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type Operation int

const (
	operationUnknown Operation = iota
	operationBackup
	operationRestore
)

func getOperation(ctx *cli.Context) (Operation, error) {
	argsOperation := ctx.String(flagOperation.Name)
	if len(argsOperation) == 0 {
		return readOperationFromCLI()
	} else {
		return stringToOperation(argsOperation)
	}
}

func readOperationFromCLI() (Operation, error) {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < retryCount; i++ {
		fmt.Printf("Enter the operation ((B)ackup / (R)restore): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return operationUnknown, err
		}
		input = strings.TrimSpace(input)
		operation, err := stringToOperation(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			return operation, err
		}
	}

	return operationUnknown, errors.New("too many failed attempts")
}

func stringToOperation(operation string) (Operation, error) {
	if strings.EqualFold(operation, "backup") || strings.EqualFold(operation, "b") {
		return operationBackup, nil
	}

	if strings.EqualFold(operation, "restore") || strings.EqualFold(operation, "r") {
		return operationRestore, nil
	}

	return operationUnknown, fmt.Errorf("unknown operation %s", operation)
}

func operationToString(operation Operation) string {
	switch operation {
	case operationBackup:
		return "backup"
	case operationRestore:
		return "restore"
	default:
		return "unknown"
	}
}
