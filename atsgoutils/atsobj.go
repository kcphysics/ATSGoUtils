package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	PARSEC               = 3085659622.014257
	LIGHTSPEED           = 29.979246
	AVG_COCHRANE_DENSITY = 1298.737508
)

var (
	Gates = map[string]*AstralBody{
		"Transwarp Gate U-02": nil,
		"Transwarp Gate T-08": nil,
		"Zausta VI":           nil,
		"Boreth":              nil,
		"Latinum Galleria":    nil,
		"Elosian City":        nil,
		"Clispau IX":          nil,
		"Kildare XI":          nil,
	}
)

type SpaceObject interface {
	CreatePoint()
	DistanceToObject(AstralBody) float64
	DistanceToPoint(Point) float64
	TimeToObject(AstralBody, float64) (float64, error)
	TimeToPoint(Point, float64, float64) (float64, error)
}

type ATSData struct {
	NavcompDB NavcompDB `json:"ATS_Navcomp_DB"`
}

func (a *ATSData) FindObject(name string) (*AstralBody, error) {
	for _, empire := range a.NavcompDB.Empires {
		body, err := empire.FindObject(name)
		if err == nil {
			return body, nil
		}
	}
	return nil, fmt.Errorf("object %s not found", name)
}

func (a *ATSData) FilterBodies(filter func(AstralBody) bool) []AstralBody {
	var bodies []AstralBody
	for _, empire := range a.NavcompDB.Empires {
		empireBodies := empire.FilterBodies(filter)
		bodies = append(bodies, empireBodies...)
	}
	return bodies
}

func (a *ATSData) ResolveObjects(source, target string) (*AstralBody, *AstralBody, error) {
	sourceObj, err := a.FindObject(source)
	if err != nil {
		return nil, nil, fmt.Errorf("error looking up source %s: %w", source, err)
	}
	targetObj, err := a.FindObject(target)
	if err != nil {
		return nil, nil, fmt.Errorf("error looking up target %s: %w", target, err)
	}
	return sourceObj, targetObj, nil
}

type NavcompDB struct {
	Version float64  `json:"version"`
	Empires []Empire `json:"empires"`
}

type Empire struct {
	Name        string       `json:"name"`
	Description string       `json:"desc"`
	Borders     []Border     `json:"borders"`
	Planets     []AstralBody `json:"planets"`
	Stations    []AstralBody `json:"stations"`
}

func (e *Empire) FilterBodies(filter func(AstralBody) bool) []AstralBody {
	var bodies []AstralBody
	for _, body := range e.Planets {
		if filter(body) {
			bodies = append(bodies, body)
		}
	}
	for _, body := range e.Stations {
		if filter(body) {
			bodies = append(bodies, body)
		}
	}
	return bodies
}

func (e Empire) FindObject(name string) (*AstralBody, error) {
	for _, body := range e.Planets {
		if body.IsObjectByName(name) {
			return &body, nil
		}
	}
	for _, body := range e.Stations {
		if body.IsObjectByName(name) {
			return &body, nil
		}
	}
	return nil, fmt.Errorf("object %s not found", name)
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

type AstralBody struct {
	Name      string  `json:"name"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Cochranes float64 `json:"cochranes"`
	Market    int64   `json:"market"`
	Point     Point   `json:"-"`
}

func (a AstralBody) IsObjectByName(name string) bool {
	return strings.Contains(strings.ToLower(a.Name), strings.ToLower(name))
}

func (a *AstralBody) CreatePoint() {
	a.Point = Point{X: a.X, Y: a.Y, Z: a.Z}
}

func (a AstralBody) DistanceToObject(target AstralBody) float64 {
	return a.Point.Distance(target.Point)
}

func (a AstralBody) DistanceToPoint(target Point) float64 {
	return a.Point.Distance(target)
}

func (a AstralBody) TimeToObject(target AstralBody, speed float64) (float64, error) {
	tCochranes := target.Cochranes
	if target.Cochranes == 0 {
		tCochranes = AVG_COCHRANE_DENSITY
	}
	return a.TimeToPoint(target.Point, tCochranes, speed)
}

func (a AstralBody) TimeToPoint(target Point, targetCochranes, speed float64) (float64, error) {
	distance := a.DistanceToPoint(target)
	if distance == 0 {
		return 0, nil
	}
	sCochranes := a.Cochranes
	if a.Cochranes == 0 {
		sCochranes = AVG_COCHRANE_DENSITY
	}
	avgCochranes := (sCochranes + targetCochranes) / 2
	velocity := math.Pow(speed, 3.33) * avgCochranes * LIGHTSPEED / PARSEC
	return distance / velocity, nil
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
		for indx, planet := range empire.Planets {
			atsData.NavcompDB.Empires[ndx].Planets[indx].CreatePoint()
			_, ok := Gates[planet.Name]
			if ok {
				Gates[planet.Name] = &atsData.NavcompDB.Empires[ndx].Planets[indx]
			}
		}
		for indx, station := range empire.Stations {
			atsData.NavcompDB.Empires[ndx].Stations[indx].CreatePoint()
			_, ok := Gates[station.Name]
			if ok {
				Gates[station.Name] = &atsData.NavcompDB.Empires[ndx].Stations[indx]
			}
		}
	}
	return &atsData, nil
}
