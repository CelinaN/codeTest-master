package autoTest_test

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var racesRepoTest db.RacesRepo

func TestMain(m *testing.M) {
	racingDB, err := sql.Open("sqlite3", "./db/racing.db")
	if err != nil {
		log.Fatal("DB table open failed", err)
	}

	racesRepoTest = db.NewRacesRepo(racingDB)
	if err := racesRepoTest.Init(); err != nil {
		log.Fatal("Failed to initialize Data in DB")
	}
	runTest := m.Run()
	after(racingDB)
	os.Exit(runTest)
}

func TestListVisible(t *testing.T) {
	filter := new(racing.ListRacesRequestFilter)

	filter.Visible = "TRUE"

	races, _ := racesRepoTest.List(filter)

	for _, race := range races {
		assert.Equal(t, race.Visible, true, "Race visible not match.")
	}
}

func TestListInvisible(t *testing.T) {
	filter := new(racing.ListRacesRequestFilter)

	filter.Visible = "FALSE"

	races, _ := racesRepoTest.List(filter)

	for _, race := range races {
		assert.Equal(t, race.Visible, "0", "Race visible not match.")
	}
}

func TestRaceStatus(t *testing.T) {
	filter := new(racing.ListRacesRequestFilter)

	races, _ := racesRepoTest.List(filter)

	for _, race := range races {

		if race.AdvertisedStartTime.AsTime().After(time.Now()) {
			assert.Equal(t, race.Status, "OPEN", "Status not match")
		} else {
			assert.Equal(t, race.Status, "CLOSED", "Status not match")
		}
	}
}

func TestRaceGetID(t *testing.T) {
	id := new(racing.GetRacesRequest)
	id.Id = []int64{1, 2}
	races, _ := racesRepoTest.Get(id)

	for _, race := range races {
		assert.Contains(t, race.Id, []int64{1, 2}, "ID not match")

	}
}

func after(db *sql.DB) {
	if _, err := db.Exec("drop table races"); err != nil {
		log.Fatal("Failed to Clean Up DB")
	}
}
