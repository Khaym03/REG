package login

import (
	"fmt"
	"regexp"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type User struct {
	Username string
	Password string
}

type CloseSession func()

func Login(page *rod.Page, user User) (CloseSession, error) {
	page.MustNavigate(c.LoginURL)
	page.MustWaitLoad()

	page.MustElementX(makeInteracteableButtonSelector).MustClick()
	page.MustElement(emailInputSelector).MustClick().MustInput(user.Username)
	page.MustElement(passwordInputSelector).MustInput(user.Password)
	page.MustElement(loginButtonSelector).MustClick()
	page.MustWaitLoad()

	errorElements := page.MustElements(`.error, .alert, .alert-danger`)
	if len(errorElements) > 0 {
		errorText := errorElements.First().MustText()
		if errorText != "" {
			return nil, fmt.Errorf("login failed: %s", errorText)
		}
	}

	// Check if verification code input is present (indicates login step succeeded)
	verifyElements := page.MustElementsX(verifyInputSelector)
	if len(verifyElements) == 0 {
		return nil, fmt.Errorf("login failed: verification input not found")
	}

	code := extractVerificationCode(page.MustElementX(codeTextSelector).MustText())
	if code == "" {
		return nil, fmt.Errorf("verification code not found")
	}
	page.MustElementX(verifyInputSelector).MustClick().MustInput(code)
	if err := page.Keyboard.Press(input.Enter); err != nil {
		return nil, err
	}

	page.MustWaitLoad()

	// Verify successful login by checking for profile dropdown
	profileElements := page.MustElementsX(profileDropdownSelector)
	if len(profileElements) == 0 {
		return nil, fmt.Errorf("login verification failed: profile dropdown not found")
	}

	return func() {
		page.MustNavigate(c.BaseURL)
		page.MustElementX(profileDropdownSelector).MustClick()
		page.MustElementX(logoutSelector).MustClick()
	}, nil
}

func extractVerificationCode(text string) string {
	// Regex: Matches exactly 6 digits isolated by word boundaries
	re := regexp.MustCompile(`\b\d{6}\b`)

	// Find the first match
	result := re.FindString(text)

	if result != "" {
		fmt.Printf("Sequence found: %s\n", result)
	} else {
		fmt.Println("No 6-digit sequence was found.")
	}
	return result
}
