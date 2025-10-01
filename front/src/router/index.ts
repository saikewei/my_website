import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue';
import UploadPhoto from '../pages/UploadPhoto.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: HomePage,
    },
    {
      path: '/admin',
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
