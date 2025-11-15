import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

export interface Res {
  code: number
  msg: string
  data: {
    [key: string]: any
  }
}

const _errors = {
  401: {
    message: '未授权',
    description: 'Status:401，未授权，请联系管理员！',
  },
  404: {
    message: '页面不存在',
    description: 'Status:404，页面不存在，请联系管理员！',
  },
  500: {
    message: '服务器错误',
    description: 'Status:500，服务器内部错误，请联系管理员！',
  },
}

axios.defaults.baseURL = import.meta.env.VITE_BASE_API || '/api' // 使用环境变量中的 API 前缀

// 请求拦截器
axios.interceptors.request.use(
  (config) => {
    // 为请求头添加token
    const useToken = config.headers.useToken ?? true
    if (useToken) {
      const token = localStorage.getItem('x-token')
      if (token) {
        config.headers['x-token'] = token
      }
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
axios.interceptors.response.use(
  (response) => {
    // 检查是否需要重新登录
    if (response.data.code === 7 && response.data.data && response.data.data.reload === true) {
      // 显示提示信息
      ElMessage({
        type: 'warning',
        message: response.data.msg || '您的帐户异地登陆或令牌失效，请重新登录',
        duration: 2000,
      })

      // 延迟跳转，让用户看到提示
      setTimeout(() => {
        // 清除所有缓存
        localStorage.clear()
        sessionStorage.clear()
        // 跳转到登录页
        window.location.href = '/login'
      }, 1800)

      return response
    }

    // 如果返回信息有message, 则显示
    if (response.data.msg) {
      if (response.data.showType === 'alert') {
        ElMessageBox.alert(
          response.data.msg,
          response.data.title || '提示',
          {
            dangerouslyUseHTMLString: true,
          },
        )
      }
      else {
        if (response.data.code !== 0) {
          ElMessage({
            type: 'error',
            dangerouslyUseHTMLString: true,
            message: response.data.msg,
            // description: response.data.description || ''
          })
        }
      }
    }
    return response
  },
  (error) => {
    const code = error.code || ''
    let status = error.status || ''
    let message = error.msg || ''

    if (error.response) {
      status = error.response.status || status
      message = error.response.statusText || message
    }

    if (error.response && error.response.data) {
      status = error.response.data.statusCode || status
      message
              = error.response.data.msg
                || error.response.data.statusMessage
                || message
    }

    ElMessage.error(
      `${status} ${code}: ${message || '未知错误，请联系管理员！'}`,
    )
    return Promise.reject(error)
  },
)

async function get(path = '', params = {}, config = {}) {
  return await new Promise((resolve, reject) => {
    axios({
      url: path,
      method: 'get',
      params,
      ...config,
    })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

// const execGet = async (path = "", params = {}, config = {}) => {
//     const loading = ref(false);
//     const data = ref({});
//     const execute = async () => {
//         loading.value = true;
//         data.value = await new Promise((resolve, reject) => {
//             axios({
//                 url: path,
//                 method: "get",
//                 params,
//                 ...config,
//             })
//                 .then((res) => {
//                     resolve(res.data);
//                 })
//                 .catch((err) => {
//                     reject(err);
//                 })
//                 .finally(() => {
//                     loading.value = false;
//                 });
//         });
//     };

//     return { data, loading, execute };
// };

async function post(path = '', data = {}, config = {}) {
  return await new Promise((resolve, reject) => {
    axios({
      url: path,
      method: 'post',
      data,
      ...config,
    })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

async function put(path = '', data = {}, config = {}) {
  return await new Promise((resolve, reject) => {
    axios({
      url: path,
      method: 'put',
      data,
      ...config,
    })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

async function del(path = '', data = {}, config = {}) {
  return await new Promise((resolve, reject) => {
    axios({
      url: path,
      method: 'delete',
      data,
      ...config,
    })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

// 文件上传方法
async function upload(path = '', formData: FormData, config: any = {}) {
  return await new Promise((resolve, reject) => {
    axios({
      url: path,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data',
        ...(config.headers || {}),
      },
      ...config,
    })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

// const execPost = async (path = "", params = {}, config = {}) => {
//     const loading = ref(false);
//     const data = ref({});
//     const execute = async () => {
//         data.value = await new Promise((resolve, reject) => {
//             axios({
//                 url: path,
//                 method: "post",
//                 data: params,
//                 ...config,
//             })
//                 .then((res) => {
//                     resolve(res.data);
//                 })
//                 .catch((err) => {
//                     reject(err);
//                 })
//                 .finally(() => {
//                     loading.value = false;
//                 });
//         });
//     };
//     return { data, loading, execute };
// };

export default {
  get,
  // execGet,
  post,
  put,
  del,
  upload,
  // execPost
}
