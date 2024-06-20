package main

import (
	"fmt"
	"go-template/cmd/seeder/utls"
	"go-template/pkg/utl/secure"
)

func main() {
	sec := secure.New(1, nil)
	query := "INSERT INTO public.authors" +
		" (first_name, last_name, created_at, email, password) VALUES" +
		" ('John', 'Doe', NOW(),'johndoe@mail.com','%s')," +
		" ('Jane', 'Doe', NOW(),'janedoe@mail.com','%s');"
	_ = utls.SeedData("authors", fmt.Sprintf(
		query, sec.Hash("password"), sec.Hash("password")),
	)
}
