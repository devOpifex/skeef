package db

import "github.com/devOpifex/skeef-app/stream"

// InsertStream Insert a new stream
func (DB *Database) InsertStream(name, follow, track, locations string) error {

	stmt, err := DB.Con.Prepare("INSERT INTO streams (name, follow, track, locations, active) VALUES (?,?,?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, follow, track, locations, 0)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) StreamExists(name string) (bool, error) {

	stmt, err := DB.Con.Prepare("SELECT COUNT(1) FROM streams WHERE name = ?;")

	if err != nil {
		return false, err
	}

	row := stmt.QueryRow(name)

	res := 0
	row.Scan(&res)

	return res > 0, nil
}

func (DB *Database) GetStreams() ([]stream.Stream, error) {
	var streams []stream.Stream

	rows, err := DB.Con.Query("SELECT name, follow, track, locations, active FROM streams;")

	if err != nil {
		return streams, err
	}

	for rows.Next() {
		var stream stream.Stream

		err := rows.Scan(&stream.Name, &stream.Follow, &stream.Track, &stream.Locations, &stream.Active)

		if err != nil {
			continue
		}

		streams = append(streams, stream)
	}

	return streams, nil
}
