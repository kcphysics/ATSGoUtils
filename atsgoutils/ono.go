package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
)

var (
	OrderedList = OrderedByDist{}
)

type OrderedByDist struct {
	HeadingResults []HeadingResult
}

func (obd *OrderedByDist) Reset() {
	obd.HeadingResults = []HeadingResult{}
}

// Sorts HeadingResults by Distance, smallest first, using the
func (obd *OrderedByDist) AddHeadingResult(hr HeadingResult) {
	if len(obd.HeadingResults) <= 0 {
		obd.HeadingResults = []HeadingResult{hr}
		return
	}
	targetNdx := -1
	for ndx, ehr := range obd.HeadingResults {
		if hr.Distance < ehr.Distance {
			targetNdx = ndx
			break
		} else if hr.Distance == ehr.Distance {
			// I want to capture all items at the same row for the same distance
			newName := fmt.Sprintf("%s, %s", ehr.BodyOfInterest.Name, hr.BodyOfInterest.Name)
			obd.HeadingResults[ndx].BodyOfInterest.Name = newName
			targetNdx = -2
			break
		}
	}
	if targetNdx == -1 {
		obd.HeadingResults = append(obd.HeadingResults, hr)
	} else if targetNdx >= 0 {
		obd.HeadingResults = append(obd.HeadingResults, HeadingResult{})
		copy(obd.HeadingResults[targetNdx+1:], obd.HeadingResults[targetNdx:])
		obd.HeadingResults[targetNdx] = hr
	}
}

func FilterBodiesByDistanceFromTarget(target *AstralBody, spaceRange *float64) func(AstralBody) bool {
	return func(astralBody AstralBody) bool {
		dist := target.DistanceToObject(astralBody)
		if dist <= *spaceRange {
			hr := HeadingResult{
				Distance:         dist,
				BodyOfInterest:   astralBody,
				Time:             -1,
				ContainingRadius: *spaceRange,
			}
			OrderedList.AddHeadingResult(hr)
			return true
		}
		return false
	}
}

func PrintBodies(target *AstralBody, numResults *int) {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	ndx := 0
	t.AppendHeader(table.Row{"#", "Object", "Distance"})
	for _, hr := range OrderedList.HeadingResults {
		ndx++
		t.AppendRow(table.Row{ndx, fmt.Sprintf("%-50s", hr.BodyOfInterest.Name), fmt.Sprintf("%4.2f", hr.Distance)})
		if ndx >= *numResults {
			break
		}
	}
	fmt.Printf("Objects near %s:\n%s", target.Name, t.Render())
}

func getNearbyObjects(target *AstralBody, spaceRange *float64, numResults *int) {
	filterFunc := FilterBodiesByDistanceFromTarget(target, spaceRange)
	OrderedList.Reset()
	NavComp.FilterBodies(filterFunc)
	PrintBodies(target, numResults)
}
