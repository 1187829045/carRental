package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func intParam(c *gin.Context, keys ...string) int {
	for _, k := range keys {
		if v := strings.TrimSpace(c.PostForm(k)); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
		if v := strings.TrimSpace(c.Query(k)); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return 0
}

func intArrayParam(c *gin.Context, key string) []int {
	values := c.PostFormArray(key)
	if len(values) == 0 {
		values = c.QueryArray(key)
	}
	out := make([]int, 0, len(values))
	for _, v := range values {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			out = append(out, n)
		}
	}
	return out
}

func stringArrayParam(c *gin.Context, key string) []string {
	values := c.PostFormArray(key)
	if len(values) == 0 {
		values = c.QueryArray(key)
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, ResultObj{Code: -1, Msg: msg})
}

