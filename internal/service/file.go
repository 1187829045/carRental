package service

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileService struct {
	UploadRoot     string // absolute/relative root, e.g. "./upload"
	ProjectRoot    string
	StaticAssets   []string
	RemoteCarImages []string
}

func NewFileService(projectRoot string) (*FileService, error) {
	// Try reading legacy file.properties; fall back for mac/Linux.
	propsPath := filepath.Join(projectRoot, "src/main/resources/file.properties")
	root := ""
	if f, err := os.Open(propsPath); err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) != 2 {
				continue
			}
			if strings.TrimSpace(kv[0]) == "path" {
				root = strings.TrimSpace(kv[1])
				break
			}
		}
		_ = f.Close()
	}

	// Environment override.
	if v := strings.TrimSpace(os.Getenv("UPLOAD_PATH")); v != "" {
		root = v
	}

	// If config points to a Windows drive or is empty, use a local folder.
	if root == "" || looksLikeWindowsDrive(root) {
		root = filepath.Join(projectRoot, "upload")
	}
	return &FileService{
		UploadRoot:  root,
		ProjectRoot: projectRoot,
		StaticAssets: []string{
			"src/main/webapp/static/images/cars/placeholder-1.svg",
			"src/main/webapp/static/images/cars/placeholder-2.svg",
			"src/main/webapp/static/images/cars/placeholder-3.svg",
		},
		RemoteCarImages: []string{
			"https://images.pexels.com/photos/170811/pexels-photo-170811.jpeg?auto=compress&cs=tinysrgb&w=1200",
			"https://images.pexels.com/photos/358070/pexels-photo-358070.jpeg?auto=compress&cs=tinysrgb&w=1200",
			"https://images.pexels.com/photos/210019/pexels-photo-210019.jpeg?auto=compress&cs=tinysrgb&w=1200",
		},
	}, nil
}

func looksLikeWindowsDrive(p string) bool {
	if len(p) < 3 {
		return false
	}
	c := p[0]
	return ((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) && p[1] == ':' && (p[2] == '\\' || p[2] == '/')
}

func DateDir(t time.Time) string {
	return t.Format("2006-01-02")
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func SafeJoin(root string, rel string) (string, error) {
	rel = strings.TrimPrefix(rel, "/")
	rel = filepath.Clean(rel)
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", errors.New("invalid path")
	}
	return filepath.Join(root, rel), nil
}

func (f *FileService) CleanupTempFiles(maxAge time.Duration) (int, error) {
	removed := 0
	now := time.Now()
	err := filepath.Walk(f.UploadRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), "_temp") {
			return nil
		}
		if now.Sub(info.ModTime()) < maxAge {
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		removed++
		return nil
	})
	return removed, err
}
