// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: false},
    modules: ["@nuxt/ui", '@nuxt/icon', '@nuxtjs/color-mode', '@vueuse/nuxt', 'dayjs-nuxt'],
    ssr: false,
    colorMode: {
        preference: 'system',
        fallback: 'light',
        classSuffix: '',
    },
    dayjs: {
        locales: ['zh'],
        defaultLocale: 'zh'
    },
    icon: {
        clientBundle: {
            scan: {
                globInclude: ['**/*.{vue,jsx,tsx}', 'node_modules/@nuxt/ui/**/*.js'],
                globExclude: ['.*', 'coverage', 'test', 'tests', 'dist', 'build'],
            },
        },
    },
    tailwindcss: {
        safelist: [
            'grid-cols-1',
            'grid-cols-3',
        ]
    },
    vue: {
        compilerOptions: {
            isCustomElement: (tag:string) => ['meting-js'].includes(tag),
        },
    },
    app: {
        head: {
            meta: [
                { name: "viewport", content: "width=device-width, initial-scale=1, user-scalable=no" },
                { charset: "utf-8" },
                { name: "color-scheme", content: "light dark" },
            ],
            script: [
                // APlayer/Meting 改为在 MusicPreview 组件挂载时按需注入,
                // 不再全局加载;main.js 提供全页通用的代码块复制功能,保留。
                {src: `/js/main.js`, type: 'text/javascript', async: true, defer: true},
            ]
        }
    },
    vite: {
        server: {
            proxy: {
                "/api": {
                    target: "http://localhost:37892",
                },
                "/upload": {
                    target: "http://localhost:37892",
                },
                "/rss": {
                    target: "http://localhost:37892",
                },
                "/swagger": {
                    target: "http://localhost:37892",
                },
            },
        },
        build: {
            rollupOptions: {
                output: {
                    hashCharacters: 'base36'
                }
            }
        }
    }
})
