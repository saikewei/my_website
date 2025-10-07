<template>
    <h1>相册管理</h1>
    <el-button type="primary">添加相册</el-button>
    <div class="scrollbar-container">
        <el-scrollbar height="600px">
            <div class="card-container">
                <el-card style="max-width: 400px" v-for="index in 10" :key="index">
                    <template #header>Yummy hamburger</template>
                    <img
                    src="https://shadow.elemecdn.com/app/element/hamburger.9cf7b091-55e9-11e9-a976-7f4d0b07eef6.png"
                    style="width: 100%"
                    />
                </el-card>
            </div>
        </el-scrollbar>
    </div>
</template>
<script setup lang="ts">

</script>
<style scoped>
.card-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    margin-top: 20px;
}

.scrollbar-container {
    margin-top: 20px;
    width: 90%;
}
</style>