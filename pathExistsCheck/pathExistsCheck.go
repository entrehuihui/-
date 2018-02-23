package pathExistsCheck

import (
	"fmt"
	"os"
	"strings"
)

//Check 检查是否存在目录
func Check(pathExists string) error {

	_, err := os.Stat(pathExists)
	if err == nil {
		fmt.Println("图片保存在：", pathExists)
		return err
	}
	err = AnalysisPath(pathExists)

	return err
}

//AnalysisPath 解析输入的路径
func AnalysisPath(path string) error {

	var err error

	pathSlice := strings.Split(path, "/")
	foderPath := "./"
	for _, date := range pathSlice {
		foderPath += date
		_, err = os.Stat(foderPath)
		if err == nil {
		} else {
			err = os.Mkdir(foderPath, 0777)
			if err != nil {
				fmt.Println("Create error:", err)
				return err
			}
		}
		foderPath += "/"

	}
	return err
}
