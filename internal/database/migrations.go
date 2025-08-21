package database

import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			email_verified BOOLEAN DEFAULT FALSE,
			email_verify_token VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS friendships (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			accepted_at TIMESTAMP,
			UNIQUE(user_id, friend_id)
		)`,
		`CREATE TABLE IF NOT EXISTS races (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			distance INTEGER NOT NULL,
			status VARCHAR(20) DEFAULT 'waiting',
			created_by INTEGER REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP,
			finished_at TIMESTAMP,
			countdown_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS race_participants (
			id SERIAL PRIMARY KEY,
			race_id INTEGER REFERENCES races(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			status VARCHAR(20) DEFAULT 'not_ready',
			current_distance INTEGER DEFAULT 0,
			finished_at TIMESTAMP,
			pace VARCHAR(10),
			position INTEGER,
			joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(race_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS race_updates (
			id SERIAL PRIMARY KEY,
			race_id INTEGER REFERENCES races(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			distance INTEGER NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_user_id ON friendships(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_friendships_friend_id ON friendships(friend_id)`,
		`CREATE INDEX IF NOT EXISTS idx_race_participants_race_id ON race_participants(race_id)`,
		`CREATE INDEX IF NOT EXISTS idx_race_participants_user_id ON race_participants(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_race_updates_race_id ON race_updates(race_id)`,
		`CREATE INDEX IF NOT EXISTS idx_race_updates_user_id ON race_updates(user_id)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			refresh_token_hash VARCHAR(255) NOT NULL,
			device_type VARCHAR(50) NOT NULL,
			user_agent TEXT,
			ip_address VARCHAR(45),
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token_hash ON sessions(refresh_token_hash)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}