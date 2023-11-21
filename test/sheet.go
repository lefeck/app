package main

import (
	"github.com/unidoc/unioffice/spreadsheet"
	"log"
)

func main() {
	ss := spreadsheet.New()

	sheet := ss.AddSheet()

	for r := 0; r < 5; r++ {
		row := sheet.AddRow()
		for c := 0; c < 5; c++ {
			cell := row.AddCell()
			cell.SetString("good")
		}
	}
	if err := ss.Validate(); err != nil {
		log.Fatalf("error validating sheet: %s", err)
	}

	ss.SaveToFile("test.xlsx")
}
