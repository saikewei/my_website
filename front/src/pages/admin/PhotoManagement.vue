<template>
    <h2>全部照片</h2>
    <el-scrollbar style="height: 75vh; width: 100%;">
        <div v-if="loading">
            <p>正在加载照片...</p>
        </div>
        <div v-else-if="photos.length === 0">
            <p>没有照片可显示。</p>
        </div>
        <div v-else>
            <div class="photo-container">
                <div v-for="photo in photos" :key="photo.id" class="photo-item">
                    <el-image
                        style="width: 150px; height: 150px"
                        :src="photo.thumbnail_base64" 
                        :alt="photo.title"
                        fit="cover" 
                        lazy
                    ></el-image>
                    <!-- <span class="photo-title">{{ photo.title }}</span> -->
                </div>
            </div>
        </div>
    </el-scrollbar>
    <el-pagination 
        layout="prev, pager, next" 
        v-model:current-page="currentPage"
        :total="total" 
        :page-size="pageSize"
        @current-change="handleOnPageChange"
    />
</template>
<script setup lang="ts">
import { ref , onBeforeMount } from 'vue';
import request from '@/utils/request';

interface Photo {
    id: number;
    album_id: number | null;
    title: string;
    description: string | null;
    file_path: string;
    file_name: string;
    file_size: number;
    width: number;
    height: number;
    is_featured: boolean;
    shot_at: string | null;
    created_at: string;
    updated_at: string;
    camera: string | null;
    lens: string | null;
    focal_length: string | null;
    aperture: string | null;
    shutter_speed: string | null;
    iso: string | null;
    exposure_bias: string | null;
    flash_fired: boolean | null;
    gps_latitude: number | null;
    gps_longitude: number | null;
    tags: string | null; // 从后端视图看，这是一个逗号分隔的字符串
    thumbnail_base64: string; // Service层添加的字段
}

interface PhotoApiResponse {
    total: number;
    photos: Photo[];
}

const total = ref(500);
const pageSize = ref(40);
const currentPage = ref(1);
const loading = ref(false);

const photos = ref<Photo[]>([]);

const fetchPhotos = async () => {
    try {
        loading.value = true;
        const response = await request.get('/photo/page', {
            params: {
                'page-num': currentPage.value,
                'page-size': pageSize.value,
            },
        }) as PhotoApiResponse;
        total.value = response.total;
        photos.value = response.photos;
        loading.value = false;
        console.log('Fetched photos:', photos.value);
    } catch (error) {
        console.error('Failed to fetch photos:', error);
    }
}

const handleOnPageChange = () => {
    fetchPhotos();
}

onBeforeMount(() => {
    fetchPhotos();
});

</script>
<style scoped>
.photo-container {
    display: flex;
    flex-wrap: wrap;
    gap: 16px; /* 图片之间的间距 */
    padding: 16px;
}

.photo-item {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.photo-title {
    margin-top: 8px;
    font-size: 14px;
    width: 150px;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>