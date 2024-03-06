package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

const (
	PARSEC               = 3085659622.014257
	LIGHTSPEED           = 29.979246
	AVG_COCHRANE_DENSITY = 1298.737508
)

var (
	Gates = []string{
		"Transwarp Gate U-02",
		"Transwarp Gate T-08",
		"Zausta VI",
		"Boreth",
		"Latinum Galleria",
		"Elosian City",
		"Clispau IX",
		"Kildare",
	}
)

type ATSData struct {
	NavcompDB NavcompDB `json:"ATS_Navcomp_DB"`
}

type NavcompDB struct {
	Version float64  `json:"version"`
	Empires []Empire `json:"empires"`
}

type Empire struct {
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Borders     []Border  `json:"borders"`
	Planets     []Planet  `json:"planets"`
	Stations    []Station `json:"stations"`
}

type Border struct {
	Name   string  `json:"name"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Z      float64 `json:"z"`
	Radius float64 `json:"radius"`
	Point  Point   `json:"-"`
}

func (b *Border) CreatePoint() {
	b.Point = Point{X: b.X, Y: b.Y, Z: b.Z}
}

type Planet struct {
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Cochranes float64 `json:"cochranes"`
	Market    int64   `json:"market"`
	Point     Point   `json:"-"`
}

func (b *Planet) CreatePoint() {
	b.Point = Point{X: b.X, Y: b.Y, Z: b.Z}
}

type Station struct {
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Cochranes float64 `json:"cochranes"`
	Market    int64   `json:"market"`
	Point     Point   `json:"-"`
}

func (b *Station) CreatePoint() {
	b.Point = Point{X: b.X, Y: b.Y, Z: b.Z}
}

type Point struct {
	X float64
	Y float64
	Z float64
}

func (p *Point) Distance(p2 Point) float64 {
	nx := p.X - p2.X
	ny := p.Y - p2.Y
	nz := p.Z - p2.Z
	return math.Sqrt(math.Pow(nx, 2) + math.Pow(ny, 2) + math.Pow(nz, 2))
}

func ParseATSDataFromFile(filename string) (*ATSData, error) {
	var atsData ATSData
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filename, err)
	}
	err = json.Unmarshal(rawBytes, &atsData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %w", err)
	}
	for ndx, empire := range atsData.NavcompDB.Empires {
		for indx := range empire.Borders {
			atsData.NavcompDB.Empires[ndx].Borders[indx].CreatePoint()
		}
		for indx := range empire.Planets {
			atsData.NavcompDB.Empires[ndx].Planets[indx].CreatePoint()
		}
		for indx := range empire.Stations {
			atsData.NavcompDB.Empires[ndx].Stations[indx].CreatePoint()
		}
	}
	return &atsData, nil
}
