package petitem

import (
	"time"

	"puppynote/internal/model"
)

type PetItemCreateRequest struct {
	PetID             int64              `json:"petId" binding:"required"`
	Name              string             `json:"name" binding:"required"`
	Category          model.ItemCategory `json:"category" binding:"required"`
	PurchaseCycleDays int                `json:"purchaseCycleDays" binding:"required,min=1"`
	PurchaseUrl       string             `json:"purchaseUrl"`
	ImageKey          string             `json:"imageKey"`
}

type PetItemUpdateRequest struct {
	Name              string             `json:"name" binding:"required"`
	Category          model.ItemCategory `json:"category" binding:"required"`
	PurchaseCycleDays int                `json:"purchaseCycleDays" binding:"required,min=1"`
	PurchaseUrl       string             `json:"purchaseUrl"`
	ImageKey          string             `json:"imageKey"`
}

type PetItemResponse struct {
	PetItemID         int64              `json:"petItemId"`
	PetID             int64              `json:"petId"`
	Name              string             `json:"name"`
	Category          model.ItemCategory `json:"category"`
	PurchaseCycleDays int                `json:"purchaseCycleDays"`
	PurchaseUrl       string             `json:"purchaseUrl"`
	ImageUrl          string             `json:"imageUrl"`
	LastPurchasedAt   *time.Time         `json:"lastPurchasedAt"`
	NextPurchaseAt    *time.Time         `json:"nextPurchaseAt"`
}

type PetItemPurchaseCreateRequest struct {
	PurchasedAt *time.Time `json:"purchasedAt"`
}

type PetItemPurchaseResponse struct {
	ID          int64     `json:"id"`
	PetItemID   int64     `json:"petItemId"`
	PurchasedAt time.Time `json:"purchasedAt"`
}
