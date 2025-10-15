<template>
    <h1>修改密码</h1>
    <div class="form-container">
        <el-form label-width="100px" :model="form" :rules="rules" ref="formRef">
            <el-form-item label="旧密码" prop="oldPassword">
                <el-input type="password" v-model="form.oldPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
                <el-input type="password" v-model="form.newPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirmNewPassword">
                <el-input type="password" v-model="form.confirmNewPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
            </el-form-item>
        </el-form>
    </div>
</template>
<script setup lang="ts">
import { ref, reactive } from 'vue';
import { ElMessage, type FormInstance, type FormRules} from 'element-plus';
import request from '@/utils/request';
import { AxiosError } from 'axios';

const formRef = ref<FormInstance>();
const form = reactive({
    oldPassword: '',
    newPassword: '',
    confirmNewPassword: ''
});
const rules = reactive<FormRules>({
    oldPassword: [
        { required: true, message: '请输入旧密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
    ],
    newPassword: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
    ],
    confirmNewPassword: [
        { required: true, message: '请确认新密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能小于6位', trigger: 'blur' },
        { validator: (rule, value) => {
            if (value !== form.newPassword) {
                return new Error('两次输入的密码不一致');
            }
            return true;
        }, trigger: 'blur' }
    ]
});

const handleChangePassword = async () => {
    try {
        if (!formRef.value) {
            console.error('表单引用未定义');
            throw new Error('表单引用未定义');
        }
        await formRef.value.validate((valid, fields) => {
            if (!valid) {
                console.log('表单验证失败:', fields);
                throw new Error('表单验证失败');
            }
        })

        await request.put('/auth/change-password', {
            old_password: form.oldPassword,
            new_password: form.newPassword
        });

        ElMessage.success('密码修改成功，请重新登录');
        form.oldPassword = '';
        form.newPassword = '';
        form.confirmNewPassword = '';
        // 清除本地存储的token
    } catch (error) {
        console.error('修改密码失败:', error);
        if(error instanceof AxiosError && error.response?.status === 400) {
            ElMessage.error('旧密码错误，请重试');
        }
        else{
            ElMessage.error('修改密码失败，请重试');
    
        }
    }
};

</script>
<style scoped>
.form-container {
    max-width: 400px;
    margin: 0 auto; 
    margin-top: 100px;
}
</style>