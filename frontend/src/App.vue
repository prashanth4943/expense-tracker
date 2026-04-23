<script setup>
import { ref } from 'vue'
import ExpenseForm from './components/ExpenseForm.vue'
import ExpenseList from './components/ExpenseList.vue'

const refreshKey = ref(0)

function handleAdded() {
  refreshKey.value++
}
</script>

<template>
  <div>
    <h1>Expense Tracker</h1>

    <ExpenseForm @added="handleAdded" />

<ExpenseList :refreshTrigger="refreshKey" />
  </div>
</template>