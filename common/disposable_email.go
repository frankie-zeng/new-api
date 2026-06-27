package common

import (
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/billionverify/disposable"
)

var (
	disposableEmailBlocklistReady int32
	customEmailDomainBlacklist    map[string]struct{}
)

func init() {
	if disposable.Count() > 0 {
		atomic.StoreInt32(&disposableEmailBlocklistReady, 1)
	}
	if raw := strings.TrimSpace(os.Getenv("EMAIL_DOMAIN_BLACKLIST")); raw != "" {
		customEmailDomainBlacklist = parseCustomEmailDomainBlacklist(raw)
	}
}

func DisposableEmailBlocklistReady() bool {
	return atomic.LoadInt32(&disposableEmailBlocklistReady) == 1
}

func DisposableEmailBlocklistCount() int {
	return disposable.Count()
}

func parseCustomEmailDomainBlacklist(raw string) map[string]struct{} {
	raw = strings.ReplaceAll(raw, "\n", ",")
	parts := strings.Split(raw, ",")
	blocklist := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		domain := strings.TrimSpace(strings.ToLower(strings.TrimPrefix(part, "@")))
		if domain == "" {
			continue
		}
		blocklist[domain] = struct{}{}
	}
	return blocklist
}

func isCustomEmailDomainBlacklisted(email string) bool {
	if len(customEmailDomainBlacklist) == 0 {
		return false
	}
	parts := strings.Split(strings.TrimSpace(email), "@")
	if len(parts) != 2 || parts[1] == "" {
		return false
	}
	domain := strings.ToLower(parts[1])
	if _, ok := customEmailDomainBlacklist[domain]; ok {
		return true
	}
	for blocked := range customEmailDomainBlacklist {
		if strings.HasSuffix(domain, "."+blocked) {
			return true
		}
	}
	return false
}

func validateReceiversNotDisposable(receivers []string) error {
	if !DisposableEmailBlocklistReady() {
		return fmt.Errorf("disposable email blocklist is not available")
	}
	for _, receiver := range receivers {
		addr := strings.TrimSpace(receiver)
		if addr == "" {
			continue
		}
		if isCustomEmailDomainBlacklisted(addr) {
			return fmt.Errorf("email domain is not allowed: %s", addr)
		}
		if disposable.IsEmail(addr) {
			return fmt.Errorf("disposable email address is not allowed: %s", addr)
		}
	}
	return nil
}
