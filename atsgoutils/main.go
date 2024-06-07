package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	NavComp *ATSData
)

func init() {
	atsDataFilename := "./atsdata.json"
	atsData, err := ParseATSDataFromFile(atsDataFilename)
	if err != nil {
		log.Printf("Error parseing ATS Data from file %s: %s", atsDataFilename, err)
		panic(err)
	}
	NavComp = atsData
	log.Printf("NavComp Loaded")
	r, err := LoadCacheFromFile("atscache.json")
	if err != nil {
		log.Printf("Error loading cache from file %s: %s", "atscache.json", err)
		panic(err)
	}
	routeCache = r
	log.Println("Route Cache Loaded")

}

func findHeading(x, y, z, pitch, yaw, speed, lineDistance, sdist *float64, empire *string) error {
	if (x != nil && *x <= -99999) || (y != nil && *y <= -99999) || (z != nil && *z <= -99999) {
		return fmt.Errorf("expected x, y and z, received %f, %f and %f respectively", *x, *y, *z)
	}
	if (pitch != nil && *pitch >= 400) || (yaw != nil && *yaw >= 400) {
		return fmt.Errorf("expected pitch and yaw, received %f and %f respectively", *pitch, *yaw)
	}
	if speed != nil && *speed == 999 {
		return fmt.Errorf("expected speed, received %f", *speed)
	}
	headings, err := FindObjectAlongLine(*x, *y, *z, *yaw, *pitch, *speed, *lineDistance, *sdist, empire)
	if err != nil {
		return fmt.Errorf("error during findobject: %w", err)
	}
	for _, heading := range headings {
		fmt.Println(heading)
	}
	return nil
}

func main() {
	brouteCmd := flag.NewFlagSet("bestroute", flag.ExitOnError)
	brouteSource := brouteCmd.String("source", "", "Source Object Name or Partial name (e.g. magna for Magna Roma)")
	brouteTarget := brouteCmd.String("target", "", "Target Object Name or Partial name (e.g. 303 for 303)")
	brouteSpeed := brouteCmd.Float64("speed", 22, "Speed in knots")
	findHeadingCmd := flag.NewFlagSet("findheading", flag.ExitOnError)
	findHeadingX := findHeadingCmd.Float64("X", -99999, "X Coordinate")
	findHeadingY := findHeadingCmd.Float64("Y", -99999, "X Coordinate")
	findHeadingZ := findHeadingCmd.Float64("Z", -99999, "X Coordinate")
	findHeadingSpeed := findHeadingCmd.Float64("speed", 999, "Warp Speed the object was travelling at")
	findHeadingPitch := findHeadingCmd.Float64("pitch", 999, "Pitch of the object")
	findHeadingYaw := findHeadingCmd.Float64("yaw", 999, "Yaw of the object")
	findHeadingEmpire := findHeadingCmd.String("frame", "grc", "Empire, by default, this is GRC")
	findHeadingSDist := findHeadingCmd.Float64("sdist", 50, "Distance to use as the 'Same Object', you don't get the station -and- planet")
	findHeadingLineDist := findHeadingCmd.Float64("dist", 1000, "Length of virtual line to use in Parsecs")
	onoCmd := flag.NewFlagSet("ono", flag.ExitOnError)
	onoSource := onoCmd.String("source", "", "Source Object to determine objects nearby")
	onoRange := onoCmd.Float64("range", 200, "Range of objects to consider")
	onoNumResults := onoCmd.Int("num-results", 20, "Number of results to display")
	if len(os.Args) < 2 {
		log.Println("Expected subcommand of bestroute, findheading, or ono")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "bestroute":
		brouteCmd.Parse(os.Args[2:])
		if *brouteSource == "" || *brouteTarget == "" {
			log.Println("Expected 'source' and 'target' flags")
			os.Exit(1)
		}
		source, target, err := NavComp.ResolveObjects(*brouteSource, *brouteTarget)
		if err != nil {
			log.Printf("Error resolving objects: %s", err)
			panic(err)
		}
		route, err := BestRoute(source, target)
		if err != nil {
			log.Printf("Error calculating best route: %s", err)
			panic(err)
		}
		_, statement := route.GetStatement(*brouteSpeed)
		fmt.Println(statement)
		// duration := route.TimeToExecute(*brouteSpeed)
		// if err != nil {
		// 	log.Printf("Error calculating best route: %s", err)
		// 	panic(err)
		// }
		// log.Printf("It should take %s to go from %s to %s at %f", duration, *brouteSource, *brouteTarget, *brouteSpeed)
	case "findheading":
		findHeadingCmd.Parse(os.Args[2:])
		err := findHeading(findHeadingX, findHeadingY, findHeadingZ, findHeadingPitch, findHeadingYaw, findHeadingSpeed, findHeadingLineDist, findHeadingSDist, findHeadingEmpire)
		if err != nil {
			log.Printf("Error finding Heading: %s", err)
			panic(err)
		}
	case "ono":
		onoCmd.Parse(os.Args[2:])
		if *onoSource == "" {
			log.Println("Source is required, none supplied")
			os.Exit(1)
		}
		sourceObject, err := NavComp.FindObject(*onoSource)
		if err != nil {
			log.Printf("Cannot locate object from string %s: %s", *onoSource, err)
			os.Exit(1)
		}
		getNearbyObjects(sourceObject, onoRange, onoNumResults)
		if err != nil {
			log.Printf("Unable to find nearby objects: %s", err)
			os.Exit(1)
		}
	default:
		log.Println("Expected subcommand of bestroute, findheading, or ono")
		os.Exit(1)
	}
}
