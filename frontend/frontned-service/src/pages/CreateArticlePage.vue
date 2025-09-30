<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import NavBarView from '@/components/NavBarView.vue'

const router = useRouter()

const title = ref('')
const content = ref('')
const category = ref('')

const handleCreateArticle = async () => {
  const token = localStorage.getItem('auth_token')
  if (!token) {
    console.error('Authentication error: No token found.')
    router.push('/login')
    return
  }

  const payload = {
    title: title.value,
    content: content.value,
    category: category.value,
  }

  const config = {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }

  try {
    const response = await axios.post('http://localhost:8081/article', payload, config)
    console.log('Article created successfully!', response.data)

    const newArticleId = response.data.ID
    if (newArticleId) {
      router.push(`/article/${newArticleId}`)
    } else {
      router.push('/')
    }
  } catch (error) {
    console.error('Failed to create article:', error)
    if (axios.isAxiosError(error) && error.response) {
      if (error.response.status === 401) {
        localStorage.removeItem('auth_token')
        router.push('/login')
      } else {
        // The error message from the backend is in error.response.data
        console.error(`Error: ${error.response.data}`)
      }
    }
  }
}
</script>

<template>
  <NavBarView />
  <div class="page-layout">
    <div class="form-container">
      <h1 class="page-title">Create a New Article</h1>
      <form id="create-article-form" @submit.prevent="handleCreateArticle">
        <label for="title"><b>Title</b></label>
        <input
          v-model="title"
          type="text"
          placeholder="Enter article title"
          name="title"
          required
        />

        <label for="category"><b>Category</b></label>
        <input
          v-model="category"
          type="text"
          placeholder="e.g., Technology, Lifestyle"
          name="category"
          required
        />

        <label for="content"><b>Content</b></label>
        <textarea
          v-model="content"
          placeholder="Write your article content here..."
          name="content"
          rows="15"
          required
        ></textarea>

        <button type="submit">Publish Article</button>
      </form>
    </div>
  </div>
</template>

<style scoped>
/* Your existing styles will work perfectly here */
.page-layout {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 2rem;
  min-height: 90vh;
}
.form-container {
  width: 100%;
  max-width: 800px;
  display: flex;
  flex-direction: column;
}
.page-title {
  font-size: 40px;
  font-weight: bold;
  text-align: center;
  margin-bottom: 2rem;
}
#create-article-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
#create-article-form input,
#create-article-form textarea,
#create-article-form button {
  padding: 12px;
  border-radius: 5px;
  font-family: inherit;
  font-size: 1rem;
}
#create-article-form input,
#create-article-form textarea {
  border: 1px solid #d9d9d9;
  transition: all 100ms;
}
#create-article-form textarea {
  resize: vertical;
}
#create-article-form button {
  margin-top: 10px;
  background-color: black;
  color: white;
  border: none;
  cursor: pointer;
  font-weight: bold;
}
#create-article-form input:focus,
#create-article-form textarea:focus {
  outline: 2px solid rgb(93, 93, 93);
}
</style>
