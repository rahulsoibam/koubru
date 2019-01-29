package googlejwt

import "time"

var (
	// MaxTokenLifetime is one day
	MaxTokenLifetime = time.Second * 86400

	// ClockSkew - five minutes
	ClockSkew = time.Minute * 5

	// Issuers is the allowed oauth token issuers
	Issuers = []string{
		"accounts.google.com",
		"https://accounts.google.com",
	}
)

// GoogleIDTokenVerifier instance
type GoogleIDTokenVerifier struct{}

// VerifyIDToken verifies the ID Token
func (v *GoogleIDTokenVerifier) VerifyIDToken(idToken string, audience []string) error {
	certs, err := GetFederatedSignonCerts()
	if err != nil {
		return err
	}
	return VerifySignedJWTWithCerts(idToken, certs, audience, Issuers, MaxTokenLifetime)
}
