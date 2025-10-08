<template>
    <h1>相册管理</h1>
    <el-button type="primary">添加相册</el-button>
    <div class="scrollbar-container">
        <el-scrollbar height="70vh">
            <div class="card-container">
                <el-card style="max-width: 240px" v-for="album in albums" :key="album.id">
                    <template #header>{{ album.title }}</template>
                    <el-image :src="album.cover_photo_url" style="width: 100%; height: 170px;" fit="cover">
                        <template #error>
                            <div class="image-viewer-slot image-slot">
                            <el-icon><icon-picture /></el-icon>
                            </div>
                        </template>
                    </el-image>
                    <p>{{ album.description }}</p>
                    <template #footer>
                        <div class="card-footer">
                            <span>创建时间: {{ album.created_at }}</span>
                        </div>
                    </template>
                    <el-button type="text" size="small">编辑</el-button>
                    <el-button type="text" size="small" style="color: red;">删除</el-button>
                </el-card>
            </div>
        </el-scrollbar>
    </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import request from '@/utils/request';
import { Picture as IconPicture } from '@element-plus/icons-vue'

type AlbumIdResponse = {
            albums_id: number[];
};

type Album = {
    id: number;
    title: string;
    description: string;
    cover_photo_id?: string;
    created_at: string;
    updated_at: string;
    cover_photo_url?: string;
};

const albumIds = ref<number[]>([]);
const albums = ref<Album[]>([]);
const fetchAlbumIds = async () => {
    try {
        
        const response = await request.get('/photo/albums-id') as AlbumIdResponse;
        albumIds.value = response.albums_id;
        console.log('Fetched album IDs:', albumIds.value);

        const albumDetailPromises = albumIds.value.map(async (id) => {
            const detailResponse = await request.get(`/photo/album/${id}`) as { album: Album };
            detailResponse.album.created_at = new Date(detailResponse.album.created_at).toLocaleDateString();
            detailResponse.album.updated_at = new Date(detailResponse.album.updated_at).toLocaleDateString();
            detailResponse.album.cover_photo_url = detailResponse.album.cover_photo_id
                ? `http://localhost:9000/api/photo/${detailResponse.album.cover_photo_id}/thumbnail`
                : undefined;
            return detailResponse.album;
        })

        const albumDetails = await Promise.all(albumDetailPromises);
        albums.value = albumDetails;
        console.log('Fetched album details:', albums.value);
    } catch (error) {
        console.error('Failed to fetch album IDs:', error);
    }
};

onMounted(() => {
    fetchAlbumIds();
});

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

.card-footer {
    text-align: right;
    font-size: 12px;
    color: #999;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  font-size: 30px;
  height: 170px;
  background: #fff;
}
</style>