package batch

import (
	"log"
	"strconv"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	"github.com/PuppyNote/puppynote-server-golang/pkg/push"
	"gorm.io/gorm"
)

type PetItemJob struct{ db *gorm.DB }

func NewPetItemJob(db *gorm.DB) *PetItemJob { return &PetItemJob{db: db} }

type upcomingPurchase struct {
	PurchaseID        int64
	PetItemID         int64
	ItemName          string
	PetID             int64
	PurchaseCycleDays int
}

func (j *PetItemJob) Run() {
	var results []upcomingPurchase
	j.db.Raw(`
		SELECT pip.id AS purchase_id, pip.pet_item_id, pi.name AS item_name,
		       pi.pet_id, pi.purchase_cycle_days
		FROM pet_item_purchases pip
		JOIN pet_items pi ON pip.pet_item_id = pi.id
		WHERE pip.id IN (SELECT MAX(id) FROM pet_item_purchases GROUP BY pet_item_id)
		AND DATE(pip.purchased_at + INTERVAL pi.purchase_cycle_days DAY)
		    = DATE(NOW() + INTERVAL 1 DAY)
	`).Scan(&results)

	if len(results) == 0 {
		return
	}

	var notifications []push.Notification
	var histories []model.AlertHistory

	for _, r := range results {
		var members []model.FamilyMember
		j.db.Where("pet_id = ? AND status = ?", r.PetID, model.FamilyStatusDone).Find(&members)
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

		body := r.ItemName + " 구매 주기가 내일 도래해요! 미리 준비해두세요 🛒"
		itemIDStr := strconv.FormatInt(r.PetItemID, 10)

		for _, uid := range userIDs {
			notifications = append(notifications, push.Notification{
				PushToken:            tokenMap[uid],
				Body:                 body,
				AlertDestinationType: model.AlertDestPetItem,
				AlertDestinationInfo: itemIDStr,
			})
			histories = append(histories, model.AlertHistory{
				UserID:               uid,
				AlertDescription:     body,
				AlertHistoryStatus:   model.AlertHistoryStatusUnchecked,
				AlertDestinationType: model.AlertDestPetItem,
				AlertDestinationInfo: itemIDStr,
			})
		}
	}

	if len(histories) > 0 {
		j.db.Create(&histories)
	}
	push.Send(notifications)
	log.Printf("[PetItemJob] %d 알림 발송 완료", len(notifications))
}
