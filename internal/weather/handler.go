package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	redispkg "github.com/PuppyNote/puppynote-server-golang/pkg/redis"
	"github.com/PuppyNote/puppynote-server-golang/pkg/middleware"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
)

const weatherCacheTTL = time.Hour

type WalkCondition string

const (
	WalkConditionGreat    WalkCondition = "GREAT"
	WalkConditionGood     WalkCondition = "GOOD"
	WalkConditionModerate WalkCondition = "MODERATE"
	WalkConditionBad      WalkCondition = "BAD"
	WalkConditionDanger   WalkCondition = "DANGER"
)

type WeatherResponse struct {
	Temperature       float64       `json:"temperature"`
	WeatherCode       int           `json:"weatherCode"`
	WeatherDesc       string        `json:"weatherDescription"`
	WindSpeed         float64       `json:"windSpeed"`
	Precipitation     float64       `json:"precipitation"`
	WalkCondition     WalkCondition `json:"walkCondition"`
	WalkMessage       string        `json:"walkMessage"`
}

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/api/v1/weather", middleware.JWTAuth(), h.getWeather)
}

func (h *Handler) getWeather(c *gin.Context) {
	var lat, lon float64
	if _, err := fmt.Sscanf(c.Query("latitude"), "%f", &lat); err != nil {
		response.BadRequest(c, "latitude가 필요합니다.")
		return
	}
	if _, err := fmt.Sscanf(c.Query("longitude"), "%f", &lon); err != nil {
		response.BadRequest(c, "longitude가 필요합니다.")
		return
	}

	cacheKey := fmt.Sprintf("weather:%.2f:%.2f",
		math.Round(lat/0.05)*0.05,
		math.Round(lon/0.05)*0.05,
	)

	ctx := context.Background()
	if cached, err := redispkg.Client.Get(ctx, cacheKey).Result(); err == nil {
		var res WeatherResponse
		if json.Unmarshal([]byte(cached), &res) == nil {
			response.OK(c, res)
			return
		}
	}

	res, err := fetchWeather(lat, lon)
	if err != nil {
		response.InternalServerError(c, "날씨 정보를 가져올 수 없습니다.")
		return
	}

	if data, err := json.Marshal(res); err == nil {
		redispkg.Client.Set(ctx, cacheKey, data, weatherCacheTTL)
	}

	response.OK(c, res)
}

func fetchWeather(lat, lon float64) (*WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weathercode,windspeed_10m,precipitation&timezone=Asia%%2FSeoul",
		lat, lon,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Current struct {
			Temperature   float64 `json:"temperature_2m"`
			WeatherCode   int     `json:"weathercode"`
			WindSpeed     float64 `json:"windspeed_10m"`
			Precipitation float64 `json:"precipitation"`
		} `json:"current"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	cur := result.Current
	desc, baseCondition := weatherCodeInfo(cur.WeatherCode)
	condition, msg := adjustForTemperature(baseCondition, cur.Temperature)

	return &WeatherResponse{
		Temperature:   cur.Temperature,
		WeatherCode:   cur.WeatherCode,
		WeatherDesc:   desc,
		WindSpeed:     cur.WindSpeed,
		Precipitation: cur.Precipitation,
		WalkCondition: condition,
		WalkMessage:   msg,
	}, nil
}

func weatherCodeInfo(code int) (string, WalkCondition) {
	switch {
	case code == 0:
		return "맑음", WalkConditionGreat
	case code == 1:
		return "대체로 맑음", WalkConditionGreat
	case code == 2 || code == 3:
		return "흐림", WalkConditionGood
	case code == 45 || code == 48:
		return "안개", WalkConditionModerate
	case code >= 51 && code <= 65:
		return "비", WalkConditionBad
	case code == 66 || code == 67:
		return "어는 비", WalkConditionDanger
	case code >= 71 && code <= 77:
		return "눈", WalkConditionModerate
	case code >= 80 && code <= 82:
		return "소나기", WalkConditionBad
	case code == 85 || code == 86:
		return "눈 소나기", WalkConditionModerate
	case code >= 95:
		return "뇌우", WalkConditionDanger
	default:
		return "알 수 없음", WalkConditionModerate
	}
}

func adjustForTemperature(base WalkCondition, temp float64) (WalkCondition, string) {
	switch {
	case temp >= 32:
		return WalkConditionDanger, "산책하기 위험한 날씨에요! 뜨거운 바닥에 발바닥 화상 주의!"
	case temp >= 28:
		return WalkConditionBad, "산책을 자제하는 것이 좋아요. 이른 아침이나 저녁 산책을 추천해요."
	case temp >= 23:
		if base == WalkConditionGreat || base == WalkConditionGood {
			return WalkConditionGood, "산책하기 좋은 날씨에요. 다소 더우니 수분 보충에 신경써주세요."
		}
		return base, walkConditionMessage(base)
	case temp <= -10:
		return WalkConditionDanger, "산책하기 위험한 날씨에요! 동상 위험이 있어요."
	case temp <= -5:
		return WalkConditionBad, "산책을 자제하는 것이 좋아요. 동상 위험이 있어요."
	case temp <= 5:
		msg := walkConditionMessage(base)
		return base, msg + " 쌀쌀한 날씨니 따뜻하게 입혀주세요."
	default:
		return base, walkConditionMessage(base)
	}
}

func walkConditionMessage(c WalkCondition) string {
	switch c {
	case WalkConditionGreat:
		return "산책하기 최적인 날씨에요!"
	case WalkConditionGood:
		return "산책하기 좋은 날씨에요."
	case WalkConditionModerate:
		return "산책은 가능하지만 주의하세요."
	case WalkConditionBad:
		return "산책을 자제하는 것이 좋아요."
	default:
		return "산책하기 위험한 날씨에요!"
	}
}
