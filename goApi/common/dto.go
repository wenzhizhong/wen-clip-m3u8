package common

// #EXT-X-VERSION:3：HLS 协议版本号（常见为 3 或以上）。
// #EXT-X-TARGETDURATION:10：所有媒体分片的最大持续时间（秒）。
// #EXT-X-MEDIA-SEQUENCE:0：第一个分片的序列号（用于直播流）。
// #EXT-X-PLAYLIST-TYPE:VOD：播放列表类型（VOD 点播或 EVENT 直播事件）。
// #EXT-X-ENDLIST：表示播放列表结束（仅限点播，直播中无此标签）。
// #EXTINF: 10.000，：媒体分片的持续时间（秒）。
// #EXT-X-DISCONTINUITY：表示分片 discontinuity（分片 discontinuity）。
// #EXT-X-KEY:METHOD=AES-128,URI="pathxxx",IV=xxx ：定义解密密钥（若流被加密）。
type M3u8Info = struct {
	ExtVersion        int
	ExtTargetduration int
	ExtMediaSequence  int
	ExtPlaylistType   string
	ExtList           map[string][]ExtListItem
	ExtListLen        int
	// ExtKey            string
	// ExtKeyMethod      string
	// ExtKeyUri         string
	// ExtKeyIv          string
	// ExtKeyTrue        string
	// ExtKeyIvTrue      string
	HasExtDiscontinuity bool
}
type ExtListItem = struct {
	Path         string
	StartSec     int64
	StartTimeStr string
	ExtDuration  float64
	ExtKeyUri    string
	ExtKeyTrue   string
	ExtKeyIvTrue string
}

type AllPathDto = struct {
	SliceMp4Path      string
	M3u8Dir           string
	UniqueName        string
	M3u8VideoBasePath string
	MergeEndPath      string
	MergeDecPath      string
	M3u8VideoPathTpl  string
	CoverImagePathTpl string
}
