import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue';
import AdminPage from '../pages/AdminPage.vue';
import UploadPhoto from '../pages/admin/UploadPhoto.vue';
import AlbumManagement from '@/pages/admin/AlbumManagement.vue';
import PhotoManagement from '@/pages/admin/PhotoManagement.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: HomePage,
    },
    {
      path: '/admin',
      component: AdminPage,
      redirect: '/admin/upload-photo',
      children: [
        {
          path: 'upload-photo',
          component: UploadPhoto,
        },
        {
          path: 'album',
          component: AlbumManagement,
        },
        {
          path: 'photo-management',
          component: PhotoManagement,
        }
      ]
    }
  ],
})

export default router
