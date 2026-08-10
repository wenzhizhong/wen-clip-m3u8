
<template>
  <div class="preload-video">
    <div class="top-tool-bar">
      <span class="c-button c-button-info" @click="onMultiSelectM3u8">
        <img src="/src/assets/images/header/upload.png" alt="upload.png"> 
        <span>多文件上传</span>
      </span> 
      <span>&nbsp;&nbsp;&nbsp;预生成预览列表</span>
    </div>
    <div class="file-list-box">
      <table class="c-table">
        <thead>
          <tr>
            <th>文件路径</th>
            <th style="width: 100px;">状态</th>
            <th style="width: 100px;">错误信息</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in pathLists" :key="index">
            <td style="text-align: left;">{{ item.path }}</td>
            <!-- <td>{{ item.status === 0 ? '未解析' : item.status === 1 ? '解析成功' : '解析失败' }}</td> -->
            <td>
              <span v-if="item.status === 0" class="c-tag c-tag-warning">未解析</span>
              <span v-if="item.status === 1" class="c-tag c-tag-success">解析成功</span>
              <span v-if="item.status === 2" class="c-tag c-tag-error">解析失败</span>
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
import { OpenM3u8File } from '../../bindings/clipM3u8Media/goApi/M3u8Handler';
import { OpenFileDialog } from '../../bindings/clipM3u8Media/goApi/Runtime';
import { uploadM3u8Interface } from '../common/types/m3u8Slice';
import { toast } from '../components/toast.vue';

const pathListsCacheKey = 'pathLists'
interface pathListsInterface {
  path :string
  status :number // 0:未解析 1:解析成功 2:解析失败
  error ?: string
}
const pathLists = ref<pathListsInterface[]>([])

defineOptions(
  {
    name: 'PreloadVideo'
  }
)
onMounted(() => {
  init()
})

function init() {
  getListPaths()
}
function getListPaths() {
  let tmpData = localStorage.getItem(pathListsCacheKey)
  if (tmpData) {
    pathLists.value = JSON.parse(tmpData)
  }
}
function setListPaths() {
  localStorage.setItem(pathListsCacheKey, JSON.stringify(pathLists.value))
}


function  onMultiSelectM3u8() {
  pathLists.value = []

  const options = {
    Title: "请选择m3u8文件",
    Filters: [
      { DisplayName: "图片文件", Pattern: "*.m3u8" },
    ],
    CanChooseFiles: true,
    AllowsMultipleSelection: true,
  };

  OpenFileDialog(options).then(async (m3u8Paths: string[])=> {
    let pathLen = m3u8Paths && m3u8Paths.length || 0;
    if (pathLen > 0) {
      toast.warning("正在解析视频，请耐心等待....", -1)

      for (let i = 0; i < pathLen; i++) {
        pathLists.value.push({
          path: m3u8Paths[i],
          status: 0
        })
      }
      for (let i = 0; i < pathLen; i++) {
        await doOpenM3u8File(m3u8Paths[i]).then((res)=>{
          let res1 = res as uploadM3u8Interface;
          toast.success(`解析成功${pathLen}/${i+1}`, 10000)
          pathLists.value[i].status = 1
        }).catch((error: any)=>{
          let msg = typeof error === 'string' ? error : error.message;
          toast.error(msg, -1)
          let item = JSON.parse(JSON.stringify(pathLists.value[i]))
          pathLists.value[i] = item as pathListsInterface
        })
        setListPaths()
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
      let msg = typeof error === 'string' ? error : error.message;
      reject(msg)
    });
  })
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
</style>