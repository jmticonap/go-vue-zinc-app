      
<script setup>
import { ChevronDoubleLeftIcon, ChevronDoubleRightIcon } from '@heroicons/vue/20/solid'
import { store } from '../store'

const goBackPage = async () => {
  if (store.currentPage > 1) {
    store.currentPage--
    await store.onSearch()
  }
}

const goForwarePage = async () => {
  if (store.currentPage * 20 < store.emailsTotal) {
    store.currentPage++
    await store.onSearch()
  }
}

</script>

<template>
  <div class="flex flex-row p-4 justify-center">
    <ul class="flex flex-row gap-4">
      <li>
        <button @click="goBackPage"
          class="rounded-md bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5  font-semibold text-white">
          <ChevronDoubleLeftIcon class="h-6 w-6" aria-hidden="true" />
        </button>
      </li>
      <li>
        <h2>
          {{ store.currentPage }} / {{ Math.ceil(store.emailsTotal / 20) }}
        </h2>
      </li>
      <li>
        <button @click="goForwarePage"
          class="rounded-md bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5  font-semibold text-white">
          <ChevronDoubleRightIcon class="h-6 w-6" aria-hidden="true" />
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>

</style>