<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AdminUserView from '@/components/AdminUserView.vue'
import AdminArticlesView from '@/components/AdminArticlesView.vue'

const router = useRouter()
const loginRedirect = () => {
  console.log('Redirecting to login page...')
  router.push('/login') // Navigate to the /login path
}

const activeSection = ref('Users')
</script>

<template>
  <div id="nav-bar-container">
    <p @click="loginRedirect" class="nav-bar-button">LOGIN</p>
  </div>

  <div id="admin-container">
    <h1 class="welcome-text">
      WELCOME <br />
      TO ADMIN PANEL
    </h1>

    <div id="admin-buttons">
      <button @click="activeSection = 'Users'">Users</button>
      <button @click="activeSection = 'Articles'">Articles</button>
    </div>

    <AdminUserView v-if="activeSection === 'Users'" />
    <AdminArticlesView v-if="activeSection === 'Articles'" />
  </div>
</template>
<style scoped>
#nav-bar-container {
  display: flex;
  flex-direction: row-reverse;
  align-items: center;
  height: 3.5rem;
  width: 100%;
  background-color: white;
  margin: 0;
  padding: 0;
}
.nav-bar-button {
  margin: 0 2rem;
  cursor: pointer;
}
.welcome-text {
  text-align: center;
}
#admin-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

#admin-buttons {
  display: flex;
  gap: 10px;
}

#admin-buttons button {
  border: 1px solid grey;
  background-color: white;
  color: black;

  padding: 8px 16px;
  cursor: pointer;
  transition:
    background-color 0.3s,
    color 0.3s;
}

#admin-buttons button:hover {
  background-color: black;
  color: white;
}
</style>
