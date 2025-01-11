package main

import (
	"context"
	_ "embed"
	"github.com/goccy/go-yaml"
	"github.com/patyhank/ptd/app"
	"github.com/patyhank/ptd/app/config"
	"os"
	"os/signal"
)

//go:embed config.example.yml
var exampleConfig []byte

var cfg config.Config

func main() {
	file, err := os.ReadFile("config.yml")
	if os.IsNotExist(err) {
		err = os.WriteFile("config.yml", exampleConfig, 0644)
		if err != nil {
			panic(err)
		}
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	instance := app.NewInstance(cfg)

	err = instance.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
