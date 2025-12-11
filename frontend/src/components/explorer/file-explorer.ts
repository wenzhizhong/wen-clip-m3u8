export interface FileItem {
  id: string | number
  name: string
  type?: 'file' | 'folder'
  size?: number
  modified?: string | Date
  [key: string]: any
}

export interface SelectedItem {
  index: number
  data: FileItem
}

export interface FileExplorerExpose {
  selectAll: () => void
  clearSelection: () => void
  selectItem: (index: number) => void
  getSelectedItems: () => SelectedItem[]
  getSelectedIds: () => (string | number)[]
  focusItem: (index: number) => void
}

// 新增：ContextMenu 相关类型
export interface ContextMenuCommand {
  id: string
  label: string
  icon?: string
  shortcut?: string
  disabled?: boolean
  danger?: boolean
  type?: 'normal' | 'separator'
  handler?: (selectedItems: SelectedItem[]) => void | Promise<void>
}

export interface ContextMenuProps {
  x: number
  y: number
  selectedItems: SelectedItem[]
  visible: boolean
  commands?: ContextMenuCommand[]
}