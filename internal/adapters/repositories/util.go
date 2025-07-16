package repositories

import (
	"github.com/google/uuid"
	"time"
	"github.com/jackc/pgx/v5/pgtype"
	"encoding/json"
)


// Helper functions for type conversion
func toPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func toPgUUIDPtr(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}

func getUUIDFromPgUUID(u pgtype.UUID) uuid.UUID {
	if !u.Valid {
		return uuid.Nil
	}
	return u.Bytes
}

func toPgJSONB(data map[string]interface{}) []byte {
	if data == nil || len(data) == 0 {
		return nil
	}
	
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// If marshaling fails, return empty JSON object
		return []byte("{}")
	}
	
	return jsonBytes
}

// getJSONBMap converts pgtype.JSONB to map[string]interface{}
func getJSONBMap(j []byte) map[string]interface{} {
	if j == nil || len(j) == 0 {
		return make(map[string]interface{})
	}
	
	var result map[string]interface{}
	err := json.Unmarshal(j, &result)
	if err != nil {
		// If unmarshaling fails, return empty map
		return make(map[string]interface{})
	}
	
	return result
}


// Helper functions for data conversion
func getTextStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func getDatePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func getTimestampPtr(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func toPgDatePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func toPgTimestampPtr(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}


func getTimestamptzPtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func getTimestamptzTime(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}


// getTimestamptz converts pgtype.Timestamptz to *time.Time (for nullable fields)
func getTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// getBool converts pgtype.Bool to bool
func getBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}


// toPgBool converts bool to pgtype.Bool
func toPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

// Helper function for timestamp pointer conversion
func toPgTimestamptzPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}