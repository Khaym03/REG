package auth

import (
	"context"

	c "github.com/Khaym03/REG/internal/constants"
	"github.com/Khaym03/REG/internal/repo"
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
		existingUser.LoginState()
		a.UpdateUser(*existingUser)
		return true
	}

	return false
}

// Login to check is user exists, the session keeps open,
// is upto the caller if want to logout inmediatly
func (a *AccountService) AuthUser(
	ctx context.Context,
	user User,
	session Session,
) (err error) {

	var existingUser *RegisterUsers = nil

	users, err := a.p.Load()
	if err != nil {
		return err
	}

	for _, knownUser := range users {
		if user.Username == knownUser.Username {
			existingUser = &knownUser
		}

	}

	if existingUser != nil {
		// dont need authentication
		existingUser.LoginState()
		a.UpdateUser(*existingUser)
		return nil
	}

	err = a.auth.Login(ctx, session, user)
	if err != nil {
		return err
	}

	// no erros means user exist
	return a.StoreUserSecret(user)
}

func (a *AccountService) StoreUserSecret(user User) (err error) {
	if err = user.Validate(); err != nil {
		return err
	}

	err = keyring.Set(c.AppName, user.Username, user.Password)
	if err != nil {
		return err
	}

	users, err := a.p.Load()
	if err != nil {
		return err
	}

	for _, u := range users {
		u.LastUse = false
	}

	usr := RegisterUsers{
		Username: user.Username,
	}

	usr.LoginState()

	users = append(users, usr)

	return a.p.Save(users)
}

func (a *AccountService) GetUserPassword(username string) (User, error) {
	secret, err := keyring.Get(c.AppName, username)
	if err != nil {
		return User{}, err
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

	for _, u := range users {
		if u.Logged {
			return &u
		}
	}
	return nil
}
