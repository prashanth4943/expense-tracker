<template>
  <div class="card">
    <h2>Add Expense</h2>

    <div v-if="error" class="alert alert--error">{{ error }}</div>
    <div v-if="success" class="alert alert--success">Expense saved!</div>

    <form @submit.prevent="submit">
      <div class="form-row">
        <label>Amount (₹)
          <input
            v-model="form.amount"
            type="number"
            step="0.01"
            min="0.01"
            placeholder="0.00"
            required
          />
        </label>

        <label>Category
          <select v-model="form.category" required>
            <option value="" disabled>Select…</option>
            <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
            <option value="__new__">+ New category</option>
          </select>
        </label>

        <label v-if="form.category === '__new__'">New Category
          <input v-model="form.newCategory" placeholder="e.g. Groceries" required />
        </label>

        <label>Date
          <input v-model="form.date" type="date" required />
        </label>
      </div>

      <label style="margin-top: 12px; display:block;">Description
        <input v-model="form.description" type="text" placeholder="What was this for?" />
      </label>

      <button type="submit" :disabled="submitting" class="btn btn--primary">
        <span v-if="submitting">Saving…</span>
        <span v-else>Add Expense</span>
      </button>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { createExpense, listCategories } from '../utils/api.js'
import { generateIdempotencyKey } from '../utils/idempotency.js'

const emit = defineEmits(['added'])

const PRESET_CATEGORIES = ['Food', 'Travel', 'Utilities', 'Entertainment', 'Health', 'Shopping', 'Other']

const form = reactive({
  amount: '',
  category: '',
  newCategory: '',
  description: '',
  date: today(),
})

const categories = ref([...PRESET_CATEGORIES])
const submitting = ref(false)
const error = ref('')
const success = ref(false)

// Each form load gets a fresh idempotency key.
// This key lives for the lifetime of this submission attempt —
// if the network fails and the user tries again, the key is refreshed
// only on a successful submission, not on network errors, so retries are safe.
let idempotencyKey = generateIdempotencyKey()

onMounted(async () => {
  try {
    const cats = await listCategories()
    // Merge server categories with presets (deduplicated)
    const merged = [...new Set([...PRESET_CATEGORIES, ...cats])]
    categories.value = merged
  } catch {
    // Non-fatal: fallback to preset list
  }
})

async function submit() {
  error.value = ''
  success.value = false

  const category = form.category === '__new__' ? form.newCategory.trim() : form.category
  if (!category) {
    error.value = 'Please enter a category name.'
    return
  }

  submitting.value = true
  try {
    await createExpense({
      amount: String(parseFloat(form.amount).toFixed(2)),
      category,
      description: form.description,
      date: form.date,
      idempotency_key: idempotencyKey,
    })

    // Success — rotate the key so the next distinct submit is a new expense.
    idempotencyKey = generateIdempotencyKey()
    success.value = true

    // Reset form
    form.amount = ''
    form.description = ''
    form.category = ''
    form.newCategory = ''
    form.date = today()

    emit('added')
    setTimeout(() => (success.value = false), 3000)
  } catch (e) {
    // Network error: keep the same idempotency key so a retry is safe.
    error.value = e.message
  } finally {
    submitting.value = false
  }
}

function today() {
  return new Date().toISOString().slice(0, 10)
}
</script>

<style>

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

label {
  font-size: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #444;
}
input,
select {
  padding: 8px 10px;
  border-radius: 6px;
  border: 1px solid #ddd;
  font-size: 14px;
  outline: none;
  transition: border 0.2s, box-shadow 0.2s;
}

input:focus,
select:focus {
  border-color: #4f46e5;
  box-shadow: 0 0 0 2px rgba(79, 70, 229, 0.1);
}

.btn {
  margin-top: 16px;
  padding: 10px 14px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s, opacity 0.2s;
}

.btn--primary {
  background: #4f46e5;
  color: white;
}

.btn--primary:hover {
  background: #4338ca;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.alert--success {
  background: #dcfce7;
  color: #166534;
}

h2 {
  margin-bottom: 12px;
}

input::placeholder {
  color: #aaa;
}

form {
  margin-top: 8px;
}

input[v-model="form.newCategory"] {
  border-color: #4f46e5;
}

</style>