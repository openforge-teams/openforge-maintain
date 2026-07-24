package utils

import (
	"errors"
	"path/filepath"
	"strings"
)

// SafeJoin 安全地拼接路径，防止路径穿越攻击
// 如果 target 试图访问 base 目录之外的文件，返回错误
func SafeJoin(base, target string) (string, error) {
	// 清理路径
	cleanBase := filepath.Clean(base)
	cleanTarget := filepath.Clean(target)

	// 如果 target 是绝对路径，直接拒绝
	if filepath.IsAbs(cleanTarget) {
		return "", errors.New("path traversal detected: absolute path not allowed")
	}

	// 检查 target 中是否包含路径穿越字符
	if strings.Contains(cleanTarget, "..") {
		return "", errors.New("path traversal detected: relative parent reference not allowed")
	}

	// 拼接路径
	joined := filepath.Join(cleanBase, cleanTarget)
	cleanJoined := filepath.Clean(joined)

	// 确保最终路径在 base 目录内
	if !strings.HasPrefix(cleanJoined, cleanBase+string(filepath.Separator)) && cleanJoined != cleanBase {
		return "", errors.New("path traversal detected: target path escapes base directory")
	}

	return cleanJoined, nil
}
