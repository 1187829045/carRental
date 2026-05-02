package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetAttachmentFilenameAddsFilenameStar(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	setAttachmentFilename(c, "客户数据.xlsx")
	h := w.Header().Get("Content-Disposition")
	if !strings.Contains(h, "filename*=") {
		t.Fatalf("expected filename* in Content-Disposition, got %q", h)
	}
	if !strings.Contains(h, "UTF-8''") {
		t.Fatalf("expected UTF-8'' prefix, got %q", h)
	}
}

