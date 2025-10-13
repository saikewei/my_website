import axios from 'axios';
import { ElMessage } from 'element-plus';

const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 100000,
});

service.interceptors.request.use(
    (config) => {
        return config;
    },
    (error) => {
        console.log('request err' + error); // for debug
        return Promise.reject(error);
    }
);

service.interceptors.response.use(
    (response) => {
        // 检查 HTTP 状态码是否为 2xx 范围
        if (response.status >= 200 && response.status < 300) {
            const res = response.data;
            // 检查后端返回的 JSON 中是否包含 error 字段
            if (res.error) {
                // 如果有 error 字段，认为是业务失败
                ElMessage({ message: res.error, type: 'error' });
                return Promise.reject(new Error(res.error));
            }
            // 没有 error 字段，认为是成功，返回整个 data
            return res;
        }
        // 非 2xx 的 HTTP 状态码直接认为是错误
        return Promise.reject(new Error(response.statusText || 'Error'));
    },
    (error)=>{
        console.log('response err' + error); // for debug
        return Promise.reject(error);
    }
);

export default service;