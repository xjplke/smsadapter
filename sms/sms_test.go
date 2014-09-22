package sms

import (
	"github.com/Unknwon/goconfig"
	"log"
	"testing"
)


func TestHttpSms(t *testing.T){
	c, err := goconfig.LoadConfigFile("../sms.conf")
        if err != nil {
                log.Fatal(err)
        }
	
	sms := NewSmsHttp()

	sms.Init(c)
	sms.Send("13426347659","123qwe")	
}

