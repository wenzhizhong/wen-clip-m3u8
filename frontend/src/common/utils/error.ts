export const parseGoApiError = (e:any, defMsg = "未知错误"): string => {
    let defErr =  {message:defMsg,cause:{},kind:""}
    if (e){
        if(typeof e === 'string'){
            defErr.message = e || defErr.message
        }else if (e.message){
            if (typeof e.message === 'string' && e.message.startsWith('{')){
                let errJsonString = e?.message ? e.message : ""
                defErr = errJsonString && JSON.parse(errJsonString) as {message: string, cause: any, kind: string} || defErr
                
            }else {
                defErr.message = e.message
            }
        }
    }
    let err = defErr.message
    if (defErr.kind) 
        err +=  "\n" + defErr.kind
    if (defErr.cause && Object.keys(defErr.cause).length > 0)
        err += "\n" + JSON.stringify(defErr.cause)
    return err
}