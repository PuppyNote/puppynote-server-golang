package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
)

func getOAuthEmail(snsType model.SnsType, token string) (string, error) {
	switch snsType {
	case model.SnsTypeKakao:
		return getKakaoEmail(token)
	case model.SnsTypeGoogle:
		return getGoogleEmail(token)
	case model.SnsTypeApple:
		return getAppleEmail(token)
	default:
		return "", pnerrors.New("지원하지 않는 SNS 타입입니다.")
	}
}

func getKakaoEmail(token string) (string, error) {
	req, _ := http.NewRequest(http.MethodPost, "https://kapi.kakao.com/v2/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		KakaoAccount struct {
			Email string `json:"email"`
		} `json:"kakao_account"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.KakaoAccount.Email == "" {
		return "", pnerrors.New("카카오 이메일을 가져올 수 없습니다.")
	}
	return result.KakaoAccount.Email, nil
}

func getGoogleEmail(token string) (string, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Email string `json:"email"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.Email == "" {
		return "", pnerrors.New("구글 이메일을 가져올 수 없습니다.")
	}
	return result.Email, nil
}

func getAppleEmail(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", pnerrors.ErrTokenInvalid
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("애플 토큰 디코딩 실패: %w", err)
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err = json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}
	if claims.Email == "" {
		return "", pnerrors.New("애플 이메일을 가져올 수 없습니다.")
	}
	return claims.Email, nil
}
