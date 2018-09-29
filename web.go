package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// the button to enable log printer
var logEnable bool
var staticURL string
var logURL string
var expiration_time int
var cache_size int
var replacement_policy string

/**************
* domainDict has 2-level indexs. The first index is domain and the second index is url(after hash-sha256)
* All data(including image pixels or text) is stored as []bytes
* e.g.
*	return data of image : domainNode['www.google.com']['256bit.jpg'].body
*	return timestamp of image : domainNode['www.google.com']['256bit.jpg'].time
 */
var domainDict map[string]domainNode

type Node struct {
	timestamp int
	body      []byte // content
}

type domainNode struct {
	urls map[string]Node
}

/*
Name: getTime
@ para: None
@ Return: string
Func: return the current time by format
*/
func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*
Name: PathExists
@ para: path string
@ Return: bool
Func: return ture if there exists a file according to path or return false if not
*/
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

/*
Name: removeFile
@ para: path string
@ Return: None
Func: remove the file according to the path
*/
func removeFile(path string) {
	if PathExists(path) == false {
		return // return directly if there exists no such file
	}
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
	if logEnable == true {
		println(getTime(), ": remove ", path)
	}
}

func readFileByte(filePath string) []byte {
	if PathExists(filePath) == false {
		_, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

/*
Name: readFile
@ para: filePath string
@ Return: string
Func: read and then return the content from file in corresponding path
*/
func readFile(filePath string) string {
	return string(readFileByte(filePath))
}

func readImageByte(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(getTime(), ": readImageByte:", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file) //解码
	if err != nil {
		fmt.Println(getTime(), ": readImageByte:", err)
	}
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

/*
Name: writeFile
@ para: filePath string
@ para: content string
@ para: appendEnable string
@ Return: None
Func: write the string content into assigned path by method of overwriting or appending
*/
func writeFile(filePath string, content string, appendEnable bool) {
	if appendEnable == false {
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
		f.WriteString(content)
	} else {
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
		f.WriteString(content)
	}
}

/*
Name: makedir
@ para: path string
@ Return: None
Func: build a new dictionary if there is no correspondinng dictionary
*/
func makedir(path string) {
	if PathExists(path) == true {
		return
	}
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	if logEnable == true {
		println(getTime(), ": Make dir in ", path)
	}
}

/*
Name: DirSize
@ para: path string
@ Return: int64, error
Func: Get the size(KB) of a Dir
*/
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size / 1000, err
}

func checkCacheSize() {
	size, _ := DirSize(staticURL)
	if int(size) > cache_size {
		println(getTime(), ": exceed cache_size", "-", cache_size, "now is ", size)
		if replacement_policy == "LRU" {
			LRUPolicy()
		} else if replacement_policy == "LFU" {
			LFUPolicy()
		}
	}
}

func LRUPolicy() {
	dir_list, e := ioutil.ReadDir(staticURL)
	if e != nil {
		fmt.Println("read dir error")
	}
	if logEnable == true {
		println("Cache Eviction and Replacement:")
	}

	lastfilepath := ""
	lasttime := 0

	for _, v := range dir_list {
		filepath := staticURL + "/" + v.Name()
		filetime := int(getUpdateUnixTime(filepath))
		if lasttime == 0 {
			lasttime = filetime
			lastfilepath = filepath
		} else {
			if filetime < lasttime {
				lasttime = filetime
				lastfilepath = filepath
			}
		}
	}

	if logEnable == true {
		println(getTime(), ": delete LRU file :", lastfilepath)
	}
	removeFile(lastfilepath)
}

func LFUPolicy() {

}

/*
Name: getSha256Code
@ para: string s
@ Return: string
Func: return the hash code of input string by the method of SHA256
*/
func getSha256Code(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func fileType(filename string) string {
	list := strings.Split(filename, ".")
	return list[len(list)-1]
}

/*
Name: splitname
@ para: url string
@ Return: string
Func: split the url and get the last element of url, then  transfer it in method of hash
*/
func hashName(url string) string {
	outpath := getSha256Code(url)
	filetype := fileType(url)
	if filetype == "jpg" ||
		filetype == "jpeg" ||
		filetype == "png" ||
		filetype == "js" ||
		filetype == "css" {
		outpath += "." + fileType(url)
	} else {
		outpath += ".html"
	}
	return outpath
}

/*
Name: touch
@ para: url string
@ Return: None
Func: It seems like the "touch" command in linux, which changes the last update time of a file
*/
func touch(url string) {
	// return directly if there exists no such file
	if PathExists(url) == false {
		return
	}
	timestamp := time.Now()
	if logEnable == true {
		println(getTime(), ": update timestamp ")
	}
	os.Chtimes(url, timestamp, timestamp)
}

/*
Name: touch
@ para: url string
@ Return: None
Func: It seems like the "touch" command in linux, which changes the last update time of a file
*/
func download(url string, outpath string) error {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	// request remote source and store it into disk
	resp, err := c.Get(url)
	if err != nil {
		fmt.Println("LALA", err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println(getTime(), ": StatusCode ", resp.StatusCode) // success
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	f, err1 := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, os.ModePerm) // make a new file
	if err1 != nil {
		// panic(err1)
		fmt.Println(err)
		log.Fatal(err)
		return err
	}
	defer f.Close()

	for {
		n, _ := resp.Body.Read(buf)
		if 0 == n {
			break
		}
		f.WriteString(string(buf[:n]))
	}
	return nil
}

/*
Name: readByte
@ para: url string
@ Return: []byte
Func: judge the type of file first and read the content in format of []byte
*/
func readByte(url string) []byte {
	var contentByte []byte
	if fileType(url) == "jpg" ||
		fileType(url) == "jpeg" ||
		fileType(url) == "png" {
		contentByte = readImageByte(url)
	} else {
		contentByte = readFileByte(url)
	}
	return contentByte
}

func addRecord(host string, outpath string, currenttime int) {
	content := host + ";" + outpath + ";" + strconv.Itoa(currenttime) + "\n"
	writeFile(logURL, content, true)
}

func addMemoryDick(host string, outpath string, currenttime int) bool {
	if PathExists(outpath) == false {
		return false
	}

	contentByte := readByte(outpath)

	_, ok := domainDict[host]
	if ok == false {
		domainDict[host] = domainNode{make(map[string]Node)}
	}

	domainDict[host].urls[outpath] = Node{currenttime, contentByte}
	return true
}

/*
Name: getRemoteContent
@ para: url string
@ Return: string
Func: firstly check whether the resouce in local disk. if so, return the url to the local resource
	otherwise, get the remote content by http request and download it into local disk
	finally return the path of corresponding resource in the disk
*/
func getRemoteContent(host string, url string) (string, error) {

	// get the url after hash sha-256
	makedir(staticURL + "/" + host)
	outpath := staticURL + "/" + host + "/" + hashName(url)

	if logEnable == true {
		println(getTime(), ": Get source in ", url)
		println(getTime(), ": Outpath is ", outpath)
	}

	// if the comtent has been in disk, then just return uri instead of requesting remote source
	if PathExists(outpath) == true {
		touch(outpath)
		if logEnable == true {
			println(getTime(), ": In webcache: True")
		}
	} else {
		err := download(url, outpath)
		if err != nil {
			return "Wrong", err
		}
		currenttime := int(time.Now().Unix())

		addRecord(host, outpath, currenttime)

		if logEnable == true {
			println(getTime(), ": Insert into memory dict", outpath, domainDict[host].urls[outpath].timestamp)
			println(getTime(), ": In webcache: False")
			println(getTime(), ": Download ", url, " into")
			println(getTime(), ": ", outpath)
		}
	}

	checkCacheSize()

	return outpath, nil
}

/*
Name: parsingHTML
@ para: content string
@ Return: string
Func: parse the whole content of html, find the link and download the linked source. Replace the remote link in
	the html with new url linked to the source in webcache. Finally, return the new html after being replaced
*/
func parsingHTML(host string, content string) string {
	exp1 := regexp.MustCompile(`(src|href)=".+?"`) // filter the link source by regular expressions
	res1 := exp1.FindAllStringSubmatch(content, -1)

	for i := 0; i < len(res1); i++ {
		// fmt.Println("%v", res1[i])
		url := res1[i][0]
		for j := 0; j < len(url); j++ {
			if string(url[j]) == "\"" {
				url = url[j+1 : len(url)-1]
				break
			}
		}

		if PathExists(url) == true {
			continue
		}

		localurl, err := getRemoteContent(host, url)
		if err != nil {
			continue
		}
		// replace the old url into new ones
		// newurl := "http://127.0.0.1:8080" + localurl[1:]
		newurl := localurl
		content = strings.Replace(content, url, newurl, -1)

	}
	return content
}

/*
Name: getHTML
@ para: htmlurl string
@ Return: string
Func: get the html by functon getRemoteContent and parse it and finally return it as string
*/
func getHTML(host string, htmlurl string) string {
	url, err := getRemoteContent(host, htmlurl)
	if err != nil {
		return "404 Bad Request"
	}
	data := readFile(url)
	data = parsingHTML(host, data)
	writeFile(url, data, false)
	return data
}

/*
Name: cache
@ para: w http.ResponseWriter
@ para: r *http.Request
@ Return: None
Func: the handler function of http.HandleFunc("/").
	When receiving the http request, the goroutine cache starts
*/
func cache(w http.ResponseWriter, r *http.Request) {

	// ignore all the methods except "GET"
	if r.Method != "GET" {
		println(getTime(), ": It's not a GET method so webcache ignores it.")
		return
	}

	host := r.URL.Hostname()
	host_port := r.URL.Host
	path := r.URL.Path
	httpurl := "http://" + host_port + path
	localurl := "." + path
	hashurl := staticURL + "/" + host + "/" + hashName(httpurl)

	if logEnable == true {
		println("--------------------------------------------")
		println("Hi, here comes a HTTP request whose method is", r.Method)
		println("HOST_PORT:", host_port)
		println("PATH:", path)
		// println("HTTPurl:", httpurl)
		// println("Localurl:", localurl)
		// println("Hashurl:", hashurl)
		println()
	}

	// check in data structure
	_, ok_1 := domainDict[host]
	if ok_1 == true {
		_, ok_2 := domainDict[host].urls[localurl]
		if ok_2 == true {
			println(getTime(), ": GET ", localurl)
			fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, getTime()+" : Read from webcache memory: True, ok_2", 0x1B)

			if fileType(localurl) == "css" {
				w.Header().Set("Content-Type", "text/css")
			}

			w.Write(domainDict[host].urls[localurl].body)

			if logEnable == true {
				println("\nAll work has been done")
				println("--------------------------------------------\n")
			}
			return
		}
		_, ok_3 := domainDict[host].urls[hashurl]
		if ok_3 == true {
			println(getTime(), ": GET ", hashurl)
			print(getTime())
			fmt.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, getTime()+" : Read from webcache memory: True, ok_3", 0x1B)

			if fileType(hashurl) == "css" {
				w.Header().Set("Content-Type", "text/css")
			}

			w.Write(domainDict[host].urls[hashurl].body)

			if logEnable == true {
				println("\nAll work has been done")
				println("--------------------------------------------\n")
			}
			return
		}
	}

	// then check disk
	if fileType(path) == "jpg" ||
		fileType(path) == "jpeg" ||
		fileType(path) == "png" ||
		fileType(path) == "css" ||
		fileType(path) == "js" {

		// if the type is css, must set the header first
		if fileType(path) == "css" {
			w.Header().Set("Content-Type", "text/css")
		}

		if PathExists(localurl) == true {
			touch(localurl)
			if logEnable == true {
				println(getTime(), ": GET ", localurl)
				println(getTime(), ": In webcache disk: True")
			}
			w.Write(readByte(localurl))
		} else {
			println(getTime(), ": try to requst remote server by http ")
			imageURL, err := getRemoteContent(host, httpurl)
			if err != nil {
				println(getTime(), ": Can't load the content")
				return
			}
			w.Write(readByte(imageURL))
		}

	} else {
		if logEnable == true {
			println(getTime(), ": GET ", httpurl)
		}
		data := getHTML(host, httpurl)
		fmt.Fprintf(w, data)
	}

	if logEnable == true {
		println("\nAll work has been done")
		println("--------------------------------------------\n")
	}
}

/*
Name: getUpdateUnixTime
@ para: url string
@ Return: int64
Func: get the last update time of a file by the url
*/
func getUpdateUnixTime(url string) int64 {
	fileInfo, _ := os.Stat(url)
	return fileInfo.ModTime().Unix()
}

/*
Name: checkFiles
@ para: url string
@ Return: None
Func: check the file in assigned dictionary and check whether they have been expired
*/
func checkFiles(url string) {
	dir_list, e := ioutil.ReadDir(url)
	if e != nil {
		fmt.Println("read dir error")
	}
	if logEnable == true {
		println("--------------------------------------------")
		println("Hi, it's time to clean exipired files:")
	}
	for _, v := range dir_list {
		filepath := url + "/" + v.Name()
		gaptime := time.Now().Unix() - getUpdateUnixTime(filepath)
		if int(gaptime) > expiration_time {
			removeFile(filepath)
		}
	}
	if logEnable == true {
		println("--------------------------------------------\n")
	}
	println()
}

/*
Name: persistenceCheck
@ para: None
@ Return: None
Func: check the files in static source dictionary and delete expired ones
*/
func persistenceCheck(interval int) {
	timer := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-timer.C:
			{
				checkFiles(staticURL)
			}
		}
	}
}

func displayDict(intercal int64) {
	timer := time.NewTicker(time.Duration(intercal) * time.Second)
	for {
		select {
		case <-timer.C:
			{
				println("~~~~~~~~~~~~~~~~~~")
				for host, v := range domainDict {
					println("HOST", host)
					for url, s := range v.urls {
						println("url", url, len(s.body))
					}
				}
				println("~~~~~~~~~~~~~~~~~~")
			}
		}
	}
}

func loadLog(logURL string) {
	lines := strings.Split(readFile(logURL), "\n")
	count := 0
	if logEnable == true {
		println(getTime(), ": restart webcache and load memory.\n")
	}

	for i := 0; i < len(lines); i++ {
		line := strings.Split(lines[i], ";")
		if len(line) != 3 {
			continue
		}
		ctime, _ := strconv.Atoi(line[2])
		readstatus := addMemoryDick(line[0], line[1], ctime)
		if readstatus == false {
			continue
		}
		count++
		fmt.Printf("%d. host:%s, url: %s\n", count, line[0], line[1])
	}

	if logEnable == true {
		println()
		println(getTime(), ": loading webcache done.")
	}
}

/*
Name: Init
@ para: None
@ Return: None
Func: Initialization of the program in the first step
*/
func Init() {

	if logEnable == true {
		println(getTime(), ": Initializing...")
	}

	staticURL = "./static"
	expiration_time = 60
	cache_size = 6000
	replacement_policy = "LRU"
	logURL = "./indexLog.txt"
	logEnable = true

	domainDict = make(map[string]domainNode)
	loadLog(logURL)

	if PathExists(staticURL) == false {
		makedir(staticURL)
	}

	if logEnable == true {
		println(getTime(), ": Initializing Done.\n")
	}
}

func main() {
	println(getTime(), ": WebCache starts monitoring port 8080")

	go displayDict(30)

	// Timer to clean the expired files
	// go persistenceCheck(10)

	Init()
	http.HandleFunc("/", cache)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
