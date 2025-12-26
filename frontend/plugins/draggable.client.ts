import draggable from '@sortablejs/vue-draggable-next'

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.component('Draggable', draggable)
})
