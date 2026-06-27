package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisposableEmailBlocklistLoaded(t *testing.T) {
	require.True(t, DisposableEmailBlocklistReady())
	assert.Greater(t, DisposableEmailBlocklistCount(), 0)
}

func TestValidateReceiversNotDisposable(t *testing.T) {
	require.True(t, DisposableEmailBlocklistReady())

	err := validateReceiversNotDisposable([]string{"user@gmail.com"})
	require.NoError(t, err)

	err = validateReceiversNotDisposable([]string{"user@mailinator.com"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disposable email address is not allowed")
}

func TestValidateReceiversNotDisposableRejectsCustomBlacklist(t *testing.T) {
	require.True(t, DisposableEmailBlocklistReady())

	customEmailDomainBlacklist = parseCustomEmailDomainBlacklist("custom-blocked.example")
	defer func() { customEmailDomainBlacklist = nil }()

	err := validateReceiversNotDisposable([]string{"user@custom-blocked.example"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email domain is not allowed")
}

func TestSendEmailRejectsDisposableAddress(t *testing.T) {
	withSMTPSettings(t)

	SMTPServer = "smtp.example.com"
	SMTPPort = 587
	SMTPAccount = "sender@example.com"
	SMTPFrom = "sender@example.com"
	SMTPToken = "secret"
	SystemName = "New API"

	err := SendEmail("Verification", "user@mailinator.com", "<p>123456</p>")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disposable email address is not allowed")
}
