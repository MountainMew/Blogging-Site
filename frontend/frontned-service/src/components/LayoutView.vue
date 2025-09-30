<!-- TODO: Add social media icons -->
<script setup lang="ts"></script>

<template>
  <div id="layout-container">
    <div id="layout-line"></div>
    <div id="icons-container">
      <span>Icon</span>
      <span>Icon</span>
      <span>Icon</span>
    </div>
  </div>
</template>

<style scoped>
#layout-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin: 0 2rem;
}

#layout-line {
  height: 50rem;
  width: 2.5px;
  background-color: black;
}
#icons-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>
