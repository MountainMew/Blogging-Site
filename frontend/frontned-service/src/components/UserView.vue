<script setup lang="ts">
// TODO: Add styling
import { ref } from 'vue'

import type { User } from '@/models/user.interface'
const props = defineProps<{
  user: User
}>()

const activeSection = ref('')
</script>

<template>
  <div class="user-wrapper">
    <div class="user-card">
      <p>{{ user.Username }}</p>
      <p>{{ user.Email }}</p>
      <p>{{ user.Role }}</p>
    </div>
    <div class="actions">
      <button @click="activeSection = ''">Delete User</button>
      <button @click="activeSection = 'Update'">Update User</button>
    </div>
    <div class="form-wrapper" v-if="activeSection !== ''">
      <form>
        <input type="text" class="action-input" name="usersname" placeholder="Type New Username" />
        <input type="text" class="action-input" name="password" placeholder="Type New Password" />
        <input type="text" class="action-input" name="email" placeholder="Type New Email" />
        <input type="text" class="action-input" name="role" placeholder="Change Role" />
      </form>
      <button @click="activeSection = ''">Cancel</button>
    </div>
  </div>
</template>

<style scoped>
.user-wrapper {
  display: flex;
  flex-direction: column;
}
.user-card {
  background-color: #f9f9f9;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
}

.user-card p {
  margin: 0;
  padding: 4px 0;
}

.user-card p:first-child {
  font-weight: 600;
  font-size: 1.1em;
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.form-wrapper {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e0e0e0;
}

.form-wrapper form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
}

button {
  font-size: 1rem;
  padding: 10px 16px;
  border-radius: 6px;
  cursor: pointer;
  background-color: #ffffff;
  border: 1px solid #cccccc;
  color: #333;
  transition:
    background-color 0.2s,
    color 0.2s;
}

button:hover {
  background-color: #212121;
  color: #ffffff;
  border-color: #212121;
}

.action-input {
  font-size: 1rem;
  padding: 10px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  background-color: #ffffff;
}

.action-input:focus {
  outline: none;
  border-color: #000000;
  box-shadow: 0 0 0 3px rgba(0, 0, 0, 0.05);
}
</style>
