/** 表单电话号码校验 */
export function useValidatePhone(rule: any, value: any, callback: any) {
  if (!value) {
    callback()
  } else {
    const phoneReg = /^1\d{10}$/
    if (!phoneReg.test(value)) {
      callback(new Error("手机号码不合规"))
    } else {
      callback()
    }
  }
}

/** 邮箱校验 */
export function useValidatePassword(rule: any, value: any, callback: any) {
  if (!value) {
    callback()
  } else if (value.length < 8) {
    callback(new Error("密码长度不能少于8位"))
  } else if (!/[a-zA-Z]/.test(value)) {
    callback(new Error("密码必须包含字母"))
  } else if (!/[0-9]/.test(value)) {
    callback(new Error("密码必须包含数字"))
  } else {
    callback()
  }
}

export function useValidateEmail(rule: any, value: any, callback: any) {
  if (!value) {
    callback()
  } else {
    const emailReg = /^[\w-]+@[\w-]+(\.[\w-]+)+$/
    if (!emailReg.test(value)) {
      callback(new Error("邮箱不合规"))
    } else {
      callback()
    }
  }
}
