<script lang="ts">
import { createVNode, h, render } from 'vue';
import { LOCAL_FILE } from '../common/request/api';
import { getPathFileName } from '../common/utils/path';

export default {
  name: 'MyVideo',
  props: {
    localPath: {
      type: String,
      default: ''
    },
    callback: {
      type: Function,
      default: () => {}
    }
  },
  data() {
    return {
      fileName: '',
      controls: false,
    };
  },
  mounted() {},
  watch: {
    localPath: {
      handler(val: string) {
        if (val){
          this.fileName = getPathFileName(this.localPath);
        }
      },
      deep: true,
      immediate: true,
    }
  },
  methods: {
    generateVideoId() {
      return `my-video-${this.fileName}`;
    },
  },
  render() {
    let id = this.generateVideoId();

    let url = LOCAL_FILE + "?path="+encodeURIComponent(this.localPath) 
    return h('div', {
      class: ['my-video-box']
    }, [
      h('video', {
        id: id,
        ref: id,
        class: ['my-video'],
        controls: this.controls,
        src: url,
        width: '100%',
        onLoadedMetadata: () => {
          this.callback();
        },
        onMouseover: () => {
          this.controls = true;
        },
        onMouseout: () => {
          this.controls = false;
        }
      }) ,
      h('div', {
        class: ['my-video-name'],
      }, this.fileName)
    ]);
  }
};
</script>

<style scoped>
.my-video-box {
  box-sizing: border-box;
  border-radius: 6px;
  
  /* padding: 10px;
  margin-bottom: 20px; */
  border: solid 1px rgba(0, 38, 255, 0.1);
  background-color: rgba(255, 255, 255, 1 );
  width: inherit;
  overflow: hidden;
  transition: all 0.2s linear;
  .my-video{
    max-width: 100%;
  }
}
</style>