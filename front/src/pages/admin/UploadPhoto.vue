<template>
    <div>
        <h1>上传照片</h1>
        <el-upload
            ref="uploadRef"
            v-model:file-list="filelist"
            class="upload-demo"
            action="#"
            :on-change="handleChange"
            :on-before-upload="handleBeforeUpload"
            :auto-upload="false"
            multiple
            accept=".jpg, .jpeg, .png"
        >
            <template #trigger>
                <el-button type="primary" class="select-file-btn">选择文件</el-button>
            </template>

            <el-button class="ml-3" type="success" @click="submitUpload">上传到服务器</el-button>

            <template #tip>
                <div class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
            </template>
        </el-upload>

    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage, type UploadInstance, type UploadProps, type UploadRawFile } from 'element-plus';
import ExifReader from 'exifreader';

const uploadRef = ref<UploadInstance>();
const filelist = ref([]);

const handleChange: UploadProps['onChange'] = async (uploadFile) => {
    if(uploadFile.raw){
        try {
            const tags = await ExifReader.load(uploadFile.raw);
            // 删除体积较大的缩略图数据，以便于在控制台查看
            if (tags.Thumbnail) {
                delete tags.Thumbnail;
            }
            
            if (Object.keys(tags).length > 0) {
                console.log(`照片[${uploadFile.name}]的EXIF信息:`, tags);
            } else {
                console.log(`照片[${uploadFile.name}]没有EXIF信息`);
            }
        } catch (error) {
            console.error(`读取照片[${uploadFile.name}]的EXIF信息失败:`, error);
        }
    }
}

const handleBeforeUpload: UploadProps['beforeUpload'] = (rawFile: UploadRawFile) => {
    const isJpgOrPng = rawFile.type === 'image/jpeg' || rawFile.type === 'image/png';
    if (!isJpgOrPng) {
        ElMessage.error('上传图片只能是 JPG 或 PNG 格式!')
        return false
    }
    return true
}

const submitUpload = () => {
  // 注意：你需要将 el-upload 的 action 属性设置为你真实的后端上传API地址
  // 这里我们调用 submit 方法来手动触发上传
  // 由于 action="#" 是一个无效地址，实际上传会失败，但会执行 before-upload 钩子
  // 在实际项目中，你需要一个真实的后端接口
  uploadRef.value!.submit()
  ElMessage.info("请查看浏览器控制台输出的元数据。实际上传需要后端接口支持。")
}
</script>

<style scoped>
.select-file-btn {
    margin-right: 10px;
}
</style>