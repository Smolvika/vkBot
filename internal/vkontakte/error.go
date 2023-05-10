package vkontakte

import "log"

var (
	placeSendMessageGreeting = "SendMessageGreeting message"
	placeSendMessage         = "SendMessage message"
)

func errorsMessage(place string, err error) {
	log.Printf("error %s: %v\n", place, err)
}

func errorsUnmarshal(err error) {
	log.Printf("error RequestUnmarshal: %v", err)
}
