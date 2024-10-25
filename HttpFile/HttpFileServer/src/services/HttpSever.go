package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"
)

var testRes = 0
var testTalkRes = 0
var PlayURL = ""
var TalkURL = ""

func TestResWebSocket(conn *websocket.Conn) {
	for {
		var msg string
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 打印接收到的消息
		fmt.Println(msg)

		if msg == "ReqTalkURL" {
			if TalkURL == "" {
				fmt.Println("recv ReqTalkURL,but TalkURL NULL")
			} else {
				websocket.Message.Send(conn, TalkURL)
			}
		}
		if msg == "ReqPlayURL" {
			if PlayURL == "" {
				fmt.Println("recv ReqTalkURL,but PlayURL NULL")
			} else {
				websocket.Message.Send(conn, PlayURL)
			}
		}

		if msg == "PlaySuccess" {
			testRes++
		}
		if msg == "PlayFail" {
			testRes = 0
		}
		if msg == "TalkSuccess" {
			testTalkRes++
		}
		if msg == "TalkFail" {
			testTalkRes = 0
		}

		// 回复消息
		//err = websocket.Message.Send(conn, msg)

	}
}

// 函数名开头大写，表示可公开的
func TestResultHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case "OPTIONS":
		{
			// 设置允许的请求方法
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// 设置允许的请求头
			w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-headers,Access-Control-Allow-Origin,content-type")
			// 设置允许跨域请求的来源
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 处理请求
			w.WriteHeader(http.StatusOK)
		}
	case "POST":
		{
			res := r.URL.Query().Get("res")
			fmt.Println("Received res:", res)

			resTalk := r.URL.Query().Get("resTalk")
			fmt.Println("Received resTalk:", resTalk)

			if res == "1" {
				testRes++
			} else {
				testRes = 0
			}
			if resTalk == "1" {
				testTalkRes++
			} else {
				testTalkRes = 0
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)

		}
	case "GET":
		{
			fmt.Println("Received a GET request")
		}
	case "DELETE":
		fmt.Println("Received a DELETE request")
	default:
		fmt.Println("Received an unknown request method")
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case "OPTIONS":
		{
			// 设置允许的请求方法
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// 设置允许的请求头
			w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-headers,Access-Control-Allow-Origin,content-type")
			// 设置允许跨域请求的来源
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 处理请求
			w.WriteHeader(http.StatusOK)
		}
	case "POST":
		{
			// 限制请求方法为 POST
			if r.Method != http.MethodPost {
				http.Error(w, "不支持此方法", http.StatusMethodNotAllowed)
				return
			}

			// 解析表单数据，限制上传文件为 10 GB
			err := r.ParseMultipartForm(10 << 30) // 10 GB ，10 << 20) // 10 MB
			if err != nil {
				http.Error(w, "解析表单失败", http.StatusBadRequest)
				return
			}

			// 获取文件
			file, fileHeader, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "获取文件失败", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// 获取文件名
			fileName := fileHeader.Filename
			fmt.Println("上传的文件名:", fileName)

			// 设置保存路径
			savePath := filepath.Join("/work/share/upload", fileName) // 在当前目录下保存文件

			// 创建目标文件
			outFile, err := os.Create(savePath)
			if err != nil {
				http.Error(w, "创建文件失败", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			// 将上传的文件内容复制到目标文件
			_, err = io.Copy(outFile, file)
			if err != nil {
				http.Error(w, "保存文件失败", http.StatusInternalServerError)
				return
			}

			// 返回成功响应
			fmt.Fprintf(w, "文件上传成功: %s\n", savePath)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
		}
	case "GET":
		{
			fmt.Println("Received a GET request")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusBadRequest)
		}
	case "DELETE":
		{
			fmt.Println("Received a DELETE request")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		{
			fmt.Println("Received an unknown request method")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
