<script setup>
import { ref, onMounted } from 'vue'
import { Dialog, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'
import { XMarkIcon, ChevronDoubleLeftIcon, ChevronDoubleRightIcon } from '@heroicons/vue/24/outline'
import { searchEmail } from './services/EmailService'

const open = ref(true)
const search = ref("")
const lastSearch = ref("")
const emails = ref([])
const emailsTotal = ref(1)
const currentPage = ref(1)
const selectedEmailID = ref('')
const selectedEmailInfo = ref(undefined)

const onSearch = async evt => {
  try {
    if (lastSearch.value != search.value) {
      emailsTotal.value = 1
      currentPage.value = 1
    }
    const { _emails, _total } = await searchEmail(
      search.value,
      {
        from: currentPage.value,
        size: 20,
        count: emailsTotal.value
      }
    )
    lastSearch.value = search.value

    //removing unuseful characters
    for (let mail of _emails) {
      for (let d in mail._source) {
        mail._source[d] = mail._source[d].trim()
      }
    }

    emails.value = _emails
    emailsTotal.value = _total

    open.value = true
  } catch (error) {
    alert(error)
  }
}

const selectRow = evt => {
  selectedEmailInfo.value = emails.value.find(email => email._id == selectedEmailID.value)
  if (!selectedEmailInfo.value._source.content)
    selectedEmailInfo.value._source.content = "---NO CONTENT HERE---"

  open.value = false
}

const goBackPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    onSearch()
  }
}

const goForwarePage = () => {
  if (currentPage.value * 20 < emailsTotal.value) {
    currentPage.value++
    onSearch()
  }
}

onMounted(async () => {
  await onSearch(undefined, {
    from: currentPage.value,
    size: 20,
    count: emailsTotal.value
  })
})

</script>

<template>
  <header>
    <nav class="fixed py-4 bg-white flex flex-row flex-nowrap justify-center items-center p-2 px-4 min-w-full top-nav-shadow">
      <form @submit.prevent="onSearch" class="flex flex-row justify-center w-full">
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
  <div class="p-8 pt-20">
    <pre>
      {{ selectedEmailInfo?._source.content }}
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
                      <div class="overflow-y-auto h-full border-2 border-dashed border-gray-200" aria-hidden="true">
                        <!--  -->
                        <table class="min-w-full">
                          <thead class="bg-slate-400 border-b border-slate-200">
                            <tr>
                              <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">Subject</th>
                              <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">From</th>
                              <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">To</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr @click.stop="() => { selectedEmailID = email._id; selectRow() }" v-for="email in emails"
                              class="odd:bg-white even:bg-slate-200 hover:bg-yellow-200">
                              <td class="cell-border">
                                {{ email._source['subject']
                                    ? email._source['subject']
                                    : '---none---'
                                }}
                              </td>
                              <td class="cell-border">
                                {{ email._source['from'] || '---none---' }}
                              </td>
                              <td class="cell-border">
                                {{ email._source['to'] || '---none---' }}
                              </td>
                            </tr>
                          </tbody>
                        </table>

                      </div>

                    </div>
                  </div>

                  <!--  -->
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
                          {{ currentPage }} / {{ Math.ceil(emailsTotal / 20) }}
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
      @click="open = true">
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
