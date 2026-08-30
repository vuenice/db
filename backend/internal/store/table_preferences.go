package store

import (
	"context"
	"database/sql"
	"errors"
)

// GetTablePreferences returns the preferences JSON string for a user, connection, schema, and table.
func (s *Store) GetTablePreferences(ctx context.Context, userID, connID int, schemaName, tableName string) (string, error) {
	var pref string
	err := s.DB.QueryRowContext(ctx, `
		SELECT preferences_json FROM table_preferences 
		WHERE user_id = ? AND connection_id = ? AND schema_name = ? AND table_name = ?`,
		userID, connID, schemaName, tableName).Scan(&pref)
	
	if errors.Is(err, sql.ErrNoRows) {
		return "{}", nil
	}
	if err != nil {
		return "", err
	}
	return pref, nil
}

// SaveTablePreferences saves the preferences JSON string for a user, connection, schema, and table.
func (s *Store) SaveTablePreferences(ctx context.Context, userID, connID int, schemaName, tableName, preferencesJSON string) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO table_preferences (user_id, connection_id, schema_name, table_name, preferences_json, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, connection_id, schema_name, table_name) DO UPDATE SET
			preferences_json = excluded.preferences_json,
			updated_at = CURRENT_TIMESTAMP
	`, userID, connID, schemaName, tableName, preferencesJSON)
	return err
}
