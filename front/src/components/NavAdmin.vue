<template>
  <el-row class="tac">
    <el-col :span="12">
      <el-menu
        active-text-color="#ffd04b"
        background-color="#545c64"
        class="el-menu-vertical-demo"
        :default-active="activeMenu"
        text-color="#fff"
        @open="handleOpen"
        @close="handleClose"
        router
      >
        <el-sub-menu index="1">
          <template #title>
            <el-icon><icon-menu /></el-icon>
            <span>照片管理</span>
          </template>
          <el-menu-item-group title="照片">
            <el-menu-item index="/admin/upload-photo">上传照片</el-menu-item>
            <el-menu-item index="/admin/photo-management">编辑照片信息</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="相册">
            <el-menu-item index="/admin/album">相册管理</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>
        <el-menu-item index="2">
          <el-icon><icon-menu /></el-icon>
          <span>修改密码</span>
        </el-menu-item>
        <el-menu-item index="3" disabled>
          <el-icon><document /></el-icon>
          <span>Navigator Three</span>
        </el-menu-item>
        <el-menu-item index="4">
          <el-icon><setting /></el-icon>
          <span>Navigator Four</span>
        </el-menu-item>
      </el-menu>
    </el-col>
  </el-row>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router'
import {
  Document,
  Menu as IconMenu,
  Setting,
} from '@element-plus/icons-vue'

import { computed } from 'vue'

const route = useRoute()

const activeMenu = computed(() => {
  const { path } = route;
  // 如果当前路径以 /admin/photo-management 开头，
  // 就返回 /admin/photo-management 作为激活的 index
  if (path.startsWith('/admin/photo-management')) {
    return '/admin/photo-management';
  }
  // 否则，直接返回当前路径
  return path;
});

const handleOpen = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
const handleClose = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}
</script>

<style scoped>
.el-menu-vertical-demo {
  width: 200px;
  height: 450px;
}
</style>
