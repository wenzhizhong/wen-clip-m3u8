<template>
  <div id="SidebarCmp">
    <div class="sidebar-count">
      <div>总分片：{{  playPathList.length }}</div>
      <div>已标记删除： {{  Object.keys(playPathListDeleted).length }}</div>
    </div>
    <section v-for="value in playPathList" @click="gotoSlice" :class="playPathListDeleted[value.path]? 'deleted-item' : ''" :key="value.name">
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
  playPathListDeleted:{
    Type: Object as () => Record<string, boolean>,
    default: ()=>{
      return {} as Record<string, boolean>
    }
  },
  callback: {
    type: Function,
    default: () => {}
  }
})

defineOptions({name: 'SidebarCmp'})

function gotoSlice(e:any) {
  if (e.target.innerText){
    window.location.href = "#"+ e.target.innerText
  }
}

</script>
<style scoped>
#SidebarCmp{
  width: 200px;
  flex-shrink: 0;
  height: calc(100vh - 60px);
  border-right: solid 2px #c9c9c9;
  overflow: hidden;
  overflow-y: scroll;
  .sidebar-count{
    padding: 10px 0 10px 10px;
    text-align: left;
    border-bottom: solid 1px #ececec;
  }
  &>section{
    cursor: pointer;
    box-sizing: border-box;
    border-bottom: solid 1px #ececec;
    text-align: left!important;
    padding: 6px 2px 6px 10px;
    transition: all .2s ease-in-out;
    &:hover{
      background-color: rgba(0, 149, 255, 0.3);
    }
    &.deleted-item{
      background-color: rgba(255, 0, 0, 0.1);
      border-bottom: solid 1px #c7c7c7;
    }
    &.deleted-item:hover{
      background-color: rgba(255, 0, 0, 0.3);
    }
  }
}
</style>