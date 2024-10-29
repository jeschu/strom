package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"strconv"
	"time"
)

//goland:noinspection GoUnhandledErrorResult
func main() {
	var (
		strom Strom
		err   error
	)
	strom, err = Load()
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	if len(os.Args) > 1 {
		var input float64
		if input, err = strconv.ParseFloat(os.Args[1], 64); err == nil {
			strom.Zaehlerstaende = append(strom.Zaehlerstaende, Zaehlerstand{Time: now, Zaehler: input})
			err = strom.Save()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var fmt = message.NewPrinter(language.German)

	var timeToZaehlerstand = strom.TimeToZaehlerstand()

	count := len(strom.Tarife)
	for idx, tarif := range strom.Tarife {
		fmt.Printf("[%s] >> %s <<\n", tarif.Start.Format("02.01.2006"), tarif.Name)
		fmt.Printf("  Grundgebühr: %5.2f€ | Arbeitspreis: %5.2fct\n", tarif.Grundgebuehr, tarif.ArbeitspreisCt)
		start := strom.GetOrCalculateZaehlerstand(tarif)
		var ende Zaehlerstand
		if idx == count-1 {
			ende = strom.GetOrCalculateZaehlerstand(Tarif{Name: "Ende", Start: now})
		} else {
			ende = strom.GetOrCalculateZaehlerstand(strom.Tarife[idx+1])
		}
		tage := float64(ende.Time.Sub(start.Time).Truncate(24*time.Hour)) / float64(24*time.Hour)
		verbrauch := ende.Zaehler - start.Zaehler
		arbeitspreis := verbrauch * tarif.ArbeitspreisCt / 100.0
		grundgebuehr := tarif.Grundgebuehr * 12.0 / 365.2 * tage
		kosten := arbeitspreis + grundgebuehr
		abschlaege := tarif.Abschlag * 12.0 / 365.2 * tage
		fmt.Printf("         Tage: %.0f\n", tage)
		fmt.Printf("        Start: %8.1f kWh %-11s | Ende: %8.1f kWh %-11s | Verbrauch: %8.1f kWh (%.1f kWh/Tag)\n",
			start.Zaehler, "("+start.Art.String()+")",
			ende.Zaehler, "("+ende.Art.String()+")",
			verbrauch, verbrauch/tage)
		fmt.Printf("       Kosten: %8.2f€ = %8.2f€ (Arbeitspreis) + %8.2f€ (Grundgebühr)                      (%2.4f€/Tag)\n", kosten, arbeitspreis, grundgebuehr, kosten/tage)
		fmt.Printf("    Abschläge: %8.2f€\n", abschlaege)
		if abschlaege < kosten {
			fmt.Printf("  Nachzahlung: %8.2f€\n", abschlaege-kosten)
		} else {
			fmt.Printf("   Erstattung: %8.2f€\n", abschlaege-kosten)
		}
		_ = idx
		_ = count
	}

	_ = timeToZaehlerstand
	/*
		zaehlerstand := (strom.Zaehlerstaende[len(strom.Zaehlerstaende)-1]).Zaehler
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
	*/
}
