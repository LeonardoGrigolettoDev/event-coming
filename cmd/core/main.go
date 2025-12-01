package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"io/ioutil"

	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/api"
	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/db"
	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/repositories"
	"github.com/LeonardoGrigolettoDev/event-coming.git/internal/services"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("MIGRATE_ON_START") == "true" {
		b, err := ioutil.ReadFile("migrations/schema.sql")
		if err != nil {
			log.Fatal(err)
		}
		if err := db.RunMigration(ctx, pool, string(b)); err != nil {
			log.Fatal(err)
		}
	}

	contactRepo := repositories.NewContactRepository(pool)
	schedulerRepo := repositories.NewSchedulerRepository(pool)
	scRepo := repositories.NewSchedulerContactsRepository(pool)
	eventRepo := repositories.NewEventRepository(pool)
	consolidatedRepo := repositories.NewConsolidatedRepository(pool)

	_ = services.NewEventService(eventRepo, consolidatedRepo)
	_ = services.NewSchedulerService(schedulerRepo, scRepo, contactRepo, eventRepo)

	fmt.Println("core up")

	r := api.NewRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
