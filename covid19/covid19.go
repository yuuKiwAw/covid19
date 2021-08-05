package covid19

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 获取百度新冠病毒原始信息
func bdcovid19_response() string {
	client := &http.Client{}
	url := "https://voice.baidu.com/act/newpneumonia/newpneumonia"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer resp.Body.Close()

	log.Println("get baidu_covid19 success")
	return string(body)
}

// io保存到本地文件
func savelocal(info string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("write io error", err)
		return
	}
	file.WriteString(info)
	log.Println("save file success")
	defer file.Close()
}

// io读取本地文件
func readlocal(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open io error", err)
		return err.Error()
	}
	defer file.Close()

	var buffer [128]byte
	var content []byte
	for {
		n, err := file.Read(buffer[:])
		if err == io.EOF {
			// 文件读取完毕
			break
		}
		if err != nil {
			log.Fatal(err.Error())
			return err.Error()
		}
		content = append(content, buffer[:n]...)
	}
	return string(content)
}

// 筛选出covid19相关的信息
func selection_result(htmlDom string) string {
	var outValue string
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlDom))
	if err != nil {
		log.Fatalln(err)
		return err.Error()
	}

	dom.Find("#captain-config").Each(func(i int, selection *goquery.Selection) {
		outValue = selection.Text()
		log.Println("selection info success")
	})
	return outValue
}

func GetCovid19info() {
	html_filePath := "./html/bdcovid19.html"
	save_jsonPath := "./data/covid19info_selected_all.json"

	// 获取response并且保存到本地html
	covid19info := bdcovid19_response()
	savelocal(covid19info, html_filePath)

	// 解析获取到的信息
	htmlDom := readlocal(html_filePath)
	covid19info_selected_all := selection_result(htmlDom)
	savelocal(covid19info_selected_all, save_jsonPath)
}

func init() {
	log_infoPath := "./log/logs.log"

	// 保存日志信息
	logFile, err := os.OpenFile(log_infoPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}

	// 同时显示在终端并写入到log文件
	writers := []io.Writer{
		logFile,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)

	log.SetOutput(fileAndStdoutWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
