<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import LayoutView from '@/components/LayoutView.vue'
import NavBarView from '@/components/NavBarView.vue'

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const handleRegister = async () => {
  if (password.value !== confirmPassword.value) {
    alert('Passwords do not match!')
    return
  }

  const payload = {
    username: username.value,
    email: email.value,
    password: password.value,
  }

  try {
    const response = await axios.post('http://localhost:8081/auth/register', payload)

    if (response.status === 200) {
      console.log('Registration successful! Response:', response)
      alert('Registration was successful!')
    }
  } catch (error) {
    console.error('Registration failed:', error)

    if (axios.isAxiosError(error) && error.response) {
      const errorMessage = error.response.data
      alert(`Error: ${errorMessage}`)
    } else {
      alert('An unexpected network error occurred.')
    }
  }
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
      <form id="login-form" @submit.prevent="handleRegister">
        <label for="uname"><b>Username</b></label>
        <input v-model="username" type="text" placeholder="Enter Username" name="uname" required />

        <label for="email"><b>Email</b></label>
        <input v-model="email" type="email" placeholder="Enter Email" name="email" required />

        <label for="psw"><b>Password</b></label>
        <input
          v-model="password"
          type="password"
          placeholder="Enter Password"
          name="psw"
          required
        />

        <label for="conf-psw"><b>Confirm Password</b></label>
        <input
          v-model="confirmPassword"
          type="password"
          placeholder="Confirm Password"
          name="conf-psw"
          required
        />

        <button type="submit">Register</button>
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
