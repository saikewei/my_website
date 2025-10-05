<template>
    <div class="header-container">
      <h1>后台</h1>
    </div>
    <div class="page-container">
        <nav-admin />
        <div class="content">
            <router-view />
        </div>
    </div>
</template>

<script setup lang="ts">
import NavAdmin from '@/components/NavAdmin.vue';
import { RouterView } from 'vue-router';
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: row;
}

.content {
  margin-left: 20px;
}

.header-container {
  background-color: #4f5d64; /* 与导航菜单颜色匹配 */
  padding: 20px;
  color: white;
  margin-bottom: 1px;
}
</style>