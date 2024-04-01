package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// This file is to help determine the heading of other ships
// For instance, often a bot net will say there's a contact or border crossing
// These functions will help in those calculations

type Heading struct {
	Yaw, Pitch float64
}

type HeadingResult struct {
	Distance, Time, ContainingRadius float64
	BodyOfInterest                   AstralBody
}

func (h HeadingResult) String() string {
	duration := time.Duration(h.Time * 1e9).Truncate(time.Second)
	return fmt.Sprintf("%-20s\t%20s\t[%.2f]", h.BodyOfInterest.Name, duration, h.Distance)
}

type ByDistance []HeadingResult

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func Rads(angle float64) float64 {
	return angle * (math.Pi / 180)
}

func ProjectHeading(heading Heading, spoint Point, d float64) *Point {
	nx := d * math.Cos(Rads(heading.Yaw)) * math.Cos(Rads(heading.Pitch))
	ny := d * math.Sin(Rads(heading.Yaw)) * math.Cos(Rads(heading.Pitch))
	nz := d * math.Sin(Rads(heading.Pitch))
	return &Point{
		X: spoint.X + nx,
		Y: spoint.Y + ny,
		Z: spoint.Z + nz,
	}
}

func ConvertToGRC(p Point, frame string) (*Point, error) {
	// Converts from a given frame to Galactic Real Coordinates
	for _, emp := range NavComp.NavcompDB.Empires {
		for _, bor := range emp.Borders {
			if strings.Contains(strings.ToLower(bor.Name), strings.ToLower(frame)) {
				return &Point{
					X: p.X + bor.X,
					Y: p.Y + bor.Y,
					Z: p.Z + bor.Z,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find frame %s", frame)
}

func bourkianDeterminant(p1, p2, s *Point, sd float64) bool {
	a := math.Pow((p2.X-p1.X), 2) + math.Pow((p2.Y-p1.Y), 2) + math.Pow((p2.Z-p1.Z), 2)
	b := 2 * ((p2.X-p1.X)*(p1.X-s.X) + (p2.Y-p1.Y)*(p1.Y-s.Y) + (p2.Z-p1.Z)*(p1.Z-s.Z))
	c := math.Pow(s.X, 2) + math.Pow(s.Y, 2) + math.Pow(s.Z, 2) + math.Pow(p1.X, 2) + math.Pow(p1.Y, 2) + math.Pow(p1.Z, 2) - 2*(s.X*p1.X+s.Y*p1.Y+s.Z*p1.Z) - math.Pow(sd, 2)
	bourke := math.Pow(b, 2) - 4*a*c
	return bourke >= 0
}

func FilterAstralBodiesByBkAndDist(source, projected *Point, rad, sdist float64) func(astralBody AstralBody) bool {
	return func(astralBody AstralBody) bool {
		bkDeterminant := bourkianDeterminant(&astralBody.Point, projected, source, rad)
		distFilter := astralBody.DistanceToPoint(*source) > sdist
		return bkDeterminant && distFilter
	}
}

func DetermineTimeAndOrder(source *Point, bodies []AstralBody, speed, rad float64) ([]HeadingResult, error) {
	var results []HeadingResult
	for _, body := range bodies {
		d := body.DistanceToPoint(*source)
		t, err := body.TimeToPoint(*source, AVG_COCHRANE_DENSITY, speed)
		if err != nil {
			continue
		}
		heading := HeadingResult{
			Distance:       d,
			Time:           t,
			BodyOfInterest: body,
		}
		results = append(results, heading)
	}
	sort.Sort(ByDistance(results))
	return results, nil
}

func FindObjectAlongLine(x, y, z, yaw, pitch, speed, distance, sdist float64, frame *string) ([]HeadingResult, error) {
	var bodies []AstralBody
	var finalContainingRadius float64
	source := &Point{X: x, Y: y, Z: z}
	if frame != nil && *frame != "grc" {
		s, err := ConvertToGRC(*source, *frame)
		if err != nil {
			return nil, fmt.Errorf("cannot find object along line: %w", err)
		}
		source = s
	}
	heading := Heading{Yaw: yaw, Pitch: pitch}
	projectedHeading := ProjectHeading(heading, *source, speed)
	for crad := float64(0); crad < 12; crad += 2 {
		filter := FilterAstralBodiesByBkAndDist(source, projectedHeading, crad, sdist)
		radialBodies := NavComp.FilterBodies(filter)
		if len(radialBodies) > 0 {
			bodies = append(bodies, radialBodies...)
			finalContainingRadius = crad
			break
		}
	}
	if len(bodies) == 0 {
		return nil, fmt.Errorf("no bodies found along line")
	}
	return DetermineTimeAndOrder(source, bodies, speed, finalContainingRadius)
}
