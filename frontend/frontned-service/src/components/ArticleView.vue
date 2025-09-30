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

function redirectToUserProfile() {
  console.log(`Redirecting to user profile with ID: ${props.article.PublisherID}`)

  router.push({
    name: 'UserDetail',
    params: {
      Id: props.article.PublisherID,
    },
  })
}
</script>
<template>
  <div class="flex-article-card">
    <img class="article-photo" src="../assets/images/Placeholder.png" alt="Placeholder Image" />
    <div class="article-content">
      <p>{{ article.Category }}</p>
      <h2 v-on:click="redirectToArticle" class="click-text">{{ article.Title }}</h2>
      <p>{{ article.Content.slice(0, 200) }}...</p>
      <p>
        Publisher:
        <span v-on:click="redirectToUserProfile" class="bold-text click-text">{{
          article.PublisherName
        }}</span>
      </p>
      <p>
        Date: <span class="bold-text">{{ article.CreatedAt.split('T')[0] }}</span>
      </p>
    </div>
  </div>
  <div id="style-line"></div>
</template>

<style scoped>
.flex-article-card {
  max-width: 1000px;
  margin: 0;
  padding: 0;
  display: flex;
  gap: 1.5rem;
}
.click-text {
  cursor: pointer;
}
.bold-text {
  font-weight: bold;
  color: #808080;
}
.article-photo {
  max-width: 300px;
  max-height: 300px;
  border-radius: 5%;
}

.article-content {
  display: flex;
  flex-direction: column;
}

#style-line {
  height: 1px;
  width: 100%;
  background-color: #808080;
  margin: 20px 0;
  padding: 0;
}
</style>
