<script setup lang="ts">
import type { Article } from '@/models/article.interface'
import { useRouter } from 'vue-router'
// Define article props
// The article interface must match the API response
const props = defineProps<{
  article: Article
}>()

const router = useRouter()

function redirectToArticle() {
  console.log(`Redirecting to article with ID: ${props.article.ID}`)

  router.push({
    name: 'ArticleDetail',
    params: {
      Id: props.article.ID,
    },
  })
}
</script>
<template>
  <div class="flex-article-card">
    <img class="article-photo" src="../assets/images/Placeholder.png" alt="Placeholder Image" />
    <div class="article-content">
      <p v-on:click="redirectToArticle" class="click-text white-text">{{ article.Title }}</p>
      <p class="white-text">
        {{ article.PublisherName }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.flex-article-card {
  margin: 0;
  padding: 0;
  display: flex;
  gap: 20px;
  align-items: center;
}
#article-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.click-text {
  cursor: pointer;
  font-weight: bold;
}
.article-photo {
  max-width: 50px;
  max-height: 50px;
  border-radius: 5%;
}
.white-text {
  color: white;
  margin: 0;
  padding: 0;
}
</style>
