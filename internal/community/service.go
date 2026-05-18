package community

import (
	"context"
	"math"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
)

type Service struct {
	repo      *Repository
	likeRedis *likeRedis
}

func NewService(repo *Repository, rdb *goredis.Client) *Service {
	return &Service{repo: repo, likeRedis: newLikeRedis(rdb)}
}

// Post Services
func (s *Service) CreatePost(userID int64, req PostCreateRequest) (int64, error) {
	p := &model.Post{
		UserID:  userID,
		Content: req.Content,
	}

	for _, tag := range req.Hashtags {
		p.Hashtags = append(p.Hashtags, model.PostHashtag{Hashtag: tag})
	}
	for i, key := range req.ImageKeys {
		p.Images = append(p.Images, model.PostImage{ImageKey: key, OrderNum: i + 1})
	}

	if err := s.repo.SavePost(p); err != nil {
		return 0, err
	}
	return p.ID, nil
}

func (s *Service) GetPost(postID, currentUserID int64) (*PostResponse, error) {
	p, err := s.repo.FindPostByID(postID)
	if err != nil {
		return nil, err
	}
	likeCount, isLiked := s.getLikeInfo(postID, currentUserID)
	return toPostResponse(p, likeCount, isLiked), nil
}

func (s *Service) GetPosts(currentUserID int64, page, size int) (*PostListResponse, error) {
	posts, total, err := s.repo.FindPosts(page, size)
	if err != nil {
		return nil, err
	}
	return s.buildListResponse(posts, total, currentUserID, page, size), nil
}

func (s *Service) GetMyPosts(userID int64, page, size int) (*PostListResponse, error) {
	posts, total, err := s.repo.FindMyPosts(userID, page, size)
	if err != nil {
		return nil, err
	}
	return s.buildListResponse(posts, total, userID, page, size), nil
}

func (s *Service) UpdatePost(postID, userID int64, req PostUpdateRequest) error {
	p, err := s.repo.FindPostByID(postID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return pnerrors.NewWithCode(403, "게시글 수정 권한이 없습니다.")
	}

	p.Content = req.Content

	if err = s.repo.DeleteHashtagsByPostID(postID); err != nil {
		return err
	}
	p.Hashtags = nil
	for _, tag := range req.Hashtags {
		p.Hashtags = append(p.Hashtags, model.PostHashtag{PostID: postID, Hashtag: tag})
	}

	if len(req.DeleteImageKeys) > 0 {
		if err = s.repo.DeleteImagesByKeys(postID, req.DeleteImageKeys); err != nil {
			return err
		}
	}

	nextOrder := s.repo.MaxOrderNum(postID)
	for _, key := range req.AddImageKeys {
		nextOrder++
		p.Images = append(p.Images, model.PostImage{PostID: postID, ImageKey: key, OrderNum: nextOrder})
	}

	return s.repo.UpdatePost(p)
}

func (s *Service) DeletePost(postID, userID int64) error {
	p, err := s.repo.FindPostByID(postID)
	if err != nil {
		return err
	}
	if p.UserID != userID {
		return pnerrors.NewWithCode(403, "게시글 삭제 권한이 없습니다.")
	}
	return s.repo.SoftDeletePost(postID)
}

func (s *Service) SearchHashtags(keyword string) ([]string, error) {
	return s.repo.SearchHashtags(keyword)
}

// Like Services
func (s *Service) ToggleLike(postID, userID int64) (*PostLikeToggleResponse, error) {
	if _, err := s.repo.FindPostByID(postID); err != nil {
		return nil, err
	}

	ctx := context.Background()

	// 캐시 미스 시 DB에서 초기화
	_, countCached := s.likeRedis.getLikeCount(ctx, postID)
	if !countCached {
		dbCount := s.repo.CountLikesByPostID(postID)
		dbLiked := s.repo.ExistsLike(postID, userID)
		s.likeRedis.initCache(ctx, postID, userID, dbCount, dbLiked)
	}

	liked, count, err := s.likeRedis.toggle(ctx, postID, userID)
	if err != nil {
		return nil, err
	}

	// DB 동기화
	if liked {
		_ = s.repo.SaveLike(&model.PostLike{PostID: postID, UserID: userID})
	} else {
		_ = s.repo.DeleteLike(postID, userID)
	}

	return &PostLikeToggleResponse{Liked: liked, LikeCount: count}, nil
}

// Helpers
func (s *Service) getLikeInfo(postID, userID int64) (int64, bool) {
	ctx := context.Background()

	count, countOk := s.likeRedis.getLikeCount(ctx, postID)
	liked, likedOk := s.likeRedis.getLikedStatus(ctx, userID, postID)

	if countOk && likedOk {
		return count, liked
	}
	return s.repo.CountLikesByPostID(postID), s.repo.ExistsLike(postID, userID)
}

func (s *Service) buildListResponse(posts []model.Post, total int64, currentUserID int64, page, size int) *PostListResponse {
	responses := make([]PostResponse, 0, len(posts))
	for _, p := range posts {
		likeCount, isLiked := s.getLikeInfo(p.ID, currentUserID)
		responses = append(responses, *toPostResponse(&p, likeCount, isLiked))
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	return &PostListResponse{
		Posts:       responses,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalCount:  total,
	}
}

func toPostResponse(p *model.Post, likeCount int64, isLiked bool) *PostResponse {
	imageKeys := make([]string, 0, len(p.Images))
	for _, img := range p.Images {
		imageKeys = append(imageKeys, img.ImageKey)
	}

	hashtags := make([]string, 0, len(p.Hashtags))
	for _, tag := range p.Hashtags {
		hashtags = append(hashtags, tag.Hashtag)
	}

	return &PostResponse{
		PostID:         p.ID,
		UserID:         p.User.ID,
		UserNickname:   p.User.NickName,
		UserProfileUrl: p.User.ProfileUrl,
		Content:        p.Content,
		ImageKeys:      imageKeys,
		Hashtags:       hashtags,
		CreatedDate:    p.CreatedDate,
		LikeCount:      likeCount,
		IsLiked:        isLiked,
	}
}
