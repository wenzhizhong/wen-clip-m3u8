<script lang="ts">
import { createVNode, h, render } from 'vue';
import ToastComponent from './toast.vue'; // 当前文件本身

// 定义 Toast 类型
type ToastType = 'info' | 'success' | 'error' | 'warning';

interface ToastOptions {
  text: string;
  type?: ToastType;
  time?: number;
}

// 创建一个容器用于挂载 toast
let toastContainer: HTMLElement | null = null;

function createToastContainer() {
  if (!toastContainer) {
    toastContainer = document.createElement('div');
    toastContainer.id = 'toast-container';
    document.body.appendChild(toastContainer);
  }
  return toastContainer;
}

function showToast(options: string | ToastOptions) {
  // 处理参数
  const opts: ToastOptions = typeof options === 'string' ? { text: options } : options;
  opts.time = opts.time || 3000
  opts.type = opts.type || 'info'
  // 创建容器
  const container = createToastContainer();
  
  // 创建 VNode
  const vnode = createVNode(ToastComponent, {
    text: opts.text,
    type: opts.type ,
    time: opts.time,
    show: true,
    callback: (data: { show: boolean }) => {
      if (!data.show) {
        // 关闭时卸载组件
        render(null, container);
      }
    }
  });
  
  if (opts.time !== undefined && opts.time > 0){
    setTimeout(() => {
      render(null, container);
    }, opts.time);
  }
  // 渲染到页面
  render(vnode, container);
  
  return vnode.component;
}

// 导出不同类型的 toast 方法
export const toast = {
  info(text: string, time: number = 3000) {
    return showToast({ text, type: 'info', time });
  },
  success(text: string, time: number = 3000) {
    return showToast({ text, type: 'success', time });
  },
  error(text: string, time: number = 3000) {
    return showToast({ text, type: 'error', time });
  },
  warning(text: string, time: number = 3000) {
    return showToast({ text, type: 'warning', time });
  }
};

export default {
  name: 'toast',
  props: {
    text: {
      type: String,
      default: ''
    },
    time: {
      type: Number,
      default: 3000
    },
    show: {
      type: Boolean,
      default: false
    },
    type: {
      type: String,
      default: 'info'
    },
    callback: {
      type: Function,
      default: () => {}
    }
  },
  data() {
    return {
      timer: 0,
      open: false,
    }
  },
  mounted() {
    // 如果初始 show 为 true，则设置 open 状态
    if (this.show) {
      this.open = true;
    }
  },
  methods: {
    close() {
      this.open = false
      this.callback({show: false})
    }
  },
  render() {
    if (!this.open) {
      return null;
    }

    return h('div', {
      class: ['toast', this.type]
    }, [
      h('span', {
        class: 'text'
      }, this.text),
      h('button', {
        class: 'close',
        onClick: this.close
      }, 'X')
    ]);
  }
};
</script>

<style>
.toast {
  height: 50px;
  line-height: 50px;
  position: fixed;
  top: 100px;
  left: 50%;
  box-sizing: border-box;
  transform: translate(-50%, -50%);
  border-radius: 5px;
  border: solid 1px rgba(0, 0, 0, 0.1);
  background-color: rgba(255, 255, 255, 1);
  z-index: 9999;
}

.toast.info {
  border: solid 1px rgb(172, 172, 172);
  background-color: rgb(247, 247, 247);
}

.toast.success {
  border: solid 1px rgb(84, 192, 93);
  background-color: rgb(225, 255, 227);
}

.toast.error {
  border: solid 1px rgb(255, 139, 139);
  background-color: rgb(253, 227, 227);
}

.toast.warning {
  border: solid 1px rgb(255, 173, 106);
  background-color: rgb(255, 243, 232);
}

.toast .text {
  display: inline-block;
  min-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding-left: 10px;
}

.toast .close {
  width: 24px;
  height: 24px;
  line-height: 20px;
  text-align: center;
  font-size: 14px;
  color: #888888;
  background-color: rgb(255, 241, 241);
  position: absolute;
  top: 11px;
  right: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: solid 1px rgba(0, 0, 0, 0.2);
  border-radius: 50%;
  user-select: none;
}

.toast .close:hover {
  color: #ffffff;
  background-color: rgb(255, 141, 141);
}
</style>