package community

import "time"

type PostCreateRequest struct {
	Content   string   `json:"content" binding:"required,max=2000"`
	Hashtags  []string `json:"hashtags"`
	ImageKeys []string `json:"imageKeys"`
}

type PostUpdateRequest struct {
	Content         string   `json:"content" binding:"required,max=2000"`
	Hashtags        []string `json:"hashtags"`
	AddImageKeys    []string `json:"addImageKeys"`
	DeleteImageKeys []string `json:"deleteImageKeys"`
}

type PostResponse struct {
	PostID         int64      `json:"postId"`
	UserID         int64      `json:"userId"`
	UserNickname   string     `json:"userNickname"`
	UserProfileUrl string     `json:"userProfileUrl"`
	Content        string     `json:"content"`
	ImageKeys      []string   `json:"imageKeys"`
	Hashtags       []string   `json:"hashtags"`
	CreatedDate    time.Time  `json:"createdDate"`
	LikeCount      int64      `json:"likeCount"`
	IsLiked        bool       `json:"isLiked"`
}

type PostListResponse struct {
	Posts       []PostResponse `json:"posts"`
	CurrentPage int            `json:"currentPage"`
	TotalPages  int            `json:"totalPages"`
	TotalCount  int64          `json:"totalCount"`
}

type PostLikeToggleResponse struct {
	Liked     bool  `json:"liked"`
	LikeCount int64 `json:"likeCount"`
}
