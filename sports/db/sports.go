package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sports/proto/sports"
	"strings"
	"sync"
	_ "syreclabs.com/go/faker"
	"time"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// Get will return a list of sport base on request IDs
	Get(id *sports.GetSportsRequest) ([]*sports.Sport, error)

	// Write new sport to sports.db
	Write(sportDetails *sports.AddSportRequest) (string, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sport repository.
func NewSportsRepo(db *sql.DB) *sportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy records.
		err = s.seed()
	})

	return err
}

// Get sport base on input id
func (s *sportsRepo) Get(id *sports.GetSportsRequest) ([]*sports.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQueries()[sportQuery]

	query, args, err = s.getFilter(query, id)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	sportList, err := s.scanSports(rows)
	if err != nil {
		return nil, err
	}
	if sportList == nil {
		return nil, nil
	}
	return sportList, err
}

// getFilter from input request
func (s *sportsRepo) getFilter(query string, filter *sports.GetSportsRequest) (string, []interface{}, error) {
	var (
		clauses string
		args    []interface{}
	)

	if filter == nil {
		return query, args, nil
	}

	// Filter by an array of sport
	if len(filter.Id) > 0 {
		clauses = " WHERE id IN (" + strings.Repeat("?,", len(filter.Id)-1) + "?)"
		for _, sportId := range filter.Id {
			args = append(args, sportId)
		}
	}

	if len(clauses) != 0 {
		query += clauses
	}

	log.Printf("The query is %v, %v", query, args)

	return query, args, nil

}

// scanSports from db table
func (s *sportsRepo) scanSports(rows *sql.Rows) ([]*sports.Sport, error) {
	var sportList []*sports.Sport

	for rows.Next() {
		var sport sports.Sport
		var advertisedStartTime time.Time

		if err := rows.Scan(&sport.Id, &sport.Name, &advertisedStartTime); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		sport.AdvertisedStartTime = timestamppb.New(advertisedStartTime)

		sportList = append(sportList, &sport)

	}

	return sportList, nil
}

// Write sport to db table
func (s *sportsRepo) Write(sportDetails *sports.AddSportRequest) (string, error) {
	var (
		txt   string
		query string
		args  []interface{}
	)

	if sportDetails.Sport.Id != 0 {

		query = addSportQueries()[addSport]
		query += "( ?,?,current_timestamp)"
		args = append(args, sportDetails.Sport.Id)
		args = append(args, sportDetails.Sport.Name)
		_, err := s.db.Exec(query, args...)
		if err != nil {
			return txt, err
		} else {
			log.Printf("The query is %v, %v", query, args)
			txt = "Add sport Success."
		}

	}

	return txt, nil
}
