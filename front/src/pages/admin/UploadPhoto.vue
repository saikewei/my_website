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
import { ElMessage, type UploadInstance, type UploadProps, type UploadRawFile, type UploadFile } from 'element-plus';
import ExifReader from 'exifreader';
import request from '@/utils/request';

interface UploadFileWithMeta extends UploadFile {
    metaData?: CleanExifData;
}

interface CleanExifData {
    width?: number;
    height?: number;
    is_featured: boolean;
    tags: string[];
    shot_at: string | null;
    focal_length?: string;
    aperture?: string;
    shutter_speed?: string;
    iso?: string;
    exposure_bias?: string;
    flash_fired: boolean;
    gps_latitude?: number;
    gps_longitude?: number;
    camera?:string;
    lens?:string;
}

const uploadRef = ref<UploadInstance>();
const filelist = ref<UploadFileWithMeta[]>([]);

const extractCleanExif = (tags: ExifReader.Tags): CleanExifData => {
    const getDesc = (tagName: string): string | undefined => tags[tagName]?.description;
    const getValue = (tagName: string): unknown => tags[tagName]?.value;

    // 处理拍摄时间，将其转换为 ISO 8601 格式
    let shotAt: string | null = null;
    const dateTimeOriginal = getDesc('DateTimeOriginal');
    if (dateTimeOriginal) {
        let timeStr = dateTimeOriginal.replace(':', '-').replace(':', '-');
        timeStr = timeStr.replace(' ', 'T'); // 将日期和时间之间的空格替换为 'T'
        
        const offset = getDesc('OffsetTimeOriginal') || '+00:00';
        shotAt = `${timeStr}${offset}`;
    }

    const isoValue = getValue('ISOSpeedRatings') as number | undefined;

    return {
        width: getValue('Image Width') as number | undefined,
        height: getValue('Image Height') as number | undefined,
        is_featured: false, // 默认值，后续可在UI中修改
        tags: [], // 默认值，后续可在UI中添加
        shot_at: shotAt,
        focal_length: getDesc('FocalLength'),
        aperture: getDesc('FNumber'),
        shutter_speed: getDesc('ExposureTime'),
        iso: isoValue !== undefined ? String(isoValue) : undefined,
        exposure_bias: getDesc('ExposureBiasValue') ? `${getDesc('ExposureBiasValue')} EV` : undefined,
        flash_fired: getDesc('Flash')?.includes('did not fire') === false,
        gps_latitude: getValue('GPSLatitude') as number | undefined,
        gps_longitude: getValue('GPSLongitude') as number | undefined,
        camera: getDesc('Model'),
        lens: getDesc('LensModel'),
    };
};


const handleChange: UploadProps['onChange'] = async (uploadFile) => {
    if(uploadFile.raw && uploadFile.status === 'ready'){
        try {
            const rawTags = await ExifReader.load(uploadFile.raw);
            // 删除体积较大的缩略图数据，以便于在控制台查看
            delete rawTags.Thumbnail;

            const tags = extractCleanExif(rawTags);

            (uploadFile as UploadFileWithMeta).metaData = tags;
            
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
        ElMessage.error('上传图片只能是 JPG 或 PNG 格式!');
        return false;
    }
    return true
}

const submitUpload = () => {
    const filesToUpload = filelist.value;

    if (!filesToUpload || filesToUpload.length === 0) {
        ElMessage.warning("没有选择任何文件进行上传。");
        return;
    }

    const uploadPromises = filesToUpload.map(file => {
        if (!file.raw) {
            return Promise.resolve();
        }

        const formData = new FormData();
        formData.append('file', file.raw);

        const meta = (file as UploadFileWithMeta).metaData || {};
        formData.append('meta', JSON.stringify(meta));

        return request({
            url: '/photo/upload',
            method: 'POST',
            data: formData,
            headers: { 'Content-Type': 'multipart/form-data' },
        });
    });
    
    Promise.all(uploadPromises)
        .then(() => {
            ElMessage.success("所有文件上传成功！");
            uploadRef.value?.clearFiles();
        })
        .catch((error) => {
            console.error("上传过程中出现错误:", error);
            ElMessage.error("上传过程中出现错误，请稍后重试。");
        });
}
</script>

<style scoped>
.select-file-btn {
    margin-right: 10px;
}
</style>