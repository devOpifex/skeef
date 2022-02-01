package db

import (
	"github.com/devOpifex/skeef/stream"
)

// InsertStream Insert a new stream
func (DB *Database) InsertStream(
	name, follow, track, locations, exclude, maxEdges, description string,
	retweetsNet, mentionsNet, hashtagsNet int,
	filterLevel string,
	min_follower_count, min_favorite_count, only_verified,
	max_hashtags, max_mentions, replyNet int) error {

	stmt, err := DB.Con.Prepare(`INSERT INTO streams 
		(
			name, follow, track, locations, active, 
			exclude, max_edges, description, retweets_net, 
			mentions_net, hashtags_net, filter_level,
			min_follower_count, min_favorite_count, only_verified,
			max_hashtags, max_mentions, reply_net
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		name,
		follow,
		track,
		locations,
		0,
		exclude,
		maxEdges,
		description,
		retweetsNet,
		mentionsNet,
		hashtagsNet,
		filterLevel,
		min_follower_count,
		min_favorite_count,
		only_verified,
		max_hashtags,
		max_mentions,
		replyNet,
	)

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

	rows, err := DB.Con.Query(`SELECT name, follow, 
		track, locations, active, max_edges, description, 
		retweets_net, mentions_net, hashtags_net, filter_level,
		min_follower_count, min_favorite_count, only_verified,
		max_hashtags, max_mentions, reply_net
		FROM streams;`)

	if err != nil {
		return streams, err
	}

	for rows.Next() {
		var stream stream.Stream

		err := rows.Scan(
			&stream.Name,
			&stream.Follow,
			&stream.Track,
			&stream.Locations,
			&stream.Active,
			&stream.MaxEdges,
			&stream.Description,
			&stream.RetweetsNet,
			&stream.MentionsNet,
			&stream.HashtagsNet,
			&stream.FilterLevel,
			&stream.MinFollowerCount,
			&stream.MinFavoriteCount,
			&stream.OnlyVerified,
			&stream.MaxHashtags,
			&stream.MaxMentions,
			&stream.ReplyNet,
		)

		if err != nil {
			continue
		}

		streams = append(streams, stream)
	}

	return streams, nil
}

func (DB *Database) DeleteStream(name string) error {
	stmt, err := DB.Con.Prepare("DELETE FROM streams WHERE name = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(&name)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) StartStream(name string) error {

	stmt, err := DB.Con.Prepare("UPDATE streams SET active = 1 WHERE name = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(&name)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) PauseStream(name string) error {

	stmt, err := DB.Con.Prepare("UPDATE streams SET active = 0 WHERE name = ?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(&name)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) PauseAllStreams() error {

	stmt, err := DB.Con.Prepare("UPDATE streams SET active = 0")

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) GetStream(name string) (stream.Stream, error) {
	var stream stream.Stream

	stmt, err := DB.Con.Prepare(`SELECT name, follow, 
		track, locations, active, max_edges, exclude, 
		description, retweets_net, mentions_net, hashtags_net,
		filter_level, min_follower_count, min_favorite_count, only_verified,
		max_hashtags, max_mentions, reply_net
		FROM streams WHERE name = ?`)

	if err != nil {
		return stream, err
	}

	row := stmt.QueryRow(name)

	row.Scan(
		&stream.Name,
		&stream.Follow,
		&stream.Track,
		&stream.Locations,
		&stream.Active,
		&stream.MaxEdges,
		&stream.Exclusion,
		&stream.Description,
		&stream.RetweetsNet,
		&stream.MentionsNet,
		&stream.HashtagsNet,
		&stream.FilterLevel,
		&stream.MinFollowerCount,
		&stream.MinFavoriteCount,
		&stream.OnlyVerified,
		&stream.MaxHashtags,
		&stream.MaxMentions,
		&stream.ReplyNet,
	)

	return stream, nil
}

func (DB *Database) UpdateStream(
	track, follow, locations, newName, currentName, exclusion, description string,
	maxEdges, retweetsNet, mentionsNet, hashtagsNet int,
	filterLevel string,
	min_follower_count, min_favorite_count, only_verified,
	max_hashtags, max_mentions, replyNet int) error {

	stmt, err := DB.Con.Prepare(`UPDATE streams SET 
		track = ?, 
		follow = ?, 
		locations = ?, 
		name = ?, 
		max_edges = ?, 
		description = ?, 
		exclude = ?, 
		retweets_net = ?, 
		mentions_net = ?, 
		hashtags_net = ?,
		filter_level = ?, 
		min_follower_count = ?, 
		min_favorite_count = ?, 
		only_verified = ?,
		max_hashtags = ?, 
		max_mentions = ?
		WHERE name = ?`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		track,
		follow,
		locations,
		newName,
		maxEdges,
		description,
		exclusion,
		retweetsNet,
		mentionsNet,
		hashtagsNet,
		filterLevel,
		min_follower_count,
		min_favorite_count,
		only_verified,
		max_hashtags,
		max_mentions,
		replyNet,
		currentName,
	)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) GetActiveStream() stream.Stream {
	var stream stream.Stream

	row := DB.Con.QueryRow(`SELECT 
		name, 
		follow, 
		track, 
		locations, 
		active, 
		max_edges, 
		exclude, 
		description, 
		retweets_net, 
		mentions_net, 
		hashtags_net, 
		filter_level,
		min_follower_count, 
		min_favorite_count, 
		only_verified,
		max_hashtags, 
		max_mentions, 
		reply_net
		FROM streams 
		WHERE active = 1;`)

	var onlyVerified int

	row.Scan(
		&stream.Name,
		&stream.Follow,
		&stream.Track,
		&stream.Locations,
		&stream.Active,
		&stream.MaxEdges,
		&stream.Exclusion,
		&stream.Description,
		&stream.RetweetsNet,
		&stream.MentionsNet,
		&stream.HashtagsNet,
		&stream.FilterLevel,
		&stream.MinFollowerCount,
		&stream.MinFavoriteCount,
		&onlyVerified,
		&stream.MaxHashtags,
		&stream.MaxMentions,
		&stream.ReplyNet,
	)

	stream.OnlyVerified = intToBool(onlyVerified)

	return stream
}

func (DB *Database) StreamOnGoing() bool {
	rows := DB.Con.QueryRow("SELECT COUNT(1) FROM streams WHERE active = 1;")

	var count int
	rows.Scan(&count)

	return count > 0
}

func intToBool(i int) bool {
	return i > 0
}
