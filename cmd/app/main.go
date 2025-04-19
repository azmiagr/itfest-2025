package main

import (
	"itfest-2025/internal/handler/rest"
	"itfest-2025/internal/repository"
	"itfest-2025/internal/service"
	"itfest-2025/pkg/config"
	"itfest-2025/pkg/database/mariadb"
	"log"
)

func main() {
	config.LoadEnvironment()

	db, err := mariadb.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = mariadb.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)

	r := rest.NewRest(svc)
	r.MountEndpoint()
	r.Run()

}
