package main

import "go-template/cmd/seeder/utls"

func main() {
	_ = utls.SeedData("authors", `INSERT INTO public.authors
	(first_name, last_name, created_at) VALUES
		('John', 'Doe', NOW()),
		('Jane', 'Doe', NOW());`,
	)
}
