<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import LayoutView from '@/components/LayoutView.vue'
import NavBarView from '@/components/NavBarView.vue'
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
</script>

<template>
  <NavBarView />
  <div id="landing-container">
    <LayoutView />
    <div id="flex-container">
      <h1 class="wellcome-text">
        WELLCOME <br />
        TO WEBAPP
      </h1>
      <div id="article-flex">
        <ArticleView v-for="article in articles" :key="article.ID" :article="article" />
      </div>
    </div>
  </div>
</template>

<style scoped>
#landing-container {
  display: flex;
}
#flex-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100vw;
  margin: 0;
  padding: 0;
}
#article-flex {
  display: flex;
  flex-direction: column;
}
.wellcome-text {
  font-size: 50px;
  font-weight: bold;
  text-align: center;
}
</style>
