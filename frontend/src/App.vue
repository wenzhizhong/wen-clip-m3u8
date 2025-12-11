<script lang="ts" setup>
import { reactive, onMounted } from 'vue';
import { operateType } from './common/constant/headerOperate';
import { uploadM3u8Key } from './common/constant/localStorageKey';
import { uploadM3u8Interface, ParseM3u8SliceInterface, PlayPathListInterface } from './common/types/m3u8Slice';
import Header from './components/Header.vue'
import SidebarCmp from './components/sidebar.vue'
import MyVideo from './components/video.vue'
import { getPathDir, getPathFileName } from './common/utils/path';


const state = reactive({
  uploadM3u8Data: {
    M3u8Info:     {} as ParseM3u8SliceInterface,
    PlayPathList: [] as PlayPathListInterface[]
  } as uploadM3u8Interface,

  uploadM3u8Dir: '',
})
onMounted(() => {
  initDataFromCache()
})
// 初始化数据
function initDataFromCache() {
  const data = localStorage.getItem(uploadM3u8Key)
  if (data) {
    const uploadM3u8Data = JSON.parse(data) as uploadM3u8Interface
    setUploadM3u8Data(uploadM3u8Data)
  }
}
const headerCallback = (type :string, data :uploadM3u8Interface) => {
  switch (type) {
    case operateType.updoad:
      localStorage.setItem(uploadM3u8Key, JSON.stringify(data))
      setUploadM3u8Data(data)
      break;
    default:
      break;
  }
}
function setUploadM3u8Data(data : uploadM3u8Interface) {
  state.uploadM3u8Dir = getPathDir(data.M3u8Path)
  state.uploadM3u8Data = data
  // console.log('data=', data)
  // console.log('state.uploadM3u8Data.PlayPathList=', state.uploadM3u8Data.PlayPathList)
  
}




import { ref, computed } from 'vue'
import VideoExplorer from './components/explorerVideo/VideoExplorer.vue'
import type {  SelectedItem, FileExplorerExpose } from './components/explorerVideo/file-explorer'

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

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date: string | Date): string => {
  return new Date(date).toLocaleDateString()
}
</script>

<template>
  <section class="container">
    <Header :callback="headerCallback"/>
    <section id="workBox">
      <SidebarCmp :playPathList="state.uploadM3u8Data.PlayPathList"/>
      <main>
        <section>
          <div id="selectFileName">
            {{ getPathFileName(state.uploadM3u8Data.M3u8Path) }}
          </div>
          <!-- <videoList :localPath="state.uploadM3u8Dir" :playPathList="state.uploadM3u8Data.PlayPathList"/> -->
          <VideoExplorer
            ref="explorerRef"
            :items="state.uploadM3u8Data.PlayPathList"
            :multi-select="true"
            class="explorer"
            @select="handleSelect"
            @focus-change="handleFocusChange"
            @item-activate="handleItemActivate"
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
    main{
      flex: 1;
      height: inherit;
      box-sizing: border-box;
      padding: 0 10px;
      height: 100vh;
      overflow: hidden;
      overflow-y: scroll;
      &>section{
        height: 100vh;
        #selectFileName{
          text-align: left;
          min-height: 30px;
          line-height: 30px;
          font-size: 16px;
          margin-bottom: 30px;
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