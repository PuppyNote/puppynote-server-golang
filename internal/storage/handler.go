package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"puppynote/config"
	"puppynote/pkg/middleware"
	"puppynote/pkg/response"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BucketKind string

const (
	BucketKindPuppyProfile  BucketKind = "PUPPY_PROFILE"
	BucketKindWalkPhoto     BucketKind = "WALK_PHOTO"
	BucketKindPetItemPhoto  BucketKind = "PET_ITEM_PHOTO"
	BucketKindUserProfile   BucketKind = "USER_PROFILE"
	BucketKindCommunityPost BucketKind = "COMMUNITY_POST"
)

var folderMap = map[BucketKind]string{
	BucketKindPuppyProfile:  "puppy-profile",
	BucketKindWalkPhoto:     "puppy-walk",
	BucketKindPetItemPhoto:  "puppy-item",
	BucketKindUserProfile:   "user-profile",
	BucketKindCommunityPost: "community",
}

type Handler struct {
	s3Client *s3.Client
}

func NewHandler() *Handler {
	cfg := config.AppConfig.AWS
	awsCfg, _ := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey, cfg.SecretKey, "",
		)),
	)
	return &Handler{s3Client: s3.NewFromConfig(awsCfg)}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/api/v1/storage", middleware.JWTAuth())
	g.POST("/:bucketKind", h.upload)
}

func (h *Handler) upload(c *gin.Context) {
	bucketKind := BucketKind(c.Param("bucketKind"))
	folder, ok := folderMap[bucketKind]
	if !ok {
		response.BadRequest(c, "지원하지 않는 bucketKind입니다.")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "파일이 필요합니다.")
		return
	}
	defer file.Close()

	if header.Size > 10<<20 {
		response.BadRequest(c, "파일 크기는 10MB를 초과할 수 없습니다.")
		return
	}

	key, err := h.uploadToS3(file, header, folder)
	if err != nil {
		response.InternalServerError(c, "파일 업로드에 실패했습니다.")
		return
	}

	cloudfrontURL := buildCloudFrontURL(key)
	response.OK(c, cloudfrontURL)
}

func (h *Handler) uploadToS3(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	cfg := config.AppConfig.AWS
	contentType := getContentType(ext)

	_, err := h.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	return key, err
}

func buildCloudFrontURL(key string) string {
	domain := config.AppConfig.AWS.CloudfrontDomain
	return fmt.Sprintf("%s/%s", strings.TrimRight(domain, "/"), key)
}

func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
