# Phase I: System Design

## High-level Architecture

```mermaid
flowchart TD
    U[User] --> CLI[Go Application Interface]
    CLI --> V[Input Validation]
    V --> S[Expense Service]
    S --> R[Expense Repository]
    R --> DB[(SQLite Database - Proposed)]
    S --> SUM[Summary and Reporting]
    SUM --> CLI
```

## Data Flow Diagram

```mermaid
flowchart TD
    U[User] -->|Expense details| P[Capture Expense]
    P --> V[Validate required fields and amount]
    V -->|Valid data| E[Expense Processing]
    V -->|Invalid data| ERR[Display validation error]
    E --> DB[(Expense Store - Proposed SQLite)]
    DB --> Q[Retrieve expenses / calculate totals]
    Q --> O[Show expense list and summaries]
    O --> U
    ERR --> U
```

## Main Components

- **Application interface:** accepts user actions and presents results.
- **Input validation:** checks required fields and sensible expense amounts.
- **Expense service:** handles adding, listing, and summarizing expenses.
- **Expense repository:** separates storage operations from application logic.
- **Storage:** SQLite is the proposed local persistence option; final implementation may refine the design.
- **Summary module:** calculates totals and groups expenses by category.

These diagrams describe the proposed Phase I design. They should not be interpreted as proof that all modules have been implemented or tested.
