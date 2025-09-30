<script setup lang="ts">
// The component now imports the corrected Comment interface
import type { Comment } from '@/models/comment.interface'

const props = defineProps<{
  comment: Comment
}>()
</script>

<template>
  <div class="comment-card">
    <p class="publisher-name">{{ comment.commenterName }}</p>

    <p class="content">{{ comment.content }}</p>

    <div class="comment-actions">
           <button @click="$emit('likeComm', comment.id)" class="like-button">
        {{ comment.liked ? "Unlike" : "Like" }} ({{ comment.likes }})
      </button>
    </div>
  </div>
</template>

<style scoped>
.comment-card {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid #e0e0e0;
  max-width: 800px;
  width: 100%;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.publisher-name {
  font-weight: 600;
  font-size: 0.9rem;
  color: #333;
  margin: 0;
}

.content {
  font-size: 1rem;
  color: #555;
  line-height: 1.5;
  margin: 0;
}

.comment-actions {
  margin-top: 12px;
  padding-top: 8px;
  border-top: 1px solid #eee;
}

.likes-count {
  font-size: 0.85rem;
  font-weight: 500;
  color: #666;
  margin: 0;
}

.like-button {
  background-color: #e7f3ff;
  color: #1c7ed6;
  border: 1px solid #a5d8ff;
  border-radius: 4px;
  padding: 6px 12px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}
.like-button:hover {
  background-color: #d0ebff;
}
</style>
