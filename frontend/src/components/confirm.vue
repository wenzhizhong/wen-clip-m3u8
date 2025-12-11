<script lang="ts">
import { createVNode, h, render } from 'vue';
import ConfirmComponent from './confirm.vue'; // 当前文件本身

interface ConfirmOptions {
  title?: string;
  content: string;
  onConfirm?: () => void;
  onCancel?: () => void;
}

// 创建一个容器用于挂载 confirm
let confirmContainer: HTMLElement | null = null;

function createConfirmContainer() {
  if (!confirmContainer) {
    confirmContainer = document.createElement('div');
    confirmContainer.id = 'confirm-container';
    document.body.appendChild(confirmContainer);
  }
  return confirmContainer;
}

function showConfirm(options: string | ConfirmOptions) {
  // 处理参数
  const opts: ConfirmOptions = typeof options === 'string' ? { content: options } : options;
  
  // 创建容器
  const container = createConfirmContainer();
  
  // 创建 VNode
  const vnode = createVNode(ConfirmComponent, {
    title: opts.title || '确认',
    content: opts.content,
    show: true,
    onConfirm: () => {
      opts.onConfirm && opts.onConfirm();
      render(null, container);
    },
    onCancel: () => {
      opts.onCancel && opts.onCancel();
      render(null, container);
    }
  });
  
  // 渲染到页面
  render(vnode, container);
  
  return vnode.component;
}

export const confirm = {
  show(options: string | ConfirmOptions) {
    return showConfirm(options);
  }
};

export default {
  name: 'confirm',
  props: {
    title: {
      type: String,
      default: '确认'
    },
    content: {
      type: String,
      required: true
    },
    show: {
      type: Boolean,
      default: false
    },
    onConfirm: {
      type: Function,
      default: () => {}
    },
    onCancel: {
      type: Function,
      default: () => {}
    }
  },
  data() {
    return {
      open: false
    }
  },
  mounted() {
    if (this.show) {
      this.open = true;
    }
  },
  methods: {
    handleConfirm() {
      this.onConfirm();
      this.open = false;
    },
    handleCancel() {
      this.onCancel();
      this.open = false;
    }
  },
  render() {
    if (!this.open) {
      return null;
    }

    return h('div', {
      class: 'confirm-overlay'
    }, [
      h('div', {
        class: 'confirm-dialog'
      }, [
        // 标题区域
        h('div', {
          class: 'confirm-header'
        }, [
          h('span', { class: 'confirm-title' }, this.title),
          h('button', {
            class: 'confirm-close',
            onClick: this.handleCancel
          }, '×')
        ]),
        
        // 内容区域
        h('div', {
          class: 'confirm-content'
        }, this.content),
        
        // 按钮区域
        h('div', {
          class: 'confirm-footer'
        }, [
          h('button', {
            class: 'confirm-btn confirm-btn-cancel',
            onClick: this.handleCancel
          }, '取消'),
          h('button', {
            class: 'confirm-btn confirm-btn-confirm',
            onClick: this.handleConfirm
          }, '确定')
        ])
      ])
    ]);
  }
};
</script>

<style scoped>
.confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9998;
  display: flex;
  justify-content: center;
  align-items: center;
}

.confirm-dialog {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 300px;
  max-width: 500px;
}

.confirm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px 10px;
  border-bottom: 1px solid #eee;
}

.confirm-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.confirm-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.confirm-close:hover {
  color: #333;
}

.confirm-content {
  padding: 20px;
  font-size: 16px;
  color: #666;
  line-height: 1.5;
}

.confirm-footer {
  display: flex;
  justify-content: flex-end;
  padding: 10px 20px 20px;
  gap: 10px;
}

.confirm-btn {
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  border: 1px solid #ddd;
}

.confirm-btn-cancel {
  background-color: #f5f5f5;
  color: #333;
}

.confirm-btn-cancel:hover {
  background-color: #e0e0e0;
}

.confirm-btn-confirm {
  background-color: #1890ff;
  color: white;
  border-color: #1890ff;
}

.confirm-btn-confirm:hover {
  background-color: #40a9ff;
}
</style>