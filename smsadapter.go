package main

import (
    "flag"
    "path/filepath"
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
    "strconv"
    "github.com/Unknwon/goconfig"
    "github.com/xjplke/smsadapter/sms"
    "fmt"
)

type Message struct {
    Body string
}

func main() {
	path := flag.String("config", "./smsadapter.conf", "设置配置文件的路径")	
        flag.Parse()
	*path = filepath.FromSlash(*path) 

	c, err := goconfig.LoadConfigFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	
	ip, iperr := c.GetValue("local","listen")
	if iperr != nil {
		log.Printf("%v use default:0.0.0.0",iperr)
		ip = "0.0.0.0"
	}
	
	port,porterr := c.Int("local","port")
	if porterr != nil {
		log.Panicf("%v,use default:8500",porterr)
		port = 8500
	}
	
    smssender := sms.NewSmsHttp()
    smssender.Init(c)

	
    fmt.Println("listen on "+ip+":"+strconv.Itoa(port))	
    handler := rest.ResourceHandler{}
    errx := handler.SetRoutes(
        &rest.Route{"GET", "/message", func(w rest.ResponseWriter, req *rest.Request) {
            w.WriteJson(&Message{
                Body: "Hello World!",
            })
        }},
	&rest.Route{"GET","/sms/send",func(w rest.ResponseWriter,req *rest.Request){
            w.Header().Set("Content-Type","text/json; charset=utf-8")
            phone := req.FormValue("phone")  
            password := req.FormValue("password")  
	    rsp := smssender.Send(phone,password)
            w.WriteJson(rsp)	
	}},
    )
    if errx != nil {
        log.Fatal(errx)
    }
    
    log.Fatal(http.ListenAndServe(ip+":"+strconv.Itoa(port), &handler))
    fmt.Println("exit1") 
}

