package model

import "time"

type ItemCategory string

const (
	ItemCategoryFood           ItemCategory = "FOOD"
	ItemCategorySnack          ItemCategory = "SNACK"
	ItemCategorySupplement     ItemCategory = "SUPPLEMENT"
	ItemCategoryDeworming      ItemCategory = "DEWORMING"
	ItemCategoryFleaPrevention ItemCategory = "FLEA_PREVENTION"
	ItemCategoryShampoo        ItemCategory = "SHAMPOO"
	ItemCategoryConditioner    ItemCategory = "CONDITIONER"
	ItemCategoryToothbrush     ItemCategory = "TOOTHBRUSH"
	ItemCategoryToothpaste     ItemCategory = "TOOTHPASTE"
	ItemCategoryDentalChew     ItemCategory = "DENTAL_CHEW"
	ItemCategoryComb           ItemCategory = "COMB"
	ItemCategoryNailClipper    ItemCategory = "NAIL_CLIPPER"
	ItemCategoryEarCleaner     ItemCategory = "EAR_CLEANER"
	ItemCategoryEyeDrops       ItemCategory = "EYE_DROPS"
	ItemCategoryGrooming       ItemCategory = "GROOMING"
	ItemCategoryDiaper         ItemCategory = "DIAPER"
	ItemCategoryPoopPad        ItemCategory = "POOP_PAD"
	ItemCategoryPoopBag        ItemCategory = "POOP_BAG"
	ItemCategoryToy            ItemCategory = "TOY"
	ItemCategoryTrainingTool   ItemCategory = "TRAINING_TOOL"
	ItemCategoryLeash          ItemCategory = "LEASH"
	ItemCategoryHarness        ItemCategory = "HARNESS"
	ItemCategoryCollar         ItemCategory = "COLLAR"
	ItemCategoryCarrier        ItemCategory = "CARRIER"
	ItemCategoryStroller       ItemCategory = "STROLLER"
	ItemCategoryClothing       ItemCategory = "CLOTHING"
	ItemCategoryShoes          ItemCategory = "SHOES"
	ItemCategoryBed            ItemCategory = "BED"
	ItemCategoryCage           ItemCategory = "CAGE"
	ItemCategoryBowl           ItemCategory = "BOWL"
	ItemCategoryPerfume        ItemCategory = "PERFUME"
	ItemCategoryCamera         ItemCategory = "CAMERA"
	ItemCategoryFeeder         ItemCategory = "FEEDER"
	ItemCategoryWaterDispenser ItemCategory = "WATER_DISPENSER"
	ItemCategoryHealthChecker  ItemCategory = "HEALTH_CHECKER"
	ItemCategoryOther          ItemCategory = "OTHER"
)

type PetItem struct {
	ID               int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	PetID            int64        `gorm:"column:pet_id;not null" json:"petId"`
	Pet              Pet          `gorm:"foreignKey:PetID" json:"-"`
	Name             string       `gorm:"column:name;not null;size:100" json:"name"`
	Category         ItemCategory `gorm:"column:category;type:varchar(50);not null" json:"category"`
	PurchaseCycleDays int         `gorm:"column:purchase_cycle_days;not null" json:"purchaseCycleDays"`
	PurchaseUrl      string       `gorm:"column:purchase_url;size:500" json:"purchaseUrl"`
	ImageKey         string       `gorm:"column:image_key;size:255" json:"imageKey"`
	BaseTimeEntity
}

func (PetItem) TableName() string { return "pet_items" }

type PetItemPurchase struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PetItemID   int64     `gorm:"column:pet_item_id;not null;index:idx_pet_item_purchases_pet_item_id" json:"petItemId"`
	PetItem     PetItem   `gorm:"foreignKey:PetItemID" json:"-"`
	PurchasedAt time.Time `gorm:"column:purchased_at;not null" json:"purchasedAt"`
	BaseTimeEntity
}

func (PetItemPurchase) TableName() string { return "pet_item_purchases" }
