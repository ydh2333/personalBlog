package util

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行后续业务逻辑（如果发生错误，会被 c.AbortWithError 或 panic 触发）
		c.Next()

		// 检查是否有错误发生（c.Errors 存储了后续逻辑中通过 c.Error() 或 AbortWithError() 抛出的错误）
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err // 获取最后一个错误（通常是关键错误）

			// 1. 处理自定义业务错误
			if bizErr, ok := err.(*BusinessError); ok {
				c.JSON(bizErr.HttpCode, &Response{
					Code:    bizErr.Code,
					Message: bizErr.Message,
					Data:    nil,
				})
				return
			}

			// 2. 处理系统错误（如数据库连接失败、空指针等）
			// 记录详细错误日志（包含堆栈信息，便于排查）
			fmt.Printf("系统错误：%v\n堆栈信息：%s\n", err, debug.Stack())
			// 返回通用 500 错误（不暴露具体细节给前端）
			c.JSON(http.StatusInternalServerError, &Response{
				Code:    ErrSystemError.Code,
				Message: ErrSystemError.Message,
				Data:    nil,
			})
			return
		}
	}
}
