package postgres

import (
	"database/sql"
	"context"
	"fmt"
	"github.com/google/uuid"
	"2025_2_404/internal/service/slot/domain/slot"
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
		return "",fmt.Errorf("invalid user ID: %w", err)
	}

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
	).Scan(&s.ID)
	
	if err != nil {
		return "", fmt.Errorf("failed to insert slot: %w", err)
	}

	return s.ID, nil
}

func (r *DB) GetByID(ctx context.Context, id slot.ID) (slot.Slot, error) {
	var s slot.Slot
	idUUID, err := uuid.Parse(string(id))
	if err != nil {
		return slot.Slot{}, fmt.Errorf("invalid slot id: %w", err)
	}
	row := r.sql.QueryRowContext(ctx, sqlTextForSelectSlotByID, idUUID)
	var idUuid, userUuid uuid.UUID 
	err = row.Scan(
		&idUuid,
		&userUuid,
		// &s.ID,
		// &s.UserID,q
		&s.SlotName,
		&s.MinCostAdv,
		&s.FormatOfBanner,
		&s.Status,
		&s.BackColor,
		&s.TextColor,
	)

	if err != nil {
		return slot.Slot{}, fmt.Errorf("failed to get slot by ID: %w", err)
	}

	s.ID = slot.ID(idUuid.String())
	s.UserID = slot.UserID(userUuid.String())

	return s, nil
}

func (r *DB) ListByUserID(ctx context.Context, userID slot.UserID) ([]slot.Slot, error) {
	userUUID, err := uuid.Parse(string(userID))
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	rows, err := r.sql.QueryContext(ctx, sqlTextForSelectSlots, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query slots: %w", err)
	}
	defer rows.Close()

	var slots []slot.Slot
	for rows.Next() {
		var s slot.Slot
		var idUuid, userUuid uuid.UUID
		err := rows.Scan(
			// &s.ID,
			&idUuid,    
			&userUuid,
			// &s.UserID,
			&s.SlotName,
			&s.MinCostAdv,
			&s.FormatOfBanner,
			&s.Status,
			&s.BackColor,
			&s.TextColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan slot: %w", err)
		}
		s.ID = slot.ID(idUuid.String())
		s.UserID = slot.UserID(userUuid.String())
		slots = append(slots, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return slots, nil
}

func (r *DB) Update(ctx context.Context, s slot.Slot) error {
	id, err := uuid.Parse(string(s.ID))
	if err != nil {
		return fmt.Errorf("invalid slot ID: %w", err)
	}

	userID, err := uuid.Parse(string(s.UserID))
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
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
		return fmt.Errorf("failed to update slot: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("slot with ID %s and user ID %s not found", s.ID, s.UserID)
	}

	return nil
}

func (r *DB) Delete(ctx context.Context, id slot.ID, userID slot.UserID) error {
	idUUID, err := uuid.Parse(string(id))
	if err != nil {
		return fmt.Errorf("invalid slot ID: %w", err)
	}

	userIDUUID, err := uuid.Parse(string(userID))
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	res, err := r.sql.ExecContext(ctx, sqlTextForDeleteSlot, idUUID, userIDUUID)
	if err != nil {
		return fmt.Errorf("failed to delete slot: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("slot with ID %s not found", id)
	}

	return nil
}