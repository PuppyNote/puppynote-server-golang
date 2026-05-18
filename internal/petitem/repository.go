package petitem

import (
	"time"

	"puppynote/internal/model"
	pnerrors "puppynote/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveItem(item *model.PetItem) error {
	return r.db.Create(item).Error
}

func (r *Repository) FindItemByID(id int64) (*model.PetItem, error) {
	var item model.PetItem
	if err := r.db.First(&item, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.New("반려동물 용품을 찾을 수 없습니다.")
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindItemsByPetID(petID int64, category *model.ItemCategory) ([]model.PetItem, error) {
	q := r.db.Where("pet_id = ?", petID)
	if category != nil {
		q = q.Where("category = ?", *category)
	}
	var items []model.PetItem
	return items, q.Find(&items).Error
}

func (r *Repository) UpdateItem(item *model.PetItem) error {
	return r.db.Save(item).Error
}

func (r *Repository) DeleteItem(id int64) error {
	return r.db.Delete(&model.PetItem{}, id).Error
}

func (r *Repository) DeleteItemsByPetID(petID int64) error {
	return r.db.Where("pet_id = ?", petID).Delete(&model.PetItem{}).Error
}

// Purchases
func (r *Repository) SavePurchase(p *model.PetItemPurchase) error {
	return r.db.Create(p).Error
}

func (r *Repository) FindLatestPurchasesByItemIDs(ids []int64) (map[int64]*time.Time, error) {
	type result struct {
		PetItemID   int64
		PurchasedAt time.Time
	}
	var results []result
	err := r.db.Model(&model.PetItemPurchase{}).
		Select("pet_item_id, MAX(purchased_at) as purchased_at").
		Where("pet_item_id IN ?", ids).
		Group("pet_item_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	m := make(map[int64]*time.Time, len(results))
	for _, res := range results {
		t := res.PurchasedAt
		m[res.PetItemID] = &t
	}
	return m, nil
}

func (r *Repository) FindPurchasesByItemID(itemID int64) ([]model.PetItemPurchase, error) {
	var purchases []model.PetItemPurchase
	err := r.db.Where("pet_item_id = ?", itemID).Order("purchased_at DESC").Find(&purchases).Error
	return purchases, err
}

func (r *Repository) DeletePurchase(id int64) error {
	return r.db.Delete(&model.PetItemPurchase{}, id).Error
}

func (r *Repository) DeletePurchasesByItemID(itemID int64) error {
	return r.db.Where("pet_item_id = ?", itemID).Delete(&model.PetItemPurchase{}).Error
}
