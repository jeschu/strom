package main

import (
	cfg "github.com/jeschu/go-config"
	"log"
	"slices"
	"time"
)

type Strom struct {
	Tarife         []Tarif        `yaml:"tarife"`
	Zaehlerstaende []Zaehlerstand `yaml:"zaehlerstaende"`
}

type Tarif struct {
	Name           string    `yaml:"name"`
	Start          time.Time `yaml:"start"`
	Grundgebuehr   float64   `yaml:"grundgebuehr"`
	ArbeitspreisCt float64   `yaml:"arbeitspreisCt"`
	Abschlag       float64   `yaml:"abschlag"`
}

type Zaehlerstand struct {
	Time    time.Time
	Zaehler float64 `yaml:"zaehler"`
	Art     Art     `yaml:"art"`
}

type Art int

const (
	ArtManuell Art = iota
	ArtBerechnet
)

func (art Art) String() string {
	switch art {
	case ArtManuell:
		return "manuell"
	case ArtBerechnet:
		return "berechnet"
	default:
		return "?"
	}
}

func Load() (Strom, error) {
	strom := Strom{}
	err := cfg.ReadConfigYaml("strom.yml", &strom)
	if err != nil {
		strom.sort()
	}
	return strom, err
}

func (strom *Strom) Save() error {
	strom.sort()
	return cfg.WriteConfigYaml("strom.yml", strom)
}

func (strom *Strom) FirstZaehlerstand() Zaehlerstand { return strom.Zaehlerstaende[0] }

func (strom *Strom) LastZaehlerstand() Zaehlerstand {
	return strom.Zaehlerstaende[len(strom.Zaehlerstaende)-1]
}

func (strom *Strom) sort() {
	slices.SortStableFunc(strom.Tarife, func(a, b Tarif) int { return int(a.Start.Sub(b.Start)) })
	slices.SortStableFunc(strom.Zaehlerstaende, func(a, b Zaehlerstand) int { return int(a.Time.Sub(b.Time)) })
}

func (strom *Strom) TimeToZaehlerstand() map[time.Time]Zaehlerstand {
	var m = make(map[time.Time]Zaehlerstand)
	for _, t := range strom.Zaehlerstaende {
		m[t.Time] = t
	}
	return m
}

func (strom *Strom) GetOrCalculateZaehlerstand(tarif Tarif) Zaehlerstand {
	var before *Zaehlerstand = nil
	var after *Zaehlerstand = nil
	for _, z := range strom.Zaehlerstaende {
		if z.Time == tarif.Start {
			return z
		}
		if z.Time.Before(tarif.Start) {
			before = &z
		}
		if z.Time.After(tarif.Start) {
			after = &z
			break
		}
	}
	if before == nil {
		log.Fatalf("Kein Zählerstand vor Tarifstart %s des Tarifs '%s'", tarif.Start.Format("02.01.2006"), tarif.Name)
	}
	if after == nil {
		now := time.Now()
		first := strom.FirstZaehlerstand()
		last := strom.LastZaehlerstand()
		verbrauch := first.Zaehler - last.Zaehler
		verbrauchDuration := float64(last.Time.Sub(first.Time))
		duration := float64(now.Sub(last.Time))
		return Zaehlerstand{
			Time:    now,
			Zaehler: last.Zaehler + (verbrauch / verbrauchDuration * duration),
			Art:     ArtBerechnet,
		}
	} else {
		return *after
	}
}
