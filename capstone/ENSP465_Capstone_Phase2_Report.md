# K.R. MANGALAM UNIVERSITY
## SCHOOL OF BASIC AND APPLIED SCIENCES / DEPARTMENT OF COMPUTER SCIENCE
### ENSP465 – CAPSTONE PROJECT

---

# PHASE II FINAL PROJECT REPORT

## **Project Title: Smart Expense Tracker Using Go**

**Submitted By:**
- **Shahbaz Ali** (Roll Number: `2401830007`)
- **Harshit** (Roll Number: `2401830008`)

**Submitted To:**
- **Dr. Hemraj Kumawat**

**Programme:** B.Sc. Cyber Security  
**Semester:** V  
**Academic Year:** 2026–2027  

---

## Abstract

Managing daily personal finances is a critical requirement for students and working professionals. Traditional manual bookkeeping and basic spreadsheets often suffer from human error, omitted entries, lack of categorization, and calculation inaccuracies. This project, **Smart Expense Tracker Using Go**, delivers a robust, secure, and modular personal financial tracking system implemented in the Go (Golang) programming language paired with an embedded SQLite relational database.

Guided by the requirements and architecture defined in Phase I, Phase II successfully realizes the complete end-to-end implementation. In strict alignment with cyber security best practices, the application incorporates defensive coding principles including zero-SQL-injection parameterized queries, rigorous boundary validation, safe exception handling, and a dedicated minor-currency integer precision engine (`Money int64`) to eliminate IEEE-754 floating-point rounding vulnerabilities. The software provides both an interactive terminal interface and an automated CLI command-line suite, complete with automated unit testing, financial reporting, ASCII distribution charts, and CSV/JSON export capabilities.

---

## 1. Introduction & Objectives Fulfillment

### 1.1 Problem Restatement
College students handle recurrent financial transactions—including tuition fees, hostel rent, textbooks, mess charges, local transport, and personal shopping. Without an automated, organized system, tracking small discretionary expenditures is difficult, resulting in budget overruns and untracked deficits.

### 1.2 Fulfillment of Phase I Objectives
All seven target milestones established in the Phase I report have been fully developed and verified:
1. **Transaction Entry**: Real-time recording of both income and expense items.
2. **Transaction Viewing & Filtering**: Formatted tabular displays with filtering by type, category, and date range.
3. **Categorization**: Multi-category tagging (e.g., *Food, Rent, Education, Transport, Utilities, Shopping, Salary*).
4. **Totals & Net Balance Calculation**: Real-time calculation based on:
   $$\text{Remaining Balance} = \text{Total Income} - \text{Total Expenses}$$
5. **Spending Review & Analytics**: Graphical category-wise percentage distribution and monthly financial summaries.
6. **Record Management**: Verified ID-based updating and deletion of transaction records.
7. **Defensive Error Handling & Validation**: Comprehensive pre-insertion checks providing user-friendly feedback without leaking internal database exceptions.

---

## 2. System Architecture & Design

The software follows a 5-tier modular architecture ensuring loose coupling and separation of concerns:

```
+--------------------------------------------------------------+
|                        User / Client                         |
+--------------------------------------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|               User Interface Layer (internal/ui)             |
|   - Interactive Terminal UI (interactive.go)                 |
|   - CLI Subcommands & Flags Engine (cli.go)                  |
|   - ASCII Table & Progress Chart Formatter (table.go)        |
+--------------------------------------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|            Application Logic Layer (internal/service)        |
|   - CRUD Operations & Financial Business Logic               |
|   - Aggregations, Trend Analysis & Export Pipeline           |
+--------------------------------------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|          Validation & Security Layer (internal/validator)    |
|   - Boundary & Type Checking, Positive Amount Enforcement    |
|   - ISO 8601 Date Parsing & String Sanitization              |
+--------------------------------------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|           Data Access Repository (internal/storage)          |
|   - Parameterized Prepared Statements (sqlite.go)            |
|   - Storage Interface Contract (storage.go)                  |
+--------------------------------------------------------------+
                               |
                               v
+--------------------------------------------------------------+
|           Embedded Database Storage (expenses.db)            |
+--------------------------------------------------------------+
```

---

## 3. Data Design & Database Schema

The database is built on **SQLite** (`modernc.org/sqlite`), operating entirely without external network dependencies or CGO requirements.

### 3.1 Relational Schema
```sql
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
    amount INTEGER NOT NULL CHECK(amount > 0),
    category TEXT NOT NULL,
    date TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
```

### 3.2 Data Fields Specification
| Field Name | Storage Type | Domain Constraint | Description |
| :--- | :--- | :--- | :--- |
| `id` | `INTEGER` | Primary Key, Auto-increment | Unique identifier for transaction management |
| `type` | `TEXT` | `CHECK(type IN ('income', 'expense'))` | Designates cash flow direction |
| `amount` | `INTEGER` | `amount > 0` (Minor units / Cents) | Monetary value scaled by $100$ |
| `category` | `TEXT` | Non-empty, $\le 50$ characters | Classification for spending analytics |
| `date` | `TEXT` | Format `YYYY-MM-DD` | Transaction occurrence date |
| `description` | `TEXT` | Optional, $\le 255$ characters | Narrative context |
| `created_at` | `DATETIME` | Server timestamp | Audit trail timestamp |

---

## 4. Cyber Security & Defensive Coding Measures

As part of the **B.Sc. Cyber Security** curriculum, specific secure programming practices were implemented:

### 4.1 SQL Injection Prevention
All database queries utilize parameterized SQL statements with placeholder bindings (`?`). Direct string interpolation and string concatenation were explicitly avoided:
```go
// Safe parameterized insertion:
query := `INSERT INTO transactions (type, amount, category, date, description) VALUES (?, ?, ?, ?, ?)`
_, err := r.db.Exec(query, string(t.Type), int64(t.Amount), t.Category, t.Date, t.Description)
```

### 4.2 Mitigation of Floating-Point Precision Exploits
Standard IEEE-754 binary floating-point numbers (`float32` / `float64`) are prone to representation inaccuracies (e.g., $0.1 + 0.2 = 0.30000000000000004$). In financial accounting, cumulative errors cause balance discrepancies.
- **Solution**: A custom `Money` type represents monetary quantities as 64-bit integers scaled by $100$ (e.g., $\$45.50 \rightarrow 4550$). All database arithmetic operations (`SUM`, `CHECK`) run in exact integer domain.

### 4.3 Input Validation & Defense-in-Depth
User input is validated before entering the service or storage layer:
- Type strictly validated against `income` and `expense`.
- Amount verified to be positive and below safe transactional thresholds ($100,000,000.00$).
- Dates strictly parsed against `2006-01-02` (ISO 8601).
- Descriptions bounded to 255 characters to prevent buffer and memory resource exhaustion.

### 4.4 Information Leakage Prevention
Database internal syntax errors and low-level stack traces are intercepted at the repository boundary and mapped to clean domain error types (e.g., `ErrNotFound`, `ValidationError`), preventing inadvertent system reconnaissance.

---

## 5. Implementation Modules

The source repository is structured as follows:

| Package | Path | Responsibility |
| :--- | :--- | :--- |
| `main` | `main.go` | Entry point, CLI router, database lifecycle management |
| `model` | `internal/model/` | Core domain structs (`Transaction`, `Money`, `Filter`, summaries) |
| `validator`| `internal/validator/` | Defensive input validation routines and error formatting |
| `storage` | `internal/storage/` | SQLite repository implementation with indexes and aggregate SQL queries |
| `service` | `internal/service/` | Business logic, calculations, data seeder, and CSV/JSON export |
| `ui` | `internal/ui/` | Interactive terminal UI, command-line flags, and ASCII table renderer |

---

## 6. Testing, Verification & Code Coverage

Automated testing was conducted using the Go standard testing framework (`go test`).

### 6.1 Test Execution Results
```
=== RUN   TestMoney_ConversionAndFormatting
--- PASS: TestMoney_ConversionAndFormatting (0.00s)
=== RUN   TestParseTransactionType
--- PASS: TestParseTransactionType (0.00s)
PASS
ok      smart-expense-tracker/internal/model    1.037s

=== RUN   TestTrackerService_AddAndGet
--- PASS: TestTrackerService_AddAndGet (0.06s)
=== RUN   TestTrackerService_InvalidInput
--- PASS: TestTrackerService_InvalidInput (0.06s)
=== RUN   TestTrackerService_BalanceCalculation
--- PASS: TestTrackerService_BalanceCalculation (0.08s)
=== RUN   TestTrackerService_UpdateAndDelete
--- PASS: TestTrackerService_UpdateAndDelete (0.08s)
=== RUN   TestTrackerService_ExportAndSeed
--- PASS: TestTrackerService_ExportAndSeed (0.20s)
PASS
ok      smart-expense-tracker/internal/service  2.386s

=== RUN   TestSQLiteRepository_CreateAndGet
--- PASS: TestSQLiteRepository_CreateAndGet (0.06s)
=== RUN   TestSQLiteRepository_UpdateAndDelete
--- PASS: TestSQLiteRepository_UpdateAndDelete (0.08s)
=== RUN   TestSQLiteRepository_TotalsAndBreakdowns
--- PASS: TestSQLiteRepository_TotalsAndBreakdowns (0.08s)
PASS
ok      smart-expense-tracker/internal/storage  2.586s

=== RUN   TestValidateTransaction_Valid
--- PASS: TestValidateTransaction_Valid (0.00s)
=== RUN   TestValidateTransaction_DefaultDate
--- PASS: TestValidateTransaction_DefaultDate (0.00s)
=== RUN   TestValidateTransaction_InvalidFields
--- PASS: TestValidateTransaction_InvalidFields (0.00s)
=== RUN   TestValidateID
--- PASS: TestValidateID (0.00s)
PASS
ok      smart-expense-tracker/internal/validator        1.210s
```
**Result**: 100% test pass rate across all modules with zero failures.

---

## 7. Execution Guide & Sample Outputs

### 7.1 Running the Application
The project includes a one-click launcher for Windows:
```cmd
run.bat
```
Alternatively, execute directly using the compiled binary or Go runtime:
```powershell
.\expense-tracker.exe
# OR
go run main.go
```

### 7.2 Sample Financial Summary Output
```
==================================================
             FINANCIAL BALANCE SUMMARY            
==================================================
 Total Transactions: 20
 Total Income:       $9050.00
 Total Expenses:     $2355.75
--------------------------------------------------
 Net Balance:        +$6694.25 (Surplus)
==================================================
 Formula: Remaining Balance = Total Income - Total Expenses
```

### 7.3 Sample Visual Category Distribution
```
================= EXPENSE BREAKDOWN BY CATEGORY =================
+------------------+------------+---------+--------------------------+
| CATEGORY         |      TOTAL |   SHARE | DISTRIBUTION             |
+------------------+------------+---------+--------------------------+
| Rent             |    1350.00 |   57.3% | ###########              |
| Education        |     390.00 |   16.6% | ###                      |
| Food             |     340.75 |   14.5% | ##                       |
| Shopping         |     140.00 |    5.9% | #                        |
| Transport        |      75.00 |    3.2% | #                        |
| Utilities        |      60.00 |    2.5% | #                        |
+------------------+------------+---------+--------------------------+
```

### 7.4 Sample Monthly Trends Table
```
================== MONTHLY FINANCIAL TRENDS ==================
+----------+--------------+--------------+----------------+
| MONTH    |       INCOME |     EXPENSES |    NET BALANCE |
+----------+--------------+--------------+----------------+
| 2026-11  |      4700.00 |       970.00 |       +3730.00 |
| 2026-10  |      2350.00 |       695.25 |       +1654.75 |
| 2026-09  |      2000.00 |       690.50 |       +1309.50 |
+----------+--------------+--------------+----------------+
```

---

## 8. Conclusion & Future Scope

### 8.1 Conclusion
The **Smart Expense Tracker Using Go** successfully transitions from the conceptual design in Phase I to a production-ready, fully verified application in Phase II. It demonstrates how modern systems languages like Go can combine high execution performance, memory safety, and clean software architecture to solve everyday financial management problems.

### 8.2 Future Scope
Potential expansions for subsequent development cycles include:
1. **Desktop Graphical User Interface (GUI)**: Implementing a native desktop GUI using Wails (Go + Web frontend) or Fyne.
2. **Multi-User Authentication**: Incorporating Argon2id password hashing and session tokens for shared family or room expense tracking.
3. **Automated Bank Statement Ingestion**: Parsing OFX and standard PDF/CSV bank statement records for automated reconciliation.
4. **Predictive Budget Alerting**: Machine learning-based budget overrun projections and Telegram/Email webhook notifications.

---

## 9. References

1. The Go Authors, *"The Go Programming Language Documentation"*, 2024. [Online]. Available: https://go.dev/doc/
2. SQLite Authors, *"SQLite Documentation and Query Language"*, 2024. [Online]. Available: https://www.sqlite.org/docs.html
3. OWASP Foundation, *"SQL Injection Prevention Cheat Sheet"*, OWASP Security Guidelines, 2023.
4. IEEE, *"IEEE Standard for Floating-Point Arithmetic (IEEE 754-2019)"*, IEEE Computer Society, 2019.

