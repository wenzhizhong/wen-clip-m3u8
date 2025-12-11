export const getSystem = ():string=> {
  let ua = window && window.navigator && window.navigator.userAgent || ''
  ua = ua.toLowerCase()
  
  let system = 'unix'
  if (ua.match(/windows|win32|win64|wow32|wow64/)) {
    system = 'windows'
  } else if (ua.match(/macintosh|macintel/)) {
    system = 'mac'
  } else if (ua.match(/x11/)) {
    system = 'linux'
  } else if (ua.match(/android/)) {
    system = 'android'
  } else if (ua.match(/iphone|ipad|ipod/)){
    system = 'ios'
  }
  return system
    
}