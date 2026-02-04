<script lang="ts" setup>
import {reactive, ref, onMounted} from 'vue'
import {MessageDialog, OpenFileDialog, } from '../../wailsjs/go/goApi/Runtime'
import {CheckFfmpeg, ClearM3u8FileJob, MergeM3u8File, OpenM3u8File, DeleteM3u8Source} from '../../wailsjs/go/goApi/M3u8Handler'

import {toast} from './toast.vue'
import { confirm } from './confirm.vue';
import { mergeSuccessInterface, uploadM3u8Interface} from '../common/types/m3u8Slice'
import { operateType } from '../common/constant/headerOperate'
import { getSystem } from '../common/utils/browser'
import { deleteTagKey, mergeSuccessKey, onSaveLockKey, uploadM3u8Key } from '../common/constant/localStorageKey';
import { getPathDir } from '../common/utils/path';

const props = defineProps({
  callback:{
    type: Function,
    default: (type :string, data :uploadM3u8Interface)=>{}
  },
})
const state = ref({
  zoomNumber: 0,
  zoomNumberOrigin: 0,
})

onMounted(() => {
  toast.info("检测ffmpeg中...") 
  CheckFfmpeg().then(()=>{
    setTimeout(() => {
      toast.success("检测成功", 10000)
    }, 500);
  }).catch(e => {
    toast.error("请先安装ffmpeg" , 10000)
  })
})

function  onSelectM3u8() {
  const options = {
    Title: "请选择m3u8文件",
    Filters: [
      { DisplayName: "图片文件", Pattern: "*.m3u8" },
    ]
  };

  OpenFileDialog(options).then((m3u8Path)=> {
    if (m3u8Path && m3u8Path.length > 0) {
      console.log('用户选择了文件：', m3u8Path);
      toast.warning("正在解析视频，请耐心等待....", -1)

      OpenM3u8File(m3u8Path).then((res :uploadM3u8Interface)=>{ 
        toast.success("解析完成", 10000)
        console.log( "res: uploadM3u8Interface=", res);
    
        res.M3u8Dir = getPathDir(m3u8Path)
        res.M3u8Path = m3u8Path
        doReset();
        props.callback(operateType.updoad, res)
      }).catch((error)=>{ 
        let msg = typeof error === 'string' ? error : error.message;
        toast.error(msg, -1)
      });
    }else{
      toast.warning("已取消选择文件" , 10000)
    }
  }).catch((error) => {
    let msg = typeof error === 'string' ? error : error.message;
    toast.error(msg, -1)
  });
}
function onClearM3u8(){
  confirm.show({
    title: '清空',
    content: '是否清空已生成数据?',
    onConfirm: () => {
      let data = localStorage.getItem( uploadM3u8Key )
      if (data) {
        const uploadM3u8Data = JSON.parse(data) as uploadM3u8Interface
        ClearM3u8FileJob(uploadM3u8Data.M3u8Path).then((res)=>{ 
          localStorage.removeItem( uploadM3u8Key )
          localStorage.removeItem( deleteTagKey )
          localStorage.removeItem( onSaveLockKey )
          localStorage.removeItem( mergeSuccessKey )
          props.callback(operateType.clear, uploadM3u8Data.M3u8Path)
          toast.success("已清空", 10000)
        }).catch((error)=>{ 
          let msg = typeof error === 'string' ? error : error.message;
          toast.error(msg, -1)
        });

        
        
      }
    },
    onCancel: () => {
      // console.log("cancel")
    }
  })
}
function onDeleteSource(){
  confirm.show({
    title: '删除',
    content: '是否删除当前数据源?',
    onConfirm: () => {
      let data = localStorage.getItem( uploadM3u8Key )
      if (data) {
        const uploadM3u8Data = JSON.parse(data) as uploadM3u8Interface
        DeleteM3u8Source(uploadM3u8Data.M3u8Path).then((res)=>{ 
          console.log(res)
          toast.success("已删除", 10000)
        }).catch((error)=>{ 
          let msg = typeof error === 'string' ? error : error.message;
          toast.error(msg, -1)
        });
      }
    },
    onCancel: () => {
      // console.log("cancel")
    }
  })
}

// onReset
function onReset(){ 
  confirm.show({
    title: '温馨提示',
    content: '是否重新编辑?',
    onConfirm: () => {
      doReset()
      props.callback(operateType.reset, {})
    }
  })
}
function doReset(){ 
    localStorage.removeItem( deleteTagKey )
    localStorage.removeItem( onSaveLockKey )
    localStorage.removeItem( mergeSuccessKey )
}
// onSave 
function onSave(){ 
  sessionStorage.getItem(onSaveLockKey)
  if (sessionStorage.getItem(onSaveLockKey)) {
    toast.error("请勿重复操作", -1)
    return
  }
  confirm.show({
    title: '合并生成',
    content: '是否并合并生成视频?',
    onConfirm: () => {

      let tmpSliceData =  localStorage.getItem( uploadM3u8Key )
      let tmpDeletedItemsMap =  localStorage.getItem( deleteTagKey )
      let sliceData = tmpSliceData ? JSON.parse(tmpSliceData) as uploadM3u8Interface : {} as  uploadM3u8Interface
      let deletedItemsMap = tmpDeletedItemsMap ? JSON.parse(tmpDeletedItemsMap)  : null
  
      let resetM3u8SliceData :string[]= [];
      if (!(sliceData && sliceData.PlayPathList && sliceData.PlayPathList.length > 0)) {
        toast.error("没有可合并文件", -1)
        return
      }
      if (!sliceData.M3u8Path){
        toast.error("没有获取到m3u8文件路径，请先选择m3u8文件重新剪辑", -1)
        return
      }
      for (let i = 0; i < sliceData.PlayPathList.length; i++) {
        const item = sliceData.PlayPathList[i];
        if (deletedItemsMap && deletedItemsMap[item.path]) {
          continue;
        }
        resetM3u8SliceData.push(item.path)
      }
      
      sessionStorage.setItem(onSaveLockKey, "1")
      toast.warning("正在合并生成视频....", -1)

      MergeM3u8File(sliceData.M3u8Path, resetM3u8SliceData).then((res :mergeSuccessInterface)=>{
        sessionStorage.removeItem(onSaveLockKey)
        toast.success("已生成", -1)
        
        res.M3u8Dir = getPathDir(sliceData.M3u8Path)
        res.M3u8Path = sliceData.M3u8Path
        props.callback(operateType.mergeSuc, res)
      }).catch((error)=>{ 
        sessionStorage.removeItem(onSaveLockKey)

        let msg = typeof error === 'string' ? error : error.message;
        toast.error(msg ||"合并视频失败", -1)
      });
    }
  })
}
// 缩放
function zoomClice(type :number){
  if (type === 1) {
    state.value.zoomNumber = state.value.zoomNumber + 20
  }else if (type === -1) {
    state.value.zoomNumber = state.value.zoomNumber - 20
  }else{
    state.value.zoomNumber = 0
  }
  let data = {zoomNumber: state.value.zoomNumber} as uploadM3u8Interface
  new Promise((resolve, reject) => { 
    resolve(true)
  }).then((res)=>{ 
    props.callback(operateType.zoom, data )
  })
}

</script>

<template>
  <header >
    <span class="opt-item" @click="onSelectM3u8">
      <img src="/src/assets/images/header/upload.png" alt="upload.png"> 
      <span>上传</span>
    </span>

    <span class="opt-item" @click="onSave">
      <img src="/src/assets/images/header/save.png" alt="save.png"> 
      <span>生成</span>
    </span>

    <span class="opt-item" @click="onReset">
      <img src="/src/assets/images/header/reset.png" alt="reset.png"> 
      <span>重置</span>
    </span>

    <span class="opt-item" @click="onClearM3u8">
      <img src="/src/assets/images/header/clear.png" alt="clear.png"> 
      <span>清空临时文件</span>
    </span>

    <span class="opt-item" @click="onDeleteSource">
      <img src="/src/assets/images/header/delete.png" alt="delete.png"> 
      <span>删除源</span>
    </span>

    <span class="p-opt-item">
      <span class="opt-item" @click="()=>{zoomClice(0)}">
        <img src="/src/assets/images/header/ic_zoom_origin.png" alt="ic_zoom_origin.png"> 
        <span>默认</span>
      </span>
      <span class="opt-item" @click="()=>{zoomClice(1)}">
        <img src="/src/assets/images/header/ic_zoom_in.png" alt="ic_zoom_in.png"> 
        <span>放大</span>
      </span>
      <span class="opt-item" @click="()=>{zoomClice(-1)}">
        <img src="/src/assets/images/header/ic_zoom_out.png" alt="ic_zoom_out.png"> 
        <span>缩小</span>
      </span>
    </span>
  </header>
</template>

<style scoped>
  *{
    color: #000;
  }
  .c-btn{
    background-color: rgba(0, 0, 0, 0);
    border: 1px solid rgba(0, 0, 0, 0.4);
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;  
    user-select: none;
    transition: all 0.2s ease-in-out;
  }
  .c-btn:hover{
    background-color: rgba(0, 0, 0, 0.07);
    border: 1px solid rgba(0, 0, 0, 0.6);
  }
  .c-btn-primary{
    color: #fff;
    border: 1px solid rgba(0, 132, 255, 0.2);
    background-color: rgba(0, 132, 255, 0.8);
  }
  .c-btn-primary:hover{
    background-color: rgba(0, 132, 255, 1);
    border: 1px solid rgba(0, 132, 255, 1);
  }
  .c-btn-blain{
    color: #000;
    background-color: transparent;
  }
  .c-btn-blain:hover{
    background-color: rgba(0, 132, 255, 0.1);
  }

  .space-wdith{
    margin-right: 50px!important;
  }

  header {
    padding: 10px;
    box-sizing: border-box;
    border-bottom: 1px solid rgba(0,0,0,0.2);
    
    display: flex;
    .p-opt-item{
      display: flex;
      margin: 0 100px;
    }
    .opt-item{
      user-select: none;
      display: flex;
      height: 24px;
      line-height: 24px;
      align-items: center;
      padding: 4px 10px;
      transition: all 0.2s ease-in-out;
      border-radius: 4px;
      &:hover{
        cursor: pointer;
        background-color: #dae8ff;
      }
      img{
        margin-right: 4px;
        width: 16px;
        height: 16px;
      }

    }
  }

</style>
