package secrets

// SecretGetter is a common interface provided for interaction with Secrets Management stores
// Refer to subpackages in this directory for service-specific implementations
type SecretGetter interface {
	GetSecret(key string) (secret string, err error)
}
