package profile

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForUpdateClient = "UPDATE client SET name = $1, email = $2, img_path = $3, user_first_name =$4, user_second_name = $5, company = $6, phone_number = $7 WHERE id = $8"
    sqlTextForShowClient = "SELECT name, email, img_path, user_first_name, user_second_name, company, phone_number FROM client WHERE id = $1"
	sqlTextForDeleteClient = "DELETE FROM client WHERE id = $1"
)

type DB struct{
	sql *sql.DB
}

func New(sql *sql.DB) *DB{
	return &DB{
		sql: sql,
	}
}

func (r *DB) Update(ctx context.Context, clientID modeluser.ID) error {
	var UserName, Email ,ImgPath, UserFirstName, UserLastName, Company, Phone sql.NullString
	var args []interface{}

	if UserName.Valid && UserName.String != "" {
		args = append(args, UserName.String)
	}

	if Email.Valid && Email.String != "" {
		args = append(args, Email.String)
	}

	if ImgPath.Valid && ImgPath.String != "" {
		args = append(args, ImgPath.String)
	} 

	if UserFirstName.Valid && UserFirstName.String != "" {
		args = append(args, UserFirstName.String)
	}

	if UserLastName.Valid && UserLastName.String != "" {
		args = append(args, UserLastName.String)
	} 

	if Company.Valid && Company.String != ""{
		args = append(args, Company.String)
	} 

	if Phone.Valid && Phone.String != ""{ 
		args = append(args, Phone.String)
	} 

	args = append(args, clientID)	
	fmt.Println("Клиент args ", args)

	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateClient, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("ad client id %v not found", clientID)
	}

	return nil
}

func (r *DB) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error){
	var client modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForShowClient, clientID).Scan(&client.UserName, &client.Email, &client.ImagePath, &client.UserFirstName, &client.UserLastName, &client.Company, &client.Phone)
	if err != nil {
		return modeluser.User{}, fmt.Errorf("failed to update user: %w", err)
	}

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
