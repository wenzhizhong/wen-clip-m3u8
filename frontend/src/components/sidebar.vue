<template>
  <div id="SidebarCmp">
    <section v-for="value in playPathList" >
      {{value.name}}
    </section>
  </div>
</template>

<script lang="ts" setup>
import { PlayPathListInterface } from '../common/types/m3u8Slice';

  
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

defineOptions({name: 'SidebarCmp'})

</script>
<style scoped>
#SidebarCmp{
  width: 200px;
  flex-shrink: 0;
  height: inherit;
  border-right: solid 2px #c9c9c9;
  overflow: hidden;
  &>section{
    box-sizing: border-box;
    border-bottom: solid 1px #ececec;
    text-align: left!important;
    padding: 6px 2px 6px 10px;
    transition: all .2s ease-in-out;
    &:hover{
      background-color: rgba(0, 149, 255, 0.3);
    }
  }
}
</style>