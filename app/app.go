package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/devOpifex/skeef-app/db"
	"github.com/devOpifex/skeef-app/graph"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/golangcollege/sessions"
	"github.com/gorilla/websocket"
	"github.com/justinas/alice"
)

// Application Application object
type Application struct {
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	Database   db.Database
	Session    *sessions.Session
	License    db.License
	Addr       string
	StopStream chan bool
	Count      int
	Stream     *twitter.Stream
}

type Setup struct {
	Tables  bool
	Admin   bool
	License bool
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan message),
	}
}

func (c *Client) Read(app *Application) {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if !app.Database.StreamOnGoing() {
			if app.Stream != nil {
				app.Stream.Stop()
			}
			continue
		}

		tokens, err := app.Database.GetTokens()

		if err != nil {
			continue
		}

		search := app.Database.GetActiveStream()

		var twitterConfig = oauth1.NewConfig(tokens.ApiKey, tokens.ApiSecret)
		var token = oauth1.NewToken(tokens.AccessToken, tokens.AccessSecret)

		// http.Client will automatically authorize Requests
		var httpClient = twitterConfig.Client(oauth1.NoContext, token)

		// Twitter client
		var client = twitter.NewClient(httpClient)

		var params = &twitter.StreamFilterParams{
			Track:         []string{search.Track},
			StallWarnings: twitter.Bool(true),
		}
		app.Stream, _ = client.Streams.Filter(params)

		var demux = twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			app.Count++
			app.InfoLog.Printf("Tweet #%v\n", app.Count)
			nodes, edges := graph.GetUserNet(*tweet)
			hashNodes, hashEdges := graph.GetHashNet(*tweet)
			nodes = append(nodes, hashNodes...)
			edges = append(edges, hashEdges...)
			g := graph.BuildGraph(nodes, edges)
			m := message{
				Graph:       g,
				TweetsCount: app.Count,
			}
			c.Pool.Broadcast <- m
		}

		for message := range app.Stream.Messages {
			demux.Handle(message)
		}

	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			for client := range pool.Clients {
				client.Conn.WriteJSON(message{})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			for client := range pool.Clients {
				client.Conn.WriteJSON(message{})
			}
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func (app *Application) socket(pool *Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := app.wsUpgrade(w, r)

	if err != nil {
		log.Println(err)
	}

	client := &Client{
		Conn: ws,
		Pool: pool,
	}

	pool.Register <- client
	client.Read(app)

}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	if !app.Database.AdminExists() {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	if !app.Database.LicenseExists() {
		http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
		return
	}

	app.render(w, r, []string{"ui/html/home.page.tmpl"}, templateData{})
}

// Handlers Returns all routes
func (app *Application) Handlers() http.Handler {

	// websocket pool
	pool := NewPool()
	go pool.Start()

	standardMiddleware := alice.New(secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/setup", dynamicMiddleware.Then(http.HandlerFunc(app.setupPage)))
	mux.Post("/setup", dynamicMiddleware.Then(http.HandlerFunc(app.setupForm)))
	mux.Get("/setup/validate", dynamicMiddleware.Then(http.HandlerFunc(app.validatePage)))
	mux.Post("/setup/validate", dynamicMiddleware.Then(http.HandlerFunc(app.validateForm)))
	mux.Get("/admin/signin", dynamicMiddleware.Then(http.HandlerFunc(app.signinPage)))
	mux.Post("/admin/signin", dynamicMiddleware.Then(http.HandlerFunc(app.signinForm)))
	mux.Get("/admin", dynamicMiddleware.Then(http.HandlerFunc(app.adminPage)))
	mux.Post("/admin", dynamicMiddleware.Then(http.HandlerFunc(app.adminForm)))
	mux.Get("/admin/signout", dynamicMiddleware.ThenFunc(app.signout))
	mux.Get("/ws", dynamicMiddleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		app.socket(pool, w, r)
	}))
	mux.Get("/admin/:stream", dynamicMiddleware.Then(http.HandlerFunc(app.streamPage)))
	mux.Post("/admin/edit", dynamicMiddleware.Then(http.HandlerFunc(app.streamForm)))

	mux.Get("/static/", app.static())

	return app.Session.Enable(standardMiddleware.Then(mux))
}
