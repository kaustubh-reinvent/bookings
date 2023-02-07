package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/kaustubh-reinvent/bookings/cmd/pkg/config"
	"github.com/kaustubh-reinvent/bookings/cmd/pkg/handlers"
	"github.com/kaustubh-reinvent/bookings/cmd/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	fmt.Println(fmt.Sprintf("Starting to listen on port %s", portNumber))

	//_ = http.ListenAndServe(portNumber, nil)

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
}
