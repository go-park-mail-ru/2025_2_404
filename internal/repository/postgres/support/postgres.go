package support

import (
	modelsup "2025_2_404/internal/domain/models/support"
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"

)

const(
	sqlTextForSelectAllSupport = "SELECT id, client_id, sup_status, category, sup_description, img_path, contact_name, contact_email FROM support WHERE client_id = $1"
    sqlTextForInsertSupport = "INSERT INTO support (client_id, sup_status, category, sup_description, img_path, contact_name, contact_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
    sqlTextForSelectSupport = "SELECT id, client_id, sup_status, category, sup_description, img_path, contact_name, contact_email FROM support WHERE id = $1 AND client_id = $2"
	sqlTextForUpdateSupport = "UPDATE support SET sup_status = $1, category = $2, sup_description = $3, img_path = $4, contact_name = $5, contact_email = $6 WHERE id = $7 AND client_id = $8"
	sqlTextForAllSupport = "SELECT id, sup_status, category, sup_description, img_path, contact_name, contact_email FROM support WHERE client_id = $1"
	sqlTextForDeleteSupport = "DELETE FROM support WHERE id = $1"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Create(ctx context.Context, sup modelsup.Support) (error) {
	res, err := r.sql.ExecContext(ctx, sqlTextForInsertSupport, sup.UserID, sup.Status, sup.Category, sup.Description, sup.ImagePath, sup.ContactName, sup.ContactEmail)
	if err != nil {
		return fmt.Errorf("failed to create ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad with id %v not found", sup.ID)
	}

	return nil
}

func (r *DB) GetOneSup(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) (modelsup.Support, error) {
	var supInfo modelsup.Support
	row := r.sql.QueryRowContext(ctx, sqlTextForSelectSupport, supID, clientID)

	err :=  row.Scan(
		&supInfo.ID,
		&supInfo.UserID,
		&supInfo.Status,
		&supInfo.Category,
		&supInfo.Description,
		&supInfo.ImagePath,
		&supInfo.ContactName,
		&supInfo.ContactEmail,
	)

	if err != nil {
		return modelsup.Support{}, fmt.Errorf("failed to find an ad: %w", err)
	}

	return supInfo, nil
}

func (r *DB) GetAllSups(ctx context.Context, superID int64) ([]modelsup.Support, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForAllSupport, superID)
	if err != nil {
		return nil, fmt.Errorf("failed to query supports by super user ID: %w", err)
	}
	defer rows.Close()

	var support []modelsup.Support
	for rows.Next() {
		var sup modelsup.Support
		err := rows.Scan(&sup.ID, &sup.Status, &sup.Category, &sup.Description, &sup.ImagePath, &sup.ContactName, &sup.ContactEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ad: %w", err)
		}
		support = append(support, sup)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return support, nil
}

func (r *DB) Update(ctx context.Context, sup modelsup.Support) error {
	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateSupport, sup.Status, sup.Category, sup.Description, sup.ImagePath, sup.ContactName, sup.ContactEmail, sup.ID, sup.UserID)
	if err != nil {
		return fmt.Errorf("failed to update ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sup with id %v and client_id %v not found", sup.ID, sup.UserID)
	}

	return nil
}

func (r *DB) Delete(ctx context.Context, supID modelsup.ID, clientID modeluser.ID) error {
	
	result, err := r.sql.ExecContext(ctx, sqlTextForDeleteSupport, supID)
	if err != nil {
		return fmt.Errorf("failed to delete ad: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get a rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Sup with ID %d not found", supID)
	}
	fmt.Printf("Пользователь с ID %d успешно удален. Затронуто строк: %d", supID, rowsAffected)
	return nil
}

func (r *DB) FindByUserID(ctx context.Context, clientID modeluser.ID) ([]modelsup.Support, error) {
	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectAllSupport, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query ads by user ID: %w", err)
	}
	defer rows.Close()

	var sups []modelsup.Support
	for rows.Next() {
		var supInfo modelsup.Support
		err :=  rows.Scan(
		&supInfo.ID,
		&supInfo.UserID,
		&supInfo.Status,
		&supInfo.Category,
		&supInfo.Description,
		&supInfo.ImagePath,
		&supInfo.ContactName,
		&supInfo.ContactEmail,
	)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ad: %w", err)
		}
		sups = append(sups, supInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return sups, nil
}
