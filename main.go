package main

import (
	"gitlab.alibaba-inc.com/zhoushua.ljj/hand/pkg/user"
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/user", user.DdLoginHandler)
	http.ListenAndServe(":8888", nil)
}