package main

import (
	"fmt"
	cfg "github.com/jeschu/go-config"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var (
		zaehlerstand float64
		err          error
	)
	strom := Strom{}
	if err = cfg.ReadConfigYaml("strom.yml", &strom); err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 1 {
		os.Exit(1)
	}
	if zaehlerstand, err = strconv.ParseFloat(os.Args[1], 64); err != nil {
		panic(err)
	}
	strom.Zaehlerstaende = append(strom.Zaehlerstaende, Zaehlerstand{Time: time.Now(), Zaehler: zaehlerstand})
	err = cfg.WriteConfigYaml("strom.yml", &strom)
	verbrauch := zaehlerstand - strom.Start.Zaehler
	now := time.Now()
	days := now.Sub(strom.Start.Datum).Hours() / 24.0
	grundgebuehr := strom.Tarif.Grundgebuehr * 12.0 / 365.2 * days
	arbeitsPreis := strom.Tarif.ArbeitspreisCt / 100.0 * verbrauch
	gesamtPreis := grundgebuehr + arbeitsPreis
	abschlaege := strom.Abschlag * 12.0 / 362.2 * days
	diff := gesamtPreis - abschlaege
	fmt.Printf("    Start: %10.1fkWh am %s\n", strom.Start.Zaehler, strom.Start.Datum.Format("02.01.2006"))
	fmt.Printf("  Aktuell: %10.1fkWh am %s\n", zaehlerstand, now.Format("02.01.2006"))
	fmt.Printf("Verbrauch: %10.1fkWh in %.0f Tagen\n", verbrauch, days)
	fmt.Printf("   Grundgebühr: %10.2f€\n", grundgebuehr)
	fmt.Printf("  Arbeitspreis: %10.2f€\n", arbeitsPreis)
	fmt.Printf("   Gesamtpreis: %10.2f€\n", gesamtPreis)
	fmt.Printf("   - Abschläge: %10.2f€\n", abschlaege)
	if diff < 0 {
		fmt.Printf("    Erstattung: %10.2f€ (%.2f€ für ein Jahr)\n", -diff, -diff/days*365.2)
	} else {
		fmt.Printf("   Nachzahlung: %10.2f€ (%.2f€ für ein Jahr)\n", diff, diff/days*365.2)
	}
}

type Strom struct {
	Start          Start          `yaml:"start"`
	Tarif          Tarif          `yaml:"tarif"`
	Abschlag       float64        `yaml:"abschlag"`
	Zaehlerstaende []Zaehlerstand `yaml:"zaehlerstaende"`
}

type Start struct {
	Datum   time.Time `yaml:"datum"`
	Zaehler float64   `yaml:"zaehler"`
}

type Tarif struct {
	Grundgebuehr   float64 `yaml:"grundgebuehr"`
	ArbeitspreisCt float64 `yaml:"arbeitspreisCt"`
}

type Zaehlerstand struct {
	Time    time.Time
	Zaehler float64 `yaml:"zaehler""`
}
