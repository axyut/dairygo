package main

import (
	"fmt"
	"net/http"

	"github.com/axyut/dairygo/db"
	"github.com/axyut/dairygo/handler"
	"github.com/axyut/dairygo/service"
)

func main() {
	mongo, err := db.NewMongo()
	if err != nil {
		panic(err)
	}
	srv := service.NewService(mongo)
	handler.RootHandler(srv)
	fmt.Println("On http://localhost:3000")
	http.ListenAndServe(":3000", nil)
	defer mongo.Close(db.Ctx)
}
