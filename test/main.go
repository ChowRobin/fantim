package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func Echo(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		//websocket接受信息
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}

		fmt.Println("reveived from client: " + reply)
		msg := "received:" + reply
		fmt.Println("send to client:" + msg)
		//这里是发送消息
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("send failed:", err)
			break
		}
	}
}

func web(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	fmt.Println("method", r.Method)
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端
		t, err := template.ParseFiles("test/websocket.html")
		if err != nil {
			log.Println(err)
			return
		}
		_ = t.Execute(w, nil)
	} else {
		//否则走打印输出post接受的参数username和password
		fmt.Println(r.PostFormValue("username"))
		fmt.Println(r.PostFormValue("password"))
	}
}

func main() {

	//接受websocket的路由地址
	http.Handle("/websocket", websocket.Handler(Echo))
	//html页面
	http.HandleFunc("/web", web)
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
