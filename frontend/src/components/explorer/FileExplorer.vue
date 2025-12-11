<template>
  <div 
    ref="containerRef"
    class="file-explorer"
    tabindex="0"
    @keydown="handleKeydown"
    @keyup="handleKeyup"
    @blur="handleBlur"
  >
    <div
      v-for="(item, index) in items"
      :key="item.id"
      ref="itemRefs"
      class="file-item"
      :class="{
        'selected': isSelected(index),
        'focused': isFocused(index),
        'anchor': isAnchor(index)
      }"
      :data-index="index"
      :data-id="item.id"
      @click="handleItemClick($event, index)"
      @mousedown="handleMouseDown($event, index)"
      @contextmenu="handleContextMenu($event, index)"
    >
      <slot name="item" :item="item" :index="index" :selected="isSelected(index)">
        <div class="item-content">
          <span class="item-name">{{ item.name }}</span>
          <span v-if="item.size" class="item-size">{{ formatSize(item.size) }}</span>
          <span v-if="item.modified" class="item-modified">{{ formatDate(item.modified) }}</span>
        </div>
      </slot>
    </div>
    
    <!-- 右键菜单（可选） -->
    <ContextMenu
      v-if="showContextMenu"
      :x="contextMenuPosition.x"
      :y="contextMenuPosition.y"
      :selected-items="selectedItems"
      :visible="showContextMenu"
      @close="closeContextMenu"
      @command="handleContextMenuCommand"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import type { FileItem, FileExplorerExpose, SelectedItem } from './file-explorer'
import ContextMenu from './ContextMenu.vue'

// Props 定义
interface FileExplorerProps {
  items: FileItem[]
  multiSelect?: boolean
  dragSelection?: boolean
}

const props = withDefaults(defineProps<FileExplorerProps>(), {
  items: () => [],
  multiSelect: true,
  dragSelection: false
})

const emit = defineEmits<{
  'select': [items: Array<{ index: number, data: FileItem }>]
  'focus-change': [index: number, item: FileItem]
  'item-activate': [index: number, item: FileItem]
}>()

// Refs
const containerRef = ref<HTMLDivElement | null>(null)
const itemRefs = ref<HTMLElement[]>([])

// 核心状态 - 使用普通变量避免响应式问题
let selectedIndexes = new Set<number>()
let anchorIndex = -1
let focusedIndex = -1
let ctrlPressed = false
let shiftPressed = false
// let showContextMenu = false
// let contextMenuPosition = { x: 0, y: 0 }

// 响应式状态用于模板渲染
const templateState = ref({
  selectedIndexes: new Set<number>(),
  anchorIndex: -1,
  focusedIndex: -1,
  showContextMenu: false,
  contextMenuPosition: { x: 0, y: 0 }
})

// 计算属性
const selectedItems = computed(() => {
  return Array.from(templateState.value.selectedIndexes)
    .sort((a, b) => a - b)
    .map(index => ({
      index,
      data: props.items[index]
    }))
})

// 辅助函数
const isSelected = (index: number): boolean => {
  return templateState.value.selectedIndexes.has(index)
}

const isFocused = (index: number): boolean => {
  return templateState.value.focusedIndex === index
}

const isAnchor = (index: number): boolean => {
  return templateState.value.anchorIndex === index
}

// 同步状态到模板
const updateTemplateState = () => {
  templateState.value = {
    selectedIndexes: new Set(selectedIndexes),
    anchorIndex,
    focusedIndex,
    showContextMenu,
    contextMenuPosition
  }
}

// 选择操作
const selectSingle = (index: number): void => {
  console.log('Select single:', index)
  
  selectedIndexes.clear()
  selectedIndexes.add(index)
  focusedIndex = index
  anchorIndex = index
  
  updateTemplateState()
  emitSelection()
  scrollToItem(index)
}

const selectRange = (start: number, end: number, merge: boolean = false): void => {
  const startIdx = Math.min(start, end)
  const endIdx = Math.max(start, end)
  
  console.log('Select range:', { start, end, startIdx, endIdx, merge })
  
  if (!merge) {
    selectedIndexes.clear()
  }
  
  for (let i = startIdx; i <= endIdx; i++) {
    selectedIndexes.add(i)
  }
  
  focusedIndex = end
  
  updateTemplateState()
  emitSelection()
  scrollToItem(end)
}

const emitSelection = () => {
  const items = selectedItems.value
  emit('select', items)
}

// 鼠标事件处理
const handleItemClick = (e: MouseEvent, index: number): void => {
  e.preventDefault()
  e.stopPropagation()
  
  console.log('Click:', {
    index,
    shiftKey: e.shiftKey,
    ctrlKey: e.ctrlKey || e.metaKey,
    anchorIndex,
    focusedIndex
  })
  
  const shiftKey = e.shiftKey
  const ctrlKey = e.ctrlKey || e.metaKey
  
  // 更新内部按键状态
  shiftPressed = shiftKey
  ctrlPressed = ctrlKey
  
  if (!props.multiSelect) {
    selectSingle(index)
    return
  }
  
  if (shiftKey) {
    // Shift + 点击：连选
    if (anchorIndex === -1) {
      // 第一次Shift点击，设置锚点为当前焦点或点击位置
      anchorIndex = focusedIndex !== -1 ? focusedIndex : index
      console.log('Set anchor to:', anchorIndex)
    }
    
    console.log('Selecting range from', anchorIndex, 'to', index, 'merge:', ctrlKey)
    selectRange(anchorIndex, index, ctrlKey)
  } else if (ctrlKey) {
    // Ctrl + 点击：切换选择
    if (selectedIndexes.has(index)) {
      selectedIndexes.delete(index)
    } else {
      selectedIndexes.add(index)
    }
    focusedIndex = index
    anchorIndex = index
    
    updateTemplateState()
    emitSelection()
    scrollToItem(index)
  } else {
    // 普通点击：单选
    selectSingle(index)
  }
}

// 键盘事件处理
const handleKeydown = (e: KeyboardEvent): void => {
  console.log('Keydown:', e.key, {
    shiftKey: e.shiftKey,
    ctrlKey: e.ctrlKey || e.metaKey,
    anchorIndex,
    focusedIndex
  })
  
  // 更新内部状态
  if (e.key === 'Shift') shiftPressed = true
  if (e.key === 'Control' || e.key === 'Meta') ctrlPressed = true
  
  // 方向键导航
  if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
    e.preventDefault()
    
    const direction = e.key === 'ArrowDown' ? 1 : -1
    let newIndex = focusedIndex
    
    if (newIndex === -1) {
      newIndex = 0
    } else {
      newIndex = Math.max(0, Math.min(props.items.length - 1, newIndex + direction))
    }
    
    if (e.shiftKey && props.multiSelect) {
      // Shift + 方向键：连选
      if (anchorIndex === -1) {
        anchorIndex = focusedIndex !== -1 ? focusedIndex : newIndex
      }
      selectRange(anchorIndex, newIndex, e.ctrlKey || e.metaKey)
    } else if (e.ctrlKey || e.metaKey) {
      // Ctrl + 方向键：只移动焦点
      focusedIndex = newIndex
      updateTemplateState()
      emit('focus-change', newIndex, props.items[newIndex])
      scrollToItem(newIndex)
    } else {
      // 方向键：单选
      selectSingle(newIndex)
    }
  } 
  // 其他快捷键
  else if (e.key === 'a' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    selectAll()
  } else if (e.key === 'Escape') {
    e.preventDefault()
    clearSelection()
  } else if (e.key === ' ' || e.key === 'Enter') {
    e.preventDefault()
    if (focusedIndex !== -1) {
      emit('item-activate', focusedIndex, props.items[focusedIndex])
    }
  }
}

const handleKeyup = (e: KeyboardEvent): void => {
  console.log('Keyup:', e.key)
  if (e.key === 'Shift') shiftPressed = false
  if (e.key === 'Control' || e.key === 'Meta') ctrlPressed = false
}

const handleBlur = (): void => {
  shiftPressed = false
  ctrlPressed = false
}

// 其他事件处理
const handleMouseDown = (e: MouseEvent, index: number): void => {
  // 拖拽选择实现
}


// 其他选择操作
const selectAll = (): void => {
  if (!props.multiSelect) return
  
  console.log('Select all')
  
  selectedIndexes.clear()
  for (let i = 0; i < props.items.length; i++) {
    selectedIndexes.add(i)
  }
  
  if (props.items.length > 0) {
    anchorIndex = 0
    focusedIndex = props.items.length - 1
  }
  
  updateTemplateState()
  emitSelection()
}

const clearSelection = (): void => {
  console.log('Clear selection')
  
  selectedIndexes.clear()
  anchorIndex = -1
  focusedIndex = -1
  
  updateTemplateState()
  emitSelection()
}

// 工具函数
const scrollToItem = (index: number): void => {
  nextTick(() => {
    const item = itemRefs.value[index]
    if (item) {
      item.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'nearest',
        inline: 'nearest' 
      })
    }
  })
}

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date: string | Date): string => {
  return new Date(date).toLocaleDateString()
}

// 在 setup 中添加右键菜单相关状态
let showContextMenu = false
let contextMenuPosition = { x: 0, y: 0 }

// 右键菜单处理
const handleContextMenu = (e: MouseEvent, index: number): void => {
  console.log("handleContextMenu========>>>")
  e.preventDefault()
  
  // 如果右键的项目不在选中列表中，先选中它
  if (!selectedIndexes.has(index)) {
    if (!e.shiftKey && !e.ctrlKey && !e.metaKey) {
      selectSingle(index)
    }
  }
  
  showContextMenu = true
  contextMenuPosition = { x: e.clientX, y: e.clientY }
  updateTemplateState()
}


const closeContextMenu = (): void => {
  showContextMenu = false
  updateTemplateState()
}

const handleContextMenuCommand = (commandId: string, items: SelectedItem[]): void => {
  console.log('Context menu command:', commandId, items)
  
  switch (commandId) {
    case 'open':
      handleOpenCommand(items)
      break
    case 'rename':
      handleRenameCommand(items)
      break
    case 'copy':
      handleCopyCommand(items)
      break
    case 'cut':
      handleCutCommand(items)
      break
    case 'delete':
      handleDeleteCommand(items)
      break
    case 'properties':
      handlePropertiesCommand(items)
      break
  }
  
  closeContextMenu()
}

// 命令处理函数
const handleOpenCommand = (items: SelectedItem[]): void => {
  console.log('打开项目:', items)
  // 实现打开逻辑
}

const handleRenameCommand = (items: SelectedItem[]): void => {
  if (items.length === 1) {
    console.log('重命名项目:', items[0])
    // 实现重命名逻辑
  }
}

const handleCopyCommand = (items: SelectedItem[]): void => {
  console.log('复制项目:', items)
  // 实现复制逻辑
}

const handleCutCommand = (items: SelectedItem[]): void => {
  console.log('剪切项目:', items)
  // 实现剪切逻辑
}

const handleDeleteCommand = (items: SelectedItem[]): void => {
  console.log('删除项目:', items)
  if (confirm(`确定要删除 ${items.length} 个项目吗？`)) {
    // 实现删除逻辑
  }
}

const handlePropertiesCommand = (items: SelectedItem[]): void => {
  console.log('显示属性:', items)
  // 实现属性显示逻辑
}
// 暴露的方法
const exposeMethods: FileExplorerExpose = {
  selectAll,
  clearSelection,
  selectItem: (index: number) => selectSingle(index),
  getSelectedItems: () => selectedItems.value,
  getSelectedIds: () => {
    return selectedItems.value.map(item => item.data.id)
  },
  focusItem: (index: number) => {
    if (index >= 0 && index < props.items.length) {
      focusedIndex = index
      updateTemplateState()
      emit('focus-change', index, props.items[index])
      scrollToItem(index)
    }
  }
}


defineExpose(exposeMethods)

// 初始化
onMounted(() => {
  if (containerRef.value) {
    containerRef.value.focus()
  }
  
  // 点击外部关闭右键菜单
  document.addEventListener('click', () => {
    if (showContextMenu) {
      showContextMenu = false
      updateTemplateState()
    }
  })
})
</script>

<style scoped>
.file-explorer {
  width: 100%;
  height: 100%;
  overflow: auto;
  outline: none;
  position: relative;
  user-select: none;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.file-item {
  padding: 8px 12px;
  margin: 1px 0;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.15s ease;
  display: flex;
  align-items: center;
  border: 2px solid transparent;
  position: relative;
}

.file-item:hover {
  background-color: rgba(0, 120, 212, 0.08);
}

.file-item.selected {
  background-color: rgba(0, 120, 212, 0.16);
}

.file-item.focused {
  border-color: #0078d4;
  background-color: rgba(0, 120, 212, 0.12);
}

.file-item.anchor::after {
  content: '';
  position: absolute;
  top: 1px;
  left: 1px;
  right: 1px;
  bottom: 1px;
  border: 1px dashed #0078d4;
  border-radius: 3px;
  pointer-events: none;
}

.file-item.selected.focused {
  background-color: rgba(0, 120, 212, 0.24);
}

.file-item.anchor.selected {
  background-color: rgba(0, 120, 212, 0.2);
}

.item-content {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-size,
.item-modified {
  margin-left: 12px;
  font-size: 0.85em;
  color: #666;
  white-space: nowrap;
}

.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 6px 0;
  min-width: 200px;
  max-width: 300px;
  z-index: 1000;
  font-size: 14px;
}
</style>