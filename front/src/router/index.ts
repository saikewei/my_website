import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue';
import AdminPage from '../pages/AdminPage.vue';
import UploadPhoto from '../pages/admin/UploadPhoto.vue';

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
      children: [
        {
          path: 'upload_photo',
          component: UploadPhoto,
        }
      ]
    }
  ],
})

export default router
