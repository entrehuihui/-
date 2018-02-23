/*
//软件参数：
// -f 保存的地址 例：/img(默认)
// -s 网址起始图片
// -g 刷新的图片数
//	-w 检索关键字
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"./pathExistsCheck"
	"./reptile"
)

//最终下载图片数量
var numEnd = 0
var num = 0

//锁
//var mu sync.Mutex

var url = "http://image.baidu.com/search/avatarjson?tn=resultjsonavatarnew&ie=utf-8&word=go图标&cg=girl&pn=10&rn=30&itg=0&z=0&fr=&width=&height=&lm=-1&ic=0&s=0&st=-1&gsm=960000000096"

func main() {
	urlBefor := "http://image.baidu.com/search/avatarjson?tn=resultjsonavatarnew&ie=utf-8&word="
	urlAfter := "&itg=0&z=0&fr=&width=&height=&lm=-1&ic=0&s=0&st=-1&gsm=960000000096"
	urlMiddle := "&rn="
	startNum := flag.Int("s", 0, "网址起始图片数")
	getNum := flag.Int("g", 20, "网址刷新的图片数（最大30）")
	searchWord := flag.String("w", "美女", "检索的关键字")
	pathExists := flag.String("f", "img", "your path exists")

	//把缓存上的数据刷新进内存
	flag.Parse()

	if *pathExists == "img" {
		*pathExists = "img/" + time.Now().Format("2006/1/2") + "/" + *searchWord
	}

	urlBefor = urlBefor + *searchWord + "&pn="
	if *getNum > 30 {
		*getNum = 30
	}

	//输入下载的图片数量
	fmt.Println("输入要下载的图片数量 ：")
	var get = 0
	fmt.Scanf("%d", &get)

	//检查文件路径
	if err := pathExistsCheck.Check(*pathExists); err != nil {
		fmt.Println("Can't find the folder:", err)
		return
	}
	//总进程等待
	var wg sync.WaitGroup
	//replite进程等待
	var childWg sync.WaitGroup
	//限制解析进程通道，最多同时五个进程解析
	chanReptile := make(chan int, 5)
LOOP:
	//创建图片地址通道
	var chanURL = make(chan reptile.Picture, 2)
	//解析网址
	wg.Add(1)
	go func() {
		//拼合网址
		urlS, _ := addURL(urlBefor, urlMiddle, urlAfter, *startNum, *getNum, get)
		for i := 0; i < len(urlS); i++ {
			url = urlS[i]
			chanReptile <- 1
			childWg.Add(1)
			func() {
				body, _ := reptile.GetUrlDate(url)
				_ = reptile.GetPictureURL(string(body), chanURL)
				<-chanReptile
				childWg.Done()
			}()
		}
		//等待所有reptile进程结束
		childWg.Wait()
		//信息解析完成，关闭通道
		close(chanURL)
		wg.Done()
	}()
	//启动下载进程
	wg.Add(1)
	go func() {
		//下载图片
		for date := range chanURL {
			date.PictureName = *pathExists + "/" + date.PictureName
			DownPicture(date)
		}
		wg.Done()
	}()

	wg.Wait()
	//等待下载进程结束
	if get > num {
		get = get - num
		num = 0
		*startNum += get
		goto LOOP
	}
	fmt.Println("======图片下载完成========")
}

//DownPicture ：下载图片并保存
func DownPicture(date reptile.Picture) {

	defer flag.Parse()
	//fmt.Println(date.PictureName)
	rep, err := http.Get(date.PictureURL)
	if err != nil {
		ret(err)
		return
	}
	body, _ := ioutil.ReadAll(rep.Body)
	defer rep.Body.Close()

	d, f := filepath.Split(date.PictureName)
	date.PictureName = d + time.Now().Format("20060102150405") + strconv.Itoa(num) + f
	file, err := os.Create(date.PictureName)
	defer file.Close()
	if err != nil {
		ret(err)
		return
	}
	//mu.Lock()
	num++
	numEnd++
	fmt.Println("正在下载第", numEnd, "张图片：", date.PictureName)
	//mu.Unlock()
	io.Copy(file, bytes.NewReader(body))
	return
}

func ret(err error) {
	fmt.Println(err)
}

func addURL(urlBefor, urlMiddle, urlAfter string, startNum, getNum, get int) (url []string, retNum int) {

	if get < getNum {
		getNum = get
		url = append(url, urlBefor+strconv.Itoa(startNum)+urlMiddle+strconv.Itoa(getNum)+urlAfter)
		retNum = startNum + getNum
		return
	}

	getQ := get / getNum
	getR := get % getNum
	var i int
	for i = 0; i < getQ; i++ {
		url = append(url, urlBefor+strconv.Itoa(startNum+getNum*i)+urlMiddle+strconv.Itoa(getNum)+urlAfter)
	}
	url = append(url, urlBefor+strconv.Itoa(startNum+getNum*i)+urlMiddle+strconv.Itoa(getR)+urlAfter)
	retNum = startNum*i + getR
	return
}
