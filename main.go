package main

import (
	"fmt"
	config "github.com/jeschu/go-config"
	"log"
	"time"
)

func main() {
	strom := Strom{}
	if err := config.LoadConfigYaml("strom.yml", &strom); err != nil {
		log.Fatal(err)
	}
	fmt.Print("Zählerstand? ")
	var z float32
	if _, err := fmt.Scanf("%f", &z); err != nil {
		log.Fatal(err)
	}
	now := time.Now()

}

type Strom struct {
	Start Start `yaml:"start"`
}

type Start struct {
	Datum   time.Time `yaml:"datum"`
	Zaehler float32   `yaml:"zaehler"`
}

type Tarif struct {
	Grundgebuehr   float32 `yaml:"grundgebuehr"`
	ArbeitspreisCt float32 `yaml:"arbeitspreisCt"`
}
