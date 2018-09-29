// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"image/jpeg"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// var domainDict map[string]domainNode

// type Node struct {
// 	timestamp int
// 	body      []byte // content
// }

// type domainNode struct {
// 	urls map[string]Node
// }

// func PathExists(path string) bool {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true
// 	}
// 	if os.IsNotExist(err) {
// 		return false
// 	}
// 	return false
// }

// func readFile(filePath string) string {
// 	return string(readFileByte(filePath))
// }

// func readFileByte(filePath string) []byte {
// 	if PathExists(filePath) == false {
// 		_, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	data, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return data
// }

// /*
// Name: writeFile
// @ para: filePath string
// @ para: content string
// @ para: appendEnable string
// @ Return: None
// Func: write the string content into assigned path by method of overwriting or appending
// */
// func writeFile(filePath string, content string, appendEnable bool) {
// 	if appendEnable == false {
// 		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		f.WriteString(content)
// 	} else {
// 		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		f.WriteString(content)
// 	}
// }

// func getTime() string {
// 	return time.Now().Format("2006-01-02 15:04:05")
// }

// func readImageByte(path string) []byte {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		fmt.Println(getTime(), ": readImageByte:", err)
// 	}
// 	defer file.Close()

// 	img, err := jpeg.Decode(file) //解码
// 	if err != nil {
// 		fmt.Println(getTime(), ": readImageByte:", err)
// 	}
// 	buf := new(bytes.Buffer)
// 	jpeg.Encode(buf, img, nil)
// 	return buf.Bytes()
// }

// func fileType(filename string) string {
// 	list := strings.Split(filename, ".")
// 	return list[len(list)-1]
// }

// func readByte(url string) []byte {
// 	var contentByte []byte
// 	if fileType(url) == "jpg" ||
// 		fileType(url) == "jpeg" ||
// 		fileType(url) == "png" {
// 		contentByte = readImageByte(url)
// 	} else {
// 		contentByte = readFileByte(url)
// 	}
// 	return contentByte
// }

// func addMemoryDick(host string, outpath string, currenttime int) {
// 	contentByte := readByte(outpath)

// 	_, ok := domainDict[host]
// 	if ok == false {
// 		domainDict[host] = domainNode{make(map[string]Node)}
// 	}

// 	domainDict[host].urls[outpath] = Node{currenttime, contentByte}
// }

// func loadLog(logURL string) {
// 	lines := strings.Split(readFile(logURL), "\n")
// 	count := 0

// 	for i := 0; i < len(lines); i++ {
// 		line := strings.Split(lines[i], ";")
// 		if len(line) != 3 {
// 			continue
// 		}
// 		ctime, _ := strconv.Atoi(line[2])
// 		addMemoryDick(line[0], line[1], ctime)
// 		count++
// 		fmt.Printf("%d. host:%s, url: %s\n", count, line[0], line[1])
// 	}

// }

// func displayDict(intercal int64) {
// 	timer := time.NewTicker(time.Duration(intercal) * time.Second)
// 	for {
// 		select {
// 		case <-timer.C:
// 			{
// 				println("~~~~~~~~~~~~~~~~~~")
// 				for host, v := range domainDict {
// 					println("HOST", host)
// 					for url, s := range v.urls {
// 						println("url", url, len(s.body))
// 					}
// 				}
// 				println("~~~~~~~~~~~~~~~~~~")
// 			}
// 		}
// 	}
// }

// func main() {
// 	domainDict = make(map[string]domainNode)
// 	loadLog("./indexLog.txt")
// 	displayDict(10)

// }
