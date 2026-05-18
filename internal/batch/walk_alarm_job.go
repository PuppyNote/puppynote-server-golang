package batch

import (
	"log"
	"strconv"
	"time"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	"github.com/PuppyNote/puppynote-server-golang/pkg/push"
	"gorm.io/gorm"
)

var seoulLoc *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		loc = time.FixedZone("KST", 9*60*60)
	}
	seoulLoc = loc
}

type WalkAlarmJob struct{ db *gorm.DB }

func NewWalkAlarmJob(db *gorm.DB) *WalkAlarmJob { return &WalkAlarmJob{db: db} }

func (j *WalkAlarmJob) Run() {
	now := time.Now().In(seoulLoc)
	timeStr := now.Format("15:04")
	dayStr := weekdayToAlarmDay(now.Weekday())

	var alarms []model.PetWalkAlarm
	j.db.Preload("AlarmDays").Preload("Pet").
		Where("alarm_status = ? AND TIME_FORMAT(alarm_time, '%H:%i') = ?", model.AlarmStatusYes, timeStr).
		Find(&alarms)

	if len(alarms) == 0 {
		return
	}

	var notifications []push.Notification
	var histories []model.AlertHistory

	for _, alarm := range alarms {
		if !containsDay(alarm.AlarmDays, dayStr) {
			continue
		}

		var members []model.FamilyMember
		j.db.Where("pet_id = ? AND status = ?", alarm.PetID, model.FamilyStatusDone).Find(&members)
		if len(members) == 0 {
			continue
		}

		userIDs := make([]int64, 0, len(members))
		for _, m := range members {
			userIDs = append(userIDs, m.UserID)
		}

		var pushList []model.Push
		j.db.Where("user_id IN ?", userIDs).Find(&pushList)

		tokenMap := make(map[int64]string, len(pushList))
		for _, p := range pushList {
			tokenMap[p.UserID] = p.PushToken
		}

		body := alarm.Pet.Name + " 산책 시간이에요! 🐾"
		petIDStr := strconv.FormatInt(alarm.PetID, 10)

		for _, uid := range userIDs {
			notifications = append(notifications, push.Notification{
				PushToken:            tokenMap[uid],
				Body:                 body,
				AlertDestinationType: model.AlertDestWalk,
				AlertDestinationInfo: petIDStr,
			})
			histories = append(histories, model.AlertHistory{
				UserID:               uid,
				AlertDescription:     body,
				AlertHistoryStatus:   model.AlertHistoryStatusUnchecked,
				AlertDestinationType: model.AlertDestWalk,
				AlertDestinationInfo: petIDStr,
			})
		}
	}

	if len(histories) > 0 {
		j.db.Create(&histories)
	}
	push.Send(notifications)
	log.Printf("[WalkAlarmJob] %d 알림 발송 완료", len(notifications))
}

func weekdayToAlarmDay(w time.Weekday) model.AlarmDay {
	switch w {
	case time.Monday:
		return model.AlarmDayMon
	case time.Tuesday:
		return model.AlarmDayTue
	case time.Wednesday:
		return model.AlarmDayWed
	case time.Thursday:
		return model.AlarmDayThu
	case time.Friday:
		return model.AlarmDayFri
	case time.Saturday:
		return model.AlarmDaySat
	default:
		return model.AlarmDaySun
	}
}

func containsDay(days []model.PetAlarmDay, day model.AlarmDay) bool {
	for _, d := range days {
		if d.AlarmDay == day {
			return true
		}
	}
	return false
}
