<template>
	<a-layout style="height: 100%;" v-if="authStore.token">
		<a-layout-sider breakpoint="lg" collapsed-width="0" @collapse="onCollapse" @breakpoint="onBreakpoint">
			<div class="logo">
				<h1>Admin</h1>
			</div>

			<a-menu v-model:selectedKeys="selectedKeys" v-model:openKeys="openKeys" theme="dark" mode="inline">
				<template v-for="menu in menus">

					<a-sub-menu :key="menu.key" v-if="menu.children">
						<template #title>
							<span>
								<component :is="icons[menu.icon]"></component>
								{{menu.name}}
							</span>
						</template>
						<a-menu-item :key="child.key" v-for="child in menu.children">
							<RouterLink :to="child.to" v-if="child.to">{{child.name}}</RouterLink>
							<a :href="child.href" v-if="child.href" :target="child.target">{{child.name}}</a>
						</a-menu-item>
					</a-sub-menu>

					<a-menu-item :key="menu.key" v-if="!menu.children">
						<component :is="icons[menu.icon]"></component>
						<span class="nav-text">
							<RouterLink :to="menu.to" v-if="menu.to">{{menu.name}}</RouterLink>
							<a :href="menu.href" v-if="menu.href" :target="menu.target">{{menu.name}}</a>
						</span>
					</a-menu-item>

				</template>
				<a-menu-item>
					<span class="nav-text">
						<a-button type="link" @click="logout">退出</a-button>
					</span>
				</a-menu-item>
			</a-menu>
		</a-layout-sider>
		<a-layout>
			<a-layout-header :style="{ background: '#fff', padding: 0 }" />
			<a-layout-content :style="{ margin: '16px' }">
				<div :style="{ padding: '16px', background: '#fff', height: '100%' }">
					<RouterView></RouterView>
				</div>
			</a-layout-content>
			<a-layout-footer style="text-align: center">@2024</a-layout-footer>
		</a-layout>
	</a-layout>
</template>

<script setup>
	import { ref, onBeforeMount } from 'vue'
	import { useRoute, useRouter } from 'vue-router';
	import { useAuthStore } from '@/store/auth'
	import { menus } from '@/routes'
	import { MenuOutlined } from '@ant-design/icons-vue';

	const icons = { MenuOutlined };
	const authStore = useAuthStore()
	const route = useRoute();
	const router = useRouter()

	const onCollapse = (collapsed, type) => {
		console.log(collapsed, type);
	};

	const onBreakpoint = (broken) => {
		console.log(broken);
	};

	const selectedKeys = ref([]);
	const openKeys = ref([])
	menus.forEach(x => {
		if ('children' in x) {
			x.children.forEach(y => {
				if (y.to != "/" && route.path == y.to) {
					selectedKeys.value.push(y.key)
					openKeys.value.push(x.key)
				}
			})
		} else if (x.to != "/" && route.path == x.to) {
			selectedKeys.value.push(x.key)
		}
	})
	if (selectedKeys.value.length == 0) selectedKeys.value = [menus[0].key]

	function logout() {
		authStore.logout()
		router.push({ path: '/login' })
	}

	onBeforeMount(_ => {
		if (!authStore.token && route.path != "/login") {
			router.push({ path: '/login' })
		}
	})
</script>

<style scoped>
	.logo {
		color: #fff;
		height: 32px;
		/* background: rgba(255, 255, 255, 0.2); */
		margin: 16px;
	}

	.site-layout-sub-header-background {
		background: #fff;
	}

	.site-layout-background {
		background: #fff;
	}

	[data-theme='dark'] .site-layout-sub-header-background {
		background: #141414;
	}
</style>