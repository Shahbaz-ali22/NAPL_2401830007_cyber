package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"smart-expense-tracker/internal/model"
)

// GenerateHTMLReport generates a modern, self-contained visual HTML dashboard.
func (s *TrackerService) GenerateHTMLReport(w io.Writer) error {
	txs, err := s.repo.List(model.Filter{})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	totals, err := s.GetBalanceSummary()
	if err != nil {
		return fmt.Errorf("failed to fetch balance: %w", err)
	}

	catExpenses, _ := s.GetExpenseBreakdownByCategory()
	monthlyTrends, _ := s.GetMonthlyBreakdown()

	// Prepare JSON data for Chart.js
	type CatChartData struct {
		Labels []string  `json:"labels"`
		Values []float64 `json:"values"`
	}
	var catData CatChartData
	for _, c := range catExpenses {
		catData.Labels = append(catData.Labels, c.Category)
		catData.Values = append(catData.Values, c.Total.ToFloat())
	}
	catJSON, _ := json.Marshal(catData)

	type TrendChartData struct {
		Months  []string  `json:"months"`
		Income  []float64 `json:"income"`
		Expense []float64 `json:"expense"`
		Balance []float64 `json:"balance"`
	}
	var trendData TrendChartData
	// Reverse order for chronological display
	for i := len(monthlyTrends) - 1; i >= 0; i-- {
		m := monthlyTrends[i]
		trendData.Months = append(trendData.Months, m.Month)
		trendData.Income = append(trendData.Income, m.Income.ToFloat())
		trendData.Expense = append(trendData.Expense, m.Expense.ToFloat())
		trendData.Balance = append(trendData.Balance, m.Balance.ToFloat())
	}
	trendJSON, _ := json.Marshal(trendData)

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Smart Expense Tracker - Capstone Dashboard</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
  <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
  <style>
    body { background-color: #f4f6f9; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; }
    .header-box { background: linear-gradient(135deg, #1e3c72 0%%, #2a5298 100%%); color: white; padding: 2rem; border-radius: 12px; margin-bottom: 2rem; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
    .stat-card { border: none; border-radius: 10px; transition: transform 0.2s; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
    .stat-card:hover { transform: translateY(-4px); }
    .chart-card { border: none; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.06); margin-bottom: 2rem; background: white; padding: 1.5rem; }
    .badge-income { background-color: #28a745; }
    .badge-expense { background-color: #dc3545; }
    .table-container { background: white; border-radius: 12px; padding: 1.5rem; box-shadow: 0 2px 10px rgba(0,0,0,0.06); }
  </style>
</head>
<body>
  <div class="container py-4">
    <!-- Header -->
    <div class="header-box text-center">
      <h5 class="text-uppercase tracking-wide opacity-75">K.R. Mangalam University &bull; ENSP465 Capstone Project</h5>
      <h1 class="display-6 fw-bold mb-2">Smart Expense Tracker Dashboard</h1>
      <p class="mb-0"><strong>Submitted By:</strong> Shahbaz Ali (2401830007) &amp; Harshit (2401830008) | <strong>Supervisor:</strong> Dr. Hemraj Kumawat</p>
    </div>

    <!-- Summary Metrics -->
    <div class="row g-3 mb-4">
      <div class="col-md-3">
        <div class="card stat-card p-3 bg-white">
          <small class="text-muted fw-bold">TOTAL TRANSACTIONS</small>
          <h2 class="fw-bold text-dark mt-1">%d</h2>
          <small class="text-secondary">Logged in SQLite</small>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stat-card p-3 bg-white border-start border-success border-4">
          <small class="text-muted fw-bold">TOTAL INCOME</small>
          <h2 class="fw-bold text-success mt-1">₹%s</h2>
          <small class="text-success">&bull; Total Inflow</small>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stat-card p-3 bg-white border-start border-danger border-4">
          <small class="text-muted fw-bold">TOTAL EXPENSES</small>
          <h2 class="fw-bold text-danger mt-1">₹%s</h2>
          <small class="text-danger">&bull; Total Outflow</small>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card stat-card p-3 bg-white border-start border-primary border-4">
          <small class="text-muted fw-bold">NET BALANCE</small>
          <h2 class="fw-bold text-primary mt-1">₹%s</h2>
          <small class="text-muted">Formula: Income - Expenses</small>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="row">
      <div class="col-lg-5">
        <div class="chart-card">
          <h5 class="fw-bold mb-3">Expenses by Category</h5>
          <canvas id="categoryChart" height="260"></canvas>
        </div>
      </div>
      <div class="col-lg-7">
        <div class="chart-card">
          <h5 class="fw-bold mb-3">Monthly Trends (Income vs Expenses)</h5>
          <canvas id="trendChart" height="260"></canvas>
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <div class="table-container">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="fw-bold mb-0">Recorded Transactions</h5>
        <input type="text" id="tableSearch" class="form-control form-control-sm w-auto" placeholder="Search transactions...">
      </div>
      <div class="table-responsive">
        <table class="table table-hover align-middle" id="txTable">
          <thead class="table-light">
            <tr>
              <th>ID</th>
              <th>Type</th>
              <th>Amount</th>
              <th>Category</th>
              <th>Date</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>`,
		totals.Count,
		totals.TotalIncome.String(),
		totals.TotalExpense.String(),
		totals.Balance.String(),
	)

	for _, tx := range txs {
		badgeClass := "badge-expense"
		sign := "-"
		if tx.Type == model.TypeIncome {
			badgeClass = "badge-income"
			sign = "+"
		}
		htmlContent += fmt.Sprintf(`
            <tr>
              <td><span class="text-muted fw-bold">#%d</span></td>
              <td><span class="badge %s text-uppercase">%s</span></td>
              <td class="fw-bold">%s₹%s</td>
              <td><span class="badge bg-secondary">%s</span></td>
              <td>%s</td>
              <td class="text-secondary">%s</td>
            </tr>`,
			tx.ID,
			badgeClass,
			tx.Type,
			sign,
			tx.Amount.String(),
			tx.Category,
			tx.Date,
			tx.Description,
		)
	}

	htmlContent += fmt.Sprintf(`
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <script>
    // Search filter
    document.getElementById('tableSearch').addEventListener('keyup', function() {
      const q = this.value.toLowerCase();
      const rows = document.querySelectorAll('#txTable tbody tr');
      rows.forEach(r => {
        r.style.display = r.innerText.toLowerCase().includes(q) ? '' : 'none';
      });
    });

    // Category Chart
    const catData = %s;
    new Chart(document.getElementById('categoryChart'), {
      type: 'doughnut',
      data: {
        labels: catData.labels,
        datasets: [{
          data: catData.values,
          backgroundColor: ['#4e73df', '#1cc88a', '#36b9cc', '#f6c23e', '#e74a3b', '#6f42c1', '#fd7e14', '#20c997']
        }]
      },
      options: { responsive: true, plugins: { legend: { position: 'bottom' } } }
    });

    // Monthly Trend Chart
    const trendData = %s;
    new Chart(document.getElementById('trendChart'), {
      type: 'bar',
      data: {
        labels: trendData.months,
        datasets: [
          { label: 'Income', data: trendData.income, backgroundColor: '#1cc88a' },
          { label: 'Expenses', data: trendData.expense, backgroundColor: '#e74a3b' }
        ]
      },
      options: { responsive: true, scales: { y: { beginAtZero: true } } }
    });
  </script>
</body>
</html>`,
		catJSON,
		trendJSON,
	)

	_, err = w.Write([]byte(htmlContent))
	return err
}

// OpenBrowser attempts to open the default web browser to view the generated file.
func OpenBrowser(target string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", target}
	case "darwin":
		cmd = "open"
		args = []string{target}
	default: // linux, freebsd
		cmd = "xdg-open"
		args = []string{target}
	}
	_ = exec.Command(cmd, args...).Start()
}
