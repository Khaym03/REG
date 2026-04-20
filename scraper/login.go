package scraper

import (
	"context"
	"fmt"
	"log"
	"regexp"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/utils"
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

func (l *LoginScraper) Login(ctx context.Context, page *rod.Page, user domain.User) error {
	page = page.Context(ctx)

	if err := page.Navigate(c.LoginURL); err != nil {
		return fmt.Errorf("failed to navigate to login: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("wait load failed: %w", err)
	}

	// We use ElementsX (plural) so it doesn't error if the modal isn't there
	els, err := page.ElementsX(makeInteracteableButtonSelector)
	if err == nil && !els.Empty() {
		log.Println("Random modal detected, dismissing...")
		modalBtn := els.First()
		if err := modalBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
			log.Printf("failed to click modal button: %v", err)
		} else {
			_ = modalBtn.WaitInvisible()
			_ = page.WaitIdle(c.TimeoutMedium)
		}
	}

	if err := utils.FillInput(page, emailInputSelector, user.Username); err != nil {
		return fmt.Errorf("email input failed: %w", err)
	}
	if err := utils.FillInput(page, passwordInputSelector, user.Password); err != nil {
		return fmt.Errorf("password input failed: %w", err)
	}

	loginBtn, err := page.Element(loginButtonSelector)
	if err != nil {
		return fmt.Errorf("login button not found: %w", err)
	}
	if err := loginBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click login: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("post-login wait failed: %w", err)
	}

	errElements, err := page.Elements(`.alert-danger`)
	if err == nil && len(errElements) > 0 {
		text, _ := errElements.First().Text()
		if text != "" {
			return fmt.Errorf("login failed: %s", text)
		}
	}

	verifyInput, err := page.ElementX(verifyInputSelector)
	if err == nil {
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
			return fmt.Errorf("verification code not found in text: %s", codeText)
		}

		if err := verifyInput.Input(code); err != nil {
			return fmt.Errorf("failed to input code: %w", err)
		}

		if err := page.Keyboard.Press(input.Enter); err != nil {
			return fmt.Errorf("failed to press enter: %w", err)
		}

		return page.WaitLoad()
	}

	return nil

}

func (l *LoginScraper) Logout(ctx context.Context, page *rod.Page) error {
	page = page.Context(ctx)

	if err := page.Navigate(c.BaseURL); err != nil {
		return err
	}

	dropdown, err := page.ElementX(profileDropdownSelector)
	if err != nil {
		return fmt.Errorf("logout dropdown not found: %w", err)
	}
	if err := dropdown.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	logoutBtn, err := page.ElementX(logoutSelector)
	if err != nil {
		return fmt.Errorf("logout button not found: %w", err)
	}

	return logoutBtn.Click(proto.InputMouseButtonLeft, 1)
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
