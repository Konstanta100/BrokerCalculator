// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: copyfrom.go

package repository

import (
	"context"
)

// iteratorForBulkInsertOperations implements pgx.CopyFromSource.
type iteratorForBulkInsertOperations struct {
	rows                 []BulkInsertOperationsParams
	skippedFirstNextCall bool
}

func (r *iteratorForBulkInsertOperations) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForBulkInsertOperations) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].Figi,
		r.rows[0].InstrumentType,
		r.rows[0].Quantity,
		r.rows[0].Payment,
		r.rows[0].Currency,
		r.rows[0].Date,
		r.rows[0].AccountID,
	}, nil
}

func (r iteratorForBulkInsertOperations) Err() error {
	return nil
}

func (q *Queries) BulkInsertOperations(ctx context.Context, arg []BulkInsertOperationsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"investing", "operations"}, []string{"id", "figi", "instrument_type", "quantity", "payment", "currency", "date", "account_id"}, &iteratorForBulkInsertOperations{rows: arg})
}
