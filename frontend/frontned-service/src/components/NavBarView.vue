<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SearchBar from './SearchBar.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 1. A reactive variable to track login status
const isLoggedIn = ref(false)

// 2. Check for the token when the component is first mounted
onMounted(() => {
  const token = localStorage.getItem('auth_token')
  // !!token converts the string (or null) to a boolean
  isLoggedIn.value = !!token
})

const landingRedirect = () => {
  router.push('/')
}

const loginRedirect = () => {
  router.push('/login')
}

const newArticleRedirect = () => {
  // You would create this route in your router configuration
  router.push('/new-article')
}

// 3. Handle the logout process
const handleLogout = () => {
  // Remove the token from storage
  localStorage.removeItem('auth_token')
  // Update the UI to show the "Login" button again
  isLoggedIn.value = false
  // Redirect to the homepage
  router.push('/login')
}
</script>

<template>
  <div id="nav-bar-container">
    <p v-on:click="landingRedirect" class="nav-bar-button">SITE NAME</p>
    <SearchBar />

    <div class="nav-actions">
      <template v-if="isLoggedIn">
        <p v-on:click="newArticleRedirect" class="nav-bar-button">New Article</p>
        <p v-on:click="handleLogout" class="nav-bar-button">Logout</p>
      </template>
      <template v-else>
        <p v-on:click="loginRedirect" class="nav-bar-button">LOGIN</p>
      </template>
    </div>
  </div>
</template>

<style scoped>
#nav-bar-container {
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 3.5rem;
  width: 100%;
  background-color: white;
  margin: 0;
  padding: 0;
}

.nav-bar-button {
  margin: 0 1rem; /* Adjusted margin for multiple buttons */
  cursor: pointer;
}

/* Added a container for the action buttons for better layout */
.nav-actions {
  display: flex;
  align-items: center;
  margin-right: 1rem; /* Align to the right */
}
</style>
