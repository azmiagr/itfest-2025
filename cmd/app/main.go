package main

import (
	"itfest-2025/internal/handler/rest"
	"itfest-2025/internal/repository"
	"itfest-2025/internal/service"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/config"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/jwt"
	"itfest-2025/pkg/middleware"
	"itfest-2025/pkg/supabase"
	"log"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal("Gagal load lokasi zona waktu:", err)
	}
	time.Local = loc
	
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
	supabase := supabase.Init()
	bcrypt := bcrypt.Init()
	jwt := jwt.Init()
	svc := service.NewService(repo, bcrypt, jwt, supabase)
	middleware := middleware.Init(svc, jwt)

	r := rest.NewRest(svc, middleware)
	r.MountEndpoint()
	r.Run()

}
