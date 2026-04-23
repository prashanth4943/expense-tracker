# 📌 Expense Tracker – Full Stack Application

## 🔗 Live Application
👉 https://expense-tracker-91g7.onrender.com

---

## 🎯 Objective

To build a simple yet **real-world resilient expense tracking system** that:

- Handles unreliable networks and retries  
- Ensures data correctness (especially for money)  
- Demonstrates clean API design and frontend-backend interaction  

---

## 🏗️ Tech Stack

- **Backend:** Go (net/http)  
- **Database:** SQLite  
- **Frontend:** Vue.js (Vite)  
- **Hosting:** Render  

---

## ✅ Features Implemented

### Core Functionality

- Add new expense (amount, category, description, date)  
- View all expenses  
- Filter by category  
- Sort by:
  - Newest first  
  - Oldest first  
  - Amount (high → low)  
  - Amount (low → high)  
- Display total for current filtered list  

---

## 🧠 Key Design Decisions

### 1. Idempotency for Safe Retries

To handle:

- Duplicate submissions  
- Network retries  
- Page refreshes  

Each request includes an **idempotency key**:

- Stored alongside expense  
- Enforced via unique constraint  
- Duplicate requests return existing record instead of inserting again  

👉 Ensures **no duplicate expenses under retry conditions**

---

### 2. Money Handling (Data Correctness)

- Amount stored as **integer (paise)** instead of float  
- Avoids floating point precision issues  

**Example:**
₹410.00 → stored as 41000

- Formatting handled on frontend  

---

### 3. Allowing Any Date (Intentional Decision)

Unlike booking systems, this is a **personal expense tracker**.

Users often:
- Log expenses later  
- Backfill historical data  

👉 Therefore:
- Past dates are allowed  
- Future dates are also allowed (intentional design choice)  

This behavior is explicitly documented.

---

### 4. Validation Strategy

Validation is enforced at **both frontend and backend**.

#### Frontend:
- Required fields  
- Minimum amount > 0  
- Basic input constraints  

#### Backend (source of truth):
- Amount must be > 0  
- Category must not be empty  
- Date must be valid  
- Description length bounded  
- Idempotency key required  

👉 Ensures correctness even if frontend is bypassed  

---

### 5. API Design

#### `POST /expenses`
- Idempotent  
- Safe under retries  
- Strict validation  

#### `GET /expenses`
Supports:
- `category` filter  
- `sort` parameter  

---

### 6. Sorting Implementation

Sorting handled at **database level using safe query mapping**:

- Prevents SQL injection  
- Supports extensibility  

---

### 7. Frontend Architecture

- Component-based (Form + List)  
- Event-driven refresh using `emit`  
- Defensive handling of API responses  

Handles:
- Loading states  
- Error states  
- Empty states  

---

## ⚠️ Trade-offs (Time Constraints)

To keep scope focused:

- SQLite used instead of managed DB (e.g., Postgres)  
- No authentication system  
- Minimal styling (focus on correctness)  
- No pagination (dataset assumed small)  
- No caching layer  

---

## 🚀 If More Time Was Available

- Add authentication & multi-user support  
- Pagination & performance optimization  
- Category analytics (charts)  
- Edit/Delete expenses  
- Rate limiting & request logging  
- Unit & integration tests  

---

## 🧪 Handling Real-World Conditions

This system is designed to handle:

- Multiple submit clicks  
- Network retries  
- Page refresh after submission  
- Partial failures  

👉 Achieved via **idempotency + backend validation**



---

## ▶️ Running Locally

```bash
# Backend
go run main.go

# Frontend
cd frontend
npm install
npm run build

screeshot


<img width="1831" height="867" alt="expense-tracker" src="https://github.com/user-attachments/assets/77b2cfc5-bc42-4fde-bab1-8ec85b025d4b" />


