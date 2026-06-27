package service

import (
	"fmt"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/model"

	"github.com/bytedance/gopkg/util/gopool"
)

func init() {
	model.RegisterNewUserCreatedHandler(NotifyNewUserRegistered)
}

func NotifyNewUserRegistered(event model.NewUserCreatedEvent) {
	barkURL := strings.TrimSpace(common.NewUserBarkURL)
	if barkURL == "" {
		return
	}

	gopool.Go(func() {
		content := fmt.Sprintf("用户名: %s", event.Username)
		if event.Email != "" {
			content += fmt.Sprintf("\n邮箱: %s", event.Email)
		}
		content += fmt.Sprintf("\n用户ID: %d", event.UserId)

		notification := dto.NewNotify(dto.NotifyTypeNewUser, "新用户注册", content, nil)
		if err := sendBarkNotify(barkURL, notification); err != nil {
			common.SysLog(fmt.Sprintf("failed to send new user bark notification for user %d: %s", event.UserId, err.Error()))
		}
	})
}
