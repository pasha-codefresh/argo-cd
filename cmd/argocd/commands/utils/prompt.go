package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/argoproj/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type Prompt interface {
	Confirm(message string) bool
}

type PromptOpts struct {
	// PromptsEnabled indicates whether prompts should be enabled
	Enabled bool
}

type prompt struct {
	enabled bool
}

func NewPrompt(opts... PromptOpts) Prompt {
	enabled := false
	if len(opts) > 0 {
		enabled = opts[0].Enabled
	}
	
	return &prompt{
		enabled: enabled,
	}
}

// Confirm prompts the user with a message (typically a yes or no question) and returns whether
// they responded in the affirmative or negative.
func (p *prompt) Confirm(message string) bool {
	if !p.enabled {
		return true
	}

	return askToProceed(message)
}


type AskToProcess interface {
	AskToProceed(message string) bool
	AskToProceedS(message string) string
}

type askToProcess struct{}

func NewAskToProcess() AskToProcess {
	return &askToProcess{}
}


// AskToProceed prompts the user with a message (typically a yes or no question) and returns whether
// they responded in the affirmative or negative.
func (a *askToProcess) AskToProceed(message string) bool {
	for {
		fmt.Print(message)
		reader := bufio.NewReader(os.Stdin)
		proceedRaw, err := reader.ReadString('\n')
		errors.CheckError(err)
		switch strings.ToLower(strings.TrimSpace(proceedRaw)) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}

// AskToProceedS prompts the user with a message (typically a yes, no or all question) and returns string
// "a", "y" or "n".
func AskToProceedS(message string) string {
	for {
		fmt.Print(message)
		reader := bufio.NewReader(os.Stdin)
		proceedRaw, err := reader.ReadString('\n')
		errors.CheckError(err)
		switch strings.ToLower(strings.TrimSpace(proceedRaw)) {
		case "y", "yes":
			return "y"
		case "n", "no":
			return "n"
		case "a", "all":
			return "a"
		}
	}
}



// PromptCredentials is a helper to prompt the user for a username and password (unless already supplied)
func PromptCredentials(username, password string) (string, string) {
	return PromptUsername(username), PromptPassword(password)
}

// PromptUsername prompts the user for a username value
func PromptUsername(username string) string {
	return PromptMessage("Username", username)
}

// PromptMessage prompts the user for a value (unless already supplied)
func PromptMessage(message, value string) string {
	for value == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message + ": ")
		valueRaw, err := reader.ReadString('\n')
		errors.CheckError(err)
		value = strings.TrimSpace(valueRaw)
	}
	return value
}

// PromptPassword prompts the user for a password, without local echo. (unless already supplied)
func PromptPassword(password string) string {
	for password == "" {
		fmt.Print("Password: ")
		passwordRaw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		errors.CheckError(err)
		password = string(passwordRaw)
		fmt.Print("\n")
	}
	return password
}

// ReadAndConfirmPassword is a helper to read and confirm a password from stdin
func ReadAndConfirmPassword(username string) (string, error) {
	for {
		fmt.Printf("*** Enter new password for user %s: ", username)
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Print("\n")
		fmt.Printf("*** Confirm new password for user %s: ", username)
		confirmPassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Print("\n")
		if string(password) == string(confirmPassword) {
			return string(password), nil
		}
		log.Error("Passwords do not match")
	}
}