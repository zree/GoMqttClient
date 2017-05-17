package example
/*

#cgo CFLAGS:-g -Wall -I./include
#cgo LDFLAGS:-L/home/xieyuan/IdeaProjects/example/src/include -lmsc

#include "vertcon.h"

*/
import "C"
import (
	"fmt"
)

const (
	isrconfig_params = "sub=iat,ptt=0,aue=speex-wb;7,result_type=plain,result_encoding=utf8,language=zh_cn,accent=mandarin,sample_rate=16000,domain=music,vad_bos=2000,vad_eos=1000"
)


//var isrParams *C.char
var text *C.char

//func SetISRParams(params string) {
//	isrParams = C.CString(params)
//}

func ISRLogin(loginParams string) error {
	ret := C.MSPLogin(nil, nil, C.CString(loginParams))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("登录失败，错误码：%d", int(ret))
	}
	return nil
}

func ISRLogout() error {
	ret := C.MSPLogout()
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("注销失败，错误码：%d", int(ret))
	}
	return nil
}


func WavToText(audiofile string) string {
	ISRLogin(Login_params)
	text = C.run_asr(C.CString(audiofile), C.CString(isrconfig_params))
	ISRLogout()
	return C.GoString(text)
}

