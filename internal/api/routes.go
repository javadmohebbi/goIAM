// Package api provides HTTP route handlers for goIAM.
// This file contains local authentication and 2FA setup routes.
package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/javadmohebbi/goIAM/internal/middleware"
)

// RegisterLocalRoutes sets up local authentication routes for the goIAM API.
//
// It defines public endpoints for registration and login, and registers protected routes
// under the /secure path using RequireAuth middleware. These protected routes include 2FA,
// profile management, and backup code functionality, registered via registerAuthRoutes().
func (a *API) registerRoutes(app *fiber.App) {
	// a.handleLogin and a.handleRegister will internally dispatch to the correct auth method
	// based on the configured precedence in a.cfg.AuthProviders
	// Register a unified login and register endpoint
	app.Post("/auth/login", a.handleLogin)
	app.Post("/auth/register", a.handleRegister)
	app.Post("/auth/reset/password/request", a.handleResetPasswordRequest)

	// token check middleware
	secure := app.Group("/s", middleware.RequireAuth(a.cfg, a.iamDB))

	// auth and profile-related routes
	a.registerAuthRoutes(secure)

	// register user-related routes
	userRoutes := secure.Group("/user")
	a.registerUserRoutes(userRoutes)
}

// handleLogin attempts login with each configured AuthProvider in order.
// It tries local, LDAP, etc., in configured order, and returns Unauthorized if all fail.
func (a *API) handleLogin(c fiber.Ctx) error {
	for _, provider := range a.cfg.AuthProviders {
		switch provider.Name {
		case "local":
			if err := a.handleLoginLocal(c); err == nil {
				return nil
			}
		case "ldap":
			// var cfg config.LDAPConfig
			// Future: implement LDAP login
		case "auth0":
			// Future: implement Auth0 login
		case "entra_id":
			// Future: implement Entra ID login
		}
	}
	return fiber.ErrUnauthorized
}

// handleRegister attempts registration with each configured AuthProvider in order.
// Only local registration is generally supported.
func (a *API) handleRegister(c fiber.Ctx) error {
	for _, provider := range a.cfg.AuthProviders {
		switch provider.Name {
		case "local":
			if err := a.handleRegisterLocal(c); err == nil {
				return nil
			}
		case "ldap":
			// var cfg config.LDAPConfig
			// Future: implement LDAP login
		case "auth0":
			// Future: implement Auth0 login
		case "entra_id":
			// Future: implement Entra ID login
		}
	}
	return fiber.ErrUnauthorized
}

// handleRegister attempts registration with each configured AuthProvider in order.
// Only local registration is generally supported.
func (a *API) handleResetPasswordRequest(c fiber.Ctx) error {
	for _, provider := range a.cfg.AuthProviders {
		switch provider.Name {
		case "local":
			if err := a.handleResetPasswordRequestLocal(c); err == nil {
				return nil
			}
		case "ldap":
			// var cfg config.LDAPConfig
			// Future: implement LDAP login
		case "auth0":
			// Future: implement Auth0 login
		case "entra_id":
			// Future: implement Entra ID login
		}
	}
	return fiber.ErrUnauthorized
}
