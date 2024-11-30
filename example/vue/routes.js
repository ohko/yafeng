import { createRouter, createWebHistory } from 'vue-router'

const routes = [
	{ path: "/login", component: () => import('./pages/user/login.vue') },
	{
		name: '首页', icon: 'MenuOutlined', path: "/", component: () => import('./layout/default.vue'),
		children: [{ path: "/", component: () => import("./pages/index.vue") }]
	},
	{
		name: '目录', icon: 'MenuOutlined', path: "/about", component: () => import('./layout/default.vue'),
		children: [
			{
				name: 'Page1', icon: 'MenuOutlined', path: "/about/page1", component: () => import("./pages/about/page1.vue"),
				children: [
					{ path: ":sub", component: () => import("./pages/about/sub.vue") },
				]
			},
			{ name: 'Page2', icon: 'MenuOutlined', path: "/about/page2", component: () => import("./pages/about/page2.vue") },
		]
	},
	{
		name: '修改密码', icon: 'MenuOutlined', path: "/user/change", component: () => import('./layout/default.vue'),
		children: [{ path: "/user/change", component: () => import("./pages/user/change.vue") }]
	},
]

let keyIndex = 0
function subm(x) {
	let children = []
	if ('children' in x) {
		x.children.forEach(y => { if ('name' in y) children.push(subm(y)) })
	}
	if (children.length > 0) return { key: ++keyIndex, to: x.path, name: x.name, icon: x.icon, children: children }
	return { key: ++keyIndex, to: x.path, name: x.name, icon: x.icon }
}

const menus = []
routes.forEach(x => {
	if ('name' in x) menus.push(subm(x))
})

const router = createRouter({
	routes,
	history: createWebHistory(import.meta.env.BASE_URL),
})

export { router, menus }