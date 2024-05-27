package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

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

	result, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	fmt.Println()

	return result, nil
}

func waitForReturn() {
	_, _ = bufio.NewReader(os.Stdin).ReadSlice('\n')
}

func readYesNo(prompt string, retryCount int) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < retryCount; i++ {
		fmt.Print(prompt)
		text, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		text = strings.TrimSpace(text)
		if strings.EqualFold(text, "yes") || strings.EqualFold(text, "y") {
			return true, nil
		}
		if strings.EqualFold(text, "no") || strings.EqualFold(text, "n") {
			return false, nil
		}
	}

	return true, errors.New("too many attempts")
}
