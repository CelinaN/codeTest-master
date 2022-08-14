package autoTest_test

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sports/db"
	"sports/proto/sports"
	"testing"
)

var sportsRepoTest db.SportsRepo

func TestMain(m *testing.M) {
	sportDB, err := sql.Open("sqlite3", ".memory:")
	if err != nil {
		log.Fatal("DB table open failed", err)
	}

	sportsRepoTest = db.NewSportsRepo(sportDB)
	sportsRepoTest.Init()
	if err := sportsRepoTest.Init(); err != nil {
		log.Fatal("Failed to initialize Data in DB")
	}
	runTest := m.Run()

	os.Exit(runTest)
}

func TestAddSport(t *testing.T) {

	//write a new sport
	sport := new(sports.Sport)
	sport.Id = 105
	sport.Name = "Test record"
	newSport := new(sports.AddSportRequest)
	newSport.Sport = sport

	_, _ = sportsRepoTest.Write(newSport)

	//get sport if match with the new add
	getSport := new(sports.GetSportsRequest)
	getSport.Id = []int64{105}
	sportList, _ := sportsRepoTest.Get(getSport)

	for _, sport := range sportList {
		assert.Equal(t, sport.Id, []int64{105}, "sport id not match.")
		assert.Equal(t, sport.Name, "Test record", "sport name not match.")
	}
}
