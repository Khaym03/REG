package auth

import (
	"context"

	c "github.com/Khaym03/REG/internal/constants"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/Khaym03/REG/internal/session"
	log "github.com/sirupsen/logrus"
	"github.com/zalando/go-keyring"
)

type RegisterUsers struct {
	Username string `json:"username"`
	LastUse  bool   `json:"last_use"`
	Logged   bool   `json:"logged"`
}

func (u *RegisterUsers) LoginState() {
	u.LastUse = true
	u.Logged = true
}

type AccountService struct {
	auth AuthService
	p    repo.Persistence[[]RegisterUsers]
}

func NewAccountService(
	auth AuthService,
	p repo.Persistence[[]RegisterUsers],
) *AccountService {
	return &AccountService{
		auth: auth,
		p:    p,
	}
}

func (a AccountService) KnownUser(user User) bool {
	var existingUser *RegisterUsers = nil

	users, err := a.p.Load()
	if err != nil {
		log.Error(err)
	}

	for _, knownUser := range users {
		if user.Username == knownUser.Username {
			existingUser = &knownUser
		}

	}

	if existingUser != nil {
		// dont need authentication
		package auth

		import (
			"context"
			"fmt"

			c "github.com/Khaym03/REG/internal/constants"
			"github.com/Khaym03/REG/internal/repo"
			"github.com/Khaym03/REG/internal/session"
			log "github.com/sirupsen/logrus"
			"github.com/zalando/go-keyring"
		)

		type RegisterUsers struct {
			Username string `json:"username"`
			LastUse  bool   `json:"last_use"`
			Logged   bool   `json:"logged"`
		}

		func (u *RegisterUsers) LoginState() {
			u.LastUse = true
			u.Logged = true
		}

		type AccountService struct {
			auth AuthService
			p    repo.Persistence[[]RegisterUsers]
		}

		func NewAccountService(
			auth AuthService,
			p repo.Persistence[[]RegisterUsers],
		) *AccountService {
			return &AccountService{
				auth: auth,
				p:    p,
			}
		}

		func (a *AccountService) KnownUser(user User) bool {
			var existingUser *RegisterUsers

			users, err := a.p.Load()
			if err != nil {
				log.Error(err)
			}

			for i := range users {
				if user.Username == users[i].Username {
					existingUser = &users[i]
					break
				}
			}

			if existingUser != nil {
				existingUser.LoginState()
				_ = a.UpdateUser(*existingUser)
				return true
			}

			return false
		}

		// Login to check is user exists, the session keeps open,
		// is upto the caller if want to logout inmediatly
		func (a *AccountService) AuthUser(
			ctx context.Context,
			user User,
			session session.Session,
		) (err error) {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			var existingUser *RegisterUsers

			users, err := a.p.Load()
			if err != nil {
				return fmt.Errorf("load users: %w", err)
			}

			for i := range users {
				if user.Username == users[i].Username {
					existingUser = &users[i]
					break
				}
			}

			if existingUser != nil {
				existingUser.LoginState()
				if err := a.UpdateUser(*existingUser); err != nil {
					return fmt.Errorf("update user: %w", err)
				}
				return nil
			}

			if err := a.auth.Login(ctx, session, user); err != nil {
				return fmt.Errorf("auth login: %w", err)
			}

			if err := a.StoreUserSecret(user); err != nil {
				return fmt.Errorf("store user secret: %w", err)
			}

			return nil
		}

		func (a *AccountService) StoreUserSecret(user User) (err error) {
			if err = user.Validate(); err != nil {
				return err
			}

			if err = keyring.Set(c.AppName, user.Username, user.Password); err != nil {
				return fmt.Errorf("set keyring secret: %w", err)
			}

			users, err := a.p.Load()
			if err != nil {
				return fmt.Errorf("load users: %w", err)
			}

			for i := range users {
				users[i].LastUse = false
			}

			usr := RegisterUsers{Username: user.Username}
			usr.LoginState()

			users = append(users, usr)

			return a.p.Save(users)
		}

		func (a *AccountService) GetUserPassword(username string) (User, error) {
			secret, err := keyring.Get(c.AppName, username)
			if err != nil {
				return User{}, fmt.Errorf("get user password: %w", err)
			}
			return User{Username: username, Password: secret}, nil
		}

		func (a *AccountService) GetRegisterUsers() ([]RegisterUsers, error) {
			return a.p.Load()
		}

		func (a *AccountService) UpdateUser(user RegisterUsers) error {
			old, err := a.p.Load()
			if err != nil {
				return err
			}

			for i, oldUser := range old {
				if user.Username == oldUser.Username {
					old[i] = user
				}
			}

			return a.p.Save(old)
		}

		func (a *AccountService) CurrentUser() *RegisterUsers {
			users, err := a.p.Load()
			if err != nil {
				log.Error(err)
			}

			for i := range users {
				if users[i].Logged {
					return &users[i]
				}
			}
			return nil
		}
