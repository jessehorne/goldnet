package util

import "fmt"

func NewSystemMessage(prefix, data string) string {
	return fmt.Sprintf("(%s) %s", prefix, data)
}

func NewChatMessage(username, data string) string {
	return fmt.Sprintf("%s - %s", username, data)
}
