package main

import (
	"fmt"
	"net/http"
	"services"

	"golang.org/x/net/websocket"
	//    "encoding/json"
	//    "bufio"
	//    "os"
)

const (
	ips = "101.69.242.10"
)

type DevInfo struct {
	SN       string `bson:"sn"`
	User     string `bson:"user"`
	Password string `bson:"password"`
	ENV      string `bson:"env"`
	Protocol string `bson:"StreamProtocol"`
}
type ResInfo struct {
	SN   string `bson:"sn"`
	Env  string `bson:"env"`
	Type string `bson:"type"`
	Url  string `bson:"url"`
	Err  string `bson:"err"`
	Time string `bson:"time"`
}
type TestCase struct {
	Env string    `json:"env"`
	Dev []DevInfo `json:"DevInfo"`
}
type TestID struct {
	Id   string     `bson:"_id"`
	Case []TestCase `bson:"case"`
}

func main() {
	// 调用外部的exe程序，并传递参数
	//cmdFile := exec.Command("C:/Users/Admin/AppData/Local/Google/Chrome/Application/chrome.exe", "D:/code/MyEye_SDK/XMediaNew/XMedia/Source/WebRTC/ClientDemo/testFile.html")
	//cmdCamera := exec.Command("C:/Users/Admin/AppData/Local/Google/Chrome/Application/chrome.exe", "D:/code/MyEye_SDK/XMediaNew/XMedia/Source/WebRTC/ClientDemo/testCamera.html")
	// 执行命令
	//err := cmdFile.Start()
	//err := cmdCamera.Start()
	http.Handle("/ws", websocket.Handler(services.TestResWebSocket))
	http.HandleFunc("/upload", services.UploadHandler)
	fmt.Println("Server is running on http://localhost:9250")
	http.ListenAndServe(":9250", nil)
	//if err != nil
	{
		//panic(err)
	}
}
