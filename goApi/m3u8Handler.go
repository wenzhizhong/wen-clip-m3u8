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
	"log"
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

var sliceMp4PathName string = common.WorkPathName + "sliceMp4Path"
var resultMp4PathName string = common.WorkPathName + "resultMp4Path"
var OptTypeVideo = "video"
var OptTypeCoverImg = "coverImg"

type M3u8Handler struct {
	Ctx *context.Context
}

// 命令执行错误
type CommandError struct {
	Cmd        string
	ExitCode   int
	Stderr     string
	Underlying error
}

// 检查环境
func (a *M3u8Handler) CheckEnv() error {
	return a.doCheckEnv()
}

// 打开m3u8文件
func (a *M3u8Handler) OpenM3u8File(path string) (data interface{}, err error) {
	return a.doOpenM3u8File(path)
}

// 清空m3u8文件作业
func (a *M3u8Handler) ClearM3u8FileJob(path string) (data bool, err error) {
	path = strings.ReplaceAll(path, common.ReCutNamePlaceholder, "")
	pathDto := a.GetGetAllPathDto(path)
	paths := []string{path, pathDto.ReCutM3u8Path}
	return a.doClearM3u8FileJob(paths)
}

// 合并每个已经生成m3u8任务文件
func (a *M3u8Handler) MergeM3u8File(path string, finalMergeFileList []string) (data interface{}, err error) {
	return a.doMergeM3u8File(path, finalMergeFileList)
}

// 删除当前作业数据源
func (a *M3u8Handler) DeleteM3u8Source(path string) (data interface{}, err error) {
	return a.doDeleteM3u8Source(path)
}

// 重新分片
func (a *M3u8Handler) ReCut(path string) (data interface{}, err error) {
	return a.doReCut(path)
}
func (a *M3u8Handler) doCheckEnv() error {
	errs := []string{}
	tmpMap := map[string][]string{
		"ffmpeg":   {"-version"},
		"openssl1": {"-version"},
	}
	for k, v := range tmpMap {
		err := a.doCheckEnvFunc(k, v...)
		if err != nil {
			errs = append(errs, k)
		}
	}

	if len(errs) > 0 {
		return errors.New("请先安装：" + strings.Join(errs, ","))
	}
	return nil
}
func (a *M3u8Handler) doCheckEnvFunc(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// 设置 Windows 下不显示窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	return cmd.Run()
}
func (a *M3u8Handler) doReCut(path string) (data interface{}, err error) {
	content, err := a.CheckM3u8File(path)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
		return data, err
	}
	m3u8Info, _, err := a.ParseM3u8File(path, &content)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
		return data, err
	}

	// 1
	pathDto := a.GetGetAllPathDto(path)
	mergeResultPathOfReCut := pathDto.MergeResultPathOfReCut
	var stderr bytes.Buffer
	reCutM3u8ContentPath := a.getM3u8ContentDir(pathDto.ReCutM3u8Path)
	reCutM3u8ContentDirName := filepath.Base(reCutM3u8ContentPath)
	mergeFromFileAbsPath := strings.ReplaceAll(pathDto.MergeFromFileAbsPath, ".m3u8", common.ReCutNamePlaceholder+".m3u8")

	mergeFileList := make([]string, 0)
	for _, listItems := range m3u8Info.ExtList {
		for _, listItem := range listItems {
			mergeFileList = append(mergeFileList, listItem.Path)
		}
	}
	err = a.generateNewM3u8File(mergeFromFileAbsPath, &content, mergeFileList)
	if err != nil {
		return data, err
	}

	args := []string{"-allowed_extensions", "ALL", "-i", mergeFromFileAbsPath, "-c", "copy", "-y", mergeResultPathOfReCut}
	fmt.Println("ffmpeg ", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return data, err
	}

	// 2
	err = os.MkdirAll(reCutM3u8ContentPath, os.ModePerm)
	if err != nil {
		fmt.Println("创建目录失败：", err)
		return data, err
	}

	fmt.Println("ffmpeg", "-i", mergeResultPathOfReCut, "-c:v", "copy", "-c:a", "copy", "-f", "hls", "-hls_time", "5", "-hls_list_size", "0", "-hls_segment_filename", strconv.Quote(filepath.Join(reCutM3u8ContentPath, "%d.ts")), mergeFromFileAbsPath)
	cmd2 := exec.Command("ffmpeg", "-i", mergeResultPathOfReCut, "-c:v", "copy", "-c:a", "copy", "-f", "hls", "-hls_time", "5", "-hls_list_size", "0", "-hls_segment_filename", (filepath.Join(reCutM3u8ContentPath, "%d.ts")), mergeFromFileAbsPath)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
	cmd2.Stderr = &stderr
	if err := cmd2.Run(); err != nil {
		fmt.Println("重新分片失败：", err, string(stderr.Bytes()))
		return data, err
	}

	// 3 read mergeFromFileAbsPath
	contentLinesBytes, err := ioutil.ReadFile(mergeFromFileAbsPath)
	if err != nil {
		fmt.Println("读取m3u8文件失败：", err)
		return data, err
	}
	contentLinesStr := string(contentLinesBytes)
	contentLines := strings.Split(contentLinesStr, "\n")
	for i, line := range contentLines {
		if !strings.Contains(line, "EXTINF") {
			continue
		}
		i++
		nextLine := reCutM3u8ContentDirName + "/" + strings.Trim(contentLines[i], "\r\n")
		contentLines[i] = nextLine
	}
	contentLinesStr = strings.Join(contentLines, "\n")
	err = ioutil.WriteFile(pathDto.ReCutM3u8Path, []byte(contentLinesStr), os.ModePerm)
	if err != nil {
		fmt.Println("重新写入m3u8文件失败：", err)
		return data, err
	}

	return map[string]string{"Path": pathDto.ReCutM3u8Path}, err
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
	if err != nil {
		common.LogToFile(path, err.Error())
		return
	}
	// fmt.Println(playPathList, err)
	// fmt.Println(m3u8Info)
	data = struct {
		M3u8Info     common.M3u8Info
		PlayPathList []map[string]interface{}
	}{
		M3u8Info:     *m3u8Info,
		PlayPathList: playPathList,
	}
	return
}

func (a *M3u8Handler) doClearM3u8FileJob(paths []string) (result bool, err error) {
	errStr := ""
	for _, p := range paths {
		if _, ok := os.Stat(p); os.IsNotExist(ok) {
			continue
		}
		tmpSliceMp4Path := a.getSliceMp4Path(p)
		fmt.Println("tmpSliceMp4Path=" + tmpSliceMp4Path)
		err1 := common.RemoveByWildcard(tmpSliceMp4Path, "*.ts")
		err2 := common.RemoveByWildcard(tmpSliceMp4Path, "*.mp4")
		err3 := common.RemoveByWildcard(tmpSliceMp4Path, "*.jpg")
		err4 := common.RemoveByWildcard(tmpSliceMp4Path, "*.key")

		tmpErr := ""
		if err1 != nil {
			tmpErr += err1.Error() + "\n"
		}
		if err2 != nil {
			tmpErr += err2.Error() + "\n"
		}
		if err3 != nil {
			tmpErr += err3.Error() + "\n"
		}
		if err4 != nil {
			tmpErr += err4.Error() + "\n"
		}
		if tmpErr != "" {
			errStr += tmpErr + "\n"
			common.LogToFile(p, tmpErr)
		}
	}
	if errStr != "" {
		err = errors.New(errStr)
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
	pathDto := a.GetGetAllPathDto(path)

	_, err = os.Stat(pathDto.ResultMp4Dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pathDto.ResultMp4Dir, os.ModePerm)
		if err != nil {
			return result, err
		}
	} else {
		os.RemoveAll(pathDto.MergeFromFileAbsPath)
		os.RemoveAll(pathDto.ResultMp4FileAbsPath)
	}

	originVideoSize, err1 := a.getM3u8ContentSize(path)
	if err1 != nil {
		return result, err1
	}

	err = a.generateNewM3u8File(pathDto.MergeFromFileAbsPath, &content, finalMergeFileList)
	if err != nil {
		return result, err
	}

	// ffmpeg -f concat -safe 0 -i list.txt -c copy output.mp4
	// fmt.Println("执行命令：", "ffmpeg", "-f", "concat", "-safe", "0", "-i", pathDto.MergeFromFileRelPath, "-c", "copy", pathDto.ResultMp4FileRelPath, "\n  ")
	// ffmpeg  -allowed_extensions ALL -i newN3u8File.m3u8 -c copy output.mp4
	fmt.Println("执行命令：", "ffmpeg", "-allowed_extensions", "ALL", "-i", pathDto.MergeFromFileRelPath, "-c", "copy", pathDto.ResultMp4FileRelPath, "\n  ")
	cmd := exec.Command("ffmpeg", "-allowed_extensions", "ALL", "-i", pathDto.MergeFromFileRelPath, "-c", "copy", pathDto.ResultMp4FileRelPath)

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = pathDto.M3u8Dir
	err = cmd.Run()
	if err != nil {
		return result, err
	}

	// resultMp4FileNameInfo, err1 := os.Stat(pathDto.ResultMp4FileAbsPath)
	// if err1 != nil {
	// 	return result, err1
	// }
	videoInfo := &common.VideoInfo{}
	videoInfo, err = common.GetVideoInfoJSON(pathDto.ResultMp4FileAbsPath)
	if err != nil {
		fmt.Println(err)
		// return result, err
	}
	fmt.Println(videoInfo)

	playPathList := []map[string]interface{}{
		{
			"path":  pathDto.ResultMp4FileRelPath,
			"name":  pathDto.ResultMp4FileName,
			"error": nil,
		},
	}
	result = struct {
		M3u8Info        common.M3u8Info
		PlayPathList    []map[string]interface{}
		MergePath       string
		M3u8Path        string
		Name            string
		VideoInfo       common.VideoInfo
		OriginVideoSize int64
	}{
		M3u8Info:        common.M3u8Info{},
		PlayPathList:    playPathList,
		MergePath:       pathDto.ResultMp4FileAbsPath,
		M3u8Path:        path,
		Name:            pathDto.ResultMp4FileName,
		VideoInfo:       *videoInfo,
		OriginVideoSize: originVideoSize,
	}
	return
}

// 删除作业数据源
func (a *M3u8Handler) doDeleteM3u8Source(path string) (result interface{}, err error) {
	errStr := ""
	path = strings.ReplaceAll(path, common.ReCutNamePlaceholder, "")
	pathDto := a.GetGetAllPathDto(path)
	paths := []string{path, pathDto.ReCutM3u8Path}

	result = struct {
		Code int
	}{
		Code: 1,
	}
	for _, p := range paths {
		if _, ok := os.Stat(p); os.IsNotExist(ok) {
			continue
		}

		_, err = a.CheckM3u8File(p)
		if err != nil {
			errStr += err.Error() + "\n"
			common.LogToFile(p, errStr)
			continue
		}

		m3u8ContentDir := a.getM3u8ContentDir(p)
		err = os.RemoveAll(m3u8ContentDir)
		if err != nil {
			errStr += err.Error() + "\n"
			common.LogToFile(p, errStr)
			continue
		}
		err = os.RemoveAll(p)
		if err != nil {
			errStr += err.Error() + "\n"
			common.LogToFile(p, errStr)
			continue
		}
	}
	if errStr != "" {
		err = errors.New(errStr)
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
func (a *M3u8Handler) ParseM3u8File(path string, content *string) (m3u8Info *common.M3u8Info, contentLines []string, err error) {
	m3u8ContentDir := a.getM3u8ContentDir(path)
	contentLines = strings.Split(*content, "\n")
	if !strings.Contains(contentLines[0], "EXTM3U") {
		return m3u8Info, contentLines, errors.New("请选择m3u8文件")
	}
	m3u8Info = &common.M3u8Info{}

	tmpKey := ""
	tmpKeyIv := ""
	tmpKeyUri := ""
	tmpKeyMethod := ""
	tmpKeyUriExists := []string{}
	extListMapKey := common.M3u8InfoConstant.ListMapDefKey //  // md5 || "none"
	extListMapKeyNumber := 0
	beginVideoLine := false
	var startSec int64 = 0
	m3u8Info.ExtList = make(map[string][]common.ExtListItem)
	m3u8Info.ExtListLen = 0
	for i := 0; i < len(contentLines); i++ {
		line := contentLines[i]
		if strings.Contains(line, "X-DISCONTINUITY") {
			tmpNextLine := strings.Trim(contentLines[i+1], "\r\n")
			if strings.Contains(tmpNextLine, "EXTINF") && !m3u8Info.HasExtDiscontinuity {
				m3u8Info.HasExtDiscontinuity = true
			}
		}
		if strings.Contains(line, "X-KEY") {
			if beginVideoLine {
				startSec = 0
				extListMapKeyNumber++
				beginVideoLine = false
			}

			if strings.Contains(line, "URI") && strings.Contains(line, "IV") {
				a.createKeyFileIfNotExist(path, line, tmpKeyUriExists)

				tmpKeyUri = line
				tmpKeyUriExists = append(tmpKeyUriExists, tmpKeyUri)
			}
		}
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
				tmpKey, tmpKeyIv, tmpKeyMethod = a.getKeyAndVi(path, line)
				extListMapKey = tmpKey // md5 || none
				if extListMapKey != "" {
					tmpExtListMapKey := md5.Sum([]byte(extListMapKey))
					extListMapKey = hex.EncodeToString(tmpExtListMapKey[:])
				}
				if strings.Contains(line, "NONE") || tmpKey == "" {
					extListMapKey = common.M3u8InfoConstant.ListMapDefKey
				}
			}
			if strings.Contains(line, "EXTINF") {
				beginVideoLine = true
			}
			if !beginVideoLine {
				continue
			}
		}
		tmpExtListMapKey := extListMapKey + "_" + strconv.Itoa(extListMapKeyNumber)
		if _, ok := m3u8Info.ExtList[tmpExtListMapKey]; !ok {
			m3u8Info.ExtList[tmpExtListMapKey] = make([]common.ExtListItem, 0)
		}

		if strings.Contains(line, "EXTINF") {
			m3u8Info.ExtListLen++
			duration, _ := a.getDurationFromExtInf(line)

			i++
			nextLine := strings.Trim(contentLines[i], "\r\n")
			contentLines[i] = nextLine + ".ts" + "\n" // 添加默认后缀.ts
			nextLineSplit := strings.Split(nextLine, "/")
			sliceFileName := nextLineSplit[len(nextLineSplit)-1]

			listItem := common.ExtListItem{
				Path:         nextLine,
				StartSec:     startSec,
				StartTimeStr: utils.MicrosecondToTime(startSec),
				ExtDuration:  duration,
				ExtKeyUri:    tmpKeyUri,
				ExtKeyTrue:   tmpKey,
				ExtKeyIvTrue: tmpKeyIv,
				ExtKeyMethod: tmpKeyMethod,
			}
			startSec += int64(listItem.ExtDuration * 1_000_000)
			if strings.Contains(nextLine, ".ts") {
				m3u8Info.ExtList[tmpExtListMapKey] = append(m3u8Info.ExtList[tmpExtListMapKey], listItem)
			} else {
				nextLine = nextLine + ".ts"
				listItem.Path = nextLine
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
						m3u8Info.ExtList[tmpExtListMapKey] = append(m3u8Info.ExtList[tmpExtListMapKey], listItem)
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
							m3u8Info.ExtList[tmpExtListMapKey] = append(m3u8Info.ExtList[tmpExtListMapKey], listItem)
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

	if m3u8Info.ExtListLen == 0 {
		return m3u8Info, contentLines, errors.New("m3u8视频分片不存在")
	}

	return
}
func (a *M3u8Handler) getM3u8Dir(path string) string {
	return filepath.Dir(path)
}

// 获取每个m3u8分片视频列表

func (a *M3u8Handler) getM3u8SliceVideo(path string, m3u8Info *common.M3u8Info, contentLines *[]string) (playPathList []map[string]interface{}, err error) {
	playPathList = make([]map[string]interface{}, 0)
	tmpPlayPathMap := make(map[string]map[string]interface{})
	extList := m3u8Info.ExtList
	listSliceLen := m3u8Info.ExtListLen

	for listMapKey, listSlice := range extList {
		listSliceChunk := utils.ArrayChunk(listSlice, 50)
		listSliceChunkLen := len(listSliceChunk)
		pathDto := a.GetGetAllPathDto(path, listMapKey)
		tmpSliceMp4Path := pathDto.SliceMp4Path
		if _, err1 := os.Stat(tmpSliceMp4Path); os.IsNotExist(err1) {
			os.MkdirAll(tmpSliceMp4Path, os.ModePerm)
		} else {
			// RemoveByWildcard(tmpSliceMp4Path, "*.mp4")
			// fmt.Println("删除文件：" + tmpSliceMp4Path + "/*.mp4")
		}

		if !strings.HasPrefix(listMapKey, common.M3u8InfoConstant.ListMapDefKey) { // "none"
			//整体合并
			err = a.mergeEncryptedSegments(pathDto.M3u8Dir, listSlice, pathDto.MergeEndPath, 3)
			if err != nil {
				return
			}
			// // isVideoPlayable() , 如果合并后正常播放，则不进行解密
			// if a.isVideoPlayable(pathDto.MergeEndPath) {
			// 	fmt.Println("合并后正常播放，不进行解密=====")
			// 	err = os.Rename(pathDto.MergeEndPath, pathDto.MergeDecPath)
			// 	if err != nil {
			// 		return
			// 	}
			// } else {
			// 	//整体解密合并文件
			// 	err1, stdErr := a.mergeDecryptedSegments(listSlice[0], pathDto.MergeEndPath, pathDto.MergeDecPath)
			// 	if err1 != nil {
			// 		err = err1
			// 		common.LogToFile(path, fmt.Sprintf("合并解密失败：%v\n%v\n", err, stdErr.String()))
			// 		return
			// 	}
			// }

			//整体解密合并文件
			err1, stdErr := a.mergeDecryptedSegments(listSlice[0], pathDto.MergeEndPath, pathDto.MergeDecPath)
			if err1 != nil {
				if a.isVideoPlayable(pathDto.MergeEndPath) {
					err = os.Rename(pathDto.MergeEndPath, pathDto.MergeDecPath)
					if err != nil {
						return
					}
				} else {
					err = err1
					common.LogToFile(path, fmt.Sprintf("合并解密失败：%v\n%v\n", err, stdErr.String()))
					return
				}
			}
		} else {
			pathDto.MergeEndPath = pathDto.MergeDecPath
			//整体合并
			err = a.mergeEncryptedSegments(pathDto.M3u8Dir, listSlice, pathDto.MergeEndPath, 3)
			if err != nil {
				return
			}
		}

		if m3u8Info.HasExtDiscontinuity {
			fixTimePath := pathDto.MergeDecPath + ".fix" + filepath.Ext(pathDto.MergeDecPath)
			tmpArgs := []string{"-i", pathDto.MergeDecPath, "-c", "copy", "-fflags", "+genpts", "-y", fixTimePath}
			cmd := exec.Command("ffmpeg", tmpArgs...)
			err = cmd.Run()
			fmt.Println("重新生成时间戳:\nffmpeg ", strings.Join(tmpArgs, " "))
			if err == nil && a.isVideoPlayable(fixTimePath) {
				err = os.Rename(fixTimePath, pathDto.MergeDecPath)
				if err != nil {
					return
				}

			}
		}

		type M3u8SliceVideo = struct {
			SliceVideo []map[string]interface{}
			Error      error
		}
		ch := make(chan []M3u8SliceVideo, listSliceChunkLen)
		for i := 0; i < listSliceChunkLen; i++ {
			go func() {
				// tmpPlayPathList, err := a.DoGetM3u8SliceVideo(path, m3u8Info, listSliceChunk[i], OptTypeCoverImg)
				tmpPlayPathList, err := a.DoGetM3u8SliceVideoV2(path, pathDto, listSliceChunk[i], OptTypeCoverImg)

				if len(tmpPlayPathList) != len(listSliceChunk[i]) {
					tmpPlayPathListIndex := make([]string, 0)
					tmpListSliceChunkIndex := make([]string, 0)
					for _, item := range tmpPlayPathList {
						tmpPlayPathListIndex = append(tmpPlayPathListIndex, item["index"].(string))
					}
					for _, item := range listSliceChunk[i] {
						tmpSliceIndex, _ := a.getSliceIndexAndName(item.Path)
						tmpListSliceChunkIndex = append(tmpListSliceChunkIndex, tmpSliceIndex)
					}
					tmpDiffRes := utils.ArrayDiff(tmpPlayPathListIndex, tmpListSliceChunkIndex, false)
					fmt.Println("切片缺失：", tmpDiffRes, tmpPlayPathListIndex, tmpListSliceChunkIndex)
				}
				ch <- []M3u8SliceVideo{{SliceVideo: tmpPlayPathList, Error: err}}
			}()
		}

		for i := 0; i < listSliceChunkLen; i++ {
			tmpM3u8SliceVideo := <-ch
			if tmpM3u8SliceVideo[0].Error != nil {
				err = tmpM3u8SliceVideo[0].Error
				common.LogToFile(path, fmt.Sprintf("获取m3u8分片视频列表失败：%v\n", err))
				continue
			}
			sliceVideoList := tmpM3u8SliceVideo[0].SliceVideo
			for j := 0; j < len(sliceVideoList); j++ {
				item := sliceVideoList[j]
				sliceIndex := item["index"].(string)
				tmpPlayPathMap[sliceIndex] = item
			}
		}
	}
	// tmpPlayPathName = utils.ArraySort(tmpPlayPathName, 1)
	for i := 0; i < listSliceLen; i++ {
		item := tmpPlayPathMap[fmt.Sprint(i)]
		if item == nil {
			continue
		}
		playPathList = append(playPathList, item)
	}

	// fmt.Println("playPathList========\n", playPathList)
	// fmt.Println("listSliceChunkLen========\n", listSliceChunkLen)
	return
}

// func (a *M3u8Handler) DoGetM3u8SliceVideo(path string, m3u8Info *M3u8Info, listSlice []string, optType string) (playPathList []map[string]interface{}, err error) {
// 	playPathList = make([]map[string]interface{}, 0)
// 	m3u8Dir := a.getM3u8Dir(path)
// 	uniqueName := a.getM3u8PathMd5(path)

// 	if optType != OptTypeVideo && optType != OptTypeCoverImg {
// 		return playPathList, errors.New("optType参数错误, 可选值：video, coverImg")
// 	}
// 	for i := 0; i < len(listSlice); i++ {
// 		sliceNameArr := strings.Split(listSlice[i], "/")
// 		sliceNameArrLen := len(sliceNameArr)
// 		re := regexp.MustCompile(`[0-9]+`)
// 		sliceIndex := re.FindString(sliceNameArr[sliceNameArrLen-1])
// 		sliceName := sliceNameArr[sliceNameArrLen-1] + ".mp4"
// 		// m3u8VideoPath := sliceMp4PathName + "/" + sliceName
// 		m3u8VideoPath := filepath.Join(sliceMp4PathName, uniqueName, sliceName)

// 		// 定义封面图路径
// 		coverImagePath := m3u8VideoPath + ".jpg"

// 		playPathListItem := map[string]interface{}{
// 			"index":      sliceIndex,
// 			"name":       sliceName,
// 			"path":       m3u8VideoPath,
// 			"cover_path": coverImagePath, // 添加封面图路径
// 			"error":      nil,
// 		}

// 		videoExists := false
// 		coverExists := false

// 		if _, err := os.Stat(filepath.Join(m3u8Dir, m3u8VideoPath)); !os.IsNotExist(err) {
// 			videoExists = true
// 		}

// 		if _, err := os.Stat(filepath.Join(m3u8Dir, coverImagePath)); !os.IsNotExist(err) {
// 			coverExists = true
// 		}

// 		// 如果视频和封面都存在，则跳过
// 		if videoExists && coverExists {
// 			playPathList = append(playPathList, playPathListItem)
// 			continue
// 		}
// 		commandArgs := []string{}
// 		if m3u8Info.ExtKeyTrue != "" && m3u8Info.ExtKeyIvTrue != "" {
// 			commandArgs = append(commandArgs, "-decryption_key", m3u8Info.ExtKeyTrue, "-decryption_iv", m3u8Info.ExtKeyIvTrue)
// 			listSlice[i] = "crypto+file:" + listSlice[i]
// 		}
// 		commandArgs = append(commandArgs, "-i", listSlice[i])

// 		if optType == OptTypeVideo {
// 			// err = exec.Command("ffmpeg", params...).Run() //  ffmpeg -allowed_extensions ALL -i "file:index.m3u8" -c copy output.mp4
// 			// err = exec.Command("ffmpeg", params2...).Run() // ffmpeg -decryption_key YOUR_KEY_HEX -decryption_iv YOUR_IV_HEX -i "crypto+file:index.m3u8_contents/0" -c copy segment_0_decrypted.ts
// 			// ffmpegStr = "ffmpeg -decryption_key f7fd2cdfb2429a9646cb69234bebc9b3 -decryption_iv 1ef58f5c956b146218c8035d387f2728 -i \"crypto+file:index.m3u8_contents/0.ts\" -c copy \"sliceMp4PathName/0.ts.mp4\""
// 			// cmd := exec.Command("cmd", "/C", ffmpegStr)

// 			commandArgs = append(commandArgs, "-c", "copy", m3u8VideoPath)
// 			cmd := exec.Command("ffmpeg",
// 				// "-decryption_key", m3u8Info.ExtKeyTrue,
// 				// "-decryption_iv", m3u8Info.ExtKeyIvTrue,
// 				// "-i", "crypto+file:"+listSlice[i],
// 				// "-c", "copy",
// 				// m3u8VideoPath,
// 				commandArgs...,
// 			)

// 			// 设置 Windows 下不显示窗口
// 			if runtime.GOOS == "windows" {
// 				cmd.SysProcAttr = &syscall.SysProcAttr{
// 					HideWindow: true,
// 				}
// 			}
// 			// 创建缓冲区用于捕获stderr
// 			// 创建多重写入器：同时写入缓冲区和终端（便于实时查看）
// 			var stderrBuf bytes.Buffer
// 			var stderrWriter io.Writer

// 			// 检查 stderr 是否可用
// 			if _, err := os.Stderr.Stat(); err == nil {
// 				// stderr 可用，创建多重写入器
// 				stderrWriter = io.MultiWriter(&stderrBuf, os.Stderr)
// 			} else {
// 				// stderr 不可用，只写入缓冲区
// 				stderrWriter = &stderrBuf
// 			}
// 			cmd.Stderr = stderrWriter
// 			cmd.Stdout = os.Stdout // 标准输出通常直接输出到终端
// 			cmd.Dir = m3u8Dir
// 			err = cmd.Run()
// 			if err != nil {
// 				cmdErr := &CommandError{
// 					Cmd:        cmd.String(),
// 					Stderr:     stderrBuf.String(),
// 					Underlying: err,
// 				}
// 				// 尝试获取退出码, 非退出错误（如命令未找到）
// 				if exitErr, ok := err.(*exec.ExitError); ok {
// 					cmdErr.ExitCode = exitErr.ExitCode()
// 				} else {
// 					cmdErr.ExitCode = -1
// 				}
// 				playPathListItem["error"] = cmdErr

// 				_, file, line, _ := runtime.Caller(0)
// 				common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, err))
// 				playPathList = append(playPathList, playPathListItem)
// 				return
// 			}
// 		}

// 		// 提取封面图
// 		if optType == OptTypeCoverImg {
// 			commandArgs = append(commandArgs, "-vframes", "1", "-an", "-sn", "-f", "image2", "-probesize", "32", "-analyzeduration", "0", "-avoid_negative_ts", "make_zero", "-fflags", "+fastseek", "-y", coverImagePath)
// 			coverCmd := exec.Command("ffmpeg",
// 				// "-decryption_key", m3u8Info.ExtKeyTrue,
// 				// "-decryption_iv", m3u8Info.ExtKeyIvTrue,
// 				// "-i", "crypto+file:"+listSlice[i],
// 				// // "-vf", "thumbnail,scale=640:-1", // 使用thumbnail过滤器提取关键帧，并缩放到宽度640
// 				// // "-vf", "scale=640:-1",  // 只做缩放，去掉thumbnail过滤器
// 				// "-vframes", "1", // 只提取一帧
// 				// "-an", // 不处理音频
// 				// "-sn", // 不处理字幕
// 				// // "-q:v", "5", // 降低质量要求以提高速度
// 				// "-f", "image2", // 图像输出格式

// 				// // "-fast", "1", // 添加快速解码参数
// 				// // "-fflags", "+fastseek", // 添加快速解码参数
// 				// // "-map_metadata", "-1", // 跳过元数据处理

// 				// "-probesize", "32", // 减少探测数据
// 				// "-analyzeduration", "0", // 快速分析
// 				// "-avoid_negative_ts", "make_zero",
// 				// "-fflags", "+fastseek", // 快速seek
// 				// "-y", // 覆盖输出文件
// 				// coverImagePath,
// 				commandArgs...,
// 			)

// 			if runtime.GOOS == "windows" {
// 				coverCmd.SysProcAttr = &syscall.SysProcAttr{
// 					HideWindow: true,
// 				}
// 			}

// 			coverCmd.Dir = m3u8Dir
// 			coverErr := coverCmd.Run()
// 			if coverErr != nil {
// 				fmt.Printf("提取封面图失败: %v\n", coverErr)
// 				playPathListItem["cover_error"] = coverErr.Error() // 记录封面提取错误
// 			}
// 		}

//			playPathList = append(playPathList, playPathListItem)
//			// fmt.Println("m3u8VideoPath=", m3u8VideoPath)
//			// fmt.Println("m3u8Dir=", m3u8Dir)
//		}
//		return playPathList, nil
//	}
func (a *M3u8Handler) DoGetM3u8SliceVideoV2(path string, pathDto *common.AllPathDto, listSlice []common.ExtListItem, optType string) (playPathList []map[string]interface{}, err error) {
	if optType != OptTypeVideo && optType != OptTypeCoverImg {
		return playPathList, errors.New("optType参数错误, 可选值：video, coverImg")
	}
	playPathList = make([]map[string]interface{}, 0)

	for i, seg := range listSlice {
		startSec := float64(seg.StartSec) / 1_000_000.0

		sliceIndex, sliceName := a.getSliceIndexAndName(seg.Path)
		m3u8VideoPath := filepath.Join(pathDto.M3u8Dir, sliceMp4PathName, pathDto.UniqueName, sliceName)

		fmt.Println(fmt.Sprintf("sliceName=%v, duration=%v, startSec=%v", sliceName, seg.ExtDuration, startSec))

		// 定义封面图路径
		coverImagePath := m3u8VideoPath + ".jpg"

		playPathListItem := map[string]interface{}{
			"index":      sliceIndex,
			"name":       sliceName,
			"time":       seg.StartTimeStr,
			"path":       m3u8VideoPath,
			"cover_path": coverImagePath, // 添加封面图路径
			"error":      nil,
		}

		videoExists := false
		coverExists := false

		if optType == OptTypeVideo {
			if _, err := os.Stat(m3u8VideoPath); !os.IsNotExist(err) {
				videoExists = true
			}
		}

		if optType == OptTypeCoverImg {
			if _, err := os.Stat(coverImagePath); !os.IsNotExist(err) {
				coverExists = true
			}
		}

		if videoExists || coverExists {
			playPathList = append(playPathList, playPathListItem)
			continue
		}
		commandArgs := []string{}
		commandPlanBArgs := []string{}
		if seg.ExtKeyTrue != "" && seg.ExtKeyIvTrue != "" {
			commandArgs = append(commandArgs, "-decryption_key", seg.ExtKeyTrue, "-decryption_iv", seg.ExtKeyIvTrue)
			seg.Path = "crypto+file:" + seg.Path
		}

		if optType == OptTypeVideo {
			// err = exec.Command("ffmpeg", params...).Run() //  ffmpeg -allowed_extensions ALL -i "file:index.m3u8" -c copy output.mp4
			// err = exec.Command("ffmpeg", params2...).Run() // ffmpeg -decryption_key YOUR_KEY_HEX -decryption_iv YOUR_IV_HEX -i "crypto+file:index.m3u8_contents/0" -c copy segment_0_decrypted.ts
			// ffmpegStr = "ffmpeg -decryption_key f7fd2cdfb2429a9646cb69234bebc9b3 -decryption_iv 1ef58f5c956b146218c8035d387f2728 -i \"crypto+file:index.m3u8_contents/0.ts\" -c copy \"sliceMp4PathName/0.ts.mp4\""
			// cmd := exec.Command("cmd", "/C", ffmpegStr)

			commandArgs = append(commandArgs, "-i", seg.Path, "-c", "copy", m3u8VideoPath)
			commandPlanBArgs = append(commandPlanBArgs, "-i", pathDto.MergeDecPath, "-ss", fmt.Sprintf("%.3f", startSec), "-t", fmt.Sprintf("%.3f", seg.ExtDuration), "-c", "copy", "-avoid_negative_ts", "1", "-y", m3u8VideoPath)

			cmd := exec.Command("ffmpeg",
				// "-decryption_key", m3u8Info.ExtKeyTrue,
				// "-decryption_iv", m3u8Info.ExtKeyIvTrue,
				// "-i", "crypto+file:"+listSlice[i],
				// "-c", "copy",
				// m3u8VideoPath,
				commandArgs...,
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
			var stderrBufPlanB bytes.Buffer
			var stderrWriter io.Writer
			var stderrWriterPlanB io.Writer

			// 检查 stderr 是否可用
			if _, err := os.Stderr.Stat(); err == nil {
				// stderr 可用，创建多重写入器
				stderrWriter = io.MultiWriter(&stderrBuf, os.Stderr)
				stderrWriterPlanB = io.MultiWriter(&stderrBufPlanB, os.Stderr)
			} else {
				// stderr 不可用，只写入缓冲区
				stderrWriter = &stderrBuf
				stderrWriterPlanB = &stderrBufPlanB
			}
			cmd.Stderr = stderrWriter
			cmd.Stdout = os.Stdout // 标准输出通常直接输出到终端
			cmd.Dir = pathDto.M3u8Dir
			err = cmd.Run()

			var cmdErr *CommandError
			// commandPlanBArgs
			if err != nil {
				cmdErr = &CommandError{
					Cmd:        cmd.String(),
					Stderr:     stderrBuf.String(),
					Underlying: err,
				}

				cmdPlanB := exec.Command("ffmpeg",
					commandPlanBArgs...,
				)
				cmdPlanB.Stderr = stderrWriterPlanB
				cmdPlanB.Stdout = os.Stdout // 标准输出通常直接输出到终端
				cmdPlanB.Dir = pathDto.M3u8Dir
				err = cmdPlanB.Run()
				if err != nil {
					cmdErr = &CommandError{
						Cmd:        cmdPlanB.String(),
						Stderr:     stderrBufPlanB.String(),
						Underlying: err,
					}
				}
			}

			if err != nil {
				// 尝试获取退出码, 非退出错误（如命令未找到）
				cmdErr.ExitCode = -1
				if err != nil {
					if exitErr, ok := err.(*exec.ExitError); ok {
						cmdErr.ExitCode = exitErr.ExitCode()
					}
				}
				playPathListItem["error"] = cmdErr

				_, file, line, _ := runtime.Caller(0)
				common.LogToFile(path, fmt.Sprintf("%s:%d %v\n", file, line, cmdErr))
				playPathList = append(playPathList, playPathListItem)
				return
			}
		}
		// 提取封面图
		if optType == OptTypeCoverImg {
			args := []string{
				"-ss", fmt.Sprintf("%.3f", startSec),
				"-i", pathDto.MergeDecPath,
				// "-vf", "thumbnail,scale=640:-1",
				// "-vcodec", "h264",
				// "-probesize", "100M", "-analyzeduration", "100M",
				// "-strict", "unofficial",
				"-vframes", "1",
				"-an", "-sn",
				"-f", "image2",
				"-y", coverImagePath,
			}
			coverCmd := exec.Command("ffmpeg", args...)
			var stderr bytes.Buffer
			coverCmd.Stderr = &stderr
			err := coverCmd.Run()
			// if err != nil && strings.Contains((stderr).String(), "dec:h263") {
			// 	args = append(args, "-vcodec", "h264")
			// 	coverCmd2 := exec.Command("ffmpeg", args...)
			// 	var stderr bytes.Buffer
			// 	coverCmd2.Stderr = &stderr
			// 	err = coverCmd2.Run()
			// }
			if err != nil {
				log.Printf("提取第 %d 个封面失败: %v \n %v", i, err, (stderr).String())
				fmt.Println("ffmpeg ", strings.Join(args, " "))
				playPathListItem["cover_error"] = err.Error() + (stderr).String() // 记录封面提取错误
			}
		}
		playPathList = append(playPathList, playPathListItem)
	}
	return
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
	fmt.Println("finalMergeFileMap	==", finalMergeFileMap)

	contentLines := strings.Split(*content, "\n")
	newContentLines := make([]string, 0)
	for i := 0; i < len(contentLines); i++ {
		line := contentLines[i]
		if strings.Contains(line, "EXTINF") {
			i++
			nextLine := strings.Trim(contentLines[i], "\r\n")
			nextLine = strings.ReplaceAll(nextLine+".ts", ".ts.ts", ".ts")
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
func (a *M3u8Handler) mergeEncryptedSegments(m3u8VideoBasePath string, segments []common.ExtListItem, mergedEncPath string, reTry int) error {
	fileInfo, err := os.Stat(mergedEncPath)
	var fileSize int64 = 0
	if err == nil && fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	if !os.IsNotExist(err) && fileSize > 0 {
		fmt.Println("合并文件已存在，跳过：" + mergedEncPath)
		return nil
	}
	if len(segments) == 0 {
		return errors.New("no segments")
	}

	out, err := os.Create(mergedEncPath)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, seg := range segments {
		in, err := os.Open(filepath.Join(m3u8VideoBasePath, seg.Path))
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, in); err != nil {
			in.Close()
			return err
		}
		in.Close()
	}
	fileSize = 0
	fileInfo, err = os.Stat(mergedEncPath)
	if err == nil && fileInfo != nil {
		fileSize = fileInfo.Size()
	}
	if reTry > 0 && fileSize == 0 {
		fmt.Println("合并文件错误，重试：" + strconv.Itoa(reTry))
		os.Remove(mergedEncPath)
		return a.mergeEncryptedSegments(m3u8VideoBasePath, segments, mergedEncPath, reTry-1)
	}
	return nil
}

func (m *M3u8Handler) mergeDecryptedSegments(listItem common.ExtListItem, mergedEncPath, mergedDecPath string) (error, bytes.Buffer) {
	var out bytes.Buffer
	if _, err := os.Stat(mergedDecPath); !os.IsNotExist(err) {
		return nil, out
	}
	args := []string{
		"-decryption_key", listItem.ExtKeyTrue,
		"-decryption_iv", listItem.ExtKeyIvTrue,
		"-i", "crypto+file:" + mergedEncPath, // 注意：输入是加密的合并文件
		"-c", "copy",
		"-y", mergedDecPath,
	}
	decryptCmd := exec.Command("ffmpeg", args...)

	decryptCmd.Stderr = &out // 捕获错误输出
	if err := decryptCmd.Run(); err != nil {
		fmt.Println("ffmpeg ", strings.Join(args, " "))
		return fmt.Errorf(err.Error() + "\n" + out.String()), out
	}
	return nil, out
}

func (a *M3u8Handler) getDurationFromExtInf(extInfStr string) (float64, error) {
	if !strings.Contains(extInfStr, "EXTINF") {
		return 0, errors.New("not extInf")
	}
	durationTemp := strings.TrimSpace(strings.Split(strings.Split(extInfStr, ":")[1], ",")[0])
	duration, err := strconv.ParseFloat(durationTemp, 64)
	return duration, err
}

func (a *M3u8Handler) getKeyAndVi(path, m3u8ItemLine string) (key, iv, method string) {
	if !strings.Contains(m3u8ItemLine, "X-KEY") {
		return
	}

	key, iv, method, _ = a.parseM3u8FileXKey(path, m3u8ItemLine)
	if !(key != "" && iv != "") {
		key = ""
		iv = ""
		common.LogToFile(path, path+"key or iv is empty")
	}
	return
}

func (a *M3u8Handler) createKeyFileIfNotExist(path string, m3u8ItemLine string, keyUriExists []string) (tryResult bool) {
	tmpKey, tmpKeyIv, tmpMethod, tmpKeyUri := a.parseM3u8FileXKey(path, m3u8ItemLine)
	if false {
		fmt.Println(tmpKey, tmpKeyIv, tmpMethod, tmpKeyUri)
	}
	if tmpKey != "" && tmpKeyIv == "" || tmpKeyUri == "" || tmpMethod == "NONE" {
		return
	}
	pathDto := a.GetGetAllPathDto(path)
	tmpKeyUriAbs := filepath.Join(pathDto.M3u8Dir, tmpKeyUri)
	if _, ok := os.Stat(tmpKeyUriAbs); !os.IsNotExist(ok) {
		return
	}
	tmpKeyFileName := filepath.Base(tmpKeyUri)
	tmpSliceName := strings.ReplaceAll(tmpKeyFileName, ".key", ".ts")
	tryParseTargetPath := filepath.Join(pathDto.SliceMp4Path, tmpSliceName)
	tryParseSourcePath := filepath.Join(a.getM3u8ContentDir(path), tmpSliceName)
	keyFilePath := filepath.Join(a.getM3u8ContentDir(path), tmpKeyFileName)

	for _, existKeyUri := range keyUriExists {
		tmpKey, tmpKeyIv, tmpMethod, tmpKeyUri = a.parseM3u8FileXKey(path, existKeyUri)
		tmpKeyUriAbs := filepath.Join(pathDto.M3u8Dir, tmpKeyUri)
		if _, ok := os.Stat(tmpKeyUriAbs); os.IsNotExist(ok) {
			continue
		}
		if tmpKey == "" || tmpKeyIv == "" || tmpMethod == "" {
			continue
		}
		tmpMethod = a.paddingKeyMethod(tmpMethod)

		var stderr bytes.Buffer
		args := []string{tmpMethod, "-d", "-in", tryParseSourcePath, "-out", tryParseTargetPath, "-nosalt", "-K", tmpKey, "-iv", tmpKeyIv}
		// openssl AES-128 -d -in 946.ts -out /dev/null -nosalt -K 30303265376262306266666164663866 -iv 00000000000000000000000000000000
		cmd := exec.Command("openssl", args...)
		cmd.Stderr = &stderr
		err := cmd.Run()
		stderrStr := strings.TrimSpace(stderr.String())
		if err != nil || stderrStr != "" {
			tmpErrStr := "openssl " + strings.Join(args, " ") + "\n验证分key、iv失败:" + err.Error() + "\n" + stderrStr
			fmt.Println(tmpErrStr)
			common.LogToFile(path, tmpErrStr)
			continue
		}

		out, err := os.Create(keyFilePath)
		if err != nil {
			common.LogToFile(path, "创建key文件失败:"+err.Error())
			out.Close()
			break
		}
		in, err := os.Open(tmpKeyUriAbs)
		if err != nil {
			common.LogToFile(path, "读取key文件失败:"+err.Error())
			in.Close()
			out.Close()
			break
		}

		if _, err := io.Copy(out, in); err != nil {
			common.LogToFile(path, "写入key文件失败:"+err.Error())
		} else {
			common.LogToFile(path, "写入key文件成功！"+m3u8ItemLine)
			tryResult = true
		}
		in.Close()
		out.Close()
		break
	}
	if !tryResult {
		common.LogToFile(path, "key file not exist:"+tmpKeyUri+";\n****** 可尝试复制key文件到m3u8文件所在目录，并重试 ******\n\n")
	}
	return
}

func (a *M3u8Handler) paddingKeyMethod(tmpMethod string) string {

	tmpMethod = strings.ToLower(tmpMethod)
	tmpMethodSlice := strings.Split(tmpMethod, "-")
	if len(tmpMethodSlice) < 3 {
		tmpMethod += "-cbc"
	}
	return tmpMethod
}
func (a *M3u8Handler) parseM3u8FileXKey(path, m3u8ItemLine string) (tmpKey, tmpKeyIv, tmpMethod, tmpKeyUri string) {
	tmpArr := strings.Split(m3u8ItemLine, ",")
	for _, v := range tmpArr {
		tmpArr2 := strings.Split(v, "=")
		if len(tmpArr2) < 2 {
			continue
		}
		if strings.Contains(v, "METHOD") {
			tmpMethod = strings.Trim(tmpArr2[1], "\r\n\"")
		}
		if strings.Contains(v, "URI") {
			tmpKeyUri = strings.Trim(tmpArr2[1], "\r\n\"")
		}
		if strings.Contains(v, "IV") {
			tmpKeyIv = strings.Trim(tmpArr2[1], "\r\n\"")
			if strings.HasPrefix(tmpKeyIv, "0x") {
				tmpKeyIv = tmpKeyIv[2:]
			}
		}
	}

	m3u8Dir := a.getM3u8Dir(path)
	keyData, _ := ioutil.ReadFile(filepath.Join(m3u8Dir, tmpKeyUri))
	tmpKey = hex.EncodeToString(keyData)
	if len(tmpKey) != 16 {

	}
	return
}
func (a *M3u8Handler) getSliceIndexAndName(slicePath string) (sliceIndex string, sliceName string) {
	if slicePath == "" {
		return
	}
	sliceNameArr := strings.Split(slicePath, "/")
	sliceNameArrLen := len(sliceNameArr)
	re := regexp.MustCompile(`[0-9]+`)
	sliceIndex = re.FindString(sliceNameArr[sliceNameArrLen-1])
	sliceName = sliceNameArr[sliceNameArrLen-1] + ".mp4"
	return
}

func (a *M3u8Handler) GetGetAllPathDto(path string, listMapKey ...string) *common.AllPathDto {
	tmpSliceMp4Path := a.getSliceMp4Path(path)

	m3u8Dir := a.getM3u8Dir(path)
	uniqueName := a.getM3u8PathMd5(path)
	m3u8VideoBasePath := filepath.Join(m3u8Dir, sliceMp4PathName, uniqueName)
	m3u8VideoPathTpl := filepath.Join(m3u8VideoBasePath, common.M3u8VideoPathTpl)
	coverImagePathTpl := filepath.Join(m3u8VideoBasePath, common.CoverImagePathTpl)
	mergeResultPathOfReCut := strings.ReplaceAll(m3u8VideoPathTpl, common.M3u8SliceNamePlaceholder, "mergeForReCut")
	reCutM3u8Path := strings.ReplaceAll(path, ".m3u8", common.ReCutNamePlaceholder+".m3u8")

	resultMp4Dir := filepath.Join(m3u8Dir, resultMp4PathName)
	resultMp4FileName := a.getM3u8PathFileName(path) + ".mp4"
	resultMp4FileName = strings.ReplaceAll(resultMp4FileName, common.ReCutNamePlaceholder, "")
	resultMp4FileRelPath := filepath.Join(resultMp4PathName, resultMp4FileName)
	resultMp4FileAbsPath := filepath.Join(resultMp4Dir, resultMp4FileName)
	mergeFromFileRelPath := common.WorkPathName + "newN3u8File.m3u8"
	mergeFromFileAbsPath := filepath.Join(m3u8Dir, mergeFromFileRelPath)

	allPathDto := &common.AllPathDto{
		SliceMp4Path:           tmpSliceMp4Path,
		UniqueName:             uniqueName,
		M3u8Dir:                m3u8Dir,
		M3u8VideoBasePath:      m3u8VideoBasePath,
		MergeEndPath:           "",
		MergeDecPath:           "",
		M3u8VideoPathTpl:       m3u8VideoPathTpl,
		CoverImagePathTpl:      coverImagePathTpl,
		MergeResultPathOfReCut: mergeResultPathOfReCut,
		ReCutM3u8Path:          reCutM3u8Path,
		ResultMp4Dir:           resultMp4Dir,
		ResultMp4FileName:      resultMp4FileName,
		ResultMp4FileRelPath:   resultMp4FileRelPath,
		ResultMp4FileAbsPath:   resultMp4FileAbsPath,
		MergeFromFileRelPath:   mergeFromFileRelPath,
		MergeFromFileAbsPath:   mergeFromFileAbsPath,
	}
	if len(listMapKey) > 0 {
		mergeEndPath := filepath.Join(m3u8VideoBasePath, listMapKey[0]+"_merged_enc.ts")
		mergeDecPath := filepath.Join(m3u8VideoBasePath, listMapKey[0]+"_merged_dec.ts")
		allPathDto.MergeEndPath = mergeEndPath
		allPathDto.MergeDecPath = mergeDecPath
	}
	return allPathDto
}

// isVideoPlayable 检测文件是否可正常播放（即含有可解码视频流）
func (a *M3u8Handler) isVideoPlayable(filePath string) bool {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_type",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)
	var out bytes.Buffer
	cmd.Stderr = &out // 捕获错误输出
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败，或者输出非空报错，认为不可播放
		return false
	}
	str := string(output)
	if strings.Contains(str, "video") {
		return true
	}
	return false
}
