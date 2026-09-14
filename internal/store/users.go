package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type UserStore struct{ DB *sql.DB }

func (s UserStore) List(ctx context.Context) ([]User, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, name, email, address FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []User{}
	for rows.Next() {
		var v User
		if err := rows.Scan(&v.ID, &v.Name, &v.Email, &v.Address); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (s UserStore) Get(ctx context.Context, id int64) (User, error) {
	var v User
	err := s.DB.QueryRowContext(ctx, "SELECT id, name, email, address FROM users WHERE id = ?", id).Scan(&v.ID, &v.Name, &v.Email, &v.Address)
	return v, err
}

func (s UserStore) Create(ctx context.Context, v User) (User, error) {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO users (name, email, address) VALUES (?, ?, ?)", v.Name, v.Email, v.Address)
	if err != nil {
		return User{}, err
	}
	v.ID, err = result.LastInsertId()
	return v, err
}

func (s UserStore) Update(ctx context.Context, v User) (User, error) {
	result, err := s.DB.ExecContext(ctx, "UPDATE users SET name = ?, email = ?, address = ? WHERE id = ?", v.Name, v.Email, v.Address, v.ID)
	if err != nil {
		return User{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return User{}, err
	}
	if count == 0 {
		return User{}, sql.ErrNoRows
	}
	return v, nil
}

func (s UserStore) Delete(ctx context.Context, id int64) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
