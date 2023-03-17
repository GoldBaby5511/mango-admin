package middleware

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error
		switch x := recovered.(type) {
		case error:
			err = x
		case runtime.Error:
			err = x
		default:
			err = fmt.Errorf("%v", x)
		}

		s := getErrorStack(err.Error(), "app")
		log.Println(s)
		response.Error(c, 500, err, "")
	})
}

// 获取错误文件和行号。 去除go自带函数和外部包。
func getErrorStack(errString string, splitDirName string) []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	idx := 0
	recorder := []string{errString}

	for i, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		filepath, line := fn.FileLine(pc)

		if strings.Contains(filepath, splitDirName) {
			if idx == 0 {
				idx = strings.Index(filepath, splitDirName)
			}
			recorder = append(recorder, fmt.Sprintf("%s:%d", filepath[idx:], line))
		}

		if i >= 20 {
			break
		}
	}

	return recorder
}

/*
func CustomError(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {

			if c.IsAborted() {
				c.Status(200)
			}
			switch errStr := err.(type) {
			case string:
				p := strings.Split(errStr, "#")
				if len(p) == 3 && p[0] == "CustomError" {
					statusCode, e := strconv.Atoi(p[1])
					if e != nil {
						break
					}
					c.Status(statusCode)
					fmt.Println(
						time.Now().Format("2006-01-02 15:04:05"),
						"[ERROR]",
						c.Request.Method,
						c.Request.URL,
						statusCode,
						c.Request.RequestURI,
						common.GetClientIP(c),
						p[2],
					)
					c.JSON(http.StatusOK, gin.H{
						"code": statusCode,
						"msg":  p[2],
					})
				}
			default:
				panic(err)
			}
		}
	}()
	c.Next()
}
*/
