<script setup>
  import { ref } from 'vue'

  defineProps({
    msg: String,
  })

  const count = ref(0)
</script>

<template>
  <div class="card">
    <h1>{{ msg }}</h1>
    <button @click="count++" class="btn">count is {{ count }}</button>
  </div>
</template>