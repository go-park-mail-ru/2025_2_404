package profile

import (
	modeluser "2025_2_404/internal/service/profile/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

const (
	sqlTextForShowClient   = "SELECT name, email, img_path, user_first_name, user_second_name, company, phone_number, created_at FROM client WHERE id = $1"
	sqlTextForDeleteClient = "DELETE FROM client WHERE id = $1"
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{
		sql: sql,
	}
}

func (r *DB) Update(ctx context.Context, client modeluser.User) error {
	var updates []string
	var args []interface{}
	argID := 1

	if client.UserName != "" {
		updates = append(updates, fmt.Sprintf("name = $%d", argID))
		args = append(args, client.UserName)
		argID++
	}

	if client.Email != "" {
		updates = append(updates, fmt.Sprintf("email = $%d", argID))
		args = append(args, client.Email)
		argID++
	}

	if client.ImagePath != "" {
		updates = append(updates, fmt.Sprintf("img_path = $%d", argID))
		args = append(args, client.ImagePath)
		argID++
	}

	if client.UserFirstName != "" {
		updates = append(updates, fmt.Sprintf("user_first_name = $%d", argID))
		args = append(args, client.UserFirstName)
		argID++
	}

	if client.UserLastName != "" {
		updates = append(updates, fmt.Sprintf("user_second_name = $%d", argID))
		args = append(args, client.UserLastName)
		argID++
	}

	if client.Company != "" {
		updates = append(updates, fmt.Sprintf("company = $%d", argID))
		args = append(args, client.Company)
		argID++
	}

	if client.Phone != "" {
		updates = append(updates, fmt.Sprintf("phone_number = $%d", argID))
		args = append(args, client.Phone)
		argID++
	}

	fmt.Println("Клиент updates ", updates)

	if len(updates) == 0 {
		return nil
	}

	updatesQuery := strings.Join(updates, ", ")
	sqlTextForUpdateClient := fmt.Sprintf("UPDATE client SET %s WHERE id = $%d", updatesQuery, argID)

	args = append(args, client.ID)

	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateClient, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad client id %v not found", client.ID)
	}

	return nil
}

func (r *DB) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error) {
	var ImagePath, UserFirstName, UserLastName, Company, Phone sql.NullString
	var CreatedAt sql.NullTime
	var client modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForShowClient, clientID).Scan(&client.UserName, &client.Email, &ImagePath, &UserFirstName, &UserLastName, &Company, &Phone, &CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return modeluser.User{}, fmt.Errorf("user with id %d not found", clientID)
		}
		return modeluser.User{}, fmt.Errorf("failed to scan user: %w", err)
	}

	client.ImagePath = ImagePath.String
	client.UserFirstName = UserFirstName.String
	client.UserLastName = UserLastName.String
	client.Company = Company.String
	client.Phone = Phone.String
	client.CreatedAt = CreatedAt.Time.String()

	fmt.Println("Клиент ", client)
	return client, nil
}

func (r *DB) Delete(ctx context.Context, clientID modeluser.ID) error {
	result, err := r.sql.ExecContext(ctx, sqlTextForDeleteClient, clientID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get a rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Profile with ID %d not found", clientID)
	}
	fmt.Printf("Пользователь с ID %d успешно удален. Затронуто строк: %d", clientID, rowsAffected)
	return nil
}
