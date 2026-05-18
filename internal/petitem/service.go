package petitem

import (
	"time"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req PetItemCreateRequest) (*PetItemResponse, error) {
	item := &model.PetItem{
		PetID:             req.PetID,
		Name:              req.Name,
		Category:          req.Category,
		PurchaseCycleDays: req.PurchaseCycleDays,
		PurchaseUrl:       req.PurchaseUrl,
		ImageKey:          req.ImageKey,
	}
	if err := s.repo.SaveItem(item); err != nil {
		return nil, err
	}
	return toItemResponse(item, nil), nil
}

func (s *Service) GetByPetID(petID int64, category *model.ItemCategory) ([]PetItemResponse, error) {
	items, err := s.repo.FindItemsByPetID(petID, category)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}

	latestMap, err := s.repo.FindLatestPurchasesByItemIDs(ids)
	if err != nil {
		return nil, err
	}

	responses := make([]PetItemResponse, 0, len(items))
	for _, item := range items {
		last := latestMap[item.ID]
		responses = append(responses, *toItemResponse(&item, last))
	}
	return responses, nil
}

func (s *Service) GetDetail(petItemID int64) (*PetItemResponse, error) {
	item, err := s.repo.FindItemByID(petItemID)
	if err != nil {
		return nil, err
	}

	latestMap, err := s.repo.FindLatestPurchasesByItemIDs([]int64{item.ID})
	if err != nil {
		return nil, err
	}

	return toItemResponse(item, latestMap[item.ID]), nil
}

func (s *Service) Update(petItemID int64, req PetItemUpdateRequest) (*PetItemResponse, error) {
	item, err := s.repo.FindItemByID(petItemID)
	if err != nil {
		return nil, err
	}

	item.Name = req.Name
	item.Category = req.Category
	item.PurchaseCycleDays = req.PurchaseCycleDays
	item.PurchaseUrl = req.PurchaseUrl
	item.ImageKey = req.ImageKey

	if err = s.repo.UpdateItem(item); err != nil {
		return nil, err
	}
	return toItemResponse(item, nil), nil
}

func (s *Service) Delete(petItemID int64) error {
	if err := s.repo.DeletePurchasesByItemID(petItemID); err != nil {
		return err
	}
	return s.repo.DeleteItem(petItemID)
}

func (s *Service) RecordPurchase(petItemID int64, req PetItemPurchaseCreateRequest) (*PetItemPurchaseResponse, error) {
	purchasedAt := time.Now()
	if req.PurchasedAt != nil {
		purchasedAt = *req.PurchasedAt
	}

	purchase := &model.PetItemPurchase{
		PetItemID:   petItemID,
		PurchasedAt: purchasedAt,
	}
	if err := s.repo.SavePurchase(purchase); err != nil {
		return nil, err
	}

	return &PetItemPurchaseResponse{
		ID:          purchase.ID,
		PetItemID:   purchase.PetItemID,
		PurchasedAt: purchase.PurchasedAt,
	}, nil
}

func (s *Service) GetPurchaseHistory(petItemID int64) ([]PetItemPurchaseResponse, error) {
	purchases, err := s.repo.FindPurchasesByItemID(petItemID)
	if err != nil {
		return nil, err
	}
	responses := make([]PetItemPurchaseResponse, 0, len(purchases))
	for _, p := range purchases {
		responses = append(responses, PetItemPurchaseResponse{
			ID:          p.ID,
			PetItemID:   p.PetItemID,
			PurchasedAt: p.PurchasedAt,
		})
	}
	return responses, nil
}

func (s *Service) DeletePurchase(purchaseID int64) error {
	return s.repo.DeletePurchase(purchaseID)
}

func toItemResponse(item *model.PetItem, lastPurchasedAt *time.Time) *PetItemResponse {
	var nextPurchaseAt *time.Time
	if lastPurchasedAt != nil {
		next := lastPurchasedAt.AddDate(0, 0, item.PurchaseCycleDays)
		nextPurchaseAt = &next
	}
	return &PetItemResponse{
		PetItemID:         item.ID,
		PetID:             item.PetID,
		Name:              item.Name,
		Category:          item.Category,
		PurchaseCycleDays: item.PurchaseCycleDays,
		PurchaseUrl:       item.PurchaseUrl,
		ImageUrl:          item.ImageKey,
		LastPurchasedAt:   lastPurchasedAt,
		NextPurchaseAt:    nextPurchaseAt,
	}
}
