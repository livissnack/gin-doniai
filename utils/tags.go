// utils/tags.go
package utils

import (
	"strings"
)

// ParseTags 处理标签字符串，支持中英文逗号分割，空值返回空数组
func ParseTags(tagStr string) []string {
	if tagStr == "" {
		return []string{}
	}

	// 替换中文逗号为英文逗号
	tagStr = strings.ReplaceAll(tagStr, "，", ",")

	// 按逗号分割并去除空格
	rawTags := strings.Split(tagStr, ",")
	var cleanedTags []string

	for _, tag := range rawTags {
		// 去除首尾空格，过滤空标签
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			cleanedTags = append(cleanedTags, trimmedTag)
		}
	}

	return cleanedTags
}
