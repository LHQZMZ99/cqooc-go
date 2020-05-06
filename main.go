package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var sessionId string

const (
	baseUrl   = "http://www.cqooc.com"
	userName  = "1111111111"
	courseId  = "2222222"
	courseCId = "3333333333"
	xsid      = "4444444444"
)

// session json 结构
type SessionStruct struct {
	Id       int64  `json:"id"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}

type Chapters struct {
	Datas []ChaptersData `json:"data"`
}

type ChaptersData struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Level string `json:"level"`
}

type StatusData struct {
	Code int    `json:"code"`
	ID   string `json:"id"`
	Msg  string `json:"msg"`
}

// 使用 https://mholt.github.io/json-to-go/ 生成
type Lessons struct {
	Meta struct {
		Total string `json:"total"`
		Start string `json:"start"`
		Size  string `json:"size"`
	} `json:"meta"`
	Data []struct {
		ID       string `json:"id"`
		ParentID string `json:"parentId"`
		Chapter  struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		} `json:"chapter"`
		CourseID string `json:"courseId"`
		Category string `json:"category"`
		Title    string `json:"title"`
		ResID    string `json:"resId"`
		Resource struct {
			ID                    string      `json:"id"`
			Title                 string      `json:"title"`
			AuthorName            string      `json:"authorName"`
			ResSort               string      `json:"resSort"`
			ResMediaType          string      `json:"resMediaType"`
			ResSize               string      `json:"resSize"`
			Viewer                string      `json:"viewer"`
			Oid                   string      `json:"oid"`
			Username              string      `json:"username"`
			ResMediaTypeLkDisplay string      `json:"resMediaType_lk_display"`
			Pages                 interface{} `json:"pages"`
			Duration              string      `json:"duration"`
			Dimension             string      `json:"dimension"`
			ResourceTypeLkDisplay string      `json:"resourceType_lk_display"`
		} `json:"resource"`
		TestID      interface{} `json:"testId"`
		Test        string      `json:"test"`
		ForumID     interface{} `json:"forumId"`
		Forum       interface{} `json:"forum"`
		OwnerID     string      `json:"ownerId"`
		Created     int64       `json:"created"`
		LastUpdated int64       `json:"lastUpdated"`
		Owner       string      `json:"owner"`
		ChapterID   string      `json:"chapterId"`
		SelfID      string      `json:"selfId"`
		IsLeader    string      `json:"isLeader"`
	} `json:"data"`
}

// 获取当前系统的时间戳
func getTs() string {
	ts := time.Now().Unix()
	return strconv.FormatInt(ts, 10)
}

// http 请求
func getHttp(url string) []byte {
	client := &http.Client{}
	request, err := http.NewRequest("GET", baseUrl+url, nil)

	if err != nil {
		return nil
	}

	request.Header.Add("Host", "www.cqooc.com")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", "player=1; xsid="+xsid)
	request.Header.Add("Referer", baseUrl+"/learn/mooc/structure?id="+courseId)
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	resp, err := client.Do(request)
	bytes, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil
	}

	return bytes
}

func postHttp(url string, postData []byte) {

}

func getSessionId() string {
	bytes := getHttp("/user/session?xsid=" + xsid + "&ts=" + getTs())

	var v SessionStruct
	err := json.Unmarshal(bytes, &v)

	if err != nil {
		return ""
	}

	sessionId := strconv.FormatInt(v.Id, 10)
	fmt.Println("成功获取 session ID: " + sessionId)

	return sessionId
}

func watchVideo(parentId string, sectionId string, chapterId string) {
	postData := []byte(fmt.Sprintf(`
    {"username": %s,
     "ownerId": %s,
     "parentId": %s,
     "action": 1,
     "courseId": %s,
     "sectionId": %s,
     "chapterId": %s,
     "category": 2
    }`, userName, sessionId, parentId, courseId, sectionId, chapterId))

	// postHttp("/learnLog/api/add", postData)
	url := baseUrl + "/learnLog/api/add"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	request.Header.Add("Host", "www.cqooc.com")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Cookie", "player=1; xsid="+xsid)
	request.Header.Add("Referer", baseUrl+"/learn/mooc/structure?id="+courseId)
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	// defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		var data StatusData
		err := json.Unmarshal(body, &data)

		if err != nil {
			return
		}

		code := data.Code

		if code == 2 {
			fmt.Println("已经添加, 跳过")
			return
		} else if code == 0 {
			fmt.Println(string(body))
			// sleep
			return
		} else if code == 3 {
			fmt.Println("非法操作")
			return
		}
	} else {
		fmt.Println("post error：" + resp.Status)
	}
}

func doChapter(id string) {
	lessonsBytes := getHttp("/json/mooc/lessons?parentId=" + id + "&limit=100&sortby=selfId&reverse=false&ts=" + getTs())
	var les Lessons
	err := json.Unmarshal(lessonsBytes, &les)

	if err != nil {
		return
	}

	// fmt.Println(string(lessonsBytes))
	// time.Sleep(time.Second * 2)
	// return

	if len(les.Data) == 0 {
		fmt.Println(id + "跳过")
		return
	}

	for i := 0; i < len(les.Data); i++ {
		data := les.Data[i]
		category := data.Category

		if category == "1" {
			// 分类 1 是视频
			watchVideo(courseCId, data.ID, data.ParentID)
		} else if category == "2" {
			// 分类 2 是单元测试
			fmt.Println("单元测试，未实现，跳过")
			continue
		} else if category == "3" {
			// 分类 3 是论坛
			fmt.Println("论坛，未实现，跳过")
			continue
		}
	}
}

func main() {
	sessionId = getSessionId()

	// 获取课程列表
	bytes := getHttp("/json/chapters?status=1&select=id,title,level&courseId=" + courseId + "&sortby=selfId&reverse=false&limit=200&start=0&ts=" + getTs())
	var v Chapters
	err := json.Unmarshal(bytes, &v)

	if err != nil {
		fmt.Println("error")
		return
	}

	for i := 0; i < len(v.Datas); i++ {
		doChapter(v.Datas[i].Id)
	}
}
