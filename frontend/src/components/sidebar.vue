<template>
  <div id="SidebarCmp">
    <section v-for="value in playPathList" @click="gotoSlice">
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
  }
}
</style>