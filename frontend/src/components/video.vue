<script lang="ts">
import { createVNode, h, render } from 'vue';
import { LOCAL_FILE } from '../common/request/api';
import { getM3u8PathFileName } from '../common/utils/path';

export default {
  name: 'MyVideo',
  props: {
    localPath: {
      type: String,
      default: ''
    },
    thumbnailUrl: {
      type: String,
      default: ''
    },
    m3u8Path: {
      type: String,
      default: ''
    },
    time: {
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
      isVideoLoaded: false, // 控制视频是否已加载
      isHovered: false      // 鼠标悬停状态
    };
  },
  mounted() {
  },
  watch: {
    localPath: {
      handler(val: string) {
        if (val) {
          this.fileName = getM3u8PathFileName(this.localPath);
        }
      },
      deep: true,
      immediate: true,
    },
  },
  methods: {
    generateVideoId() {
      return `my-video-${this.fileName}`;
    },
    
    
    // 加载视频
    loadVideo() {
      this.isVideoLoaded = true;
    },
    
    // 卸载视频（释放资源）
    unloadVideo() {
      this.isVideoLoaded = false;
      // 获取视频元素并释放资源
      const videoElement = document.getElementById(this.generateVideoId()) as HTMLVideoElement;
      if (videoElement) {
        videoElement.pause();
        videoElement.removeAttribute('src');
        videoElement.load();
      }
    },
    
    // 处理鼠标悬停
    handleMouseOver() {
      this.isHovered = true;
      // 可以选择在此处加载视频
      if (!this.isVideoLoaded) {
        this.loadVideo();
      }
    },
    
    // 处理鼠标离开
    handleMouseOut() {
      this.isHovered = false;
      // 可以选择在此处卸载视频以节省资源
      // this.unloadVideo(); // 根据需求决定是否卸载
    },
    
    // 视频加载完成回调
    handleLoadedMetadata() {
      this.callback();
    },
    
    // 视频结束播放
    handleVideoEnded() {
      // 可以选择在此处卸载视频以节省资源
      this.unloadVideo();
    },
    getSrcLink(targetPath :string) {
      return LOCAL_FILE + "?path=" + encodeURIComponent(targetPath) + "&m3u8Path=" + encodeURIComponent(this.m3u8Path) ;
    },
  },
  render() {
    let id = this.generateVideoId();
    // console.log("this.localPath=", this.localPath)
    // console.log("this.thumbnailUrl=", this.thumbnailUrl)
    let url = this.getSrcLink(this.localPath);
    let thumbnailUrl = this.getSrcLink(this.thumbnailUrl);

    // 根据 isVideoLoaded 状态决定渲染哪个元素
    let videoElement;
    
    if (this.isVideoLoaded) {
      // 渲染视频元素
      videoElement = h('video', {
        id: id,
        ref: id,
        class: ['my-video'],
        controls: true,
        src: url,
        width: '100%',
        onLoadedmetadata: this.handleLoadedMetadata,
        onEnded: this.handleVideoEnded,
        // onMouseover: this.handleMouseOver,
        // onMouseout: this.handleMouseOut
      });
    } else {
      // 渲染缩略图元素
      videoElement = h('div', {
        class: ['my-video-cover-box'],
        onClick: this.loadVideo, // 点击加载视频
        // onMouseover: this.handleMouseOver,
        // onMouseout: this.handleMouseOut,
      }, [
        h('div', {
          class: ['my-video-cover-img'],
        }, [
          h('img', {
            src: `${thumbnailUrl}`,
          }, ''),
        ]),
        // 播放按钮
        h('div', {
          class: ['my-video-play-icon'],
        }, '▶')
      ]);
    }

    return h('div', {
      class: ['my-video-box']
    }, [
      videoElement,
      h('div', {
        class: ['my-video-info'],
      }, [
        h('div', {
          class: ['my-video-name'],
        }, this.fileName),
        h('div', {
          class: ['my-video-time'],
        }, this.time),
      ]),
     
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
  background-color: rgba(255, 255, 255, 1);
  width: inherit;
  overflow: hidden;
  transition: all 0.2s linear;
  
  .my-video {
    min-height: 190px;
    max-width: 100%;
  }
  
  .my-video-cover-box {
    cursor: pointer;
    transition: opacity 0.3s ease;
    
    width: 100%;
    min-height: 190px;
    cursor: pointer;
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    .my-video-cover-img{
      img{
        width: 100%;
      }
    }
    .my-video-play-icon{
      width: 50px;
      height: 50px;
      line-height: 50px;
      background-color: rgba(0, 0, 0, 1);
      opacity: 0.2;
      border-radius: 50%;
      color: white;
      font-size: 20px;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      transition: all 0.3s ease;
    }
  }
  
  .my-video-cover-box:hover {
    .my-video-play-icon{
      opacity: 0.6;
    }
    /* opacity: 0.9; */
  }
  .my-video-info{
    display: flex;
    box-sizing: border-box;
    padding: 2px 6px;
    .my-video-name{
      /* 自动伸缩 */
      flex: 1;
      text-align: left;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .my-video-time{
      /* 禁止伸缩 */
      flex: 0;
    }
  }
}
</style>