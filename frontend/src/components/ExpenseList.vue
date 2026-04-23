<template>
  <div class="card">
    <div class="list-header">
      <h2>Expenses</h2>
      <div class="total-badge">Total: {{ total }}</div>
    </div>

    <!-- Filters -->
    <div class="filters">
      <label>
        Category
        <select v-model="filterCategory" @change="load">
          <option value="">All</option>
          <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
        </select>
      </label>

      <label>
        Sort
        <!-- <select v-model="sort" @change="load">
          <option value="date_desc">Newest first</option>
        </select> -->
        <select v-model="sort" @change="load">
  <option value="date_desc">Newest first</option>
  <option value="date_asc">Oldest first</option>
  <option value="amount_desc">Amount (high → low)</option>
  <option value="amount_asc">Amount (low → high)</option>
</select>
      </label>
    </div>

    <!-- States -->
    <div v-if="loading" class="state-msg">Loading…</div>
    <div v-else-if="error" class="alert alert--error">{{ error }}</div>
    <div v-else-if="expenses.length === 0" class="state-msg empty">No expenses yet.</div>

    <!-- Table -->
    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Category</th>
            <th>Description</th>
            <th class="right">Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in expenses" :key="e.id">
            <td>{{ formatDate(e.date) }}</td>
            <td><span class="chip">{{ e.category }}</span></td>
            <td class="desc">{{ e.description || '—' }}</td>
            <td class="right amount">{{ e.amount }}</td>
          </tr>
        </tbody>
        <tfoot>
          <tr>
            <td colspan="3"><strong>Total</strong></td>
            <td class="right amount"><strong>{{ total }}</strong></td>
          </tr>
        </tfoot>
      </table>
    </div>

    <!-- Category summary (nice-to-have) -->
    <div v-if="categorySummary.length" class="summary">
      <h3>By Category</h3>
      <div class="summary-grid">
        <div v-for="s in categorySummary" :key="s.category" class="summary-item">
          <span class="chip">{{ s.category }}</span>
          <span class="amount">{{ s.total }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { listExpenses, listCategories } from '../utils/api.js'

const props = defineProps({
  refreshTrigger: Number,
})

const expenses = ref([])
const categories = ref([])
const total = ref('₹0.00')
const loading = ref(false)
const error = ref('')
const filterCategory = ref('')
const sort = ref('date_desc')

onMounted(() => {
  load()
  loadCategories()
})

// watch(() => props.refreshTrigger, load)
watch(() => props.refreshTrigger, () => {
  load()
})
watch([filterCategory, sort], () => {
  load()
})


async function load() {
  loading.value = true
  error.value = ''
  try {
    const data = await listExpenses({
      category: filterCategory.value,
      sort: sort.value,
    })

    expenses.value = data.expenses || []
    // total.value = formatPaise(data.total || 0)
    if (typeof data.total_paise === 'number') {
  total.value = formatPaise(data.total_paise)
} else if (data.total) {
  total.value = data.total
} else {
  total.value = '₹0.00'
}

  } catch (e) {
    error.value = 'Failed to load expenses. Please try again.'
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    categories.value = await listCategories()
  } catch {
    // Non-fatal
  }
}

// Compute per-category totals from current visible list
const categorySummary = computed(() => {
  const map = {}
  for (const e of expenses.value) {
    if (!map[e.category]) map[e.category] = 0
    // map[e.category] += e.amount_paise
    map[e.category] += e.amount_paise || 0
  }
  return Object.entries(map).map(([category, paise]) => ({
    category,
    total: formatPaise(paise),
  })).sort((a, b) => a.category.localeCompare(b.category))
})

function formatDate(d) {
  if (!d) return ''
  const [y, m, day] = d.split('-')
  return `${day}/${m}/${y}`
}

function formatPaise(paise) {
  const rupees = Math.floor(paise / 100)
  const cents = paise % 100
  return `₹${rupees}.${String(cents).padStart(2, '0')}`
}
</script>

<style>

body {
  font-family: system-ui, -apple-system, sans-serif;
  background: #f6f7fb;
  margin: 0;
  padding: 20px;
  color: #222;
}

.card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  max-width: 900px;
  margin: 0 auto 20px;
}
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.total-badge {
  background: #eef2ff;
  color: #4f46e5;
  padding: 6px 12px;
  border-radius: 999px;
  font-weight: 600;
}

.filters {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filters label {
  font-size: 14px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

select {
  padding: 6px 10px;
  border-radius: 6px;
  border: 1px solid #ddd;
}

.table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 10px;
  border-bottom: 1px solid #eee;
  text-align: left;
}

th {
  font-size: 13px;
  color: #666;
  text-transform: uppercase;
}

.right {
  text-align: right;
}

.amount {
  font-weight: 600;
  color: #111;
}

.chip {
  background: #f1f5f9;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.state-msg {
  text-align: center;
  padding: 20px;
  color: #666;
}

.state-msg.empty {
  color: #999;
}

.alert {
  padding: 10px;
  border-radius: 6px;
  margin-bottom: 12px;
}

.alert--error {
  background: #fee2e2;
  color: #b91c1c;
}

.summary {
  margin-top: 20px;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 10px;
}

.summary-item {
  background: #f9fafb;
  padding: 10px;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
}
</style>