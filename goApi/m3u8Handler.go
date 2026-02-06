package goApi

import (
	"bytes"
	"clipM3u8Media/backend/apps/common/utils"
	"clipM3u8Media/goApi/common"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	sysRuntime "runtime"
	"strconv"
	"strings"
	"syscall"
)

var sliceMp4PathName string = "sliceMp4Path"
var resultMp4PathName string = "resultMp4Path"
var OptTypeVideo = "video"
var OptTypeCoverImg = "coverImg"

type M3u8Handler struct {
	Ctx *context.Context
}

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
	ExtList           []string
	ExtKey            string
	ExtKeyMethod      string
	ExtKeyUri         string
	ExtKeyIv          string
	ExtKeyTrue        string
	ExtKeyIvTrue      string
}

// 命令执行错误
type CommandError struct {
	Cmd        string
	ExitCode   int
	Stderr     string
	Underlying error
}

// 检查ffmpeg是否安装
func (a *M3u8Handler) CheckFfmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	// 设置 Windows 下不显示窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	return cmd.Run()
}

// 打开m3u8文件
func (a *M3u8Handler) OpenM3u8File(path string) (data interface{}, err error) {
	return a.doOpenM3u8File(path)
}

// 清空m3u8文件作业
func (a *M3u8Handler) ClearM3u8FileJob(path string) (data bool, err error) {
	return a.doClearM3u8FileJob(path)
}

// 合并每个已经生成m3u8任务文件
func (a *M3u8Handler) MergeM3u8File(path string, finalMergeFileList []string) (data interface{}, err error) {
	return a.doMergeM3u8File(path, finalMergeFileList)
}

// 删除当前作业数据源
func (a *M3u8Handler) DeleteM3u8Source(path string) (data interface{}, err error) {
	return a.doDeleteM3u8Source(path)
}

func (a *M3u8Handler) doOpenM3u8File(path string) (data interface{}, err error) {
	playPathList := make([]map[string]interface{}, 0)
	content, err := a.CheckM3u8File(path)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
		return data, err
	}
	m3u8Info, contentLines, err := a.ParseM3u8File(path, &content)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
		return data, err
	}

	playPathList, err = a.getM3u8SliceVideo(path, m3u8Info, &contentLines)
	// fmt.Println(playPathList, err)
	// fmt.Println(m3u8Info)
	data = struct {
		M3u8Info     M3u8Info
		PlayPathList []map[string]interface{}
	}{
		M3u8Info:     *m3u8Info,
		PlayPathList: playPathList,
	}
	return
}

func (a *M3u8Handler) doClearM3u8FileJob(path string) (result bool, err error) {
	_, err = a.CheckM3u8File(path)
	if err != nil {
		return result, err
	}

	tmpSliceMp4Path := a.getSliceMp4Path(path)
	err = common.RemoveByWildcard(tmpSliceMp4Path, "*.mp4")
	err = common.RemoveByWildcard(tmpSliceMp4Path, "*.jpg")
	if err != nil {
		return result, err
	}
	result = true

	return
}

func (a *M3u8Handler) doMergeM3u8File(path string, finalMergeFileList []string) (result interface{}, err error) {
	if len(finalMergeFileList) == 0 {
		return result, errors.New("没有可处理文件")
	}
	content := ""
	content, err = a.CheckM3u8File(path)
	if err != nil {
		return result, err
	}

	m3u8Dir := a.getM3u8Dir(path)
	resultMp4Dir := filepath.Join(m3u8Dir, resultMp4PathName)
	resultMp4FileName := a.getM3u8PathFileName(path) + ".mp4"
	resultMp4FileRelPath := filepath.Join(resultMp4PathName, resultMp4FileName)
	resultMp4FileAbsPath := filepath.Join(resultMp4Dir, resultMp4FileName)

	mergeFromFileRelPath := "newN3u8File.m3u8"
	mergeFromFileAbsPath := filepath.Join(m3u8Dir, mergeFromFileRelPath)

	_, err = os.Stat(resultMp4Dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(resultMp4Dir, os.ModePerm)
		if err != nil {
			return result, err
		}
	} else {
		os.RemoveAll(mergeFromFileAbsPath)
		os.RemoveAll(resultMp4FileAbsPath)
	}

	originVideoSize, err1 := a.getM3u8ContentSize(path)
	if err1 != nil {
		return result, err1
	}

	err = a.generateNewM3u8File(mergeFromFileAbsPath, &content, finalMergeFileList)
	if err != nil {
		return result, err
	}

	// ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp4
	// fmt.Println("执行命令：", "ffmpeg", "-f", "concat", "-safe", "0", "-i", mergeFromFileRelPath, "-c", "copy", resultMp4FileRelPath, "\n  ")
	// ffmpeg  -allowed_extensions ALL -i newN3u8File.m3u8 -c copy output.mp4
	fmt.Println("执行命令：", "ffmpeg", "-allowed_extensions", "ALL", "-i", mergeFromFileRelPath, "-c", "copy", resultMp4FileRelPath, "\n  ")
	cmd := exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", mergeFromFileRelPath, "-c", "copy", resultMp4FileRelPath)

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = m3u8Dir
	err = cmd.Run()
	if err != nil {
		return result, err
	}

	// resultMp4FileNameInfo, err1 := os.Stat(resultMp4FileAbsPath)
	// if err1 != nil {
	// 	return result, err1
	// }
	videoInfo := &common.VideoInfo{}
	videoInfo, err = common.GetVideoInfoJSON(resultMp4FileAbsPath)
	if err != nil {
		fmt.Println(err)
		// return result, err
	}
	fmt.Println(videoInfo)

	playPathList := []map[string]interface{}{
		{
			"path":  resultMp4FileRelPath,
			"name":  resultMp4FileName,
			"error": nil,
		},
	}
	result = struct {
		M3u8Info        M3u8Info
		PlayPathList    []map[string]interface{}
		MergePath       string
		M3u8Path        string
		Name            string
		VideoInfo       common.VideoInfo
		OriginVideoSize int64
	}{
		M3u8Info:        M3u8Info{},
		PlayPathList:    playPathList,
		MergePath:       resultMp4FileAbsPath,
		M3u8Path:        path,
		Name:            resultMp4FileName,
		VideoInfo:       *videoInfo,
		OriginVideoSize: originVideoSize,
	}
	return
}

// 删除作业数据源
func (a *M3u8Handler) doDeleteM3u8Source(path string) (result interface{}, err error) {
	result = struct {
		Code int
	}{
		Code: 1,
	}

	_, err = a.CheckM3u8File(path)
	if err != nil {
		return result, err
	}

	m3u8ContentDir := a.getM3u8ContentDir(path)
	err = os.RemoveAll(m3u8ContentDir)
	if err != nil {
		return result, err
	}
	err = os.RemoveAll(path)
	if err != nil {
		return result, err
	}
	result = struct {
		Code int
	}{
		Code: 0,
	}
	return
}

// 检测m3u8文件
func (a *M3u8Handler) CheckM3u8File(path string) (content string, err error) {
	if path == "" {
		return content, errors.New("请选择m3u8文件")
	}
	if !strings.HasSuffix(path, ".m3u8") {
		return content, errors.New("参数异常，不是m3u8文件路径")
	}
	fileInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return content, err
	}
	fmt.Println("文件" + fileInfo.Name())

	tmpContent, err := os.ReadFile(path)
	if err != nil {
		return content, err
	}
	if len := len(tmpContent); len == 0 {
		return content, errors.New("请选择m3u8文件")
	}
	return string(tmpContent), nil
}

// 解析m3u8文件
func (a *M3u8Handler) ParseM3u8File(path string, content *string) (m3u8Info *M3u8Info, contentLines []string, err error) {
	m3u8ContentDir := a.getM3u8ContentDir(path)
	contentLines = strings.Split(*content, "\n")
	if !strings.Contains(contentLines[0], "EXTM3U") {
		return m3u8Info, contentLines, errors.New("请选择m3u8文件")
	}
	m3u8Info = &M3u8Info{}

	beginVideoLine := false
	for i := 0; i < len(contentLines); i++ {
		line := contentLines[i]
		if !beginVideoLine {
			if strings.Contains(line, "VERSION") {
				m3u8Info.ExtVersion, _ = strconv.Atoi(strings.Split(line, ":")[1])

			} else if strings.Contains(line, "TARGETDURATION") {
				m3u8Info.ExtTargetduration, _ = strconv.Atoi(strings.Split(line, ":")[1])

			} else if strings.Contains(line, "MEDIA-SEQUENCE") {
				m3u8Info.ExtMediaSequence, _ = strconv.Atoi(strings.Split(line, ":")[1])

			} else if strings.Contains(line, "PLAYLIST-TYPE") {
				m3u8Info.ExtPlaylistType = line

			} else if strings.Contains(line, "X-KEY") {
				m3u8Info.ExtKey = line

				tmpArr := strings.Split(line, ",")
				for _, v := range tmpArr {
					tmpArr2 := strings.Split(v, "=")
					if len(tmpArr2) < 2 {
						continue
					}
					if strings.Contains(v, "URI") {
						m3u8Info.ExtKeyUri = strings.Trim(tmpArr2[1], "\r\n\"")
					}
					if strings.Contains(v, "METHOD") {
						m3u8Info.ExtKeyMethod = strings.Trim(tmpArr2[1], "\r\n\"")
					}
					if strings.Contains(v, "IV") {
						m3u8Info.ExtKeyIv = strings.Trim(tmpArr2[1], "\r\n\"")
						if strings.HasPrefix(m3u8Info.ExtKeyIv, "0x") {
							m3u8Info.ExtKeyIv = m3u8Info.ExtKeyIv[2:]
						}
					}
				}
			}
			if strings.Contains(line, "EXTINF") {
				beginVideoLine = true
			}
			if !beginVideoLine {
				continue
			}
		}

		if strings.Contains(line, "EXTINF") {
			i++
			nextLine := strings.Trim(contentLines[i], "\r\n")
			contentLines[i] = nextLine + ".ts" + "\n" // 添加默认后缀.ts
			nextLineSplit := strings.Split(nextLine, "/")
			sliceFileName := nextLineSplit[len(nextLineSplit)-1]
			if strings.Contains(nextLine, ".ts") {
				m3u8Info.ExtList = append(m3u8Info.ExtList, nextLine)
			} else {
				nextLine = nextLine + ".ts"
				tmpPath := filepath.Join(m3u8ContentDir, sliceFileName)
				oldPath := tmpPath
				newPath := tmpPath + ".ts"

				// 添加物理文件默认后缀
				_, err1 := os.Stat(strings.Trim(oldPath, "\""))
				_, err2 := os.Stat(strings.Trim(newPath, "\""))
				if err1 != nil && err2 != nil {
					err = err1
					fmt.Println("没有找到分片：" + sliceFileName + "\n" + oldPath)
					fmt.Println(err1)
					return
				} else {
					if err2 == nil {
						m3u8Info.ExtList = append(m3u8Info.ExtList, nextLine)
					} else {
						var err3 error
						var cmd *exec.Cmd
						switch sysRuntime.GOOS {
						case "windows":
							cmd = exec.Command("cmd", "/C", "move", oldPath, newPath)
							cmd.SysProcAttr = &syscall.SysProcAttr{
								HideWindow: true,
							}
						default:
							cmd = exec.Command("mv", oldPath, newPath)
						}
						err3 = cmd.Run()

						if err3 != nil {
							fmt.Println("复制分片失败：" + sliceFileName + "\n" + oldPath + "=>" + newPath + " \n")
							fmt.Println(err3)
						} else {
							m3u8Info.ExtList = append(m3u8Info.ExtList, nextLine)
						}
					}
				}
			}
		}
	}
	// jsonStr, err1 := json.Marshal(m3u8Info)
	// contentLinesStr, err2 := json.Marshal(contentLines)
	// fmt.Println("\n", string(jsonStr), err1)
	// fmt.Println("\n", string(contentLinesStr), err2)

	if len(m3u8Info.ExtList) == 0 {
		return m3u8Info, contentLines, errors.New("m3u8视频分片不存在")
	}

	m3u8Dir := a.getM3u8Dir(path)
	keyData, _ := ioutil.ReadFile(filepath.Join(m3u8Dir, m3u8Info.ExtKeyUri))
	m3u8Info.ExtKeyTrue = hex.EncodeToString(keyData)
	m3u8Info.ExtKeyIvTrue = m3u8Info.ExtKeyIv

	return
}
func (a *M3u8Handler) getM3u8Dir(path string) string {
	return filepath.Dir(path)
}

// 获取每个m3u8分片视频列表

func (a *M3u8Handler) getM3u8SliceVideo(path string, m3u8Info *M3u8Info, contentLines *[]string) (playPathList []map[string]interface{}, err error) {
	playPathList = make([]map[string]interface{}, 0)
	tmpPlayPathMap := make(map[string]map[string]interface{})
	listSlice := m3u8Info.ExtList
	listSliceChunk := utils.ArrayChunk(listSlice, 50)
	listSliceChunkLen := len(listSliceChunk)

	tmpSliceMp4Path := a.getSliceMp4Path(path)
	if _, err1 := os.Stat(tmpSliceMp4Path); os.IsNotExist(err1) {
		os.MkdirAll(tmpSliceMp4Path, os.ModePerm)
	} else {
		// RemoveByWildcard(tmpSliceMp4Path, "*.mp4")
		// fmt.Println("删除文件：" + tmpSliceMp4Path + "/*.mp4")
	}

	type M3u8SliceVideo = struct {
		SliceVideo []map[string]interface{}
		Error      error
	}
	ch := make(chan []M3u8SliceVideo, listSliceChunkLen)
	for i := 0; i < listSliceChunkLen; i++ {
		go func() {
			tmpPlayPathList, err := a.DoGetM3u8SliceVideo(path, m3u8Info, listSliceChunk[i], OptTypeCoverImg)
			ch <- []M3u8SliceVideo{{SliceVideo: tmpPlayPathList, Error: err}}
		}()
	}

	for i := 0; i < listSliceChunkLen; i++ {
		tmpM3u8SliceVideo := <-ch
		sliceVideoList := tmpM3u8SliceVideo[0].SliceVideo
		for j := 0; j < len(sliceVideoList); j++ {
			item := sliceVideoList[j]
			sliceIndex := item["index"].(string)
			tmpPlayPathMap[sliceIndex] = item
		}
	}
	// tmpPlayPathName = utils.ArraySort(tmpPlayPathName, 1)
	for i := 0; i < len(listSlice); i++ {
		item := tmpPlayPathMap[fmt.Sprint(i)]
		if item == nil {
			continue
		}
		playPathList = append(playPathList, item)
	}

	fmt.Println("playPathList========\n", playPathList)
	fmt.Println("listSliceChunkLen========\n", listSliceChunkLen)
	return
}

func (a *M3u8Handler) DoGetM3u8SliceVideo(path string, m3u8Info *M3u8Info, listSlice []string, optType string) (playPathList []map[string]interface{}, err error) {
	playPathList = make([]map[string]interface{}, 0)
	m3u8Dir := a.getM3u8Dir(path)
	uniqueName := a.getM3u8PathMd5(path)

	if optType != OptTypeVideo && optType != OptTypeCoverImg {
		return playPathList, errors.New("optType参数错误, 可选值：video, coverImg")
	}
	for i := 0; i < len(listSlice); i++ {
		sliceNameArr := strings.Split(listSlice[i], "/")
		sliceNameArrLen := len(sliceNameArr)
		re := regexp.MustCompile(`[0-9]+`)
		sliceIndex := re.FindString(sliceNameArr[sliceNameArrLen-1])
		sliceName := sliceNameArr[sliceNameArrLen-1] + ".mp4"
		// m3u8VideoPath := sliceMp4PathName + "/" + sliceName
		m3u8VideoPath := filepath.Join(sliceMp4PathName, uniqueName, sliceName)

		// 定义封面图路径
		coverImagePath := m3u8VideoPath + ".jpg"

		playPathListItem := map[string]interface{}{
			"index":      sliceIndex,
			"name":       sliceName,
			"path":       m3u8VideoPath,
			"cover_path": coverImagePath, // 添加封面图路径
			"error":      nil,
		}

		videoExists := false
		coverExists := false

		if _, err := os.Stat(filepath.Join(m3u8Dir, m3u8VideoPath)); !os.IsNotExist(err) {
			videoExists = true
		}

		if _, err := os.Stat(filepath.Join(m3u8Dir, coverImagePath)); !os.IsNotExist(err) {
			coverExists = true
		}

		// 如果视频和封面都存在，则跳过
		if videoExists && coverExists {
			playPathList = append(playPathList, playPathListItem)
			continue
		}

		if optType == OptTypeVideo {
			// err = exec.Command("ffmpeg", params...).Run() //  ffmpeg -allowed_extensions ALL -i "file:index.m3u8" -c copy output.mp4
			// err = exec.Command("ffmpeg", params2...).Run() // ffmpeg -decryption_key YOUR_KEY_HEX -decryption_iv YOUR_IV_HEX -i "crypto+file:index.m3u8_contents/0" -c copy segment_0_decrypted.ts
			// ffmpegStr = "ffmpeg -decryption_key f7fd2cdfb2429a9646cb69234bebc9b3 -decryption_iv 1ef58f5c956b146218c8035d387f2728 -i \"crypto+file:index.m3u8_contents/0.ts\" -c copy \"sliceMp4PathName/0.ts.mp4\""
			// cmd := exec.Command("cmd", "/C", ffmpegStr)

			cmd := exec.Command("ffmpeg",
				"-decryption_key", m3u8Info.ExtKeyTrue,
				"-decryption_iv", m3u8Info.ExtKeyIvTrue,
				"-i", "crypto+file:"+listSlice[i],
				"-c", "copy",
				m3u8VideoPath,
			)

			// 设置 Windows 下不显示窗口
			if runtime.GOOS == "windows" {
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true,
				}
			}
			// 创建缓冲区用于捕获stderr
			// 创建多重写入器：同时写入缓冲区和终端（便于实时查看）
			var stderrBuf bytes.Buffer
			var stderrWriter io.Writer

			// 检查 stderr 是否可用
			if _, err := os.Stderr.Stat(); err == nil {
				// stderr 可用，创建多重写入器
				stderrWriter = io.MultiWriter(&stderrBuf, os.Stderr)
			} else {
				// stderr 不可用，只写入缓冲区
				stderrWriter = &stderrBuf
			}
			cmd.Stderr = stderrWriter
			cmd.Stdout = os.Stdout // 标准输出通常直接输出到终端
			cmd.Dir = m3u8Dir
			err = cmd.Run()
			if err != nil {
				cmdErr := &CommandError{
					Cmd:        cmd.String(),
					Stderr:     stderrBuf.String(),
					Underlying: err,
				}
				// 尝试获取退出码, 非退出错误（如命令未找到）
				if exitErr, ok := err.(*exec.ExitError); ok {
					cmdErr.ExitCode = exitErr.ExitCode()
				} else {
					cmdErr.ExitCode = -1
				}
				playPathListItem["error"] = cmdErr

				_, file, line, _ := runtime.Caller(0)
				common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
				playPathList = append(playPathList, playPathListItem)
				return
			}
		}

		// 提取封面图
		if optType == OptTypeCoverImg {

			coverCmd := exec.Command("ffmpeg",
				"-decryption_key", m3u8Info.ExtKeyTrue,
				"-decryption_iv", m3u8Info.ExtKeyIvTrue,
				"-i", "crypto+file:"+listSlice[i],
				// "-vf", "thumbnail,scale=640:-1", // 使用thumbnail过滤器提取关键帧，并缩放到宽度640
				// "-vf", "scale=640:-1",  // 只做缩放，去掉thumbnail过滤器
				"-vframes", "1", // 只提取一帧
				"-an", // 不处理音频
				"-sn", // 不处理字幕
				// "-q:v", "5", // 降低质量要求以提高速度
				"-f", "image2", // 图像输出格式

				// "-fast", "1", // 添加快速解码参数
				// "-fflags", "+fastseek", // 添加快速解码参数
				// "-map_metadata", "-1", // 跳过元数据处理

				"-probesize", "32", // 减少探测数据
				"-analyzeduration", "0", // 快速分析
				"-avoid_negative_ts", "make_zero",
				"-fflags", "+fastseek", // 快速seek
				"-y", // 覆盖输出文件
				coverImagePath,
			)
			if runtime.GOOS == "windows" {
				coverCmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true,
				}
			}

			coverCmd.Dir = m3u8Dir
			coverErr := coverCmd.Run()
			if coverErr != nil {
				fmt.Printf("提取封面图失败: %v\n", coverErr)
				playPathListItem["cover_error"] = coverErr.Error() // 记录封面提取错误
			}
		}

		playPathList = append(playPathList, playPathListItem)
		// fmt.Println("m3u8VideoPath=", m3u8VideoPath)
		// fmt.Println("m3u8Dir=", m3u8Dir)
	}
	return playPathList, nil
}

func (a *M3u8Handler) generateNewM3u8File(newPath string, content *string, finalMergeFileList []string) error {
	finalMergeFileListLen := len(finalMergeFileList)
	finalMergeFileMap := make(map[string]string)
	for i := 0; i < len(finalMergeFileList); i++ {
		item := finalMergeFileList[i]
		item = strings.ReplaceAll(item, "\\", "/")
		sliceName := strings.ReplaceAll(filepath.Base(item), ".mp4", "")
		finalMergeFileMap[sliceName] = ""
	}
	fmt.Println("	==", finalMergeFileMap)

	contentLines := strings.Split(*content, "\n")
	newContentLines := make([]string, 0)
	for i := 0; i < len(contentLines); i++ {
		line := contentLines[i]
		if strings.Contains(line, "EXTINF") {
			i++
			nextLine := strings.Trim(contentLines[i], "\r\n")
			nextLine = nextLine + ".ts"
			sliceName := filepath.Base(nextLine)
			if finalMergeFileListLen > 0 {
				if _, ok := finalMergeFileMap[sliceName]; ok {
					newContentLines = append(newContentLines, line)
					newContentLines = append(newContentLines, nextLine)
				}
			} else {
				newContentLines = append(newContentLines, line)
				newContentLines = append(newContentLines, nextLine)
			}
		} else {
			newContentLines = append(newContentLines, line)
		}
	}
	newContent := strings.Join(newContentLines, "\n")
	err := os.WriteFile(newPath, []byte(newContent), os.ModePerm)
	return err
}
func (a *M3u8Handler) getM3u8ContentSize(path string) (int64, error) {
	dir := a.getM3u8ContentDir(path)
	return common.GetSize(dir)
}
func (a *M3u8Handler) getM3u8ContentDir(path string) string {
	path = path[:len(path)-5]
	return path + ".m3u8_contents"
}

func (a *M3u8Handler) getSliceMp4Path(path string) string {
	m3u8Dir := a.getM3u8Dir(path)
	uniqueName := a.getM3u8PathMd5(path)

	return filepath.Join(m3u8Dir, sliceMp4PathName, uniqueName)
}
func (a *M3u8Handler) getM3u8PathMd5(path string) string {
	uniqueName := md5.Sum([]byte(a.getM3u8PathFileName(path)))
	return hex.EncodeToString(uniqueName[:])
}

func (a *M3u8Handler) getM3u8PathFileName(path string) string {
	return filepath.Base(path)[:len(filepath.Base(path))-5]
}
