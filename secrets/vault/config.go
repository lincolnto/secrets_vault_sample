package vault

import (
	"fmt"

	"github.com/riotgames/vault-go-client"
)

// Client Configuration
const (
	// Required configuration
	envSpaceName  = "SPACE_NAME"
	envSecretPath = "SECRET_PATH"

	// Auth-related configuration
	envAWSIAMRoleARN   = "IAM_ROLE_ARN"
	envAppRoleID       = "APP_ROLE_ID"
	envAppRoleSecretID = "APP_ROLE_SECRET_ID"

	// Vault Space Path Formats
	secretsMountPathFormat = "%s/secrets"
	awsLoginPathFormat     = "%s/aws"
	appRoleLoginPathFormat = "%s/approle"
)

const (
	errFmtConfigValidation = "failed to validate Vault Provider config, %s is empty"
)

// Config fetches relevant component configuration values from the Glue property registry
type Config struct {
	// The Vault Space Name configured for the client instance (e.g. "MyTeamVaultSpace")
	SpaceName string
	// The Secret Path (i.e. secret key) to fetch secrets from (e.g. "mysecret-dev" in MyTeamVaultSpace/secrets)
	SecretPath string
	// The Vault Secret Mount Path where secrets are stored (e.g. "MyTeamVaultSpace/secrets")
	SecretMountPath string
	// The Vault Login Path to for Vault AWS Auth (e.g. "MyTeamVaultSpace/aws")
	AWSLoginPath string
	// The Vault Login Path to for Vault AppRole Auth (e.g. "MyTeamVaultSpace/approle")
	AppRoleLoginPath string
	// The AWS IAM Role ARN used for Vault AWS Auth
	AWSIAMRoleArn string
	// The Vault AppRole ID used Vault Auth
	AppRoleID string
	// The Vault AppRole Secret ID used for Vault Auth
	AppRoleSecretID string
}

// NewConfig returns a config object, with getters for config values relevant to the Vault Provider
func NewConfig() *Config {
	spaceName := os.Getenv(envSpaceName)

	return &Config{
		SpaceName:  spaceName,
		SecretPath: os.Getenv(envSecretPath),

		SecretMountPath:  fmt.Sprintf(secretsMountPathFormat, spaceName),
		AWSLoginPath:     fmt.Sprintf(awsLoginPathFormat, spaceName),
		AppRoleLoginPath: fmt.Sprintf(appRoleLoginPathFormat, spaceName),

		AWSIAMRoleArn:   os.Getenv(envAWSIAMRoleARN),
		AppRoleID:       os.Getenv(envAppRoleID),
		AppRoleSecretID: os.Getenv(envAppRoleSecretID),
	}
}

// getIAMLoginOptions creates an IAMLoginOptions object with values from the Component config
func (c *Config) getIAMLoginOptions() vault.IAMLoginOptions {
	return vault.IAMLoginOptions{
		Role:      c.AWSIAMRoleArn,
		MountPath: c.AWSLoginPath,
	}
}

// getAppRoleLoginOptions creates an AppRoleLoginOptions object with values from the Component config
func (c *Config) getAppRoleLoginOptions() vault.AppRoleLoginOptions {
	return vault.AppRoleLoginOptions{
		RoleID:    c.AppRoleID,
		SecretID:  c.AppRoleSecretID,
		MountPath: c.AppRoleLoginPath,
	}
}

// Validate checks the Config to confirm that all required values are initialized
func (c *Config) Validate() error {
	if c.SpaceName == emptyStr {
		return fmt.Errorf(errFmtConfigValidation, propertySpaceName)
	}
	if c.SecretPath == emptyStr {
		return fmt.Errorf(errFmtConfigValidation, propertySecretPath)
	}

	return nil
}
