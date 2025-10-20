<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <div class="header-container">
          <h1>后台</h1>
          <div class="header-button-container">
            <el-button type="text" style="color: white;" @click="handleOnLogout">
              退出登录
            </el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px">
          <div class="aside">
            <nav-admin />
          </div>
        </el-aside>
        <el-main>
          <div class="content">
            <router-view />
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import NavAdmin from '@/components/NavAdmin.vue';
import { RouterView, useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';

const router = useRouter();

const handleOnLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      localStorage.removeItem('authToken');
      router.push({ name: 'Login' });
    })
}
</script>

<style>
/* 添加全局样式，移除 body 的默认边距 */
body {
  margin: 0;
}
</style>

<style scoped>
.common-layout, .el-container {
  height: 100vh;
}

.el-header {
  padding: 0;
}

.aside {
  margin-top: 49px;
}

.content {
  margin-left: 10px;
  margin-top: 50px;
}

.header-container {
  display: flex;
  background-color: #4f5d64; /* 与导航菜单颜色匹配 */
  padding: 10px;
  padding-left: 35px;
  color: white;
}

.header-button-container {
  margin-left: auto;
  margin-top: 25px;
  margin-right: 20px;
}
</style>