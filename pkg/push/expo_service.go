package push

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
)

const expoAPI = "https://exp.host/--/api/v2/push/send"

type PushMessage struct {
	To    string            `json:"to"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Sound string            `json:"sound"`
	Data  map[string]string `json:"data,omitempty"`
}

type Notification struct {
	PushToken            string
	Body                 string
	AlertDestinationType model.AlertDestinationType
	AlertDestinationInfo string
}

var titleMap = map[model.AlertDestinationType]string{
	model.AlertDestWalk:        "산책 알림",
	model.AlertDestPetItem:     "용품 구매 알림",
	model.AlertDestFriend:      "친구 추가 알림",
	model.AlertDestFamilyInvite: "가족 초대 알림",
	model.AlertDestDailyReport: "일일 리포트",
}

func Send(notifications []Notification) {
	if len(notifications) == 0 {
		return
	}

	messages := make([]PushMessage, 0, len(notifications))
	for _, n := range notifications {
		if n.PushToken == "" {
			continue
		}
		title := titleMap[n.AlertDestinationType]
		if title == "" {
			title = "퍼피노트 알림"
		}
		messages = append(messages, PushMessage{
			To:    n.PushToken,
			Title: title,
			Body:  n.Body,
			Sound: "default",
			Data: map[string]string{
				"alert_destination_type": string(n.AlertDestinationType),
				"alert_destination_info": n.AlertDestinationInfo,
			},
		})
	}

	if len(messages) == 0 {
		return
	}

	payload, err := json.Marshal(messages)
	if err != nil {
		log.Printf("FCM 메시지 직렬화 실패: %v", err)
		return
	}

	resp, err := http.Post(expoAPI, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("Expo 푸시 발송 실패: %v", err)
		return
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}
