package sms

import (
	"launchpad.net/xmlpath"
	"io"
)

type httphfy struct{}

func (hfy *httphfy)Paser(reader io.Reader)(*Response){
	rsp := new(Response)	
	root, err := xmlpath.Parse(reader)
	if err != nil {
		rsp.Status = FAILED
		rsp.Msg = "短信平台服务错误"
		return rsp	
	}	

	pathstatus := xmlpath.MustCompile("//returnsms/returnstatus")
	pathmessage := xmlpath.MustCompile("//returnsms/message")
	
	if status, ok := pathstatus.String(root); ok {
		if status == "Success" {
			rsp.Status = SUCCESS	
		}else{	
			rsp.Status = FAILED	
		}
	}else{
		rsp.Status = FAILED	
		rsp.Msg = "短信平台服务错误"	
		return rsp	
	}
	
	if message,ok := pathmessage.String(root);ok {
		rsp.Msg = "短信平台:"+message	
	}
	return rsp
}


