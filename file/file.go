package file

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//GetAllMDFileName 获取当前目录下的所有指定格式的文件
func GetAllMDFileName() []string {
	pwd, _ := os.Getwd()
	filePathNames, err := filepath.Glob(filepath.Join(pwd, "*.md"))
	if err != nil {
		log.Fatal(err)
	}
	count := len(filePathNames)
	if count < 1 {
		fmt.Println("当前目录下没有markdown文件")
	} else {
		fmt.Printf("该目录下共检测到 %d 个 markdown 文件\n", count-1)
	}

	for i := range filePathNames {
		fmt.Println(filePathNames[i])
	}
	return filePathNames
}

//ReadAllMDFile 读取所有获取到的文件
func ReadAllMDFile(filePathNames []string) {
	newDir := fmt.Sprintf("newPost%d", time.Now().Unix())
	err := os.Mkdir(newDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//遍历读取所有的文件
	for k, v := range filePathNames {
		file, err := os.Open(v)
		if err != nil {
			log.Printf("can't open file %s , err : %s ", v, err)
		}
		fmt.Printf("\n正在读取第 %d 个文件 %s\n", k+1, v)
		result := ""
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			//如果行中包含 date 标签,则替换掉
			if strings.Contains(line, "date:") {
				line = HandleDate(line)
			}
			if strings.Contains(line, "tags:") {
				line = HandleTags(line)
			}
			if strings.Contains(line, "categories") {
				line = HandleCategories(line)
			}

			result = result + line + "\n"
		}
		fmt.Printf("正在更改第 %d 个文件 %s\n", k+1, v)

		if err := HandleContent(newDir, v, result); err != nil {
			log.Fatal("写入文件失败")
		}
		fmt.Printf("%s写入完成\n", v)
	}

}

//HandleContent 将处理完的内容覆盖进去
func HandleContent(newDir string, filePathNames string, result string) error {
	//如果路径分割符为 \ 则替换为 /
	if strings.Contains(filePathNames, "\\") {
		filePathNames = strings.Replace(filePathNames, "\\", "/", -1)
	}
	//获取最后一个分割符后的名称
	index := strings.LastIndex(filePathNames, "/")
	fileName := filePathNames[index+1:]
	//创建文件
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

//HandleTags 处理标签
func HandleTags(tags string) string {
	//查找tas后的空格,将所有的 tags 用 [] 包裹起来
	index := strings.Index(tags, " ")
	tags = fmt.Sprintf("%s [%s]", tags[:index], tags[index+1:])
	return tags
}

//HandleCategories 处理分类
func HandleCategories(categories string) string {
	//查找到分类的空格,整理成 hugo 的格式
	index := strings.Index(categories, " ")
	categories = fmt.Sprintf("%s [%s]", categories[:index], categories[index+1:])
	return categories
}

//HandleDate 处理日期
func HandleDate(date string) string {
	//查找到年月日后的空格,根据 hugo date 的格式进行更改
	index := strings.LastIndex(date, " ")
	date = fmt.Sprintf("%sT%s+08:00", date[:index], date[index+1:])
	return date
}
