package data

import (
	_ "embed"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

//go:embed boundingBoxes.csv
var boxData string

type BoundingBox struct {
	Country string
	LongMin float32
	LatMin  float32
	LongMax float32
	LatMax  float32
}

func ReadBoundingBox() ([]BoundingBox, error) {
	var boxes []BoundingBox

	reader := csv.NewReader(strings.NewReader(boxData))

	index := 0
	for {

		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		index++

		if index == 1 {
			continue
		}

		longMin, _ := strconv.ParseFloat(row[1], 32)
		latMin, _ := strconv.ParseFloat(row[2], 32)
		longMax, _ := strconv.ParseFloat(row[3], 32)
		latMax, _ := strconv.ParseFloat(row[4], 32)

		box := BoundingBox{
			Country: row[0],
			LongMin: float32(longMin),
			LatMin:  float32(latMin),
			LongMax: float32(longMax),
			LatMax:  float32(latMax),
		}

		boxes = append(boxes, box)

	}

	return boxes, nil
}
