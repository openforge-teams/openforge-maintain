<template>
  <div class="page-header">
    <a-breadcrumb>
      <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
        <router-link v-if="item.path" :to="item.path">{{ item.title }}</router-link>
        <span v-else>{{ item.title }}</span>
      </a-breadcrumb-item>
    </a-breadcrumb>
    <h2 class="page-title">{{ title }}</h2>
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title?: string
  breadcrumbs?: { title: string; path?: string }[]
}>(), {
  title: '',
  breadcrumbs: () => [],
})
</script>

<style scoped>
.page-header {
  margin-bottom: 24px;
}
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 12px 0 0;
  color: rgba(0, 0, 0, 0.85);
}
</style>
