# Smart Expense Tracker Using Go

**ENSP465 – Capstone Project | Phase I Implementation**  
**Programme:** B.Sc. Cyber Security (Semester V)  
**Institution:** K.R. Mangalam University  
**Supervisor:** Dr. Hemraj Kumawat  

**Authors / Submitted By:**
- **Shahbaz Ali** (Roll Number: `2401830007`)
- **Harshit** (Roll Number: `2401830008`)

---

## 1. Project Overview

The **Smart Expense Tracker** is a robust, modular personal finance management application written in **Go (Golang)** with an embedded **SQLite** storage engine. Designed specifically for students and individuals, it provides a reliable, secure digital system to record financial transactions, categorize expenditures, prevent calculation errors, and visualize spending habits through interactive and CLI-driven interfaces.

This implementation directly fulfills all requirements, architectural layers, data schemas, and security controls defined in the **ENSP465 Capstone Phase I Report**.

---

## 2. Key Features

1. **Add Transactions**: Record new income and expense entries with amount, category, date, and optional description.
2. **View & Filter Transactions**: Display transactions in clean ASCII tables with multi-criteria filtering (by Type: income/expense, Category, and Date range).
3. **Categorise Transactions**: Organize spending into custom or standard categories (*Food, Transport, Education, Shopping, Utilities, Salary, etc.*).
4. **Calculate Totals & Net Balance**:
   $$\text{Remaining Balance} = \text{Total Income} - \text{Total Expenses}$$
   Displays financial position with surplus or deficit indicators.
5. **Spending & Trend Review**:
   - Visual category distribution breakdown with dynamic ASCII bar charts and percentage calculations.
   - Monthly trend reports grouping income, expenses, and net balance by month (`YYYY-MM`).
6. **Record Management**: Update or delete any existing transaction by its unique ID with verification.
7. **Dual Interface**:
   - **Interactive Terminal UI**: Intuitive numeric menu for live demonstrations, testing, and daily tracking.
   - **CLI Command / Flag Mode**: Direct commands for fast workflows and automation scripts.
8. **Security & Data Integrity** (B.Sc. Cyber Security Focus):
   - **Zero SQL Injection Risk**: 100% parameterized SQL queries (`?` prepared statements).
   - **Monetary Precision Model**: Financial amounts stored internally as minor currency units (`Money int64` / cents/paise) to prevent IEEE-754 floating-point drift and rounding exploits.
   - **Input Sanitization & Validation**: Rejects non-positive amounts, malformed dates, overflow amounts, and oversized strings.
   - **Safe Error Handling**: User-facing errors are sanitized without leaking database internals or stack traces.

---

## 3. System Architecture

The project follows the layered architecture defined in Figure 1 and Figure 2 of the Phase I report:

```
+-------------------------------------------------------+
|                    User / Client                      |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|   User Interface Layer (internal/ui)                  |
|   - Interactive Terminal UI (interactive.go)          |
|   - CLI Subcommands & Flags (cli.go)                  |
|   - Formatted Tables & ASCII Visualizations (table.go)|
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|   Application Logic Layer (internal/service)          |
|   - Add, View, Filter, Update, Delete Transactions    |
|   - Balance Calculation & Category/Monthly Summary    |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|   Validation & Error Handling (internal/validator)    |
|   - Strict Type & Range Checks                        |
|   - Date & String Bounds Enforcement                  |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|   Data Access Layer (internal/storage)                |
|   - Storage Interface (storage.go)                    |
|   - Parameterized SQLite Engine (sqlite.go)           |
+-------------------------------------------------------+
                           |
                           v
+-------------------------------------------------------+
|   Persistent Storage (expenses.db)                    |
+-------------------------------------------------------+
```

---

## 4. Technology Stack

| Component | Technology | Rationale |
| :--- | :--- | :--- |
| **Language** | Go (Golang 1.22+) | Statically typed, fast compilation, built-in concurrency, high reliability |
| **Storage Engine** | SQLite (`modernc.org/sqlite`) | Serverless, zero-configuration, cross-platform pure-Go (no CGO required) |
| **Data Handling** | Standard Library (`database/sql`, `fmt`, `time`) | High performance, minimal external dependencies |
| **Testing** | Go Testing Package (`go test`) | Built-in unit and integration test framework |

---

## 5. Project Structure

```
capstone/
├── cmd/
├── internal/
│   ├── model/
│   │   ├── transaction.go         # Domain models, Money representation & math
│   │   └── transaction_test.go    # Money precision & type parsing unit tests
│   ├── validator/
│   │   ├── validator.go           # Input validation rules & error handling
│   │   └── validator_test.go      # Boundary, invalid input & ID unit tests
│   ├── storage/
│   │   ├── storage.go             # Repository interface definition
│   │   ├── sqlite.go              # SQLite repository with parameterized queries
│   │   └── sqlite_test.go         # Storage CRUD & aggregate unit tests
│   ├── service/
│   │   ├── service.go             # Core application & financial calculation logic
│   │   └── service_test.go        # Business logic & balance calculation tests
│   └── ui/
│       ├── cli.go                 # Non-interactive CLI command flags handler
│       ├── interactive.go         # Interactive terminal menu loop
│       └── table.go               # Formatted table and ASCII bar chart renderer
├── go.mod                         # Go module definitions
├── go.sum                         # Cryptographic checksums of dependencies
├── main.go                        # Application entry point & router
└── README.md                      # Complete documentation & user guide
```

---

## 6. Getting Started

### Prerequisites
- [Go (1.20 or newer)](https://go.dev/dl/) installed.

### Build Executable
```bash
# Compile binary
go build -o expense-tracker.exe main.go
```

---

## 7. Usage Guide

### Mode 1: Interactive Menu (Default)
To launch the interactive terminal interface, run the application with no arguments:

```bash
.\expense-tracker.exe
# OR
go run main.go
```

**Interactive Menu Options:**
```
==================================================
         SMART EXPENSE TRACKER (Go CLI)           
     K.R. Mangalam University - Capstone Project  
==================================================
 [1] Add Transaction (Income / Expense)
 [2] View All Transactions
 [3] Filter Transactions (Type / Category / Date)
 [4] Calculate Totals & Remaining Balance
 [5] Review Spending (Category & Monthly Summaries)
 [6] Update an Existing Transaction
 [7] Delete a Transaction
 [8] View Registered Categories
 [9] Export Transactions (CSV / JSON)
 [10] Load Sample Demo Data
 [0] Exit Application
--------------------------------------------------
```

---

### Mode 2: Command-Line Subcommands (CLI Mode)

#### 1. Add Transactions
```bash
# Add Income
.\expense-tracker.exe add --type income --amount 2500.00 --category "Salary" --desc "Monthly scholarship/stipend"

# Add Expenses
.\expense-tracker.exe add --type expense --amount 45.50 --category "Food" --desc "Groceries"
.\expense-tracker.exe add --type expense --amount 30.00 --category "Transport" --desc "Metro card"
.\expense-tracker.exe add --type expense --amount 120.00 --category "Education" --desc "Security textbook"
```

#### 2. View Transactions & Filtering
```bash
# List all transactions
.\expense-tracker.exe list

# Filter by type
.\expense-tracker.exe list --type expense
.\expense-tracker.exe list --type income

# Filter by category
.\expense-tracker.exe list --category Food

# Filter by date range
.\expense-tracker.exe list --from 2026-11-01 --to 2026-11-10
```

#### 3. Calculate Balance & Totals
```bash
.\expense-tracker.exe balance
```
*Output:*
```
==================================================
             FINANCIAL BALANCE SUMMARY            
==================================================
 Total Transactions: 4
 Total Income:       $2500.00
 Total Expenses:     $195.50
--------------------------------------------------
 Net Balance:        +$2304.50 (Surplus)
==================================================
 Formula: Remaining Balance = Total Income - Total Expenses
```

#### 4. Review Spending & Visual Distribution
```bash
.\expense-tracker.exe summary
```
*Includes category percentage breakdown, ASCII progress bars, and monthly comparison.*

#### 5. Update a Transaction
```bash
.\expense-tracker.exe update --id 2 --amount 55.00 --desc "Updated grocery bill"
```

#### 6. Delete a Transaction
```bash
.\expense-tracker.exe delete --id 3
```

#### 7. View Registered Categories
```bash
.\expense-tracker.exe categories
```

#### 8. Export Transactions (CSV / JSON)
```bash
# Export to CSV (Excel compatible)
.\expense-tracker.exe export --format csv --file my_expenses.csv

# Export to JSON
.\expense-tracker.exe export --format json --file my_expenses.json
```

#### 9. Load Sample Demo Data
```bash
# Instantly populates realistic student expenses across 3 months
.\expense-tracker.exe seed
```

---

## 8. Running Automated Unit Tests

Comprehensive unit tests cover all layers (validation, model arithmetic, storage CRUD, and service logic):

```bash
# Run all project tests
go test -v ./...
```

*All test suites run in memory or isolated temporary files to ensure reproducible, clean test runs.*

---

## 9. Alignment with Capstone Phase I Deliverables

| Phase I Specification | Implementation Detail | Status |
| :--- | :--- | :--- |
| **Income / Expense Entry** | `service.AddTransaction` with full validation | Complete |
| **Category Grouping** | `storage.GetCategoryBreakdown` with percentage & visual distribution | Complete |
| **Totals & Net Balance** | $\text{Balance} = \text{Total Income} - \text{Total Expenses}$ | Complete |
| **Update & Delete Records** | Verified ID-based updating and deletion | Complete |
| **Error & Input Validation** | `internal/validator` bounds checks and non-zero amount rules | Complete |
| **Relational Data Storage** | SQLite with table indexing and parameterized queries | Complete |
| **Security Best Practices** | SQL injection protection, safe error messages, integer currency units | Complete |
| **Automated Testing** | Comprehensive `go test` coverage across all packages | Complete |

