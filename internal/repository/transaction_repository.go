package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"finapp/internal/model"
)

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Upsert(ctx context.Context, tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, account_id, pluggy_transaction_id, description, amount, date, type, category_id, pluggy_category, notes, tags, is_recurring, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,NOW(),NOW())
		ON CONFLICT (pluggy_transaction_id) DO UPDATE
		SET description      = EXCLUDED.description,
		    amount           = EXCLUDED.amount,
		    date             = EXCLUDED.date,
		    pluggy_category  = EXCLUDED.pluggy_category,
		    updated_at       = NOW()
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query,
		tx.ID, tx.UserID, tx.AccountID, tx.PluggyTransactionID,
		tx.Description, tx.Amount, tx.Date, tx.Type,
		tx.CategoryID, tx.PluggyCategory, tx.Notes, tx.Tags, tx.IsRecurring,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
}

func (r *transactionRepository) BulkUpsert(ctx context.Context, txs []model.Transaction) error {
	for i := range txs {
		if err := r.Upsert(ctx, &txs[i]); err != nil {
			return fmt.Errorf("upsert transaction %s: %w", txs[i].ID, err)
		}
	}
	return nil
}

func (r *transactionRepository) FindByFilter(ctx context.Context, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	args := []interface{}{filter.UserID}
	conditions := []string{"user_id = $1"}
	idx := 2

	if filter.AccountID != nil {
		conditions = append(conditions, fmt.Sprintf("account_id = $%d", idx))
		args = append(args, *filter.AccountID)
		idx++
	}
	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", idx))
		args = append(args, *filter.CategoryID)
		idx++
	}
	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", idx))
		args = append(args, *filter.Type)
		idx++
	}
	if filter.From != nil {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", idx))
		args = append(args, *filter.From)
		idx++
	}
	if filter.To != nil {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", idx))
		args = append(args, *filter.To)
		idx++
	}

	where := strings.Join(conditions, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions WHERE %s`, where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	// Pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	dataArgs := append(args, limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT id, user_id, account_id, pluggy_transaction_id, description, amount, date, type,
		       category_id, pluggy_category, notes, tags, is_recurring, created_at, updated_at
		FROM transactions
		WHERE %s
		ORDER BY date DESC, created_at DESC
		LIMIT $%d OFFSET $%d`, where, idx, idx+1)

	rows, err := r.db.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var txList []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.PluggyTransactionID,
			&t.Description, &t.Amount, &t.Date, &t.Type,
			&t.CategoryID, &t.PluggyCategory, &t.Notes, &t.Tags,
			&t.IsRecurring, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan transaction: %w", err)
		}
		txList = append(txList, t)
	}
	return txList, total, rows.Err()
}

func (r *transactionRepository) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, pluggy_transaction_id, description, amount, date, type,
		       category_id, pluggy_category, notes, tags, is_recurring, created_at, updated_at
		FROM transactions WHERE id = $1 AND user_id = $2`
	var t model.Transaction
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.AccountID, &t.PluggyTransactionID,
		&t.Description, &t.Amount, &t.Date, &t.Type,
		&t.CategoryID, &t.PluggyCategory, &t.Notes, &t.Tags,
		&t.IsRecurring, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find transaction: %w", err)
	}
	return &t, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *model.Transaction) error {
	_, err := r.db.Exec(ctx,
		`UPDATE transactions SET category_id=$1, notes=$2, tags=$3, updated_at=NOW() WHERE id=$4 AND user_id=$5`,
		tx.CategoryID, tx.Notes, tx.Tags, tx.ID, tx.UserID,
	)
	return err
}

func (r *transactionRepository) SumByPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) (model.PeriodSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'DEBIT'  THEN amount ELSE 0 END), 0) AS expenses
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3`
	var s model.PeriodSummary
	err := r.db.QueryRow(ctx, query, userID, from, to).Scan(&s.TotalIncome, &s.TotalExpenses)
	return s, err
}

func (r *transactionRepository) SumByCategoryAndPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CategorySpending, error) {
	query := `
		SELECT
			t.category_id::text,
			COALESCE(c.name, 'Sem categoria') AS category_name,
			SUM(t.amount) AS amount,
			COUNT(*) AS count
		FROM transactions t
		LEFT JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1 AND t.date >= $2 AND t.date <= $3 AND t.type = 'DEBIT'
		GROUP BY t.category_id, c.name
		ORDER BY amount DESC`

	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var total float64
	var items []model.CategorySpending
	for rows.Next() {
		var cs model.CategorySpending
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.Amount, &cs.Count); err != nil {
			return nil, err
		}
		total += cs.Amount
		items = append(items, cs)
	}

	for i := range items {
		if total > 0 {
			items[i].Percentage = (items[i].Amount / total) * 100
		}
	}
	return items, rows.Err()
}

func (r *transactionRepository) GetMonthlyAggregates(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]model.CashFlowPoint, error) {
	query := `
		SELECT
			TO_CHAR(DATE_TRUNC('month', date), 'YYYY-MM') AS month,
			COALESCE(SUM(CASE WHEN type = 'CREDIT' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'DEBIT'  THEN amount ELSE 0 END), 0) AS expenses
		FROM transactions
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY DATE_TRUNC('month', date)
		ORDER BY DATE_TRUNC('month', date)`

	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []model.CashFlowPoint
	for rows.Next() {
		var p model.CashFlowPoint
		if err := rows.Scan(&p.Month, &p.Income, &p.Expenses); err != nil {
			return nil, err
		}
		p.Net = p.Income - p.Expenses
		points = append(points, p)
	}
	return points, rows.Err()
}

func (r *transactionRepository) FindRecurring(ctx context.Context, userID uuid.UUID, lookbackMonths int) (float64, float64, error) {
	// Detect recurring income and expense averages from last N months
	query := `
		WITH monthly AS (
			SELECT
				DATE_TRUNC('month', date) AS month,
				type,
				SUM(amount) AS total
			FROM transactions
			WHERE user_id = $1
			  AND date >= NOW() - ($2::int * INTERVAL '1 month')
			GROUP BY DATE_TRUNC('month', date), type
		)
		SELECT
			COALESCE(AVG(CASE WHEN type = 'CREDIT' THEN total END), 0) AS avg_income,
			COALESCE(AVG(CASE WHEN type = 'DEBIT'  THEN total END), 0) AS avg_expenses
		FROM monthly`

	var avgIncome, avgExpenses float64
	err := r.db.QueryRow(ctx, query, userID, lookbackMonths).Scan(&avgIncome, &avgExpenses)
	return avgIncome, avgExpenses, err
}
