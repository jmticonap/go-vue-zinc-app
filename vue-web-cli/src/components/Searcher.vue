<script setup>
import { store } from '../store'
</script>

<template>
  <nav
    class="fixed py-4 bg-white flex flex-row flex-nowrap justify-center items-center p-2 px-4 min-w-full top-nav-shadow">
    <form @submit.prevent="store.onSearch" class="flex flex-row justify-center w-full">
      <ul class="flex flex-row flex-nowrap gap-1 w-10/12">
        <li class="w-full">
          <input v-model="store.search" class="border-gray-900 border-solid border rounded-md p-1 w-full"
            type="store.search" name="search_field" placeholder="Term for search" />
        </li>
        <li>
          <button @click="store.onSearch"
            class="bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5 rounded-full font-semibold text-white"
            type="button">search</button>
        </li>
      </ul>
    </form>

  </nav>
</template>

<style scoped>

</style>