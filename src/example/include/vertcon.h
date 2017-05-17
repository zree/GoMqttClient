
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "qisr.h"
#include "msp_cmn.h"
#include "msp_errors.h"


#define	BUFFER_SIZE 2048
#define HINTS_SIZE  100
#define GRAMID_LEN	128
#define FRAME_LEN	640 
 
//int get_grammar_id(char* grammar_id, unsigned int id_len)
//{
//	FILE*			fp				=	NULL;
//	char*			grammar			=	NULL;
//	unsigned int	grammar_len		=	0;
//	unsigned int	read_len		=	0;
//	const char*		ret_id			=	NULL;
//	unsigned int	ret_id_len		=	0;
//	int				ret				=	-1;
//
//	if (NULL == grammar_id)
//		goto grammar_exit;
//
//	fp = fopen("gm_continuous_digit.abnf", "rb");
//	if (NULL == fp)
//	{
//		goto grammar_exit;
//	}
//
//	fseek(fp, 0, SEEK_END);
//	grammar_len = ftell(fp); //��ȡ�﷨�ļ���С
//	fseek(fp, 0, SEEK_SET);
//
//	grammar = (char*)malloc(grammar_len + 1);
//	if (NULL == grammar)
//	{
//		goto grammar_exit;
//	}
//
//	read_len = fread((void *)grammar, 1, grammar_len, fp); //��ȡ�﷨����
//	if (read_len != grammar_len)
//	{
//		goto grammar_exit;
//	}
//	grammar[grammar_len] = '\0';
//
//	ret_id = MSPUploadData("usergram", grammar, grammar_len, "dtt = abnf, sub = asr", &ret); //�ϴ��﷨
//	if (MSP_SUCCESS != ret)
//	{
//		goto grammar_exit;
//	}
//
//	ret_id_len = strlen(ret_id);
//	if (ret_id_len >= id_len)
//	{
//		goto grammar_exit;
//	}
//	strncpy(grammar_id, ret_id, ret_id_len);
//
//grammar_exit:
//	if (NULL != fp)
//	{
//		fclose(fp);
//		fp = NULL;
//	}
//	if (NULL!= grammar)
//	{
//		free(grammar);
//		grammar = NULL;
//	}
//	return ret;
//}
char*		grammar_id				=	NULL;
char			rec_result[BUFFER_SIZE]		 	= {'\0'};
char* run_asr(const char* audio_file, const char* params)
{
	const char*		session_id						= NULL;
	char			hints[HINTS_SIZE]				= {'\0'}; //hintsΪ�������λỰ��ԭ�����������û��Զ���
	unsigned int	total_len						= 0;
	int 			aud_stat 						= MSP_AUDIO_SAMPLE_CONTINUE;		//��Ƶ״̬
	int 			ep_stat 						= MSP_EP_LOOKING_FOR_SPEECH;		//�˵���
	int 			rec_stat 						= MSP_REC_STATUS_SUCCESS;			//ʶ��״̬	
	int 			errcode 						= MSP_SUCCESS;

	FILE*			f_pcm 							= NULL;
	char*			p_pcm 							= NULL;
	long 			pcm_count 						= 0;
	long 			pcm_size 						= 0;
	long			read_size						= 0;

    memset(rec_result,0,sizeof(rec_result)/sizeof(char));
	if (NULL == audio_file)
		goto asr_exit;

	f_pcm = fopen(audio_file, "rb");
	if (NULL == f_pcm) 
	{
		goto asr_exit;
	}
	
	fseek(f_pcm, 0, SEEK_END);
	pcm_size = ftell(f_pcm); //��ȡ��Ƶ�ļ���С 
	fseek(f_pcm, 0, SEEK_SET);		

	p_pcm = (char*)malloc(pcm_size);
	if (NULL == p_pcm)
	{
		goto asr_exit;
	}

	read_size = fread((void *)p_pcm, 1, pcm_size, f_pcm); //��ȡ��Ƶ�ļ�����
	if (read_size != pcm_size)
	{
		goto asr_exit;
	}

	session_id = QISRSessionBegin(grammar_id, params, &errcode);
	if (MSP_SUCCESS != errcode)
	{
		goto asr_exit;
	}
	
	while (1) 
	{
		unsigned int len = 10 * FRAME_LEN; // ÿ��д��200ms��Ƶ(16k��16bit)��1֡��Ƶ20ms��10֡=200ms��16k�����ʵ�16λ��Ƶ��һ֡�Ĵ�СΪ640Byte
		int ret = 0;

		if (pcm_size < 2 * len) 
			len = pcm_size;
		if (len <= 0)
			break;
		
		aud_stat = MSP_AUDIO_SAMPLE_CONTINUE;
		if (0 == pcm_count)
			aud_stat = MSP_AUDIO_SAMPLE_FIRST;

		ret = QISRAudioWrite(session_id, (const void *)&p_pcm[pcm_count], len, aud_stat, &ep_stat, &rec_stat);
		if (MSP_SUCCESS != ret)
		{
			goto asr_exit;
		}
			
		pcm_count += (long)len;
		pcm_size  -= (long)len;
		
		if (MSP_EP_AFTER_SPEECH == ep_stat)
			break;
		usleep(200);
	}
	errcode = QISRAudioWrite(session_id, NULL, 0, MSP_AUDIO_SAMPLE_LAST, &ep_stat, &rec_stat);
	if (MSP_SUCCESS != errcode)
	{
		goto asr_exit;	
	}

	while (MSP_REC_STATUS_COMPLETE != rec_stat) 
	{
		const char *rslt = QISRGetResult(session_id, &rec_stat, 0, &errcode);
		if (MSP_SUCCESS != errcode)
		{
			goto asr_exit;
		}
		if (NULL != rslt)
		{
			unsigned int rslt_len = strlen(rslt);
			total_len += rslt_len;
			if (total_len >= BUFFER_SIZE)
			{
				goto asr_exit;
			}
			strncat(rec_result, rslt, rslt_len);
		}
		usleep(150);
	}
asr_exit:
	if (NULL != f_pcm)
	{
		fclose(f_pcm);
		f_pcm = NULL;
	}
	if (NULL != p_pcm)
	{	
		free(p_pcm);
		p_pcm = NULL;
	}

	QISRSessionEnd(session_id, hints);
return rec_result;
}