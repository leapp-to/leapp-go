package main

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"text/template"
)

func main() {
	stdoutSockPath := os.Getenv("LEAPP_ACTOR_STDOUT_SOCK")
	if len(os.Args) > 1 && os.Args[1] == "server" {
		os.Remove(stdoutSockPath)
		if listener, err := net.Listen("unix", stdoutSockPath); err == nil {
			defer listener.Close()
			for {
				if sock, err := listener.Accept(); err == nil {
					defer sock.Close()
					io.Copy(os.Stdout, sock)
				}
			}
		} else {
			os.Stderr.Write([]byte(err.Error()))
		}
	} else {
		var tmpl *template.Template
		if len(os.Args) > 1 {
			funcs := template.FuncMap{"json": func(v interface{}) string {
				a, _ := json.Marshal(v)
				return string(a)
			}}
			tmpl, _ = template.New("temp").Funcs(funcs).Parse(os.Args[1])
		}
		if sock, err := net.Dial("unix", stdoutSockPath); err == nil {
			defer sock.Close()
			if tmpl == nil {
				io.Copy(sock, os.Stdin)
			} else {
				inputData := map[string]interface{}{}
				if json.NewDecoder(os.Stdin).Decode(&inputData) == nil {
					tmpl.Execute(sock, inputData)
				}
			}
		}
	}
}
