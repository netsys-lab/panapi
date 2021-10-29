package main

import (
	"fmt"
	"os"
)

type CSVPoint struct {
	Mibps float64
}

func WriteCSV(data []CSVPoint) {
	Info("Writing CSV dump...")
	f, err := os.OpenFile("spate/pid.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		Error("Failed to create spate/pid.csv: %v", err)
	}
	defer f.Close()
	f.WriteString("Mibps\n")

	for _, item := range data {
		f.WriteString(fmt.Sprintf("%v\n", item.Mibps))
	}
}
