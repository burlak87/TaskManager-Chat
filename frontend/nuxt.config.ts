import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  ssr: false,
  devtools: { enabled: true },
  modules: [
    '@pinia/nuxt',
    '@vueuse/nuxt',
    'reka-ui/nuxt',
    '@nuxtjs/i18n',
  ],
  css: [
    '~/assets/css/main.css',
  ],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  i18n: {
    defaultLocale: 'ru',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
    langDir: 'locales/',
    locales: [
      {
        code: 'ru',
        file: 'ru.json',
        iso: 'ru-RU',
        name: 'Русский',
      },
    ],
    strategy: 'prefix_except_default',
  },
});