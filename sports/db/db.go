package db

import (
	"syreclabs.com/go/faker"
	"time"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INTEGER PRIMARY KEY, name TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO sports(id, name, advertised_start_time) VALUES (?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name(),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}
