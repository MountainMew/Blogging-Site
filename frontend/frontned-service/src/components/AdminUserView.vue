<script setup lang="ts">
import { ref } from 'vue'
import type { User } from '@/models/user.interface'
import UserView from './UserView.vue'

const apiResponse = [
  {
    Id: '101',
    Username: 'alex.smith',
    Email: 'alex.smith@example.com',
    Role: 'admin',
  },
  {
    Id: '102',
    Username: 'jane.doe',
    Email: 'jane.d@example.com',
    Role: 'user',
  },
  {
    Id: '103',
    Username: 'brian_c',
    Email: 'b.clark@example.com',
    Role: 'user',
  },
  {
    Id: '104',
    Username: 'samantha.jones',
    Email: 'sjones@example.com',
    Role: 'admin',
  },
  {
    Id: '105',
    Username: 'mike_p',
    Email: 'mike.p@example.com',
    Role: 'user',
  },
]

const users = ref<User[]>(apiResponse)
</script>

<template>
  <div id="admin-user-container">
    <div id="filter-user-bar">
      <label for="users">Search User By:</label>

      <select name="users" id="users">
        <option value="Id">Id</option>
        <option value="Username">Username</option>
        <option value="Role">Role</option>
      </select>

      <form @submit.prevent>
        <input
          type="text"
          id="search-input"
          name="search-input"
          placeholder="Enter search term..."
        />
      </form>
    </div>

    <div id="users-grid">
      <UserView v-for="user in users" :key="user.Id" :user="user" />
    </div>
  </div>
</template>

<style scoped>
#admin-user-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem;
  background-color: #fdfdfd;
  width: 100%;
  box-sizing: border-box;
}

#filter-user-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  width: 100%;
  max-width: 1200px;
  background-color: #ffffff;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-sizing: border-box;
}

#filter-user-bar label {
  font-weight: 500;
  color: #555;
  white-space: nowrap;
}

#filter-user-bar select,
#filter-user-bar input {
  padding: 0.6rem 0.8rem;
  font-size: 1rem;
  color: #333;
  border: 1px solid #ced4da;
  border-radius: 6px;
  background-color: #fff;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

#filter-user-bar select:focus,
#filter-user-bar input:focus {
  outline: none;
  border-color: #80bdff;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.2);
}

#filter-user-bar form {
  flex-grow: 1;
}
#filter-user-bar #search-input {
  width: 100%;
  box-sizing: border-box;
}

#users-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  width: 100%;
  max-width: 1200px;
}
</style>
