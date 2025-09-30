<script setup lang="ts">
import SearchIcon from '../assets/svg/search-icon.svg'
//TODO: implemenet API Call logic for search bar
</script>

<template>
  <div id="search-bar-container">
    <div id="search-input-container">
      <form>
        <input type="text" id="serach-input" name="serach-input" placeholder="Search" />
      </form>
      <img :src="SearchIcon" alt="Search Icon" id="search-icon" />
    </div>
    <div id="style-line"></div>
  </div>
</template>

<style scoped>
#search-bar-container {
  display: flex;
  flex-direction: column;
  gap: 5px;
  margin: 0;
  padding: 0;
}
#style-line {
  height: 1px;
  width: 100%;
  background-color: #808080;
  margin: 0;
  padding: 0;
}
#search-input-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
#serach-input {
  width: 20rem;
  border: none;
  outline: none;
}
#search-icon {
  height: 1.5rem;
  background-color: 808080;
  margin: 0;
  padding: 0;
}
</style>
