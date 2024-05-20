package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	client "github.com/axyut/dairygo/client"
	"github.com/axyut/dairygo/db"
	"github.com/axyut/dairygo/services"
	"golang.org/x/exp/slog"
)

func main() {
	db := db.NewDb()
	logg := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	services.GlobalService(logg, db)

	component := client.Index()
	http.Handle("/", templ.Handler(component))

	fmt.Println("On http://localhost:3000")
	http.ListenAndServe(":3000", nil)

}
