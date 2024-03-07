// Package vault provides the Vault implementation for the SecretGetter interface
package vault

import (
	"fmt"

	"github.com/riotgames/vault-go-client"
)

// Misc Consts
const (
	emptyStr = ""
)

// Vault Consts
const (
	defaultURL = "https://your.vault.server.url/"
)

// Component provides a lightweight wrapper around the public riotgames/vault-go-client package for integration with
// internally hosted Vault offerings
type Component struct {
	config      *Config
	vaultClient *vault.Client
}

// NewComponent provides a method to initialize a base Component struct.
//
// Refer to options.go for auth methods. The recommended method is WithDefaultChainAuthProvider, which attempts to
// authenticate with Vault using both IAM login and AppRole login. This enables Vault auth when running your
// service in AWS (IAM login) or on your local machine (AppRole for testing)
//
// This initializes a Vault client pointing towards the Vault server URL given by the defaultURL const defined in this package
// Refer to the config method Config.Validate() for required client configuration values
//
// Refer to the Hashicorp documentation and vault-go-client for details on the various auth methods available
// https://developer.hashicorp.com/vault/docs/auth
// https://github.com/RiotGames/vault-go-client
func NewComponent(config *Config, options ...Option) *Component {
	opts := createDefaultOptions(config)
	for _, option := range options {
		option(opts)
	}

	if validateErr := config.Validate(); validateErr != nil {
		opts.log.Fatal(validateErr.Error())
	}

	vaultClientConfig := vault.DefaultConfig()
	vaultClientConfig.Address = opts.url
	vaultClient, err := vault.NewClient(vaultClientConfig)
	if err != nil {
		opts.log.Fatal(fmt.Sprintf("failed to initialize Vault client, err: %s", err.Error()))
	}
	if authErr := opts.authProvider(vaultClient); authErr != nil {
		opts.log.Fatal(fmt.Sprintf("failed to authenticate Vault client, err: %s", authErr.Error()))
	}

	return &Component{
		config:      config,
		vaultClient: vaultClient,
	}
}

// GetSecret fetches a secret from the Vault server
func (c *Component) GetSecret(key string) (secret string, err error) {
	secretResp, err := c.vaultClient.KV2.Get(vault.KV2GetOptions{
		MountPath:  c.config.SecretMountPath,
		SecretPath: c.config.SecretPath,
	})
	if err != nil {
		return emptyStr, err
	}

	secretMap := secretResp.Data["data"].(map[string]interface{})
	return secretMap[key].(string), nil
}
