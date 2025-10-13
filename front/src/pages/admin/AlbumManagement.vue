<template>
    <h1>相册管理</h1>
    <el-button type="primary" @click="onClickAddAlbum">添加相册</el-button>
    <div class="scrollbar-container">
        <el-scrollbar height="70vh">
            <div class="card-container">
                <el-card style="width: 240px" v-for="album in albums" :key="album.id">
                    <template #header>{{ album.title }}</template>
                    <el-image :src="album.cover_photo_url" style="width: 100%; height: 170px;" fit="cover">
                        <template #error>
                            <div class="image-viewer-slot image-slot">
                            <el-icon><icon-picture /></el-icon>
                            </div>
                        </template>
                    </el-image>
                    <p class="card-description">{{ album.description }}</p>
                    <template #footer>
                        <div class="card-footer">
                            <span>创建时间: {{ album.created_at }}</span>
                        </div>
                    </template>
                    <el-button type="text" size="small" @click="onClickEditAlbum(album.id)">编辑</el-button>
                    <el-button type="text" size="small" style="color: red;" @click="onClickDeleteAlbum(album.id)">删除</el-button>
                </el-card>
            </div>
        </el-scrollbar>
    </div>

    <el-dialog v-model="dialogFormVisible" :title="dialogFormTitle" width="500">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="auto">
        <el-form-item label="相册名称" :label-width="formLabelWidth" prop="name">
            <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" :label-width="formLabelWidth" prop="description">
            <el-input
            type="textarea"
            v-model="form.description"
            autocomplete="off"
            />
        </el-form-item>
        </el-form>
        <template #footer>
        <div class="dialog-footer">
            <el-button type="primary" @click="onClickSubmit(formRef)">
            确定
            </el-button>
            <el-button @click="dialogFormVisible = false">取消</el-button>
        </div>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { onBeforeMount, ref, reactive } from 'vue';
import request from '@/utils/request';
import { Picture as IconPicture } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage, ElMessageBox } from 'element-plus'

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

interface CreateForm {
    name: string;
    description: string;
}

const formRef = ref<FormInstance>();

const form = reactive<CreateForm>({
    name: '',
    description: '',
});

const currentAlbumId = ref<number | null>(null);

const rules = reactive<FormRules<CreateForm>>({
    name: [
        { required: true, message: '请输入相册名称', trigger: 'blur' },
        { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' },
    ],
    description: [
        { max: 200, message: '长度在 0 到 100 个字符', trigger: 'blur' },
    ],
});


const dialogFormVisible = ref(false)
const dialogFormTitle = ref('创建相册')
const formLabelWidth = '140px'

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

const onClickAddAlbum = () => {
    dialogFormTitle.value = '创建相册';
    dialogFormVisible.value = true;
};

const onClickSubmit = async (formEl: FormInstance | undefined) => {
    try {
        if (!formEl) return;
        await formEl.validate((valid, fileds) => {
            if (!valid) {
                console.log('error submit!', fileds);
                throw new Error('表单验证失败');
            }
        });

        if (dialogFormTitle.value === '创建相册') {
            const response = await request.post('/photo/create-album', {
                title: form.name,
                description: form.description === '' ? "无描述" : form.description,
            });
            ElMessage({
                message: '相册创建成功',
                type: 'success',
            })
            console.log('Album created successfully:', response);
        }
        else if (dialogFormTitle.value === '编辑相册') {
            const response = await request.put('/photo/edit-album', {
                id: currentAlbumId.value,
                title: form.name,
                description: form.description === '' ? "无描述" : form.description,
            });
            ElMessage({
                message: '相册编辑成功',
                type: 'success',
            })
            console.log('Album edited successfully:', response);
            currentAlbumId.value = null;
        }
        dialogFormVisible.value = false;
        
        form.name = '';
        form.description = '';
        fetchAlbumIds();
    } catch (error) {
        dialogFormVisible.value = false;
        currentAlbumId.value = null;
        
        form.name = '';
        form.description = '';
        ElMessage({
            message: '相册创建失败',
            type: 'error',
        })
        console.error('Failed to create album:', error);
    }
};

const onClickEditAlbum = (albumId: number) => {
    dialogFormTitle.value = '编辑相册';
    dialogFormVisible.value = true;
    currentAlbumId.value = albumId;
    console.log('Edit album with ID:', albumId);

    form.name = albums.value.find(album => album.id === albumId)?.title || '';
    form.description = albums.value.find(album => album.id === albumId)?.description || '';
};

const onClickDeleteAlbum = (albumId: number) => {
    ElMessageBox.confirm('确定要删除该相册吗？此操作不可逆！', '警告', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        
        type: 'warning',
    }).then(() => {
        deleteAlbum(albumId);
    }).catch(() => {
        ElMessage({
            type: 'info',
            message: '已取消删除'
        });
    });
}

const deleteAlbum = async (albumId: number) => {
    try {
        const response = await request.delete(`/photo/album/${albumId}`);
        ElMessage({
            message: '相册删除成功',
            type: 'success',
        })
        console.log('Album deleted successfully:', response);
        fetchAlbumIds();
    } catch (error) {
        ElMessage({
            message: '相册删除失败',
            type: 'error',
        })
        console.error('Failed to delete album:', error);
    }
};

onBeforeMount(() => {
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

.card-description {
    /* 1. 设置一个固定的高度，确保所有卡片的这部分高度都相同 */
    height: 63px; /* 这个高度大约可以容纳3行文字，您可以根据需要微调 */
    line-height: 1.5; /* 设置行高，便于计算 */

    /* 2. 以下是实现多行文本溢出显示省略号的关键样式 */
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all; /* 允许在单词内换行 */
}
</style>