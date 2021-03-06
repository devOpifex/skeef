package stream

type Stream struct {
	Name             string
	Follow           string
	Track            string
	Locations        string
	Exclusion        string
	MaxEdges         int
	Active           int
	Description      string
	RetweetsNet      int
	MentionsNet      int
	HashtagsNet      int
	ReplyNet         int
	FilterLevel      string
	MinFollowerCount int
	MinFavoriteCount int
	OnlyVerified     bool
	MaxHashtags      int
	MaxMentions      int
}
