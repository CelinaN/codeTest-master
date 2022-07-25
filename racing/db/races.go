package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"strings"
	"sync"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)

	// Get will return a list of races base on request IDs
	Get(id *racing.GetRacesRequest) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses      []string
		args         []interface{}
		visibleValue string
	)

	if filter == nil {
		return query, args
	}
	// filter meetingId
	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	// populate pass in filter value
	if filter.Visible == "TRUE" {
		visibleValue = "1"
	}
	if filter.Visible == "FALSE" {
		visibleValue = "0"
	}

	//filter by visibility
	if visibleValue == "1" || visibleValue == "0" {
		if len(clauses) != 0 {
			query += " AND visible = " + visibleValue
		} else {
			query += " WHERE visible = " + visibleValue

		}

	}

	// add order by field
	orderByField := getOrderByfield(filter.OrderBy)
	query += " ORDER BY " + orderByField

	log.Printf("The query is %v, %v", query, args)

	return query, args
}

func (r *racesRepo) scanRaces(rows *sql.Rows) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		race.AdvertisedStartTime = timestamppb.New(advertisedStart)

		// update status field base on the advertised start time
		if advertisedStart.After(time.Now()) {
			race.Status = "OPEN"
		} else {
			race.Status = "CLOSED"
		}

		races = append(races, &race)

	}

	return races, nil
}

func getOrderByfield(inputField string) string {

	var orderByField string

	//switch pass in order by fields
	switch inputField {
	case "NAME":
		orderByField = "name"
	case "NUMBER":
		orderByField = "number"
	case "ID":
		orderByField = "id"
	case "MEETING_ID":
		orderByField = "meeting_id"
	default:
		orderByField = "advertised_start_time"
	}
	return orderByField
}

func (r *racesRepo) Get(id *racing.GetRacesRequest) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args, err = r.getFilter(query, id)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}
	if races == nil {
		return nil, nil
	}
	return races, err
}

func (r *racesRepo) getFilter(query string, filter *racing.GetRacesRequest) (string, []interface{}, error) {
	var (
		clauses string
		args    []interface{}
	)

	if filter == nil {
		return query, args, nil
	}

	// Filter by an array of sport ids
	if len(filter.Id) > 0 {
		clauses = " WHERE id IN (" + strings.Repeat("?,", len(filter.Id)-1) + "?)"
		for _, raceId := range filter.Id {
			args = append(args, raceId)
		}
	}

	if len(clauses) != 0 {
		query += clauses
	}

	log.Printf("The query is %v, %v", query, args)

	return query, args, nil

}
