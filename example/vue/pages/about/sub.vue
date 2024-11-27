<template>
	sub:{{sub}}
</template>

<script setup>
	import { ref } from 'vue';
	import { useRoute } from 'vue-router';

	const route = useRoute();
	const sub = ref(route.params.sub)
</script>