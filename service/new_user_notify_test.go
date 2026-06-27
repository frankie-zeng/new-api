package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/stretchr/testify/require"
)

func TestNotifyNewUserRegisteredSendsBarkRequest(t *testing.T) {
	InitHttpClient()
	originalURL := common.NewUserBarkURL
	originalAllowPrivateIP := system_setting.GetFetchSetting().AllowPrivateIp
	originalSSRF := system_setting.GetFetchSetting().EnableSSRFProtection
	defer func() {
		common.NewUserBarkURL = originalURL
		system_setting.GetFetchSetting().AllowPrivateIp = originalAllowPrivateIP
		system_setting.GetFetchSetting().EnableSSRFProtection = originalSSRF
	}()
	system_setting.GetFetchSetting().AllowPrivateIp = true
	system_setting.GetFetchSetting().EnableSSRFProtection = false

	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"message":"success"}`))
	}))
	defer server.Close()

	common.NewUserBarkURL = server.URL + "/{{title}}/{{content}}"

	NotifyNewUserRegistered(model.NewUserCreatedEvent{
		UserId:   42,
		Username: "alice",
		Email:    "alice@example.com",
	})

	require.Eventually(t, func() bool {
		return receivedPath != ""
	}, 5*time.Second, 10*time.Millisecond)

	require.Contains(t, receivedPath, "新用户注册")
	require.Contains(t, receivedPath, "alice")
	require.Contains(t, receivedPath, "alice@example.com")
}

func TestNotifyNewUserRegisteredSkipsWhenURLUnset(t *testing.T) {
	originalURL := common.NewUserBarkURL
	defer func() {
		common.NewUserBarkURL = originalURL
	}()

	common.NewUserBarkURL = ""
	NotifyNewUserRegistered(model.NewUserCreatedEvent{
		UserId:   1,
		Username: "bob",
	})
}
