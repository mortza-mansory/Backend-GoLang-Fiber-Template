package app

import (
	"github.com/yourorg/go-fiber-template/internal/modules/auth"
)

// AuthModule holds the auth module's dependencies as wired at startup.
type AuthModule struct {
	Handler *auth.Handler
}

// Modules aggregates all registered feature modules.
type Modules struct {
	Auth *AuthModule
}

// newModules constructs every feature module.
//
// TODO: add new features here, e.g.:
//
//	modules.Payment = &PaymentModule{ Handler: payment.NewHandler(...) }
func newModules(d *Deps) *Modules {
	authService := auth.NewService(&authRepository{db: d.DB}, d.Logger)
	authHandler := auth.NewHandler(authService, d.Validator)

	return &Modules{
		Auth: &AuthModule{Handler: authHandler},
	}
}
