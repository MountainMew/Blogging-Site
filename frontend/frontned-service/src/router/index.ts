import AboutPage from '@/pages/AboutPage.vue'
import AdminPage from '@/pages/AdminPage.vue'
import ArticlePage from '@/pages/ArticlePage.vue'
import LandingPage from '@/pages/LandingPage.vue'
import LoginPage from '@/pages/LoginPage.vue'
import RegisterPage from '@/pages/RegisterPage.vue'
import CreateArticlePage from '@/pages/CreateArticlePage.vue'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: LandingPage },
  { path: '/register', component: RegisterPage },
  { path: '/login', component: LoginPage },
  { path: '/admin', component: AdminPage },
  {
    path: '/about/:Id',
    name: 'UserDetail',
    component: AboutPage,
  },
  {
    path: '/article/:Id',
    name: 'ArticleDetail',
    component: ArticlePage,
  },
  {
    path: '/new-article',
    name: 'CreateArticle',
    component: CreateArticlePage,
    meta: { requiresAuth: true }, // 3. Add meta field for authentication
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})

export default router
