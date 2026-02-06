
export interface uploadM3u8Interface {
	M3u8Dir: 			string // m3u8文件所在目录
	M3u8Path: 		string // m3u8文件路径
  M3u8Info:     ParseM3u8SliceInterface,
  PlayPathList: PlayPathListInterface[],
	zoomNumber ?: number
}
interface VideoInfoFormatInterface{
	duration: string,
	bit_rate: string,
	size: string
}
// streams
interface VideoInfoStreamsInterface{
	codec_type: string,
	codec_name: string,
	width: number,
	height: number,
	duration: string,
	bit_rate: string,
	r_frame_rate: string
}

export interface VideoInfoInterface{
	format: VideoInfoFormatInterface,
	streams: VideoInfoStreamsInterface[],
}


export interface mergeSuccessInterface {
	M3u8Dir: 			string // m3u8文件所在目录
	M3u8Path: 		string // m3u8文件路径
	MergePath: 		string // 合并后路径
	Name: 				string // 剧名
  M3u8Info:     ParseM3u8SliceInterface,
  PlayPathList: PlayPathListInterface[],
	VideoInfo: VideoInfoInterface
	OriginVideoSize: Number
}

export interface ParseM3u8SliceInterface {
	ExtVersion       : Number
	ExtTargetduration: Number
	ExtMediaSequence : Number
	ExtPlaylistType  : string
	ExtList          : string[]
	ExtKey           : string
	ExtKeyMethod     : string
	ExtKeyUri        : string
	ExtKeyIv         : string
}

export interface PlayPathListInterface{
  name:string,
  path:string,
  cover_path:string,
  error:any
}
