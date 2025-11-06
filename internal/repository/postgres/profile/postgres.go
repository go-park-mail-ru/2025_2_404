package profile

import (
	modeluser "2025_2_404/internal/domain/models/user"
	"context"
	"database/sql"
	"fmt"
)

const(
	sqlTextForUpdateClient = "UPDATE client SET user_login = $1, email = $2, img_path = $3, user_name =$4, user_subname = $5, company = $6, phone_number = $7 WHERE id = $8"
    sqlTextForShowClient = "SELECT user_login, email, img_path, user_name, user_subname, company, phone_number FROM client WHERE id = $1"
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

func (r *DB) Update(ctx context.Context, client modeluser.User) error {
	res, err := r.sql.ExecContext(ctx, sqlTextForUpdateClient, client.UserName, client.Email, client.ImagePath, client.UserFirstName, client.UserLastName, client.Company, client.Phone, client.ID)
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

func (r *DB) Show(ctx context.Context, clientID modeluser.ID) (modeluser.User, error){
	var imgPath, UserFirstName, UserLastName, Company,Phone sql.NullString
	var client modeluser.User
	err := r.sql.QueryRowContext(ctx, sqlTextForShowClient, clientID).Scan(&client.UserName, &client.Email, &imgPath, &UserFirstName, &UserLastName, &Company, &Phone)
	if err != nil {
		return modeluser.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	if imgPath.Valid {
		client.ImagePath = imgPath.String
	} else {
		client.ImagePath = ""
	}

	if UserFirstName.Valid {
		client.UserFirstName = UserFirstName.String
	} else {
		client.UserFirstName = ""
	}

	if UserLastName.Valid {
		client.UserLastName = UserLastName.String
	} else {
		client.UserLastName = ""
	}

	if Company.Valid {
		client.Company = Company.String
	} else {
		client.Company = ""
	}

	if Phone.Valid {
		client.Phone = Phone.String
	} else {
		client.Phone = ""
	}

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
