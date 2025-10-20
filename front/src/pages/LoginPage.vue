<template>
    <div class="login-container">
        <el-card class="login-card">
            <template #header>
                <div class="card-header">
                    <span>管理员登录</span>
                </div>
            </template>
            <el-form @submit.prevent="handleLogin">
                <el-form-item label="密码">
                    <el-input 
                        v-model="password"
                        type="password" 
                        placeholder="请输入密码"
                        size="large"
                        show-password
                    />
                </el-form-item>
                <el-form-item>
                    <el-button 
                        type="primary" 
                        style="width: 100%;"
                        :loading="loading"
                        size="large"
                        @click="handleLogin"
                    >登录</el-button>
                </el-form-item>
            </el-form>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import request from '@/utils/request';

interface LoginResponse {
    token: string;
}

const router = useRouter();
const password = ref('');
const loading = ref(false);

const handleLogin = async () => {
    if (!password.value) {
        ElMessage.error('请输入密码');
        return;
    }
    loading.value = true;
    try {
        const response = await request.post('/auth/login', {
            password: password.value
        }) as LoginResponse;

        localStorage.setItem('authToken', response.token);
        ElMessage.success('登录成功');
        router.push({ name: 'UploadPhoto' });
    } catch (error) {
        ElMessage.error('登录失败，请检查密码后重试');
        console.error('Login error:', error);
    } finally {
        loading.value = false;
    }
};

</script>
<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f2f6fc; /* Element Plus 的浅灰色背景 */
}

.login-card {
    width: 400px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 20px;
    font-weight: bold;
}
</style>