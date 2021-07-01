package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store interface 代替结构体
type Store interface {
	// 将 Querier interface 插入：使 Store interface 具有其所有功能
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

/// SQLStore 提供了执行数据库查询和事务的所有函数
// 这将是 Store interface 的真正实现,与 SQL 数据库对话
type SQLStore struct {
	db *sql.DB
	*Queries
}

/// NewStore creates a new store
/// 该函数不应该返回一个指针，而只是一个 Store interface: 即SQLStore需实现此接口所有功能
func NewStore(db *sql.DB) Store {
	// SQLStore 结构 实现 Store interface 的所有功能
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

/// ExecTx 数据库事务通用函数，具体业务函数在回调内执行
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		var transferParams CreateTransferParams
		transferParams.FromAccountID = arg.FromAccountID
		transferParams.ToAccountID = arg.ToAccountID
		transferParams.Amount = arg.Amount

		result.Transfer, err = q.CreateTransfer(ctx, transferParams)
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      account1.ID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      account2.ID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
