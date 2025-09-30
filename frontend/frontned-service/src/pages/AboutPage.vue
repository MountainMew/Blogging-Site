<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import LayoutView from '@/components/LayoutView.vue'
import NavBarView from '@/components/NavBarView.vue'
import type { Article } from '@/models/article.interface'
import ArticleView from '@/components/ArticleView.vue'
import AboutArticleView from '@/components/AboutArticleView.vue'

// 1. Initialize the refs as empty arrays. They will be filled by the API call.
const articles = ref<Article[]>([])
const latestArticles = ref<Article[]>([])

// onMounted runs its code as soon as the component is added to the page.
onMounted(async () => {
  // 2. Get the current route object to access URL parameters.
  const route = useRoute()
  const publisherId = route.params.Id // This 'Id' must match the parameter name in your router config.

  // 3. Add a guard clause in case the ID is missing from the URL.
  if (!publisherId) {
    console.error('Publisher ID is missing from the route.')
    return
  }

  // 4. Use a try...catch block to handle potential network errors.
  try {
    // 5. Make the API call to your new endpoint, passing the publisherId.
    const response = await axios.get<Article[]>(
      `http://localhost:8081/blog/by-publisher?id=${publisherId}`,
    )
    console.log(response)

    // 6. Update your reactive refs with the data from the API.
    articles.value = response.data
    latestArticles.value = response.data.slice(0, 3)

    console.log(
      `Successfully fetched ${response.data.length} articles for publisher ID: ${publisherId}`,
    )
  } catch (error) {
    console.error(`Failed to fetch articles for publisher ${publisherId}:`, error)
  }
})
</script>
<template>
  <NavBarView />
  <div id="about-container">
    <LayoutView />
    <div id="article-flex">
      <ArticleView v-for="article in articles" :key="article.ID" :article="article" />
    </div>
    <div id="about-section">
      <h1 class="about-text">John Doe</h1>
      <p class="about-text">placeholder@email.com</p>
      <div id="about-article-flex">
        <h1 class="about-text">Latest Articles</h1>
        <AboutArticleView v-for="article in latestArticles" :key="article.ID" :article="article" />
      </div>
    </div>
  </div>
</template>

<style scoped>
#about-container {
  display: flex;
  justify-content: space-between;
}
#about-section {
  background-color: #1a1a18;
  width: 30%;
}
#about-article-flex {
  display: flex;
  flex-direction: column;
  margin: 0 50px;
  gap: 1rem;
}
.about-text {
  color: white;
  text-align: center;
}
</style>
