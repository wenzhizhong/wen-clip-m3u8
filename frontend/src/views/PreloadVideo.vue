
<template>
  <div class="preload-video">
    <div class="top-tool-bar">
      <span class="c-button c-button-info" @click="onMultiSelectM3u8">
        <img src="/src/assets/images/header/add.png" alt="add.png"> 
        <span>添加文件</span>
      </span> 
      <span class="c-button" @click="()=>{doParseMultiple(0)}">
        <img src="/src/assets/images/header/play-white.png" alt="play-white.png"> 
        <span>开始</span>
      </span> 
      <span class="c-button c-button-warn" @click="()=>{doParseMultiple(2)}" v-show="preloadResultMsg.failNum">
        <img src="/src/assets/images/header/play-orange.png" alt="play-orange.png"> 
        <span>重试失败</span>
      </span> 
      <span class="">&nbsp;&nbsp;&nbsp;预生成预览列表</span>
      <span class="preload-result-msg" v-if="preloadResultMsg.total !== undefined">&nbsp;{{ `解析: 总数${preloadResultMsg.total}  成功${preloadResultMsg.successNum}  失败${preloadResultMsg.failNum}` }}</span>
    </div>
    <div class="file-list-box">
      <table class="c-table">
        <thead>
          <tr>
            <th>文件路径</th>
            <th style="width: 120px;">状态</th>
            <th style="width: 100px;">错误信息</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in pathLists" :key="index">
            <td style="text-align: left;">{{ index + 1}}&nbsp;&nbsp;{{ item.path }}</td>
            <!-- <td>{{ item.status === 0 ? '未解析' : item.status === 1 ? '解析成功' : '解析失败' }}</td> -->
            <td>
              <span v-if="item.tmpCurStatus === 3" class="c-tag c-tag-warning">解析中..</span>
              <span v-else-if="item.status === 0" class="c-tag c-tag-warning">未解析</span>
              <span v-else-if="item.status === 1" class="c-tag c-tag-success">解析成功</span>
              <span v-else-if="item.status === 2" class="c-tag c-tag-error c-cursor-pointer" @click="reParse(index)" title="点击重新解析">失败&重试</span>
            </td>
            <td>{{ item.error }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script  lang="ts" setup>
import { onMounted, ref } from 'vue';
import { OpenM3u8File, ReCut } from '../../bindings/clipM3u8Media/goApi/M3u8Handler';
import { OpenFileDialog } from '../../bindings/clipM3u8Media/goApi/Runtime';
import { uploadM3u8Interface } from '../common/types/m3u8Slice';
import { toast } from '../components/toast.vue';
import { confirm } from '../components/confirm.vue';
import { parseGoApiError } from '../common/utils/error';
import { getM3u8PathFileName } from '../common/utils/path';

const pathListsCacheKey = 'pathLists'
interface pathListsInterface {
  path :string
  name ?: string
  status :number // 0:未解析 1:解析成功 2:解析失败 
  tmpCurStatus :number // 0:未解析 1:解析成功 2:解析失败 3：解析中
  error ?: string
}
interface preloadResultInterface{
  total :number
  successNum :number
  failNum :number
}
const pathLists = ref<pathListsInterface[]>([])
const preloadResultMsg = ref<preloadResultInterface>({} as preloadResultInterface)
const lockKey =  ref<string>("doParseMultipleLock")

defineOptions(
  {
    name: 'PreloadVideo'
  }
)
onMounted(() => {
  init()
})

function init() {
  getListPathsStorage()
}
function getListPathsStorage() {
  let tmpData = localStorage.getItem(pathListsCacheKey)
  if (tmpData) {
    pathLists.value = JSON.parse(tmpData)
    refreshPreloadResultMsg()
  }
}
function setListPathsStorage() {
  refreshPreloadResultMsg()
  for (const element of pathLists.value) {
    if (element.tmpCurStatus === 3) {
      element.tmpCurStatus = 0
      element.tmpCurStatus = 0
    }
  }
  localStorage.setItem(pathListsCacheKey, JSON.stringify(pathLists.value))
}
function refreshPreloadResultMsg(){
  let pathLen = pathLists.value.length
  let successNum = pathLists.value.filter((item)=>item.status === 1).length
  let failNum = pathLists.value.filter((item)=>item.status === 2).length
  preloadResultMsg.value.total = pathLen
  preloadResultMsg.value.successNum = successNum
  preloadResultMsg.value.failNum =failNum   
}

function  onMultiSelectM3u8() {
  pathLists.value = []
  preloadResultMsg.value = {} as preloadResultInterface

  if(window[lockKey.value]){
    toast.warning("请勿等待操作完成！", 3000)
    return
  }

  const options = {
    Title: "请选择m3u8文件",
    Filters: [
      { DisplayName: "m3u8文件", Pattern: "*.m3u8" },
    ],
    CanChooseFiles: true,
    AllowsMultipleSelection: true,
  };

  OpenFileDialog(options).then(async (m3u8Paths: string[])=> {
    let pathLen = m3u8Paths && m3u8Paths.length || 0;
    if (pathLen > 0) {
      for (let i = 0; i < pathLen; i++) {
        if (m3u8Paths[i].indexOf("___new___")!==-1) continue;

        let fileName= getM3u8PathFileName(m3u8Paths[i])
        pathLists.value.push({
          path: m3u8Paths[i],
          name: fileName,
          status: 0,
          tmpCurStatus: 0,
        })
      }
    }else{
      toast.warning("已取消选择文件" , 10000)
    }
  }).catch((error: any)=>{
    let msg = typeof error === 'string' ? error : error.message;
    toast.error(msg, -1)
  });
}
async function  doOpenM3u8File(m3u8Path: string){
  return new Promise((resolve, reject) =>{
    OpenM3u8File(m3u8Path).then((res :uploadM3u8Interface)=>{ 
      resolve(res)
    }).catch((error: any)=>{ 
      reject(error)
    });
  })
}
async function doReCut(m3u8Path: string){
  return new Promise((resolve, reject) =>{
    ReCut (m3u8Path).then((res:any)=>{
      resolve(res)
    }).catch((error: any)=>{
      reject(error)
    })
  })
}
function reParse(index: number) {
  
  confirm.show({
    title: '提示',
    content: '是否重新预生成',
    onConfirm: async () => {
      let path = pathLists.value[index].path
      pathLists.value[index].tmpCurStatus = 3

      toast.warning("正在解析视频，请耐心等待....", -1)
      await doReCut(path).then(async (res :any)=>{
        console.log(res)
        if (res && res.Path){
            await doOpenM3u8File(res.Path).then((res)=>{
              let res1 = res as uploadM3u8Interface;
              updatePathListsItem(index,  res1.HasCoverImgError? 2 : 1, res1.HasCoverImgError? "存在报错": '')
            }, (error: any)=>{
              let errMsg = parseGoApiError( error )
              toast.error(errMsg, -1)
              updatePathListsItem(index,  2, errMsg)
            }).catch((error: any)=>{
              let errMsg = parseGoApiError( error )
              toast.error(errMsg, -1)
              updatePathListsItem(index,  2, errMsg)
            })
            
            setListPathsStorage()
        }
      }, (error: any)=>{
        let errMsg = parseGoApiError( error )
        toast.error(errMsg, -1)
        updatePathListsItem(index,  2, errMsg)
      })
      
      setTimeout(() => {
        toast.close()
      }, 1000);
    }
  })
}
async function doParseMultiple(optStatus : number = -1 ) {
  let pathListsLen = pathLists.value.length
  if (pathListsLen === 0) {
    toast.warning("请先选择m3u8文件", -1)
    return
  }
  if(window[lockKey.value]){
    toast.warning("请勿重复操作！", 3000)
    return
  }
  if (preloadResultMsg.value.total && preloadResultMsg.value.total === preloadResultMsg.value.successNum) {
    toast.warning("所有文件已解析成功，无需再次解析", -1)
    return
  }

  confirm.show({
    title: '提示',
    content: '是否开始解析',
    onConfirm: async () => {
      window[lockKey.value] = true
      for (let i = 0; i < pathLists.value.length; i++) {
        let path = pathLists.value[i].path
        if (pathLists.value[i].status != optStatus) continue;
        pathLists.value[i].tmpCurStatus = 3
        toast.warning("正在解析【"+pathLists.value[i].name+"】，请耐心等待....", -1)

        await doOpenM3u8File(path).then((res)=>{
          console.log("i=", i, "res=", res)

          let res1 = res as uploadM3u8Interface;
          updatePathListsItem(i,  res1.HasCoverImgError? 2 : 1, res1.HasCoverImgError? "请查看日志": '')
        }, (error: any)=>{
          console.log("i=", i, "error=", error)
          updatePathListsItem(i,  2, parseGoApiError( error ))
        }).catch((error: any)=>{
          console.log("i=", i, "catch error=", error)
          updatePathListsItem(i,  2, parseGoApiError( error ))
        })
        setListPathsStorage()
      }
      toast.close()
      window[lockKey.value] = false
    }
  })
  
}
function updatePathListsItem(index :number, status: number, error: string){
  pathLists.value[index].error = error
  pathLists.value[index].status = status
  pathLists.value[index].tmpCurStatus = status
}

</script>

<style scoped>
.top-tool-bar{
  margin-top: 20px;
  text-align: left;
  padding: 4px 20px;
}
.file-list-box{
  box-sizing: border-box;
  padding: 4px 20px;
  height: calc(100vh - 40px);
  overflow-y: scroll;
  table{
    width:calc(100vw - 40px);
  }
}
.preload-result-msg{
  color: rgb(0, 115, 247);
  display: inline-block;
  margin-left: 30px;
}
</style>