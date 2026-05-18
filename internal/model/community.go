package model

import "time"

type Post struct {
	ID        int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64       `gorm:"column:user_id;not null" json:"userId"`
	User      User        `gorm:"foreignKey:UserID" json:"-"`
	Content   string      `gorm:"column:content;type:text" json:"content"`
	Hashtags  []PostHashtag `gorm:"foreignKey:PostID" json:"hashtags,omitempty"`
	Images    []PostImage `gorm:"foreignKey:PostID" json:"images,omitempty"`
	DeletedAt *time.Time  `gorm:"column:deleted_at" json:"deletedAt,omitempty"`
	BaseTimeEntity
}

func (Post) TableName() string { return "posts" }

type PostHashtag struct {
	PostID  int64  `gorm:"column:post_id;not null" json:"postId"`
	Hashtag string `gorm:"column:hashtag" json:"hashtag"`
}

func (PostHashtag) TableName() string { return "post_hashtags" }

type PostImage struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID   int64  `gorm:"column:post_id;not null" json:"postId"`
	Post     Post   `gorm:"foreignKey:PostID" json:"-"`
	ImageKey string `gorm:"column:image_key;not null" json:"imageKey"`
	OrderNum int    `gorm:"column:order_num;not null" json:"orderNum"`
	BaseTimeEntity
}

func (PostImage) TableName() string { return "post_images" }

type PostLike struct {
	ID     int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID int64 `gorm:"column:post_id;not null;uniqueIndex:idx_post_user" json:"postId"`
	UserID int64 `gorm:"column:user_id;not null;uniqueIndex:idx_post_user" json:"userId"`
	BaseTimeEntity
}

func (PostLike) TableName() string { return "post_likes" }
