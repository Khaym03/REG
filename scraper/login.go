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
	"github.com/go-rod/rod/lib/proto"
)

var _ domain.AuthService = (*LoginScraper)(nil)

type LoginScraper struct {
}

func NewLoginScraper() *LoginScraper {
	return &LoginScraper{}
}

func (l *LoginScraper) Login(ctx context.Context, page *rod.Page, user domain.User) (err error) {
	page = page.Context(ctx)

	if err = navigate(page, c.LoginURL); err != nil {
		return err
	}

	if err = dismissOptionalModal(page, makeInteracteableButtonSelector, c.TimeoutShort); err != nil {
		return err
	}

	if err := fillInput(page, emailInputSelector, user.Username); err != nil {
		return fmt.Errorf("email input failed: %w", err)
	}
	if err := fillInput(page, passwordInputSelector, user.Password); err != nil {
		return fmt.Errorf("password input failed: %w", err)
	}

	if err = click(page, loginButtonSelector); err != nil {
		return err
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("post-login wait failed: %w", err)
	}

	if err = checkLoginError(page); err != nil {
		return err
	}

	if err = handleVerificationStep(page); err != nil {
		return err
	}

	return nil
}

func (l *LoginScraper) Logout(ctx context.Context, page *rod.Page) (err error) {
	page = page.Context(ctx)

	if err = navigate(page, c.BaseURL); err != nil {
		return err
	}

	if err = click(page, profileDropdownSelector); err != nil {
		return err
	}

	return click(page, logoutSelector)
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

// dismissOptionalModal attempts to clear a popup if it exists.
func dismissOptionalModal(page *rod.Page, selector string, idleTimeout time.Duration) error {
	elements, err := page.ElementsX(selector)
	if err != nil {
		return fmt.Errorf("failed to search for modal elements: %w", err)
	}

	// If no elements are found, we assume there is no modal to dismiss.
	if elements.Empty() {
		return nil
	}

	log.Println("Random modal detected, dismissing...")
	modalBtn := elements.First()

	// Click the button
	if err := modalBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click modal button: %w", err)
	}

	// Wait for the element to disappear
	if err := modalBtn.WaitInvisible(); err != nil {
		return fmt.Errorf("modal button did not disappear: %w", err)
	}

	// Ensure the page is idle before continuing
	if err := page.WaitIdle(idleTimeout); err != nil {
		return fmt.Errorf("error waiting for page idleness after dismissal: %w", err)
	}

	return nil
}

func checkLoginError(page *rod.Page) error {
	elements, err := page.Elements(".alert-danger")
	if err != nil {
		return fmt.Errorf("failed to query alert elements: %w", err)
	}

	if elements.Empty() {
		return nil
	}

	// Attempt to extract the text from the first alert found
	errorMessage, err := elements.First().Text()
	if err != nil {
		return fmt.Errorf("login alert detected, but failed to read text: %w", err)
	}

	// If the alert exists but is empty, provide a generic failure message
	if errorMessage == "" {
		return fmt.Errorf("login failed with an empty alert message")
	}

	return fmt.Errorf("login failed: %s", errorMessage)
}

func handleVerificationStep(page *rod.Page) error {
	verifyInput, err := page.ElementX(verifyInputSelector)
	if err != nil {
		return nil
	}

	log.Println("Verification step triggered")

	codeEl, err := page.ElementX(codeTextSelector)
	if err != nil {
		return fmt.Errorf("verification code element not found: %w", err)
	}

	codeText, err := codeEl.Text()
	if err != nil {
		return fmt.Errorf("failed to get verification text: %w", err)
	}

	code := extractVerificationCode(codeText)
	if code == "" {
		return fmt.Errorf("verification code not found in text: %q", codeText)
	}

	if err := verifyInput.Input(code); err != nil {
		return fmt.Errorf("failed to input code: %w", err)
	}

	if err := page.Keyboard.Press(input.Enter); err != nil {
		return fmt.Errorf("failed to press enter: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("page failed to load after verification: %w", err)
	}

	return nil
}
