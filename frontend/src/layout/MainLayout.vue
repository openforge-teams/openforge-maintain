<template>
  <a-layout class="main-layout">
    <Sidebar />
    <a-layout>
      <Header />
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
      <a-layout-footer class="footer">
        <span>openforge-maintain</span>
      </a-layout-footer>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
</script>

<style scoped>
.main-layout {
  height: 100vh;
}
.content {
  margin: 16px;
  padding: 24px;
  background: #f0f2f5;
  overflow-y: auto;
  min-height: 0;
}
.footer {
  text-align: center;
  padding: 12px 24px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
}
</style>
