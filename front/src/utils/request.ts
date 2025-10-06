import axios from 'axios';
import { ElMessage } from 'element-plus';

const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 10000,
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
        return response;
    },
    (error)=>{
        console.log('response err' + error); // for debug
        ElMessage({
            message: error.message,
            type: 'error',
            duration: 5 * 1000,
        });
        return Promise.reject(error);
    }
);

export default service;