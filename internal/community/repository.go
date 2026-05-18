package community

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SavePost(p *model.Post) error {
	return r.db.Create(p).Error
}

func (r *Repository) FindPostByID(id int64) (*model.Post, error) {
	var p model.Post
	err := r.db.Preload("User").Preload("Images").Preload("Hashtags").
		Where("id = ? AND deleted_at IS NULL", id).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.New("게시글을 찾을 수 없습니다.")
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) FindPosts(page, size int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	r.db.Model(&model.Post{}).Where("deleted_at IS NULL").Count(&total)

	err := r.db.Preload("User").Preload("Images").Preload("Hashtags").
		Where("deleted_at IS NULL").
		Order("created_date DESC").
		Offset(page * size).Limit(size).
		Find(&posts).Error
	return posts, total, err
}

func (r *Repository) FindMyPosts(userID int64, page, size int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	r.db.Model(&model.Post{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&total)

	err := r.db.Preload("User").Preload("Images").Preload("Hashtags").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_date DESC").
		Offset(page * size).Limit(size).
		Find(&posts).Error
	return posts, total, err
}

func (r *Repository) UpdatePost(p *model.Post) error {
	return r.db.Save(p).Error
}

func (r *Repository) SoftDeletePost(id int64) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *Repository) DeleteHashtagsByPostID(postID int64) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.PostHashtag{}).Error
}

func (r *Repository) DeleteImagesByKeys(postID int64, keys []string) error {
	return r.db.Where("post_id = ? AND image_key IN ?", postID, keys).Delete(&model.PostImage{}).Error
}

func (r *Repository) MaxOrderNum(postID int64) int {
	var maxOrder int
	r.db.Model(&model.PostImage{}).Where("post_id = ?", postID).
		Select("COALESCE(MAX(order_num), 0)").Scan(&maxOrder)
	return maxOrder
}

func (r *Repository) SearchHashtags(keyword string) ([]string, error) {
	var hashtags []string
	err := r.db.Model(&model.PostHashtag{}).
		Select("DISTINCT hashtag").
		Where("hashtag LIKE ?", "%"+keyword+"%").
		Limit(200).
		Pluck("hashtag", &hashtags).Error
	return hashtags, err
}

// Likes
func (r *Repository) CountLikesByPostID(postID int64) int64 {
	var count int64
	r.db.Model(&model.PostLike{}).Where("post_id = ?", postID).Count(&count)
	return count
}

func (r *Repository) ExistsLike(postID, userID int64) bool {
	var count int64
	r.db.Model(&model.PostLike{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count)
	return count > 0
}

func (r *Repository) SaveLike(like *model.PostLike) error {
	return r.db.Create(like).Error
}

func (r *Repository) DeleteLike(postID, userID int64) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostLike{}).Error
}
