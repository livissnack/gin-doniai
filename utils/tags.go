// utils/tags.go
package utils

import (
	"fmt"
	"strings"
	"encoding/json"
)

// ParseTags 处理标签字符串，支持中英文逗号分割和JSON数组格式，空值返回空数组
func ParseTags(tagStr string) []string {
    // 打印输入值用于调试
    fmt.Println(tagStr, "-------------")

    if tagStr == "" {
        return []string{}
    }

    // 尝试解析JSON数组格式
    if strings.HasPrefix(tagStr, "[") && strings.HasSuffix(tagStr, "]") {
        var tags []string
        if err := json.Unmarshal([]byte(tagStr), &tags); err == nil {
            // 过滤空标签
            var cleanedTags []string
            for _, tag := range tags {
                trimmedTag := strings.TrimSpace(tag)
                if trimmedTag != "" {
                    cleanedTags = append(cleanedTags, trimmedTag)
                }
            }
            return cleanedTags
        }
    }

    // 处理普通字符串格式
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

