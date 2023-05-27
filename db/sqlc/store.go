package db

import (
	"context"
	"database/sql"
)

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(queries *Queries) error {

		var err error
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		// creating entry for FromAccount
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		// creating entry for toAccount
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		//  Update the account 1 for deduction
		acc1, err := queries.GetAccountForUpdate(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.FromAccountID,
			Balance: acc1.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		//  Update the account 2 for addition
		acc2, err := queries.GetAccountForUpdate(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
			ID:      args.ToAccountID,
			Balance: acc2.Balance + args.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
