package vault

import (
	"fmt"

	"github.com/riotgames/vault-go-client"
	"go.uber.org/zap"
)

// Option is a function that sets a configuration option for the Component
type Option func(*options)

// AuthProvider is a function that accepts a Vault client and authenticates the client with the Vault server
type AuthProvider func(client *vault.Client) error

type options struct {
	log          *zap.Logger
	url          string
	authProvider AuthProvider
}

func createDefaultOptions(config *Config) *options {
	opts := &options{
		log: zap.L(),
		url: defaultURL,
	}
	WithDefaultChainAuthProvider(config)(opts)
	return opts
}

// WithLogger sets the logger for the Component
func WithLogger(l *zap.Logger) Option {
	return func(o *options) {
		o.log = l
	}
}

// WithURL sets the Vault Server URL for the Vault Client in Component
func WithURL(url string) Option {
	return func(o *options) {
		o.url = url
	}
}

// WithDefaultChainAuthProvider attempts to authenticate to Vault with the default auth chain. This method is preferred to
// authenticate the Vault Client both when running within AWS resources. and during local testing
//
// The default auth chain will read the following config values to attempt to authenticate with Riot Vault servers:
//   - vault.iamRoleARN
//   - vault.appRoleID
//   - vault.appRoleSecretID
//
// These authentication methods are attempted in the following order:
//   - IAM Login
//   - AppRole Login
//
// LDAP authentication is intentionally omitted to mitigate risk for exposing LDAP credentials during local testing.
func WithDefaultChainAuthProvider(config *Config) Option {
	return func(o *options) {
		o.authProvider = func(client *vault.Client) error {
			var errs []error

			_, iamErr := client.Auth.IAM.Login(config.getIAMLoginOptions())
			if iamErr == nil {
				o.log.Info("Vault Client authed with AWS IAM method")
				return nil
			}
			o.log.Debug(fmt.Sprintf("Vault Client IAM auth failed with err: %s", iamErr.Error()))
			errs = append(errs, iamErr)

			_, appRoleErr := client.Auth.AppRole.Login(config.getAppRoleLoginOptions())
			if appRoleErr == nil {
				o.log.Info("Vault Client authed with AppRole method")
				return nil
			}
			o.log.Debug(fmt.Sprintf("Vault Client AppRole auth failed with err: %s", appRoleErr.Error()))
			errs = append(errs, appRoleErr)

			return fmt.Errorf("%w; %w", iamErr, appRoleErr)
		}
	}
}

// WithIAMAuthProvider attempts to authenticate to Vault with the IAM login method
// https://developer.hashicorp.com/vault/docs/auth/aws
func WithIAMAuthProvider(loginOptions vault.IAMLoginOptions) Option {
	return func(o *options) {
		o.authProvider = func(client *vault.Client) error {
			_, authErr := client.Auth.IAM.Login(loginOptions)
			return authErr
		}
	}
}

// WithAppRoleAuthProvider attempts to authenticate to Vault with the AppRole login method:
// https://developer.hashicorp.com/vault/docs/auth/approle
func WithAppRoleAuthProvider(loginOptions vault.AppRoleLoginOptions) Option {
	return func(o *options) {
		o.authProvider = func(client *vault.Client) error {
			_, authErr := client.Auth.AppRole.Login(loginOptions)
			return authErr
		}
	}
}
