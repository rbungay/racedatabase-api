package storage

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
)

// SupabaseStorage handles database operations with Supabase
type SupabaseStorage struct {
    db *sql.DB
}

// NewSupabaseStorage creates a new instance of SupabaseStorage
func NewSupabaseStorage(dbURL string) (*SupabaseStorage, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, err
    }
    
    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return &SupabaseStorage{db: db}, nil
}

// cleanFeeString converts "$20.00" to 20.00
func cleanFeeString(fee string) float64 {
    // Remove "$" and any whitespace
    fee = strings.TrimSpace(strings.TrimPrefix(fee, "$"))
    // Convert to float64
    val, err := strconv.ParseFloat(fee, 64)
    if err != nil {
        return 0.0
    }
    return val
}

// nullTime returns a sql.NullTime for a given time string
func nullTime(t string) sql.NullTime {
    if t == "" {
        return sql.NullTime{Valid: false}
    }
    parsed, err := time.Parse(time.RFC3339, t)
    if err != nil {
        return sql.NullTime{Valid: false}
    }
    return sql.NullTime{Time: parsed, Valid: true}
}

// SaveRace stores a race and its associated events in the database
func (s *SupabaseStorage) SaveRace(race *models.RaceDetails) error {
    tx, err := s.db.BeginTx(context.Background(), nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert race
    _, err = tx.Exec(`
        INSERT INTO races (id, name, url, external_url, logo_url, timezone)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            url = EXCLUDED.url,
            external_url = EXCLUDED.external_url,
            logo_url = EXCLUDED.logo_url,
            timezone = EXCLUDED.timezone,
            updated_at = NOW()
    `, race.ID, race.Name, race.URL, race.ExternalURL, race.LogoURL, race.Timezone)
    if err != nil {
        return err
    }

    // Insert events
    for _, event := range race.Events {
        _, err = tx.Exec(`
            INSERT INTO events (event_id, race_id, name, start_time, end_time, event_type, distance, registration_opens, category)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            ON CONFLICT (event_id) DO UPDATE SET
                name = EXCLUDED.name,
                start_time = EXCLUDED.start_time,
                end_time = EXCLUDED.end_time,
                event_type = EXCLUDED.event_type,
                distance = EXCLUDED.distance,
                registration_opens = EXCLUDED.registration_opens,
                category = EXCLUDED.category,
                updated_at = NOW()
        `, event.EventID, race.ID, event.Name, 
           nullTime(event.StartTime), nullTime(event.EndTime),
           event.EventType, event.Distance, 
           nullTime(event.RegOpens), event.Category)
        if err != nil {
            return err
        }

        // Insert registration periods with cleaned fee values
        for _, period := range event.RegPeriods {
            _, err = tx.Exec(`
                INSERT INTO registration_periods (event_id, opens_at, closes_at, race_fee, processing_fee)
                VALUES ($1, $2, $3, $4, $5)
            `, event.EventID, 
               nullTime(period.Opens), nullTime(period.Closes),
               cleanFeeString(period.Fee), cleanFeeString(period.ProcFee))
            if err != nil {
                return err
            }
        }
    }

    return tx.Commit()
}