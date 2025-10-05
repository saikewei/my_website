<template>
    <h1>后台</h1>
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
</style>