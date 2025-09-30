<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import CommentView from '@/components/CommentView.vue'
import LayoutView from '@/components/LayoutView.vue'
import NavBarView from '@/components/NavBarView.vue'
import type { Comment } from '@/models/comment.interface'
import type { ArticleDTO } from '@/models/articleDTO.interface'

const comments = ref<Comment[]>([])
const article = ref<ArticleDTO | null>(null)
const newCommentContent = ref('') // State for the new comment input
const isLoggedIn = ref(false) // State to track login status

const route = useRoute()
const router = useRouter()
const articleId = route.params.Id as string


//TODO: use this function to communicate with the like route
const handleLikeComment = async (commentIdToLike: number) => {
  console.log("Not yet implemented in the backend")
  const token = localStorage.getItem('auth_token');
  if (!token) {
    console.error("Cannot like comment: user is not logged in.");
    // Optionally redirect to login
    // router.push('/login');
    return;
  }

  try {
    //TODO: ADD the corect api URL
    // MAKE THE API CALL to your future backend route.
    // A POST request is standard for creating a new "like".
    await axios.post(`http://localhost:8081/comment/${commentIdToLike}/like`, 
      {}, // Empty body, the important info is in the URL and headers
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );

    // INSTANT UI UPDATE: Find the comment in the local array and update its likes.
    // This provides great feedback to the user without a page refresh.
    const likedComment = comments.value.find(c => c.id === commentIdToLike);
    if (likedComment) {
      // For simplicity, we can just increment the count.
      // A more robust solution would be to refetch the comment or have the API
      // return the new list of likes.
      //likedComment.likes.push({} as any); // Add a placeholder like to update the count
      if (likedComment.liked === true) {
        likedComment.likes -= 1;
        likedComment.liked = false;
      }else {
        likedComment.likes += 1;
        likedComment.liked = true
      }
    }

  } catch (error) {
    console.error("Failed to like the comment:", error);
    // You could show an error message to the user here.
  }
};

/**
 * Fetches the article and its comments from the backend.
 */
const fetchArticleAndComments = async () => {
  if (!articleId) return

  try {
    //TODO: Add the username in the Response
    // Check the console.log
    const token = localStorage.getItem('auth_token');
    if (!token) {
      console.error("Cannot like comment: user is not logged in.");
      return;
    }
    const response = await axios.get(`http://localhost:8081/article/${articleId}`, 
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    );

    article.value = response.data
    console.log("======Cooments Data========")
    console.log(response.data)
    comments.value = response.data.comments || []
  } catch (error) {
    console.error('Failed to fetch the article:', error)
  }
}

/**
 * Handles the submission of the new comment form.
 */
const handlePostComment = async () => {
  const token = localStorage.getItem('auth_token')
  if (!token) {
    console.error('Not logged in. Cannot post comment.')
    router.push('/login')
    return
  }

  // Basic validation
  if (newCommentContent.value.trim() === '') {
    return
  }

  const payload = {
    content: newCommentContent.value,
    articleId: articleId,
  }

  const config = {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }

  try {
    // The endpoint matches your Go handler for creating a comment
    await axios.post('http://localhost:8081/comment', payload, config)
    // On success, clear the input and refresh the comments
    newCommentContent.value = ''
    await fetchArticleAndComments()
  } catch (error) {
    console.error('Failed to post comment:', error)
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      // If unauthorized, token is bad. Log the user out.
      localStorage.removeItem('auth_token')
      isLoggedIn.value = false
      router.push('/login')
    }
  }
}

// When the component mounts, check login status and fetch data
onMounted(async () => {
  const token = localStorage.getItem('auth_token')
  isLoggedIn.value = !!token
  await fetchArticleAndComments()
  console.log(await fetchArticleAndComments())
})
</script>

<template>
  <NavBarView />
  <div id="article-section">
    <LayoutView />
    <div v-if="article" id="article">
      <h1>{{ article.title }}</h1>
      <p class="article-content">{{ article.content }}</p>
    </div>
    <div v-else>
      <p>Loading article...</p>
    </div>
  </div>

  <div id="comment-section">
    <div id="style-line"></div>

    <!-- ADD COMMENT FORM (Visible only if logged in) -->
    <div v-if="isLoggedIn" class="comment-form-container">
      <h3>Leave a Comment</h3>
      <form @submit.prevent="handlePostComment" class="comment-form">
        <textarea
          v-model="newCommentContent"
          placeholder="Write your comment here..."
          rows="4"
          required
        ></textarea>
        <button type="submit">Post Comment</button>
      </form>
    </div>
    <CommentView v-for="comment in comments" :key="comment.id" :comment="comment"  @like-comm="handleLikeComment"  />
  </div>
</template>

<style scoped>
#article-section {
  display: flex;
}
#article {
  text-align: center;
  width: 100%;
  margin: 0 100px;
}
.article-content {
  text-align: left;
  line-height: 1.6;
  white-space: pre-wrap; /* Preserves formatting like line breaks */
}
#style-line {
  height: 1px;
  width: 100%;
  max-width: 800px;
  background-color: #dcdcdc;
  margin: 40px 0;
}
#comment-section {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: 0 2rem;
  padding-bottom: 4rem;
}

.comment-form-container {
  width: 100%;
  max-width: 800px;
  margin-bottom: 40px;
}

.comment-form-container h3 {
  margin-bottom: 1rem;
  font-size: 1.5rem;
  text-align: center;
}

.comment-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.comment-form textarea {
  padding: 12px;
  border-radius: 5px;
  border: 1px solid #d9d9d9;
  font-family: inherit;
  font-size: 1rem;
  resize: vertical;
}

.comment-form textarea:focus {
  outline: 2px solid rgb(93, 93, 93);
}

.comment-form button {
  padding: 12px;
  border-radius: 5px;
  background-color: black;
  color: white;
  border: none;
  cursor: pointer;
  font-weight: bold;
  align-self: flex-end; 
  width: 150px;
}


</style>
