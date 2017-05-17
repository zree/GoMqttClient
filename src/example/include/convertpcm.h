//#include <stdio.h>
//#include <string.h>
//#include <stdlib.h>
//#include <unistd.h>
//#include <malloc.h>
//
//#include "qtts.h"
//#include "msp_cmn.h"
//#include "msp_errors.h"
//
//int textS_to_speech(char* result, const char* src_text, const char* params, int sleep_time)
//{
//	int          ret          = -1;
//	const char*  sessionID    = NULL;
//	unsigned int audio_len    = 0;
//	int          synth_status = MSP_TTS_FLAG_STILL_HAVE_DATA;
//
//
//	if (NULL == src_text)
//	{
//		return ret;
//	}
//
//	//开始合成
//	sessionID = QTTSSessionBegin(params, &ret);
//	if (MSP_SUCCESS != ret)
//	{
//		return ret;
//	}
//	ret = QTTSTextPut(sessionID, src_text, (unsigned int)strlen(src_text), NULL);
//	if (MSP_SUCCESS != ret)
//	{
//		QTTSSessionEnd(sessionID, "TextPutError");
//		return ret;
//	}
//
//	while (1)
//        {
//                /* 获取合成音频 */
//                const void* data = QTTSAudioGet(sessionID, &audio_len, &synth_status, &ret);
//                if (MSP_SUCCESS != ret)
//                        break;
//                if (NULL != data)
//                {
//
//                        char *temp =(char *)malloc((strlen(result)+audio_len+1)*sizeof(char));
//                        strcpy(temp,result);
//                        strcat(temp,data);
//                        result = temp;
//
//                }
//                if (MSP_TTS_FLAG_DATA_END == synth_status)
//                        break;
//                usleep(sleep_time); //防止频繁占用CPU
//        }//合成状态synth_status取值请参阅《讯飞语音云API文档》
//
//	if (MSP_SUCCESS != ret)
//	{
//		QTTSSessionEnd(sessionID, "AudioGetError");
//		return ret;
//	}
//
//	ret = QTTSSessionEnd(sessionID, "Normal");
//	return ret;
//}
