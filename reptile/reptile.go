package reptile

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//Picture 创建图片名和地址结构体
type Picture struct {
	PictureURL  string
	PictureName string
}

var existsURL = regexp.MustCompile(`"objURL":"http://[\S\s]*?"fromPageTitle":"\S*<strong>\S*?</strong>\S*?"`)
var existsURL1 = regexp.MustCompile(`http://[\S]*?\.[A-Za-z]{3,9}"`)
var existsName = regexp.MustCompile(`"fromPageTitle":"\S*[^"]`)
var existsname1 = regexp.MustCompile(`\.\S{3,5}$`)

//GetPictureURL 从下载的网页数据中解析出来 return：URL(string)，name(string)
func GetPictureURL(date string, chanURL chan (Picture)) int {

	//解析计数
	var num = 0

	pictureURL := existsURL.FindAllStringSubmatch(date, 30)
	for _, rangeDate := range pictureURL {
		defer func() {
			if err := recover(); err != nil {
				file, err := os.OpenFile("panic.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
				if err != nil {
					return
				}
				_, _ = io.WriteString(file, date)
			}
		}()
		var f Picture
		url := existsURL1.FindAllStringSubmatch(rangeDate[0], 1)[0][0]
		//fmt.Println(url)
		f.PictureURL = string([]byte(url)[:len(url)-1])

		name := existsName.FindAllStringSubmatch(rangeDate[0], 1)[0][0]
		name = string([]byte(name)[17:len(name)])
		name = strings.Replace(name, "<strong>", "", -1)
		name = strings.Replace(name, "</strong>", "", -1)
		name = strings.Replace(name, "\\", "", -1)
		name = strings.Replace(name, "/", "", -1)
		name = strings.Replace(name, ">", "", -1)
		name = strings.Replace(name, "<", "", -1)
		//name = strings.Replace(name, ":", "", -1)
		//name = strings.Replace(name, "|", "", -1)
		f.PictureName = name + existsname1.FindAllStringSubmatch(f.PictureURL, 1)[0][0]
		//f.PictureName = time.Now().Format("20060102150405") + name
		//fmt.Println(f.PictureName)
		chanURL <- f
		num++
	}
	return num
}

/*
//解析图片地址的规则
var existsURL = regexp.MustCompile(`"objURL":"http://\S*?\.\S{3,4}"`)
//解析图片名称的规则
var existsFile = regexp.MustCompile(`(\S[^\.])*\.`)
//GetPictureURL 从下载的网页数据中解析出来 return：URL(string)，name(string)
func GetPictureURL(date string, chanURL chan (Picture)) int {
	//粗略提取网址
	pictureURL := existsURL.FindAllStringSubmatch(date, 35)

	//计数
	var num = 0
	for _, urlDate := range pictureURL {

		var f Picture
		//去掉多余的部分，提取真正的网址
		f.PictureURL = string([]byte(urlDate[0])[10 : len(urlDate[0])-1])

		ti := time.Now().Unix()

		f.PictureName = strings.Replace(f.PictureURL, "/", "_", -1)
		f.PictureName = strings.Replace(f.PictureName, ":", "@", -1)
		num1 := strings.Count(f.PictureName, ".")
		f.PictureName = strings.Replace(f.PictureName, ".", "@", num1-1)
		f.PictureName = strconv.FormatInt(ti, 10) + f.PictureName
		//检测是否空名
		if f.PictureName != "" {
			chanURL <- f
		}
	}

	return num
}
*/

//GetUrlDate 把指定网页的内容下载下来
func GetUrlDate(url string) ([]byte, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Http err:", err)
		return nil, err
	}
	resq, err := client.Do(req)
	if err != nil {
		fmt.Println("Client.Do err:", err)
		return nil, err
	}
	defer resq.Body.Close()

	body, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		fmt.Println("ReadAll err:", err)
		return nil, err
	}
	return body, err
}
