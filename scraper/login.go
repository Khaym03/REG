package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

var _ domain.AuthService = (*LoginScraper)(nil)

type LoginScraper struct {
	browser *rod.Browser
}

func NewLoginScraper(browser *rod.Browser) *LoginScraper {
	return &LoginScraper{browser: browser}
}

func (l *LoginScraper) Login(ctx context.Context, user domain.User) error {
	// Use Try to catch panics from "Must" calls and return them as errors
	return rod.Try(func() {
		page := l.browser.MustPage()
		defer page.Close()

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
}

func (l *LoginScraper) Logout(ctx context.Context) error {
	return rod.Try(func() {
		page := l.browser.MustPage()
		defer page.Close()

		page.MustNavigate(c.BaseURL)
		page.MustElementX(profileDropdownSelector).MustClick()
		page.MustElementX(logoutSelector).MustClick()
	})
}

const (
	emailInputSelector    = `#exampleInputEmail`
	passwordInputSelector = `#passwordInput`

	loginButtonSelector = `body > div.container-scroller > div > div > div > div.col-lg-5.d-flex.flex-column.align-items-center.justify-content-center.vh-100 > div.auth-form-transparent.text-left.p-3 > form > div.my-3 > button`

	verifyInputSelector = `/html/body/div[1]/div/div/div/div[1]/div[1]/div/form[1]/div/input`

	codeTextSelector                = `/html/body/div[1]/div/div/div/div[1]/div[1]/div/p[1]`
	makeInteracteableButtonSelector = `//*[@id="modal_sesion_notificacion_login"]/div/div/div/button`

	profileDropdownSelector = `//*[@id="profileDropdown"]`
	logoutSelector          = `/html/body/div[1]/nav/div[2]/ul/li[3]/div/a`
)

// Regex: Matches exactly 6 digits isolated by word boundaries
var verificationRegex = regexp.MustCompile(`\b\d{6}\b`)

func extractVerificationCode(text string) string {
	result := verificationRegex.FindString(text)

	if result != "" {
		log.Printf("Verification sequence found: %s", result)
	} else {
		log.Println("No 6-digit sequence was found in the provided text.")
	}
	return result
}
