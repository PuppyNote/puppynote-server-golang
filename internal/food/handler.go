package food

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/PuppyNote/puppynote-server-golang/config"
	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	"github.com/PuppyNote/puppynote-server-golang/pkg/middleware"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FoodAskRequest struct {
	Question string `json:"question" binding:"required"`
}

type FoodResponse struct {
	ID          int64             `json:"id"`
	Question    string            `json:"question"`
	Answer      string            `json:"answer"`
	SafetyLevel model.SafetyLevel `json:"safetyLevel"`
}

type FoodListResponse struct {
	Content    []FoodResponse `json:"content"`
	Page       int            `json:"page"`
	TotalCount int64          `json:"totalCount"`
}

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/api/v1/foods", middleware.JWTAuth())
	{
		g.GET("", h.list)
		g.POST("/ai", h.askAI)
	}
}

func (h *Handler) list(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	var histories []model.FoodChatHistory
	var total int64

	q := h.db.Model(&model.FoodChatHistory{})
	if keyword := c.Query("question"); keyword != "" {
		q = q.Where("question LIKE ?", "%"+keyword+"%")
	}
	q.Count(&total)
	q.Order("created_date DESC").Offset(page * size).Limit(size).Find(&histories)

	content := make([]FoodResponse, 0, len(histories))
	for _, h := range histories {
		content = append(content, FoodResponse{
			ID:          h.ID,
			Question:    h.Question,
			Answer:      h.Answer,
			SafetyLevel: h.SafetyLevel,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	_ = totalPages
	response.OK(c, FoodListResponse{Content: content, Page: page, TotalCount: total})
}

func (h *Handler) askAI(c *gin.Context) {
	var req FoodAskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var existing model.FoodChatHistory
	if err := h.db.Where("question = ?", req.Question).First(&existing).Error; err == nil {
		response.OK(c, FoodResponse{
			ID:          existing.ID,
			Question:    existing.Question,
			Answer:      existing.Answer,
			SafetyLevel: existing.SafetyLevel,
		})
		return
	}

	analysis, err := callOllama(req.Question)
	if err != nil {
		response.InternalServerError(c, "AI 분석에 실패했습니다.")
		return
	}

	if !analysis.IsFood {
		_ = c.Error(pnerrors.ErrFoodNotRelated)
		return
	}

	history := &model.FoodChatHistory{
		Question:    req.Question,
		Answer:      analysis.Answer,
		SafetyLevel: model.SafetyLevel(analysis.SafetyLevel),
	}
	h.db.Create(history)

	response.OK(c, FoodResponse{
		ID:          history.ID,
		Question:    history.Question,
		Answer:      history.Answer,
		SafetyLevel: history.SafetyLevel,
	})
}

type ollamaAnalysis struct {
	IsFood      bool   `json:"isFood"`
	SafetyLevel string `json:"safetyLevel"`
	Answer      string `json:"answer"`
}

func callOllama(question string) (*ollamaAnalysis, error) {
	host := config.AppConfig.Ollama.Host
	prompt := buildOllamaPrompt(question)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  "exaone3.5:7.8b",
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.3,
		},
	})

	resp, err := http.Post(host+"/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Response string `json:"response"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var analysis ollamaAnalysis
	if err = json.Unmarshal([]byte(result.Response), &analysis); err != nil {
		return nil, fmt.Errorf("ollama 응답 파싱 실패: %w", err)
	}
	return &analysis, nil
}

func buildOllamaPrompt(question string) string {
	return fmt.Sprintf(`당신은 20년 경력의 수의사 영양 전문가입니다. 강아지 음식 안전성에 대한 질문에 답변해주세요.

반드시 아래 JSON 형식으로만 응답하세요:
{
  "isFood": true/false,
  "safetyLevel": "GOOD|NOTION|BAD",
  "answer": "상세 답변"
}

- isFood: 음식 관련 질문이면 true, 아니면 false
- safetyLevel: GOOD(안전), NOTION(주의필요), BAD(위험)
- answer: 마크다운 형식으로 영양정보, 섭취방법, 주의사항 포함

질문: %s`, question)
}
