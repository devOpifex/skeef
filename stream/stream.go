package stream

type Stream struct {
	Name        string
	Follow      string
	Track       string
	Locations   string
	Exclusion   string
	MaxEdges    int
	Active      int
	Description string
	RetweetsNet int
	MentionsNet int
	HashtagsNet int
	FilterLevel string
}
