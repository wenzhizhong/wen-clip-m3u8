<script lang="ts" setup>
import { reactive, onMounted } from 'vue';
import { operateType } from './common/constant/headerOperate';
import { deleteTagKey, mergeSuccessKey, uploadM3u8Key } from './common/constant/localStorageKey';
import { uploadM3u8Interface, ParseM3u8SliceInterface, PlayPathListInterface, mergeSuccessInterface } from './common/types/m3u8Slice';
import Header from './components/Header.vue'
import SidebarCmp from './components/sidebar.vue'
import MyVideo from './components/video.vue'
import MyVideoInfo from './components/videoInfo.vue'
import { getPathDir, getPathFileName } from './common/utils/path';


const state = reactive({
  uploadM3u8Data: {
    M3u8Info:     {} as ParseM3u8SliceInterface,
    PlayPathList: [] as PlayPathListInterface[],
    M3u8Path: '',
  } as uploadM3u8Interface,

  mergeSuccessData: {
    M3u8Info:     {} as ParseM3u8SliceInterface,
    PlayPathList: [] as PlayPathListInterface[],
    M3u8Path: '',
    MergePath: '',
    Name : '',
  } as mergeSuccessInterface,

  uploadM3u8Dir: '',
  playPathListDeletedTag: {} as Record<string, boolean>,
})
onMounted(() => {
  Init()
})

function Init(){
  initDataFromCache()
  initDeleteTag()
  initMergeSucCacheData()
}
// 初始化数据
function initDataFromCache() {
  const data = localStorage.getItem(uploadM3u8Key)
  if (data) {
    const uploadM3u8Data = JSON.parse(data) as uploadM3u8Interface
    setUploadM3u8Data(uploadM3u8Data)
  }
}
const headerHeaderCallback = (type :string, data :any) => {
  switch (type) {
    case operateType.updoad:
      data = data as uploadM3u8Interface
      localStorage.setItem(uploadM3u8Key, JSON.stringify(data))
      setUploadM3u8Data(data)
      break;
    case operateType.mergeSuc:
      data = data as mergeSuccessInterface
      if (data && data.PlayPathList && data.PlayPathList.length){
        localStorage.setItem(mergeSuccessKey, JSON.stringify(data))
        setMergeSuccessData(data)
      }
      break;
    case operateType.clear:
      window.location.reload()
      break;
    case operateType.reset:
      window.location.reload()
      console.log("============= window.location.reload()")
      break;
    default:
      break;
  }
}
function setUploadM3u8Data(data : uploadM3u8Interface) {
  state.uploadM3u8Dir = getPathDir(data.M3u8Path)
  state.uploadM3u8Data = data
}
function setMergeSuccessData(data : mergeSuccessInterface) {
  state.mergeSuccessData = data
}





/** 连选、单选处理- 开始*/
import { ref, computed } from 'vue'
import VideoExplorer from './components/explorerVideo/VideoExplorer.vue'
import type {  SelectedItem, FileExplorerExpose } from './components/explorerVideo/file-explorer'
import VideoList from './components/videoList.vue';

const explorerRef = ref<FileExplorerExpose | null>(null)

// 响应式状态（来自子组件事件）
const selectedItems = ref<SelectedItem[]>([])
const focusedIndex = ref<number>(-1)
const anchorIndex = ref<number>(-1)
const selectedCount = computed(() => selectedItems.value.length)

// 事件处理
const handleSelect = (items: SelectedItem[]) => {
  selectedItems.value = items
  console.log('选中项目:', items)
}

const handleFocusChange = (index: number, item: PlayPathListInterface) => {
  focusedIndex.value = index
  console.log('焦点变化:', index, item.name)
}

const handleItemActivate = (index: number, item: PlayPathListInterface) => {
  console.log('激活项目:', index, item.name)
  // 可以在这里处理打开文件/文件夹的逻辑
}

const deletedTagCallback = (items: SelectedItem[])=> {
  if (items && items.length > 0){
    console.log('删除项目:', items)
    for (let i = 0; i < items.length; i++) {
      state.playPathListDeletedTag[items[i].data.path] = true
    }
    cacheDeleteTag()
  }
}
const deletedTagRecverCallback = (items: SelectedItem[])=> {
  if (items && items.length > 0){
    console.log('接收到删除项目:', items)
    for (let i = 0; i < items.length; i++) {
        delete state.playPathListDeletedTag[items[i].data.path]
    }
    cacheDeleteTag()
  }
}
const cacheDeleteTag = ()=>{ 
  localStorage.setItem(deleteTagKey, JSON.stringify(state.playPathListDeletedTag))
}

const initDeleteTag = ()=>{ 
  const data = localStorage.getItem(deleteTagKey)
  if (data){
    state.playPathListDeletedTag = data ? JSON.parse(data) as Record<string, boolean> : {}
  }
}

const initMergeSucCacheData = ()=>{ 
  const data = localStorage.getItem(mergeSuccessKey)
  if (data){
    let tmpData = data ? JSON.parse(data): {} 
    state.mergeSuccessData = tmpData as mergeSuccessInterface
  }
}


/** 连选、单选处理- 结束*/


</script>

<template>
  <section class="container">
    <!-- header -->
    <Header :callback="headerHeaderCallback"/>

    <section id="workBox">
      <!-- sidebar -->
      <SidebarCmp :playPathList="state.uploadM3u8Data.PlayPathList"/>
      <main>
        <section>
          <div id="selectFileName" v-if="state.uploadM3u8Data.M3u8Path">
            <h3>{{ getPathFileName(state.uploadM3u8Data.M3u8Path) }}</h3>
          </div>
          <div id="successMergeBox" v-if="state.mergeSuccessData.PlayPathList.length">
            <div >
              <VideoList :localPath="state.uploadM3u8Dir" :playPathList="state.mergeSuccessData.PlayPathList"/> 
            </div>
            <MyVideoInfo :merge-success-data="state.mergeSuccessData"></MyVideoInfo>
          </div>
          <VideoExplorer
            ref="explorerRef"
            :items="state.uploadM3u8Data.PlayPathList"
            :multi-select="true"
            class="explorer"
            @select="handleSelect"
            @focus-change="handleFocusChange"
            @item-activate="handleItemActivate"
            
            :itemsDeletedSet="state.playPathListDeletedTag"
            :deletedTagCallback="deletedTagCallback"
            :deletedTagRecverCallback="deletedTagRecverCallback"
          >
            <template #item="{ item, index, selected }">
              <div class="custom-item">
                <MyVideo :localPath="state.uploadM3u8Dir + '/' + item.path"/>
              </div>
            </template>
          </VideoExplorer>
        </section>
      </main>
    </section>
  </section>
</template>

<style >
  html,body{
    padding: 0;
    text-align: left!important;
    color: #000!important;
    box-sizing: border-box;
    background-color: #fff;
    overflow: hidden;
  }
  
  #workBox{
    display: flex;
    flex-direction: row; 
    height: 100vh;
    box-sizing: border-box;
    main{
      flex: 1;
      height: inherit;
      box-sizing: border-box;
      padding: 0 10px;
      overflow: hidden;
      overflow-y: scroll;
      &>section{
        box-sizing: border-box;
        #selectFileName{
          box-sizing: border-box;
          text-align: left;
          min-height: 30px;
          line-height: 30px;
          font-size: 16px;
          /* margin-bottom: 30px; */
        }
        #successMergeBox{ 
          box-sizing: border-box;
          display: flex;
          flex-direction: row;
          margin-bottom: 50px;
          &>div:first-child{ 
            width: 1200px;
            flex-shrink: 0;
          }
          &>div:last-child{ 
            flex: 1;
            text-align: left;
            box-sizing: border-box;
            margin-left: 20px;
            ul{
              list-style: none;
              padding: 0;
              padding-left: 20px;
              margin: 0;
              li{
                margin-bottom: 10px;
                margin-bottom: 20px;
              }
            }
          }
        }
      }
    }
  }
</style>
<style scoped>
.explorer {
  height: calc(100vh - 120px);
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: auto;
  margin-bottom: 20px;
  background: white;
}
/* file-explorer */
.custom-item {
  display: flex;
  align-items: center;
  /* padding: 10px 12px; */
  border-bottom: 1px solid #eee;
  transition: background-color 0.2s;
}

</style>