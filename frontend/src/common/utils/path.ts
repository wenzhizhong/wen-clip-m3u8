import { getSystem } from "./browser"

// 获取路径最后一级目录分隔符索引
export const getPathLastDirIndex = (path :string):number=>{
  path = path || ""
  let index = -1
  let os = getSystem()
  if(os.indexOf("win") > -1){
    path = path.replace(/\\/g, "/").replace(/\/\//g, "/")
  }
  index = path.lastIndexOf("/")
  return index
}
// 获取路径目录
export const getPathDir = (path :string):string=>{
  path = path || ""
  let index = getPathLastDirIndex(path)
  return path.substring(0, index)
}
// 获取路径文件名
export const getM3u8PathFileName = (path :string):string=>{
  path = path || ""
  let index = getPathLastDirIndex(path)
  return path.substring(index+1)
}