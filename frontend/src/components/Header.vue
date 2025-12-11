<script lang="ts" setup>
import {reactive, ref, onMounted} from 'vue'
import {MessageDialog, OpenFileDialog, } from '../../wailsjs/go/goApi/Runtime'
import {CheckFfmpeg, OpenM3u8File} from '../../wailsjs/go/goApi/M3u8Handler'

import {toast} from './toast.vue'
import { confirm } from './confirm.vue';
import { uploadM3u8Interface} from '../common/types/m3u8Slice'
import { operateType } from '../common/constant/headerOperate'
import { getSystem } from '../common/utils/browser'
import { uploadM3u8Key } from '../common/constant/localStorageKey';
import { getPathDir } from '../common/utils/path';

const props = defineProps({
  callback:{
    type: Function,
    default: (type :string, data :uploadM3u8Interface)=>{}
  }
})
const data = reactive({
  name: "",
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

  OpenFileDialog(options).then((filePath)=> {
    if (filePath && filePath.length > 0) {
      console.log('用户选择了文件：', filePath);
      toast.warning("正在解析视频....", -1)

      OpenM3u8File(filePath).then((res :uploadM3u8Interface)=>{ 
        toast.success("解析完成", 10000)
        console.log( "res: uploadM3u8Interface=", res);
    
        res.M3u8Dir = getPathDir(filePath)
        res.M3u8Path = filePath
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
    content: '是否清空?',
    onConfirm: () => {
      let data = localStorage.getItem( uploadM3u8Key )
      if (data) {
        const uploadM3u8Data = JSON.parse(data) as uploadM3u8Interface
        

        localStorage.removeItem( uploadM3u8Key )
        props.callback(operateType.clear, uploadM3u8Data.M3u8Path)
        toast.success("已清空", 10000)
      }

    },
    onCancel: () => {
      // console.log("cancel")
    }
  })
}

</script>

<template>
  <header >
    <button class="c-btn c-btn-blain header-button" title="保存">
      <img src="/src/assets/images/header/save.png" alt="save.png">
    </button>
    <button class="c-btn c-btn-blain header-button" @click="onClearM3u8">
      <img src="/src/assets/images/header/clear.png" alt="clear.png">
    </button>

    <button class="c-btn c-btn-blain header-button" @click="onSelectM3u8">
      <img src="/src/assets/images/header/upload.png" alt="upload.png">
    </button>


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


  header {
    padding: 10px;
    box-sizing: border-box;
    border-bottom: 1px solid rgba(0,0,0,0.2);
    
    display: flex;
    button{
      margin-right: 10px;
    }
    .header-button{
      width: 36px;
      height: 36px;
      margin: 0;
      margin-right: 10px;
      padding: 0;
      img{
        width: 24px;
        height: 24px;
      }
    }
  }
</style>
