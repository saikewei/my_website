<template>
    <h2>全部照片</h2>
    <el-scrollbar style="height: 70vh; width: 100%;">
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
                        @click="selectPhoto(photo.id)"
                    ></el-image>
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
    <el-dialog
        title="编辑照片信息"
        v-model="dialogFormVisible"
        width="700px"
    >
        <div class="dialog-body-wrapper">
            <div class="dialog-image-container">
                <el-image
                    class="dialog-image"
                    :src="dialogImageUrl"
                    :alt="photos.find(p => p.id === selectedPhotoId)?.title"
                    fit="scale-down"
                >
                    <template #placeholder>
                        <div class="image-slot">
                            <el-icon><icon-picture /></el-icon>
                        </div>
                    </template>
                </el-image>
            </div>
            <div v-if="selectedPhoto" class="photo-details-grid">
                <!-- 2. 将每个详情项包裹起来，并添加标签 -->
                <div class="detail-item">
                    <span class="detail-label">相机:</span>
                    <span class="detail-value">{{ selectedPhoto.camera || 'N/A' }}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">镜头:</span>
                    <span class="detail-value">{{ selectedPhoto.lens || 'N/A' }}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">拍摄参数:</span>
                    <span class="detail-value">{{ formatShootingParams(selectedPhoto) }}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-label">尺寸:</span>
                    <span class="detail-value">{{ selectedPhoto.width }} x {{ selectedPhoto.height }}</span>
                </div>
                 <div class="detail-item">
                    <span class="detail-label">拍摄时间:</span>
                    <span class="detail-value">{{ formatShotAt(selectedPhoto.shot_at) }}</span>
                </div>
            </div>
            <el-form
                ref="formRef"
                :model="form"
                :rules="rules"
                label-width="100px"
            >
                <el-form-item label="标题" prop="title">
                    <el-input v-model="form.title"></el-input>
                </el-form-item>
                <el-form-item label="描述" prop="description">
                    <el-input 
                        type="textarea" 
                        :rows="3" 
                        v-model="form.description"
                    ></el-input>
                </el-form-item>
                <el-form-item label="是否精选" prop="is_featured">
                    <el-switch v-model="form.is_featured"></el-switch>
                </el-form-item>
                <el-form-item label="标签" prop="tags">
                    <el-input-tag
                        v-model="form.tags"
                        placeholder="请输入标签"
                        aria-label="Please click the Enter key after input"
                    />
                </el-form-item>
                <el-form-item label="相册" prop="album">
                    <el-select 
                        v-model="form.album"
                        placeholder="选择相册" 
                        value-key="id"
                    >
                        <el-option
                            v-for="album in albumOptions"
                            :key="album.id"
                            :label="album.title"
                            :value="album"
                        ></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
        </div>
        <template #footer>
            <el-button type="danger" style="margin-right: 450px" @click="handleOnDeleteSubmit">删 除</el-button>
            <el-button type="primary" @click="handleOnEditSubmit">确 定</el-button>
            <el-button @click="dialogFormVisible = false">取 消</el-button>    
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { ref , watch, reactive, computed, onBeforeMount } from 'vue';
import request from '@/utils/request';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus';
import { Picture as IconPicture } from '@element-plus/icons-vue'

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

interface PhotoEditForm {
    title: string;
    description: string;
    is_featured: boolean;
    tags: string[] | null; // 前端使用数组
    album: AlbumOption;
}

interface PhotoEditRequest {
    id: number;
    title: string;
    description: string | null;
    is_featured: boolean;
    tags: string[] | null;
    album_id: number;
}

interface AlbumResponse {
    albums: AlbumOption[];
}

interface AlbumOption {
    id: number;
    title: string;
    description: string | null;
    created_at: string;
    updated_at: string;
}

const route = useRoute();
const router = useRouter();

const total = ref(500);
const pageSize = ref(40);
const currentPage = ref(1);
const loading = ref(false);

const photos = ref<Photo[]>([]);

const selectedPhotoId = ref<number | null>(null);

const formRef = ref<FormInstance>();
const form = reactive<PhotoEditForm>({
    title: '',
    description: '',
    is_featured: false,
    tags: null,
    album: { id: 0, title: '', description: null, created_at: '', updated_at: '' },
});
const rules = reactive<FormRules<PhotoEditForm>>({
    title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 1, max: 50, message: '标题长度在 1 到 50 个字符之间', trigger: 'blur' },
    ],
    description: [
        { max: 300, message: '描述长度不能超过 300 个字符', trigger: 'blur' },
    ],
    tags: [
        { type: 'array', message: '标签必须是一个数组', trigger: 'change' },
    ],
});
const albumOptions = ref<AlbumOption[]>([]);

const selectedPhoto = computed(() => {
    if (selectedPhotoId.value === null) {
        return null;
    }
    return photos.value.find(p => p.id === selectedPhotoId.value);
});

const dialogFormVisible = ref(false);
const dialogImageUrl = computed(() => {
    if (selectedPhotoId.value !== null) {
        return `${import.meta.env.VITE_API_BASE_URL}/photo/${selectedPhotoId.value}/thumbnail?size=600`;
    }
    return '';
})

const formatShotAt = (dateString: string | null) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleString();
};

const formatShootingParams = (photo: Photo) => {
    const params = [photo.focal_length, photo.aperture, photo.shutter_speed, photo.iso ? `ISO ${photo.iso}` : null];
    return params.filter(p => p).join(' | ') || 'N/A';
}

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
    // console.log('Page changed to:', currentPage.value);
    router.push({ 
        name: 'PhotoManagement',
        params: { page: currentPage.value } 
    });
}

const selectPhoto = (photoId: number) => {
    selectedPhotoId.value = photoId;

    form.title = selectedPhoto.value?.title || '';
    form.description = selectedPhoto.value?.description || '';
    form.is_featured = selectedPhoto.value?.is_featured || false;
    form.tags = selectedPhoto.value?.tags ? selectedPhoto.value.tags.split(',').map(tag => tag.trim()).filter(tag => tag) : [];
    const album = albumOptions.value.find(a => a.id === selectedPhoto.value?.album_id);
    form.album = album || { id: 0, title: '无专辑', description: null, created_at: '', updated_at: '' };

    dialogFormVisible.value = true;
    console.log('Selected photo ID:', selectedPhotoId.value);
}

const fetchAlbums = async () => {
    try {
        const response = await request.get('/photo/album/details') as AlbumResponse;
        albumOptions.value = response.albums;
        albumOptions.value.unshift({ id: 0, title: '无专辑', description: null, created_at: '', updated_at: '' });

        console.log('Fetched albums:', albumOptions.value);
    } catch (error) {
        console.error('Failed to fetch albums:', error);
    }
};

const handleOnEditSubmit = async () => {
    if (!formRef.value) return;

    try {
        await formRef.value.validate((valid, fileds) => {
            if (!valid) {
                console.log('error submit!', fileds);
                throw new Error('表单验证失败');
            }
        });

        const payload: PhotoEditRequest = {
            id: selectedPhotoId.value as number,
            title: form.title,
            description: form.description || null,
            is_featured: form.is_featured,
            tags: form.tags && form.tags.length > 0 ? form.tags : null,
            album_id: form.album.id,
        };
        dialogFormVisible.value = false;

        await request.put(`/photo/edit/photo`, payload);

        ElMessage.success('照片信息更新成功!');
    } catch (error) {
        console.error('Failed to submit edit form:', error);
        ElMessage.error('照片信息更新失败!');
        dialogFormVisible.value = false;
        return;
    } finally {
        fetchPhotos();
    }
}

const handleOnDeleteSubmit = async () => {
    ElMessageBox.confirm('确定要删除该照片吗？此操作不可逆！', '警告', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        
        type: 'warning',
    }).then(() => {
        if (selectedPhotoId.value === null) {
            ElMessage.error('未选择照片');
            return;
        }

        request.delete(`/photo/${selectedPhotoId.value}`).then(() => {
            ElMessage.success('照片删除成功');
            dialogFormVisible.value = false;
            fetchPhotos();
        }).catch((error) => {
            console.error('Failed to delete photo:', error);
            ElMessage.error('照片删除失败');
        });
    }).catch(() => {
        ElMessage.info('已取消删除');
    });
}

onBeforeMount(() => {
    fetchAlbums();
});

watch(
    () => route.fullPath,
    () => {
        const pageFromUrl = parseInt(route.params.page as string) || 1;

        currentPage.value = pageFromUrl;
        fetchPhotos();
    },
    { immediate: true}
);

</script>
<style scoped>
.photo-details-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 12px 24px; /* 垂直间距12px，水平间距24px */
    padding: 12px;
    background-color: #f9f9f9;
    border-radius: 4px;
    margin-top: 16px;
}

.detail-item {
    display: flex;
    align-items: center;
    font-size: 14px;
}

.detail-label {
    color: #606266;
    margin-right: 8px;
    font-weight: bold;
}

.detail-value {
    color: #303133;
}


.photo-container {
    display: flex;
    flex-wrap: wrap;
    gap: 8px; /* 图片之间的间距 */
    padding: 16px;
}

.photo-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
}

.photo-item:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
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

.dialog-image {
    width: 55%;
    max-height: 600px; /* 可以根据需要调整最大高度 */
    border-radius: 4px;
}

.dialog-image-container {
    /* 使用 Flexbox 居中，这是最现代和可靠的方式 */
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 24px; /* 替代原来图片上的 margin-bottom */
}

.image-slot {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background: var(--el-fill-color-light); /* 使用 Element Plus 的主题色 */
    color: var(--el-text-color-secondary);
}

.image-slot .el-icon {
    font-size: 30px;
}

.dialog-body-wrapper {
    height: 65vh; /* 或者一个固定的像素值，如 600px */
    overflow-y: auto; /* 这是关键，当内容超出时，出现垂直滚动条 */
    padding: 0 10px; /* 为滚动条留出一点空间，避免内容紧贴边缘 */
}
</style>