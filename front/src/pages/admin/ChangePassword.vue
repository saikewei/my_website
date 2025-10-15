<template>
    <h1>修改密码</h1>
    <div class="form-container">
        <el-form label-width="100px">
            <el-form-item label="旧密码">
                <el-input type="password" v-model="oldPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="新密码">
                <el-input type="password" v-model="newPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="确认新密码">
                <el-input type="password" v-model="confirmNewPassword" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
            </el-form-item>
        </el-form>
    </div>
</template>
<script setup lang="ts">

</script>
<style scoped>
.form-container {
    max-width: 400px;
    margin: 0 auto; 
    margin-top: 100px;
}
</style>