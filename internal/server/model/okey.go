package model

import (
	"fmt"
)

func GetMsgOffKey(mid int32) string {
	return fmt.Sprintf("MESSAGES_OFF_LIST:%d", mid)
}

func GetBanKey(mid int32) string {
	return fmt.Sprintf("GAG|BAN|%d", mid)
}

func GetCDKey(mid int32) string {
	return fmt.Sprintf("MARCH|USER_LEFT_TIME|%d", mid)
}

func GetTimesKey(mid int32) string {
	return fmt.Sprintf("MARCH|USER_USE_TIMES|%d", mid)
}

func GetVipKey(mid int32) string {
	return fmt.Sprintf("MARCH|USER_VIP_LEVEL|%d", mid)
}

func GetRuleKey() string {
	return "MARCH|CHAT_RULES"
}
