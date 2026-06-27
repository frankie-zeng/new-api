package common

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/billionverify/disposable"
)

var disposableEmailBlocklistReady int32

func init() {
	if disposable.Count() > 0 {
		atomic.StoreInt32(&disposableEmailBlocklistReady, 1)
	}
}

func DisposableEmailBlocklistReady() bool {
	return atomic.LoadInt32(&disposableEmailBlocklistReady) == 1
}

func DisposableEmailBlocklistCount() int {
	return disposable.Count()
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
		if disposable.IsEmail(addr) {
			return fmt.Errorf("disposable email address is not allowed: %s", addr)
		}
	}
	return nil
}
