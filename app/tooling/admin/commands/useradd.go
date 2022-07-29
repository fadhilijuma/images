package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/fadhilijuma/images/business/core/user"
	"github.com/fadhilijuma/images/business/sys/database"
	"github.com/fadhilijuma/images/business/web/auth"
	"go.uber.org/zap"
)

// UserAdd adds new users into the database.
func UserAdd(log *zap.SugaredLogger, cfg database.Config, name, email, agency, password string) error {
	if name == "" || email == "" || password == "" || agency == "" {
		fmt.Println("help: useradd <name> <email> <agency> <password>")
		return ErrHelp
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	core := user.NewCore(log, db)

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Agency:          agency,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	usr, err := core.Create(ctx, nu, time.Now())
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
