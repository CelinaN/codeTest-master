package db

const (
	sportQuery = "list"
	addSport   = "list"
)

func getSportQueries() map[string]string {
	return map[string]string{
		sportQuery: `SELECT * FROM sports`,
	}
}

func addSportQueries() map[string]string {
	return map[string]string{
		addSport: `INSERT INTO sports VALUES`,
	}
}
