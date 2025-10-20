import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue';
import AdminPage from '../pages/AdminPage.vue';
import UploadPhoto from '../pages/admin/UploadPhoto.vue';
import AlbumManagement from '@/pages/admin/AlbumManagement.vue';
import PhotoManagement from '@/pages/admin/PhotoManagement.vue';
import ChangePassword from '@/pages/admin/ChangePassword.vue';
import LoginPage from '@/pages/LoginPage.vue';
import { jwtDecode } from 'jwt-decode';

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
      meta: { requiresAuth: true },
      children: [
        {
          path: 'upload-photo',
          name: 'UploadPhoto',
          component: UploadPhoto,
        },
        {
          path: 'album',
          component: AlbumManagement,
        },
        {
          path: 'photo-management',
          redirect: '/admin/photo-management/1',
        },
        {
          path: 'photo-management/:page',
          name: 'PhotoManagement',
          component: PhotoManagement,
        },
        {
          path: 'change-password',
          component: ChangePassword,
        }
      ]
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginPage,
    }
  ],
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('authToken');
  let isAuthenticated = false;

  if (token) {
    try {
      const decodedToken = jwtDecode(token);

      if (decodedToken.exp && decodedToken.exp * 1000 > Date.now()) {
        isAuthenticated = true;
      } else {
        // Token 已过期
        localStorage.removeItem('authToken'); // 清除过期的 token
      }
    } catch (error) {
      console.error('无效的 token:', error);
      localStorage.removeItem('authToken'); // 清除无效的 token
    }
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);

  if (requiresAuth && !isAuthenticated) {
    // 需要登录但未通过验证，跳转到登录页
    next({ name: 'Login' });
  } else if (to.name === 'Login' && isAuthenticated) {
    // 已登录但访问登录页，重定向到后台首页
    next({ path: '/admin' });
  } else {
    // 其他情况，直接放行
    next();
  }
})

export default router
