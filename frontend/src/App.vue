<template>
  <div>
    <h1>Expense Tracker</h1>
    <ExpenseForm />
    <ExpenseList />
  </div>
</template>

<script setup>
import ExpenseForm from './components/ExpenseForm.vue'
import ExpenseList from './components/ExpenseList.vue'
</script>