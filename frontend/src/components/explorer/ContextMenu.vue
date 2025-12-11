<template>
  <Teleport to="body">
    <div 
      v-if="visible"
      ref="menuRef"
      class="context-menu"
      :style="menuStyle"
      @click.stop="handleMenuClick"
      @mousedown.stop
      @contextmenu.stop.prevent
    >
      <div 
        v-for="command in filteredCommands" 
        :key="command.id"
        class="context-menu-item"
        :class="{
          'disabled': command.disabled,
          'danger': command.danger,
          'separator': command.type === 'separator'
        }"
        @mousedown="handleCommandMousedown(command, $event)"
        @click="handleCommandClick(command)"
      >
        <template v-if="command.type !== 'separator'">
          <span class="context-menu-icon" v-if="command.icon">
            {{ command.icon }}
          </span>
          <span class="context-menu-label">{{ command.label }}</span>
          <span class="context-menu-shortcut" v-if="command.shortcut">
            {{ command.shortcut }}
          </span>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import type { SelectedItem } from './file-explorer'

// 命令类型定义
export interface ContextMenuCommand {
  id: string
  label?: string
  icon?: string
  shortcut?: string
  disabled?: boolean
  danger?: boolean
  type?: 'normal' | 'separator'
  handler?: (selectedItems: SelectedItem[]) => void | Promise<void>
}

// Props 定义
interface ContextMenuProps {
  x: number
  y: number
  selectedItems: SelectedItem[]
  visible: boolean
  commands?: ContextMenuCommand[]
}

const props = withDefaults(defineProps<ContextMenuProps>(), {
  selectedItems: () => [],
  commands: undefined,
  visible: false
})

// Emits 定义
const emit = defineEmits<{
  'close': []
  'command': [commandId: string, selectedItems: SelectedItem[]]
}>()

// Refs
const menuRef = ref<HTMLDivElement | null>(null)
const adjustedX = ref(0)
const adjustedY = ref(0)

// 默认命令集
const defaultCommands: ContextMenuCommand[] = [
  {
    id: 'open',
    label: '打开',
    icon: '📂',
    disabled: false,
    type: 'normal'
  },
  {
    id: 'rename',
    label: '重命名',
    icon: '✏️',
    disabled: true,
    type: 'normal'
  },
  {
    id: 'separator-1',
    type: 'separator'
  },
  {
    id: 'copy',
    label: '复制',
    icon: '📋',
    shortcut: 'Ctrl+C',
    type: 'normal'
  },
  {
    id: 'cut',
    label: '剪切',
    icon: '✂️',
    shortcut: 'Ctrl+X',
    type: 'normal'
  },
  {
    id: 'paste',
    label: '粘贴',
    icon: '📝',
    shortcut: 'Ctrl+V',
    disabled: true,
    type: 'normal'
  },
  {
    id: 'separator-2',
    type: 'separator'
  },
  {
    id: 'delete',
    label: '删除',
    icon: '🗑️',
    shortcut: 'Del',
    danger: true,
    type: 'normal'
  },
  {
    id: 'separator-3',
    type: 'separator'
  },
  {
    id: 'properties',
    label: '属性',
    icon: '📊',
    disabled: false,
    type: 'normal'
  }
]

// 计算属性
const usedCommands = computed<ContextMenuCommand[]>(() => {
  return props.commands || defaultCommands
})

const filteredCommands = computed<ContextMenuCommand[]>(() => {
  return usedCommands.value.map(command => {
    const cmd = { ...command }
    
    if (cmd.id === 'rename' || cmd.id === 'properties') {
      cmd.disabled = props.selectedItems.length !== 1
    } else if (cmd.id === 'open') {
      cmd.disabled = props.selectedItems.length === 0
    } else if (cmd.id === 'delete') {
      cmd.disabled = props.selectedItems.length === 0
    }
    
    return cmd
  })
})

const menuStyle = computed(() => {
  return {
    left: `${adjustedX.value}px`,
    top: `${adjustedY.value}px`,
    display: props.visible ? 'block' : 'none'
  }
})

// 方法
const handleCommandMousedown = (command: ContextMenuCommand, e: MouseEvent): void => {
  // 防止事件冒泡
  e.stopPropagation()
  if (command.disabled || command.type === 'separator') {
    e.preventDefault()
  }
}

const handleCommandClick = (command: ContextMenuCommand): void => {
  console.log('Command clicked:', command.id)
  
  if (command.disabled || command.type === 'separator') return
  
  emit('command', command.id, props.selectedItems)
  emit('close')
}

const handleMenuClick = (e: MouseEvent): void => {
  e.stopPropagation()
}

// 调整菜单位置
const adjustMenuPosition = (): void => {
  if (!menuRef.value) {
    // 如果没有 ref，直接使用 props 位置
    adjustedX.value = props.x
    adjustedY.value = props.y
    return
  }
  
  nextTick(() => {
    if (!menuRef.value) return
    
    const rect = menuRef.value.getBoundingClientRect()
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    
    let newX = props.x
    let newY = props.y
    
    // 水平方向调整
    if (newX + rect.width > windowWidth) {
      newX = windowWidth - rect.width - 10
    }
    
    // 垂直方向调整
    if (newY + rect.height > windowHeight) {
      newY = windowHeight - rect.height - 10
    }
    
    // 确保位置不小于0
    newX = Math.max(10, newX)
    newY = Math.max(10, newY)
    
    adjustedX.value = newX
    adjustedY.value = newY
  })
}

// 点击外部关闭菜单（只监听左键点击）
const handleClickOutside = (e: MouseEvent): void => {
  if (!props.visible) return
  
  // 只处理左键点击
  if (e.button !== 0) return
  
  // 如果点击的是菜单本身，不关闭
  if (menuRef.value && menuRef.value.contains(e.target as Node)) {
    return
  }
  
  emit('close')
}

// 键盘事件处理
const handleKeydown = (e: KeyboardEvent): void => {
  if (!props.visible) return
  
  switch(e.key) {
    case 'Escape':
      e.preventDefault()
      emit('close')
      break
    case 'ArrowUp':
      e.preventDefault()
      // 这里可以添加键盘导航
      break
    case 'ArrowDown':
      e.preventDefault()
      // 这里可以添加键盘导航
      break
    case 'Enter':
      e.preventDefault()
      // 这里可以添加键盘导航
      break
  }
}

// 监听 visible 变化
watch(() => props.visible, (newVisible, oldVisible) => {
  console.log('ContextMenu visibility changed:', { newVisible, oldVisible })
  
  if (newVisible) {
    // 立即设置初始位置
    adjustedX.value = props.x
    adjustedY.value = props.y
    
    // 延迟调整位置，确保 DOM 已渲染
    setTimeout(() => {
      adjustMenuPosition()
    }, 10)
    
    // 只添加左键点击监听，不监听右键
    document.addEventListener('mousedown', handleClickOutside)
    document.addEventListener('keydown', handleKeydown)
    window.addEventListener('resize', adjustMenuPosition)
    
    // 添加全局样式防止滚动
    document.body.style.overflow = 'hidden'
  } else {
    // 移除事件监听
    document.removeEventListener('mousedown', handleClickOutside)
    document.removeEventListener('keydown', handleKeydown)
    window.removeEventListener('resize', adjustMenuPosition)
    
    // 恢复滚动
    document.body.style.overflow = ''
  }
}, { immediate: true })

// 监听位置变化
watch(() => props.x, () => {
  if (props.visible) {
    adjustMenuPosition()
  }
})

watch(() => props.y, () => {
  if (props.visible) {
    adjustMenuPosition()
  }
})

// 清理
onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', adjustMenuPosition)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.15);
  padding: 6px 0;
  min-width: 220px;
  max-width: 320px;
  z-index: 9999;
  font-size: 14px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.98);
}

.context-menu-item {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.15s ease;
  position: relative;
  color: #333;
  user-select: none;
}

.context-menu-item:not(.separator):not(.disabled):hover {
  background-color: rgba(0, 120, 212, 0.1);
}

.context-menu-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.context-menu-item.disabled:hover {
  background-color: transparent;
}

.context-menu-item.danger {
  color: #d13438;
}

.context-menu-item.danger:hover {
  background-color: rgba(209, 52, 56, 0.1);
}

.context-menu-item.separator {
  height: 1px;
  padding: 0;
  margin: 6px 0;
  background-color: #eee;
  cursor: default;
}

.context-menu-icon {
  margin-right: 12px;
  width: 18px;
  text-align: center;
  font-size: 16px;
  flex-shrink: 0;
}

.context-menu-label {
  flex: 1;
  font-weight: 400;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.context-menu-shortcut {
  margin-left: 16px;
  color: #888;
  font-size: 0.9em;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  flex-shrink: 0;
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.context-menu {
  animation: fadeIn 0.15s ease-out;
  transform-origin: top left;
}

/* 暗色主题支持 */
@media (prefers-color-scheme: dark) {
  .context-menu {
    background: rgba(30, 30, 30, 0.98);
    border-color: #444;
  }
  
  .context-menu-item {
    color: #e0e0e0;
  }
  
  .context-menu-item.separator {
    background-color: #444;
  }
  
  .context-menu-shortcut {
    color: #aaa;
  }
  
  .context-menu-item:not(.separator):not(.disabled):hover {
    background-color: rgba(0, 120, 212, 0.2);
  }
  
  .context-menu-item.danger {
    color: #ff6b6b;
  }
}
</style>