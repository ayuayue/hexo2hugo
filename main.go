package main

import (
	"hexo2hugo/file"
)

func main() {
	filePathNames := file.GetAllMDFileName()
	file.ReadAllMDFile(filePathNames)
	file.Waiter.Wait()
}
