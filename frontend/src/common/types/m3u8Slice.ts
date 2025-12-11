
export interface uploadM3u8Interface {
	M3u8Dir: 			string // m3u8文件所在目录
	M3u8Path: 		string // m3u8文件路径
  M3u8Info:     ParseM3u8SliceInterface,
  PlayPathList: PlayPathListInterface[]
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
  error:any
}
