<script setup>
import { ref } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { XMarkIcon, ChevronDoubleLeftIcon } from '@heroicons/vue/24/outline'
import { searchEmail } from './services/EmailService'

const open = ref(false)
const search = ref("")
const emails = ref([])
const selectedEmailID = ref('')
const selectedEmailInfo = ref(undefined)

const onSearch = async evt => {
  emails.value = await searchEmail(search.value)
  open.value = true
}

const selectRow = evt => {
  selectedEmailID.value = evt.target.id
  selectedEmailInfo.value = emails.value.find(email => email._id == selectedEmailID.value)

  open.value = false
}
</script>

<template>
  <header>
    <nav class="flex flex-row flex-nowrap justify-center items-center p-2 px-4 min-w-full">
      <form @submit.default="onSearch" class="flex flex-row justify-center w-full">
        <ul class="flex flex-row flex-nowrap gap-1 w-10/12">
          <li class="w-full">
            <input v-model="search" class="border-gray-900 border-solid border rounded-md p-1 w-full" type="search"
              name="search_field" placeholder="Term for search" />
          </li>
          <li>
            <button @click="onSearch"
              class="bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5 rounded-full font-semibold text-white"
              type="button">search</button>
          </li>
        </ul>
      </form>

    </nav>
  </header>
  <div class="p-8">
    <pre>
      {{ selectedEmailInfo?._source.Content}}
    </pre>
  </div>
  <TransitionRoot as="template" :show="open">
    <Dialog as="div" class="relative z-10" @close="open = false">
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
                      @click="open = false">
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
                      <div class="h-full border-2 border-dashed border-gray-200" aria-hidden="true">
                        <!--
                        <table class="min-w-full">
                          
                          <thead class="bg-slate-50 border-b border-slate-200">
                            <tr>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">Subject</th>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">From</th>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">To</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr @click.stop="selectRow" v-for="email in emails" class="odd:bg-white even:bg-slate-50">
                              <td v-bind:id="email._id"
                                class="my-colum px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-900">
                                {{ email._source['Subject'] }}
                              </td>
                              <td class="my-colum px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-900">
                                {{ email._source['From'] }}
                              </td>
                              <td class="my-colum px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-900">
                                {{ email._source['To'] }}
                              </td>

                            </tr>

                          </tbody>
                        </table>
                      -->
                        <table class="min-w-full">
                          <thead class="bg-slate-400 border-b border-slate-200">
                            <tr>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">Subject</th>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">From</th>
                              <th class="my-colum px-6 py-3 text-left text-sm font-medium text-slate-900">To</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr @click.stop="selectRow" v-for="email in emails"
                              class="odd:bg-white even:bg-slate-200 hover:bg-yellow-200">
                              <td v-bind:id="email._id">
                                {{ email._source['Subject'] }}
                              </td>
                              <td>
                                {{ email._source['From'] }}
                              </td>
                              <td>
                                {{ email._source['To'] }}
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                      <!--
                        <Paginator />
                      -->
                    </div>
                    <!-- /End replace -->
                  </div>
                </div>
              </DialogPanel>
            </TransitionChild>
          </div>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
  <div class="bg-red-700 relative top_button">
    <button type="button"
      class="btn_pnl_shower rounded-md bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5  font-semibold text-white absolute right-0 top-1/2"
      @click="open = true">
      <ChevronDoubleLeftIcon class="h-6 w-6" aria-hidden="true" />
    </button>
  </div>
</template>

<style scoped>
.top_button {
  position: fixed;
  top: 50%;
  right: 0;
}

.btn_pnl_shower {
  height: 3rem;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.my-colum {}

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
