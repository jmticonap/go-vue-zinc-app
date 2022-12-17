<script setup>
import { store } from '../store'

const selectRow = evt => {
  store.selectedEmailInfo = store.emails.find(email => email._id == store.selectedEmailID)
  if (!store.selectedEmailInfo._source.content)
    store.selectedEmailInfo._source.content = "---NO CONTENT HERE---"

    store.open = false
}
</script>

<template>
  <table class="min-w-full">
    <thead class="bg-slate-400 border-b border-slate-200">
      <tr>
        <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">Subject</th>
        <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">From</th>
        <th class="px-6 py-3 text-left text-sm font-medium text-slate-900">To</th>
      </tr>
    </thead>
    <tbody>
      <tr @click.stop="() => { store.selectedEmailID = email._id; selectRow() }" v-for="email in store.emails"
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
</template>

<style>

</style>