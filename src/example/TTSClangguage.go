package example
//
//
///*
//
//#cgo CFLAGS:-g -Wall -I./include
//#cgo LDFLAGS:-L/home/xieyuan/IdeaProjects/example/src/include -lmsc
//
//#include "qtts.h"
//#include "msp_cmn.h"
//#include "msp_errors.h"
//#include "msp_types.h"
//*/
//import "C"
//import (
//	"unsafe"
//	"bytes"
//)
//
//func TextTToSpeech( result *bytes.Buffer, text string)(retnum int){
//	var ret = C.int(-1)
//	var sessionID *C.char
//	var dataLen C.uint
//	var synth_status = C.MSP_TTS_FLAG_STILL_HAVE_DATA
//	var buffer bytes.Buffer
//
//	if (text==""){
//		return int(ret);
//	}
//	sessionID = C.QTTSSessionBegin(C.CString(ttsconfig_params), &ret)
//	if (C.MSP_SUCCESS != ret) {
//		return int(ret);
//	}
//	ret = C.QTTSTextPut(sessionID, C.CString(text), C.uint(len(text)), nil);
//	if (C.MSP_SUCCESS != ret) {
//	C.QTTSSessionEnd(sessionID, C.CString("TextPutError"));
//		return int(ret);
//	}
//	for true {
//		/* 获取合成音频 */
//
//		data := unsafe.Pointer(C.QTTSAudioGet(sessionID, &dataLen, &synth_status, &ret));
//		if (C.MSP_SUCCESS != ret){break;}
//
//
//		if (data!=nil) {
//			buffer.WriteString(string(*data))
//		}
//		if (C.MSP_TTS_FLAG_DATA_END == synth_status) {
//			break;
//		}
//	}//合成状态synth_status取值请参阅《讯飞语音云API文档》
//	result = &buffer
//
//	if (C.MSP_SUCCESS != ret) {
//		C.QTTSSessionEnd(sessionID, C.CString("AudioGetError"));
//		return int(ret);
//	}
//
//	ret = C.QTTSSessionEnd(sessionID, C.CString("Normal"));
//	return int(ret);
//}