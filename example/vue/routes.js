import { createRouter, createWebHistory } from 'vue-router'

const data = [
	{ name: '首页', icon: 'MenuOutlined', to: "/", path: "/", layout: './layout/default.vue', component: "./pages/index.vue" },
	{
		name: '目录', icon: 'MenuOutlined', path: "/about", layout: './layout/default.vue',
		children: [
			{ name: 'Page1', icon: 'MenuOutlined', to: "/about", path: "", component: "./pages/about/index.vue" },
			{ name: 'Page2', icon: 'MenuOutlined', to: "/about/sub", path: ":sub", component: "./pages/about/sub.vue" },
		]
	},
	{ name: '修改密码', icon: 'MenuOutlined', to: "/user/change", path: "/user/change", layout: './layout/default.vue', component: "./pages/user/change.vue" },
]

function subr(x) {
	let children = []
	if ('children' in x) { x.children.forEach(y => { children.push(subr(y)) }) }

	if ('layout' in x && 'component' in x) {
		children.unshift({ path: "", component: () => import(/* @vite-ignore */x.component) })
	}
	return { path: x.path, component: () => import(/* @vite-ignore */'layout' in x ? x.layout : x.component), children: children }
}

let keyIndex = 0
function subm(x) {
	if ('children' in x) {
		let children = []
		x.children.forEach(y => { children.push(subm(y)) })
		return { key: ++keyIndex, name: x.name, icon: x.icon, children: children }
	} else {
		return { key: ++keyIndex, to: x.to, name: x.name, icon: x.icon }
	}
}

const routes = [{ path: "/login", component: () => import('./pages/user/login.vue') }]
const menus = []
data.forEach(x => { routes.push(subr(x)); menus.push(subm(x)) })

const router = createRouter({
	routes,
	history: createWebHistory(import.meta.env.BASE_URL),
})

export { router, menus }