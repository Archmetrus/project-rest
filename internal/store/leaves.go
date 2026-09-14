package store

import (
	"context"
	"database/sql"
)

type Leave struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type LeaveStore struct{ DB *sql.DB }

func (s LeaveStore) List(ctx context.Context) ([]Leave, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, user_id, start_date, end_date FROM leaves ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Leave{}
	for rows.Next() {
		var v Leave
		if err := rows.Scan(&v.ID, &v.UserID, &v.StartDate, &v.EndDate); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (s LeaveStore) Get(ctx context.Context, id int64) (Leave, error) {
	var v Leave
	err := s.DB.QueryRowContext(ctx, "SELECT id, user_id, start_date, end_date FROM leaves WHERE id = ?", id).Scan(&v.ID, &v.UserID, &v.StartDate, &v.EndDate)
	return v, err
}

func (s LeaveStore) Create(ctx context.Context, v Leave) (Leave, error) {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO leaves (user_id, start_date, end_date) VALUES (?, ?, ?)", v.UserID, v.StartDate, v.EndDate)
	if err != nil {
		return Leave{}, err
	}
	v.ID, err = result.LastInsertId()
	return v, err
}

func (s LeaveStore) Update(ctx context.Context, v Leave) (Leave, error) {
	result, err := s.DB.ExecContext(ctx, "UPDATE leaves SET user_id = ?, start_date = ?, end_date = ? WHERE id = ?", v.UserID, v.StartDate, v.EndDate, v.ID)
	if err != nil {
		return Leave{}, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Leave{}, err
	}
	if count == 0 {
		return Leave{}, sql.ErrNoRows
	}
	return v, nil
}

func (s LeaveStore) Delete(ctx context.Context, id int64) error {
	result, err := s.DB.ExecContext(ctx, "DELETE FROM leaves WHERE id = ?", id)
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
