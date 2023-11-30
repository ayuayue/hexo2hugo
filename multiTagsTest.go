package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	input := "- 算法a - test - tag3 - tag4"

	// 使用正则表达式匹配标签
	re := regexp.MustCompile(`- ([\p{Han}a-zA-Z0-9]+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	// 提取匹配的标签
	var tags []string
	for _, match := range matches {
		tags = append(tags, match[1])
	}

	// 将标签转化为字符串
	result := "[" + strings.Join(tags, ", ") + "]"

	fmt.Println(result)
}
