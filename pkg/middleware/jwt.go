package middleware

import (
	"strings"

	pnjwt "github.com/PuppyNote/puppynote-server-golang/pkg/jwt"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
)

const LoginUserKey = "loginUser"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "인증 토큰이 필요합니다.")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		info, err := pnjwt.Parse(tokenString)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set(LoginUserKey, info)
		c.Next()
	}
}

func GetLoginUser(c *gin.Context) *pnjwt.LoginUserInfo {
	val, exists := c.Get(LoginUserKey)
	if !exists {
		return nil
	}
	info, ok := val.(*pnjwt.LoginUserInfo)
	if !ok {
		return nil
	}
	return info
}
