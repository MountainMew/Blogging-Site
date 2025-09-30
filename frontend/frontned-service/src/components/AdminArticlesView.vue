<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import type { Article } from '@/models/article.interface'
import ArticleView from '@/components/ArticleView.vue'

const articles = ref<Article[]>([])

onMounted(async () => {
  try {
    const response = await axios.get('http://localhost:8081/blog')

    console.log('Raw data from backend:', response.data)

    articles.value = response.data
  } catch (error) {
    console.error('There was an error fetching the articles:', error)
  }
})

const handleDelete = (id: string) => {
  console.log(id)
}
</script>

<template>
  <div id="admin-user-container">
    <div id="filter-user-bar">
      <label for="cars">Search Article By:</label>

      <select name="users" id="users">
        <option value="Id">Id</option>
        <option value="Username">Publisher</option>
      </select>
      <form>
        <input type="text" id="serach-input" name="serach-input" placeholder="Search" />
      </form>
      <div id="article-flex">
        <div v-for="article in articles" :key="article.ID" class="article-wrapper">
          <div>
            <ArticleView :article="article" />
          </div>

          <button @click="handleDelete(article.ID)">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.article-wrapper {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}
</style>
