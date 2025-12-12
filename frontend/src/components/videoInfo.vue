<script lang="ts" setup>
import {  PropType} from 'vue';
import { mergeSuccessInterface } from '../common/types/m3u8Slice';
import { formatSize, formatDuration } from '../common/utils/file';

defineOptions({ name: 'MyVideoInfo' })
const props = defineProps({
  mergeSuccessData:{
    type: Object as PropType<mergeSuccessInterface>,
    default: () => ({} as mergeSuccessInterface)
  }
})
// defineExpose({
//   formatSize
// })
</script>
<template>
  <div>
      <h3>剪辑后信息</h3>
      <ul>
        <li>
          <b>剧名：</b> <span>{{ mergeSuccessData.Name }}</span>
        </li>
        <li>
          <b>文件位置：</b> <span class="c-clolor-blue">{{ mergeSuccessData.MergePath }}</span>
        </li>
        <li>
          <b>总时长：</b> <span>{{ formatDuration(Number(mergeSuccessData.VideoInfo?.format?.duration) ) }}</span>
        </li>
        <li>
          <b>大小：</b> <span>{{ formatSize(Number(mergeSuccessData.VideoInfo?.format?.size || 0)) }}</span>
        </li>
        <li>
          <b>分辨率：</b> <span>{{ mergeSuccessData.VideoInfo?.streams?.[0]?.width || "" }} x {{ mergeSuccessData.VideoInfo?.streams?.[0]?.height || "" }}</span>
        </li>
        <li>
          <b>码率：</b> <span>{{ mergeSuccessData.VideoInfo?.streams?.[0]?.bit_rate || "" }}</span>
        </li>
        <li>
          <b>格式：</b> <span>{{ mergeSuccessData.VideoInfo?.streams?.[0]?.codec_name || "" }}</span>
        </li>
        <li>
          <b>音频：</b> <span>{{ mergeSuccessData.VideoInfo?.streams?.[1]?.codec_name || "" }}</span>
        </li>
      </ul>
    </div>
</template>

<style scoped>
ul, li{
  list-style: none;
  padding: 0;
  margin: 0;
}
ul{
  li{
    margin-bottom: 10px;
    display: flex;
    b{
      display: inline-block;
      width: 80px;
      flex-shrink: 0;
    }
    span{
      display: inline-block;
      flex: 1;
    }
    span.c-clolor-blue{
      color: #0767c1;
    }
  }
  
}
</style>