package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// 设置 gin mode 之后再执行测试用例
	// TestMode在测试时不会出现重复的调试日志
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
