<template>
    <h1>Admin Page</h1>
    <div class="page-container">
        <nav-admin />
        <div class="content">
            <h1>Upload Photo Page</h1>
            This is a page to upload photo.
        </div>
    </div>
</template>

<script setup lang="ts">
import NavAdmin from '@/components/NavAdmin.vue';
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