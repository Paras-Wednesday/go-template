package main

import (
	"context"
	"fmt"
	"go-template/cmd/seeder/utls"
	"go-template/internal/config"
	"go-template/internal/postgres"
	"go-template/models"
	"log"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("couldn't load the evn: %s", err)
	}
	db, err := postgres.Connect()
	if err != nil {
		log.Fatal(err, "pop")
	}
	author, err := models.Authors(qm.OrderBy("id")).One(context.Background(), db)
	if err != nil {
		log.Fatalf("Couldn't query author: %s", err)
	}
	_ = utls.SeedData("posts", fmt.Sprintf(`
	INSERT INTO public.posts(author_id, content)
	VALUES (%d, 'This is my first post');`, author.ID))
}
