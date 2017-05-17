package example

/*

#cgo CFLAGS:-g -Wall -I./include
#cgo LDFLAGS:-L/home/xieyuan/IdeaProjects/example/src/include -lmsc

#include "convert.h"

*/
import "C"
import "fmt"

/*
* rdn:           合成音频数字发音方式
* volume:        合成音频的音量
* pitch:         合成音频的音调
* speed:         合成音频对应的语速
* voice_name:    合成发音人
* sample_rate:   合成音频采样率
* text_encoding: 合成文本编码格式
*
* 详细参数说明请参阅《iFlytek MSC Reference Manual》
 */

const (
	voice_kind    = "yanping"
	Login_params  = "appid = 586f0a5a, work_dir =."
	ttsconfig_params = "voice_name=" + voice_kind + ", text_encoding=UTF8, sample_rate=16000"
)

//var ttsParams *C.char
//var sleep C.int = C.int(0)
//
//func SetTTSParams(params string) {
//	ttsParams = C.CString(params)
//}
//
//func SetSleep(t int) {
//	sleep = C.int(t)
//}

func TTSLogin(loginParams string) error {
	ret := C.MSPLogin(nil, nil, C.CString(loginParams))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("登录失败，错误码：%d", int(ret))
	}
	return nil
}

func TTSLogout() error {
	ret := C.MSPLogout()
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("注销失败，错误码：%d", int(ret))
	}
	return nil
}

func TextToSpeech(text, outPath string) error {
	TTSLogin(Login_params)
	ret := C.text_to_speech(C.CString(text), C.CString(outPath), C.CString(ttsconfig_params), C.int(0))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("音频生成失败，错误码：%d", int(ret))
	}

	TTSLogout()
	return nil
}
