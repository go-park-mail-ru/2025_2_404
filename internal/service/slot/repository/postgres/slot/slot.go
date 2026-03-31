// Package postgres provides a PostgreSQL implementation of the slot repository.
package postgres

import (
	"2025_2_404/internal/service/slot/domain/slot"
	"2025_2_404/pkg/globalerrors"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
)

const (
	sqlTextForSelectSlots = `
		SELECT id, user_id, slot_name, min_cost_adv, format_of_banner, status, back_color, text_color 
		FROM slots WHERE user_id = $1
	`

	sqlTextForInsertSlot = `
		INSERT INTO slots (user_id, slot_name, min_cost_adv, format_of_banner, status, back_color, text_color) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	sqlTextForSelectSlotByID = `
		SELECT id, user_id, slot_name, min_cost_adv, format_of_banner, status, back_color, text_color 
		FROM slots WHERE id = $1
	`

	sqlTextForUpdateSlot = `
		UPDATE slots SET 
			slot_name = $1, 
			min_cost_adv = $2, 
			format_of_banner = $3, 
			status = $4, 
			back_color = $5, 
			text_color = $6 
		WHERE id = $7 AND user_id = $8
	`

	sqlTextForDeleteSlot = `
		DELETE FROM slots WHERE id = $1 AND user_id = $2
	`
)

type DB struct {
	sql *sql.DB
}

func New(sql *sql.DB) *DB {
	return &DB{sql: sql}
}

func (r *DB) Create(ctx context.Context, s slot.Slot) (slot.ID, error) {
	userID, err := uuid.Parse(string(s.UserID))
	if err != nil {
		return "", globalerrors.ErrInvalidQuery
	}

	var newID uuid.UUID
	err = r.sql.QueryRowContext(
		ctx,
		sqlTextForInsertSlot,
		userID,
		s.SlotName,
		s.MinCostAdv,
		s.FormatOfBanner,
		s.Status,
		s.BackColor,
		s.TextColor,
	).Scan(&newID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503": // foreign_key_violation (несуществующий user_id)
				return "", globalerrors.ErrInvalidQuery
			case "23514": // check_violation (неверный format_of_banner, цвет и т.д.)
				return "", globalerrors.ErrInvalidQuery
			case "23505": // unique_violation (маловероятно, но возможно)
				return "", globalerrors.ErrInternal
			}
		}
		return "", globalerrors.ErrInternal
	}

	return slot.ID(newID.String()), nil
}

func (r *DB) GetByID(ctx context.Context, id slot.ID) (slot.Slot, error) {
	idUUID, err := uuid.Parse(string(id))
	if err != nil {
		return slot.Slot{}, globalerrors.ErrInvalidQuery
	}

	var s slot.Slot
	var userUUID uuid.UUID
	row := r.sql.QueryRowContext(ctx, sqlTextForSelectSlotByID, idUUID)
	err = row.Scan(
		&idUUID,
		&userUUID,
		&s.SlotName,
		&s.MinCostAdv,
		&s.FormatOfBanner,
		&s.Status,
		&s.BackColor,
		&s.TextColor,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slot.Slot{}, globalerrors.ErrSlotNotFound
		}
		return slot.Slot{}, globalerrors.ErrInternal
	}

	s.ID = slot.ID(idUUID.String())
	s.UserID = slot.UserID(userUUID.String())
	return s, nil
}

func (r *DB) ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error) {
	userUUID, err := uuid.Parse(string(userID))
	if err != nil {
		return nil, globalerrors.ErrInvalidQuery
	}

	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectSlots, userUUID)
	if err != nil {
		return nil, globalerrors.ErrInternal
	}
	defer func() {
		_ = rows.Close()
	}()

	var slots []slot.Slot
	for rows.Next() {
		var s slot.Slot
		var idUUID, userUUID uuid.UUID
		if err := rows.Scan(
			&idUUID,
			&userUUID,
			&s.SlotName,
			&s.MinCostAdv,
			&s.FormatOfBanner,
			&s.Status,
			&s.BackColor,
			&s.TextColor,
		); err != nil {
			return nil, globalerrors.ErrInternal
		}
		s.ID = slot.ID(idUUID.String())
		s.UserID = slot.UserID(userUUID.String())
		slots = append(slots, s)
	}

	if err = rows.Err(); err != nil {
		return nil, globalerrors.ErrInternal
	}

	return slots, nil
}

func (r *DB) Update(ctx context.Context, s slot.Slot) error {
	id, err := uuid.Parse(string(s.ID))
	if err != nil {
		return globalerrors.ErrInvalidQuery
	}

	userID, err := uuid.Parse(string(s.UserID))
	if err != nil {
		return globalerrors.ErrInvalidQuery
	}

	res, err := r.sql.ExecContext(
		ctx,
		sqlTextForUpdateSlot,
		s.SlotName,
		s.MinCostAdv,
		s.FormatOfBanner,
		s.Status,
		s.BackColor,
		s.TextColor,
		id,
		userID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514": // check_violation
				return globalerrors.ErrInvalidQuery
			}
		}
		return globalerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return globalerrors.ErrInternal
	}
	if rowsAffected == 0 {
		return globalerrors.ErrAccessDenied // слот не найден ИЛИ нет прав
	}

	return nil
}

func (r *DB) Delete(ctx context.Context, id slot.ID, userID slot.UserID) error {
	idUUID, err := uuid.Parse(string(id))
	if err != nil {
		return globalerrors.ErrInvalidQuery
	}

	userIDUUID, err := uuid.Parse(string(userID))
	if err != nil {
		return globalerrors.ErrInvalidQuery
	}

	res, err := r.sql.ExecContext(ctx, sqlTextForDeleteSlot, idUUID, userIDUUID)
	if err != nil {
		return globalerrors.ErrInternal
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return globalerrors.ErrInternal
	}
	if rowsAffected == 0 {
		return globalerrors.ErrAccessDenied
	}

	return nil
}
