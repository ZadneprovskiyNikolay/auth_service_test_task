package auth

import (
	"fmt"
	"time"
)

func RefreshRequestNewIPEmail(newIP string) []byte {
	timeStr := time.Now().In(time.UTC).Format("2006-01-02 15:04:05") + " (UTC)"
	return []byte(fmt.Sprintf(
		`Обнаружен вход в ваш аккаунт с нового ip-адреса.
Время: %v
IP-адрес: %s`,
		timeStr, newIP))
}
