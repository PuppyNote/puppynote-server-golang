package main

import (
	"github.com/gin-gonic/gin"
)

// registerRoutes는 각 도메인 핸들러를 라우터에 등록합니다.
// 도메인 구현 시 순차적으로 추가됩니다.
func registerRoutes(base *gin.RouterGroup) {
	_ = base
}
