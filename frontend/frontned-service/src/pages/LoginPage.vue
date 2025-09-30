<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import LayoutView from '@/components/LayoutView.vue'
import NavBarView from '@/components/NavBarView.vue'

const router = useRouter()

const username = ref('')
const password = ref('')

const handleLogin = async () => {
  const payload = {
    username: username.value,
    password: password.value,
  }

  try {
    const response = await axios.post('http://localhost:8081/auth/login', payload)

    console.log('Login successful! Response:', response)

    const token = response.data.token

    if (token) {
      localStorage.setItem('auth_token', token)
      router.push('/')
    } else {
      console.error('Login successful, but no token was provided in the response.')
    }
  } catch (error) {
    // --- ⬇️ CHANGE IS HERE ⬇️ ---

    // Log a user-friendly error message to the console.
    if (axios.isAxiosError(error) && error.response) {
      console.error(`Login failed: ${error.response.data}`)
    } else {
      console.error('An unexpected network error occurred.')
    }

    // You can also log the full error object for more detailed debugging.
    console.error('Full error object:', error)
  }
}

const registerRedirect = () => {
  console.log('Redirecting to register page...')
  router.push('/register')
}
</script>

<template>
  <NavBarView />
  <div class="auth-layout">
    <LayoutView />
    <div class="auth-container">
      <h1 class="wellcome-text">
        WELLCOME <br />
        TO WEBAPP
      </h1>
      <form id="login-form" @submit.prevent="handleLogin">
        <label for="uname"><b>Username</b></label>
        <input v-model="username" type="text" placeholder="Enter Username" name="uname" required />

        <label for="psw"><b>Password</b></label>
        <input
          v-model="password"
          type="password"
          placeholder="Enter Password"
          name="psw"
          required
        />

        <button type="submit">Login</button>
        <p v-on:click="registerRedirect" class="signup-link">Don't have an account?</p>
      </form>
    </div>
    <img class="auth-image" src="../assets/images/Placeholder.png" alt="Placeholder Image" />
  </div>
</template>

<style scoped>
.auth-container {
  height: 90vh;
  width: 80%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}
.auth-layout {
  display: flex;
  justify-content: space-between;
}
.auth-image {
  margin: 0 2rem;
  width: 45%;
  height: 90vh;
}
.wellcome-text {
  font-size: 50px;
  font-weight: bold;
  text-align: center;
}
.signup-link {
  cursor: pointer;
}
#login-form {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  border: 1px solid #d9d9d9;
  padding: 10px;
}

/* Selects both input fields and the button within the form */
#login-form input,
#login-form button {
  padding: 12px;
  border-radius: 5px; /* Optional: softens the corners */
}

/* Specifically targets just the input fields */
#login-form input {
  border: 1px solid #d9d9d9;
  transition: all 100ms;
}

/* A little extra styling for the button to make it stand out */
#login-form button {
  margin-top: 5px;
  background-color: black; /* A nice blue background */
  color: white;
  border: none; /* Removes the default border */
  cursor: pointer;
  font-weight: bold;
}

/* Style for the signup link */
#login-form .signup-link {
  text-align: center;
  cursor: pointer;
  color: black;
  text-decoration: underline;
}

#login-form input:focus {
  outline: 2px solid rgb(93, 93, 93); /* Sets a solid black outline when the input is selected */
}
</style>
