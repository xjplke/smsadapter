package sms

import (
	"strings"
	"github.com/Unknwon/goconfig"
	"log"
	"net/http"
	"net/url"
	"io"
)

type EStatus string

const (
	SUCCESS EStatus = "success"
	FAILED	EStatus = "failed"
)

type Response struct {
	Status EStatus
	Msg string 
}

type Smssender interface {
	Init(c *goconfig.ConfigFile) error
	Send(phone,message string) (Response)
}	

type RspPaser interface {
	Paser(reader io.Reader) (*Response)
}


type SmsHttp struct {
       	url,vendor,phonekey,messagekey string 
        params map[string] string  //url param	
	paser RspPaser		
}

var conf *goconfig.ConfigFile


func NewSmsHttp()(* SmsHttp){
	return new(SmsHttp)	
}

func getValue(s,k string)(string){
	ret,err := conf.GetValue(s,k)
	if err != nil {
		log.Fatal(err)
	}	
	return ret
}


func (sh *SmsHttp)Init(c *goconfig.ConfigFile){
	conf = c	
	sh.url = getValue("http","url")
	sh.vendor = getValue("http","vendor")	
	sh.phonekey = getValue("http","phonekey")
	sh.messagekey = getValue("http","messagekey")	

	s,err := conf.GetSection(sh.vendor)
	if err != nil {
		log.Fatal(err)
	}	
	sh.params = s
	
	sh.paser = new(httphfy)
}

func (sh *SmsHttp)Send(phone,password string)(r *Response){
	params := url.Values{}  
	rsp := new(Response)
		
	for k,v := range sh.params {
		params.Add(k,v);	
	}		
	params.Set(sh.phonekey,phone)
	params.Set(sh.messagekey,strings.Replace(params.Get(sh.messagekey),"%password%",password,-1))
	postrsp,err := http.PostForm(sh.url,params)
	if err != nil {
		log.Print(err)	
		rsp.Status = FAILED;
		rsp.Msg = "短信服务器请求错误"	
		return rsp		
	}
	defer postrsp.Body.Close()
//    	body, err := ioutil.ReadAll(postrsp.Body)		
//	s := string(body)
//	fmt.Println(s)

	rsp = sh.paser.Paser(postrsp.Body)	
	return rsp
}



