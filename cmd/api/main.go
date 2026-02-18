package main

import (
	"context"
	"fmt"
	"log"
	"pasengger_service/config"
	"pasengger_service/ent"

	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Load()
    dsn := fmt.Sprintf(
	"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	cfg.Database.Host,
	cfg.Database.Port,
	cfg.Database.User,
	cfg.Database.Name,
	cfg.Database.Password,
)
 client, err :=  ent.Open("postgres",dsn)
    if err != nil {
        log.Fatalf("failed opening connection to postgres: %v", err)
    }
    defer client.Close()
    // Run the auto migration tool.
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }
}
