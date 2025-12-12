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
      <ul>
        <li>
          <b>剧名：</b> {{ mergeSuccessData.Name }}
        </li>
        <li>
          <b>文件位置：</b> <span class="c-clolor-blue">{{ mergeSuccessData.MergePath }}</span>
        </li>
        <li>
          <b>总时长：</b> {{ formatDuration(Number(mergeSuccessData.VideoInfo?.format?.duration) ) }}
        </li>
        <li>
          <b>大小：</b> {{ formatSize(Number(mergeSuccessData.VideoInfo?.format?.size || 0)) }}
        </li>
        <li>
          <b>分辨率：</b> {{ mergeSuccessData.VideoInfo?.streams?.[0]?.width || "" }} x {{ mergeSuccessData.VideoInfo?.streams?.[0]?.height || "" }}
        </li>
        <li>
          <b>码率：</b> {{ mergeSuccessData.VideoInfo?.streams?.[0]?.bit_rate || "" }}
        </li>
        <li>
          <b>格式：</b> {{ mergeSuccessData.VideoInfo?.streams?.[0]?.codec_name || "" }}
        </li>
        <li>
          <b>音频：</b> {{ mergeSuccessData.VideoInfo?.streams?.[1]?.codec_name || "" }}
        </li>
      </ul>
    </div>
</template>

<style scoped>
ul>li{
  b{
    display: inline-block;
    width: 80px;
  }
  .c-clolor-blue{
    color: #0767c1;
  }
}
</style>