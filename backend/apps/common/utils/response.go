package utils

import "clipM3u8Media/backend/apps/common/dto"

func Response(code int, msg string, data interface{}) dto.ResponseDto {
	return dto.ResponseDto{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
