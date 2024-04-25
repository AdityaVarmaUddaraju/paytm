package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDuplicateAccount = errors.New("account already exists")
	ErrNoAccount = errors.New("account does not exist")
	ErrInsuffientBalance = errors.New("insufficient funds")
)


type AccountModel struct {
	DB *sql.DB
}

type Account struct {
	UserID int64 `json:"user_id"`
	Balance int `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *AccountModel) CreateAccount(user_id int64) error {
	query := `
		INSERT INTO accounts(user_id)
		VALUES ($1)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, user_id)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"accounts_pkey\"":
			return ErrDuplicateAccount
		default:
			return err
		}
	}

	return nil
}

func (m *AccountModel) AddMoney(user_id int64,amount int) error {
	query := `
	UPDATE accounts
	SET balance = balance + $1
	WHERE user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, amount, user_id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoAccount
	}

	return nil
}

func (m *AccountModel) TransferMoney(fromUserID, toUserID int64, amount int) error {
	
	toQuery := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE user_id = $2`

	fromQuery := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE user_id = $2`

	fromBalanceQuery := `
		SELECT user_id, balance
		FROM accounts
		where user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var fromAccount Account

	err = tx.QueryRowContext(ctx, fromBalanceQuery, fromUserID).Scan(&fromAccount.UserID, &fromAccount.Balance)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoAccount
		default:
			return err
		}
	}

	if fromAccount.Balance < amount {
		return ErrInsuffientBalance
	}

	fromResult, err := tx.ExecContext(ctx, fromQuery, amount, fromUserID)

	if err != nil {
		return err
	}

	fromRowsAffected, err := fromResult.RowsAffected()

	if err != nil {
		return err
	}

	if fromRowsAffected == 0 {
		return ErrNoAccount
	}

	toResult, err := tx.ExecContext(ctx, toQuery, amount, toUserID)

	if err != nil {
		return err
	}

	toRowsAffected, err := toResult.RowsAffected()

	if err != nil {
		return err
	}

	if toRowsAffected == 0 {
		return ErrNoAccount
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	
	return nil
}

func (m *AccountModel) CheckIfUserExists(userID int64) (bool, error) {
	query := `
		SELECT user_id, balance
		FROM accounts
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	var account Account

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&account.UserID, &account.Balance)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return false, ErrNoAccount
		default:
			return false, err
		}
	}

	return true, nil
}