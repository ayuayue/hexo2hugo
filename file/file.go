package file

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GetAllMDFileName 获取当前目录下的所有指定格式的文件
func GetAllMDFileName(dirpath string) []string {
	filePathNames, err := filepath.Glob(filepath.Join(dirpath, "*.md"))
	if err != nil {
		log.Fatal(err)
	}
	count := len(filePathNames)
	if count < 1 {
		fmt.Println("当前目录下没有markdown文件")
	} else {
		fmt.Printf("该目录下共检测到 %d 个 markdown 文件\n", count)
	}

	return filePathNames
}

// processTags 处理多行tags压扁成的单行 tags
// NOTICE: 这里非常容易产生 null 字符，所以做了一些特殊判断
func processTags(inputText string) (string, error) {
	if len(inputText) == 0 {
		return "", nil
	}

	// 使用正则表达式匹配标签
	re := regexp.MustCompile(`- ([\p{Han}a-zA-Z0-9]+)`)
	matches := re.FindAllStringSubmatch(inputText, -1)

	// 提取标签名称到字符串数组
	var tags []string
	for _, match := range matches {
		tags = append(tags, match[1])
	}

	// 将字符串数组转换为 JSON 字符串
	jsonStr, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}

	return string(jsonStr), nil
}

// removeBOM 检测并去除 UTF-8 文件的 BOM
func removeBOM(s string) string {
	// Check for UTF-8 BOM (EF BB BF) and remove it
	if strings.HasPrefix(s, "\xEF\xBB\xBF") {
		return s[3:]
	}
	return s
}

// ReadAllMDFile 读取所有获取到的文件
func ReadAllMDFile(filePathNames []string) {
	newDir := fmt.Sprintf("newPost%d", time.Now().Unix())
	err := os.Mkdir(newDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	totalFileNum := len(filePathNames)
	// 遍历读取所有的文件
	for k, v := range filePathNames {
		file, err := os.Open(v)
		if err != nil {
			log.Printf("can't open file %s , err : %s ", v, err)
		}
		defer file.Close()
		fmt.Printf("\n正在读取第 %d/%d 个文件 %s\n", k+1, totalFileNum, v)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		// 使用 strings.Builder 保存结果
		var builder bytes.Buffer

		// 通过 start 检测 `---` 保证不影响文本内容
		start := 0
		multi_line := false
		buf := ""

		for scanner.Scan() {
			line := scanner.Text()
			if start < 2 {
				if multi_line {
					if strings.Contains(line, " - ") {
						buf += line
						continue
					} else {
						multi_line = false
						// 处理输入并打印结果
						fmt.Printf("buf: %q\n", buf)
						output, err := processTags(buf)

						fmt.Println("output: ", output)
						if err != nil {
							fmt.Println("Error:", err)
							log.Fatal("can't handle: ", buf)
						}
						// 检查 output 是否以换行符结尾
						if !strings.HasSuffix(output, "\n") {
							output += "\n"
						}
						builder.WriteString(output)
					}
				}
				if strings.Contains(line, "---") {
					start++
					if start == 1 {
						fmt.Printf("开始处理文件[%s]的文件头\n", v)
					}
					if start == 2 {
						fmt.Printf("文件头即将[%s]处理完成\n", v)
					}
				}
				//如果行中包含 date 标签,则替换掉
				if strings.Contains(line, "date:") {
					line = HandleDate(line)
				}
				if strings.Contains(line, "tags:") {
					// 处理多行逻辑 && 还有空的逻辑！
					if len(line) <= len("tags: ") {
						multi_line = true
						buf = ""
						// 检查 line 是否以空格符结尾, 多行需要加空格
						if !strings.HasSuffix(line, " ") {
							line += " "
						}
						builder.WriteString(line)
						continue
					}
					line = HandleTags(line)
				}
				if strings.Contains(line, "categories") {
					// 处理多行逻辑 && 还有空的逻辑！
					if len(line) <= len("categories: ") {
						fmt.Println("multi get line: ", line)
						multi_line = true
						buf = ""
						// 检查 line 是否以空格符结尾, 多行需要加空格
						if !strings.HasSuffix(line, " ") {
							line += " "
						}
						builder.WriteString(line)
						continue
					}
					line = HandleCategories(line)
					fmt.Println("After handle categories: ", line)
				}
				if strings.Contains(line, "updated:") {
					if len(line) <= len("updated: ") {
						continue
					}
					line = HandleUpdated(line)
				}
			}
			// fmt.Printf("cur line: %q \n", line)
			// 检查 line 是否以换行符结尾
			if !strings.HasSuffix(line, "\n") {
				line += "\n"
			}

			// 进行字符串拼接
			builder.WriteString(line)

			// fmt.Printf("cur res: %q \n", builder.String())
		}
		// 获取最终拼接的结果
		result := builder.String()

		// 检测和去除 BOM
		result = removeBOM(result)

		fmt.Printf("正在更改第 %d 个文件 %s\n", k+1, v)

		if err := HandleContent(newDir, v, result); err != nil {
			log.Fatal("写入文件失败")
		}
		fmt.Printf("第 %d/%d 个文件 %s 写入完成\n", k+1, totalFileNum, v)
	}

}

// HandleContent 将处理完的内容覆盖进去
func HandleContent(newDir string, filePathNames string, result string) error {
	// 如果路径分割符为 \ 则替换为 /
	filePathNames = strings.Replace(filePathNames, "\\", "/", -1)
	// 获取最后一个分割符后的名称
	index := strings.LastIndex(filePathNames, "/")
	fileName := filePathNames[index+1:]
	// 创建文件
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", newDir, fileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, result)
	if err != nil {
		return err
	}
	return nil
}

func replacePrefix(input string, oldPrefix string, newPrefix string) string {
	if strings.HasPrefix(input, oldPrefix) {
		return newPrefix + input[len(oldPrefix):]
	}
	return input
}

// HandleUpdated 处理修改时间
func HandleUpdated(updated string) string {
	//查找到年月日后的空格,根据 hugo updated 的格式进行更改
	index := strings.LastIndex(updated, " ")
	updated = fmt.Sprintf("%sT%s+08:00", updated[:index], updated[index+1:])
	// 替换前缀
	lastmod := replacePrefix(updated, "updated", "lastmod")
	return lastmod
}

// HandleTags 处理标签
// NOTICE: 不对对已经有 `[]` 的标记特殊判断
func HandleTags(tags string) string {
	// 可能的去除换行
	tags = strings.TrimRight(tags, "\n")
	//查找tags后的空格,将所有的 tags 用 [] 包裹起来
	index := strings.Index(tags, " ")
	// 处理防止 tag 后为空的情况
	if index != -1 {
		tags = fmt.Sprintf("%s [%s]", tags[:index], tags[index+1:])
	}
	return tags
}

// HandleCategories 处理分类
// NOTICE: 不对对已经有 `[]` 的标记特殊判断
func HandleCategories(categories string) string {
	// 可能得去除行尾的换行符
	categories = strings.TrimRight(categories, "\n")
	//查找到分类的空格,整理成 hugo 的格式
	index := strings.Index(categories, " ")
	// 处理防止 categories 后为空的情况
	if index != -1 {
		categories = fmt.Sprintf("%s [%s]", categories[:index], categories[index+1:])
	}
	return categories
}

// HandleDate 处理日期
func HandleDate(date string) string {
	//查找到年月日后的空格,根据 hugo date 的格式进行更改
	index := strings.LastIndex(date, " ")
	date = fmt.Sprintf("%sT%s+08:00", date[:index], date[index+1:])
	return date
}
