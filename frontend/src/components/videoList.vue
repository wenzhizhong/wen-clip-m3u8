<template>
  <section id="sliceMp4List">
    <template v-for="value in playPathList">
      <MyVideo :localPath="localPath + '/' + value.path"  :thumbnailUrl="localPath + '/' + value.cover_path"/>
    </template>
  </section>
</template>

<script lang="ts" setup>
import { defineComponent } from 'vue'
import { PlayPathListInterface } from '../common/types/m3u8Slice';
import MyVideo from './video.vue'

const props = defineProps({
  localPath:{
    type: String,
    default: ''
  },
  playPathList:{
    Type: Array as () => PlayPathListInterface[],
    default: ()=>{
      return [] as PlayPathListInterface[]
    }
  },
  callback: {
    type: Function,
    default: () => {}
  }
})

defineOptions({ name: 'videoList' })

</script>

<style scoped>
  
  #sliceMp4List{
    box-sizing: border-box;
    padding: 10px;
    background-color: rgb(0, 0, 0, 0.08);
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
  }
</style>