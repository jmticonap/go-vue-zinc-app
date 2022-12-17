<script setup>
import { onMounted } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { XMarkIcon, ChevronDoubleLeftIcon } from '@heroicons/vue/24/outline'
import Searcher from './components/Searcher.vue'
import Paginator from './components/Paginator.vue'
import TableEmails from './components/TableEmails.vue'
import { store } from './store'

onMounted(async () => {
  await store.onSearch(undefined, {
    from: store.currentPage,
    size: 20,
    count: store.emailsTotal
  })
})

</script>

<template>
  <header>
    <!-- Searcher -->
    <Searcher />
    <!--==========-->
  </header>
  <div class="p-8 pt-20">
    <pre>
      {{ store.selectedEmailInfo?._source.content }}
    </pre>
  </div>
  <TransitionRoot as="template" :show="store.open">
    <Dialog as="div" class="relative z-10" @close="store.open = false">
      <TransitionChild as="template" enter="ease-in-out duration-500" enter-from="opacity-0" enter-to="opacity-100"
        leave="ease-in-out duration-500" leave-from="opacity-100" leave-to="opacity-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-hidden">
        <div class="absolute inset-0 overflow-hidden">
          <div class="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10">
            <TransitionChild as="template" enter="transform transition ease-in-out duration-500 sm:duration-700"
              enter-from="translate-x-full" enter-to="translate-x-0"
              leave="transform transition ease-in-out duration-500 sm:duration-700" leave-from="translate-x-0"
              leave-to="translate-x-full">
              <DialogPanel class="pointer-events-auto relative w-screen max-w-4xl">
                <TransitionChild as="template" enter="ease-in-out duration-500" enter-from="opacity-0"
                  enter-to="opacity-100" leave="ease-in-out duration-500" leave-from="opacity-100" leave-to="opacity-0">
                  <div class="absolute top-0 left-0 -ml-8 flex pt-4 pr-2 sm:-ml-10 sm:pr-4">
                    <button type="button"
                      class="rounded-md text-gray-300 hover:text-white focus:outline-none focus:ring-2 focus:ring-white"
                      @click="store.open = false">
                      <span class="sr-only">Close panel</span>
                      <XMarkIcon class="h-6 w-6" aria-hidden="true" />
                    </button>
                  </div>
                </TransitionChild>
                <div class="flex h-full flex-col overflow-y-scroll bg-white py-6 shadow-xl">
                  <div class="px-4 sm:px-6">
                    <DialogTitle class="text-lg font-medium text-gray-900">Search Result</DialogTitle>
                  </div>
                  <div class="relative mt-6 flex-1 px-4 sm:px-6">
                    <!-- Replace with your content -->
                    <div class="absolute inset-0 px-4 sm:px-6">
                      <div class="overflow-y-auto h-full border-2 border-dashed border-gray-200" aria-hidden="true">
                        <!-- TableEmails -->
                        <TableEmails />
                        <!--=============-->
                      </div>

                    </div>
                  </div>

                  <!-- Paginator -->
                  <Paginator />
                  <!--===========-->
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
  <div class="bg-red-700 fixed top-1/2 top_button">
    <button type="button"
      class="btn_pnl_shower rounded-md bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5  font-semibold text-white absolute right-0 top-1/2"
      @click="store.open = true">
      <ChevronDoubleLeftIcon class="h-6 w-6" aria-hidden="true" />
    </button>
  </div>
</template>

<style scoped>
.top-nav-shadow {
  filter: drop-shadow(0 0 0.25rem #000);
}

.top_button {
  right: 0;
}

.btn_pnl_shower {
  height: 3rem;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.cell-border {
  border-left: solid #9f9f9f 0.1rem;
  padding: 0.25rem 0.5rem;
}

.cell-border:nth-child(1) {
  border-left: none
}

pre {
  white-space: pre-wrap;
  /* Since CSS 2.1 */
  white-space: -moz-pre-wrap;
  /* Mozilla, since 1999 */
  white-space: -pre-wrap;
  /* Opera 4-6 */
  white-space: -o-pre-wrap;
  /* Opera 7 */
  word-wrap: break-word;
  /* Internet Explorer 5.5+ */
}
</style>
