package login

import (
	"fmt"
	"log"
	"regexp"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// Regex: Matches exactly 6 digits isolated by word boundaries
var verificationRegex = regexp.MustCompile(`\b\d{6}\b`)

type User struct {
	Username string
	Password string
}

type CloseSession func()

func Login(page *rod.Page, user User) (CloseSession, error) {
	// Use Try to catch panics from "Must" calls and return them as errors
	err := rod.Try(func() {
		page.MustNavigate(c.LoginURL).MustWaitLoad()

		// Handle optional modal
		if el, err := page.Timeout(1 * time.Second).ElementX(makeInteracteableButtonSelector); err == nil {
			log.Println("Random modal detected, dismissing...")
			el.MustClick()
		}

		page.MustElement(emailInputSelector).MustClick().MustInput(user.Username)
		page.MustElement(passwordInputSelector).MustInput(user.Password)
		page.MustElement(loginButtonSelector).MustClick()
		page.MustWaitLoad()

		// Check for explicit error messages from the UI
		if errorElements := page.MustElements(`.alert-danger`); len(errorElements) > 0 {
			errorText := errorElements.First().MustText()
			if errorText != "" {
				panic(fmt.Errorf("login failed: %s", errorText))
			}
		}

		// Handle 2FA/Verification
		verifyElements := page.MustElementsX(verifyInputSelector)
		if len(verifyElements) > 0 {
			log.Println("Verification step triggered")
			code := extractVerificationCode(page.MustElementX(codeTextSelector).MustText())
			if code == "" {
				panic("verification code not found in text")
			}

			page.MustElementX(verifyInputSelector).MustClick().MustInput(code)
			page.Keyboard.MustType(input.Enter)
			page.MustWaitLoad()
		}
	})

	if err != nil {
		return nil, err
	}

	// Return a closure that handles its own potential failures
	return func() {
		_ = rod.Try(func() {
			page.MustNavigate(c.BaseURL)
			page.MustElementX(profileDropdownSelector).MustClick()
			page.MustElementX(logoutSelector).MustClick()
		})
	}, nil
}

func extractVerificationCode(text string) string {
	result := verificationRegex.FindString(text)

	if result != "" {
		log.Printf("Verification sequence found: %s", result)
	} else {
		log.Println("No 6-digit sequence was found in the provided text.")
	}
	return result
}
