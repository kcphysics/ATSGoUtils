package main

import (
	"flag"
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
}

func main() {
	brouteCmd := flag.NewFlagSet("bestroute", flag.ExitOnError)
	brouteSource := brouteCmd.String("source", "", "Source Object Name or Partial name (e.g. magna for Magna Roma)")
	brouteTarget := brouteCmd.String("target", "", "Target Object Name or Partial name (e.g. 303 for 303)")
	brouteSpeed := brouteCmd.Float64("speed", 22, "Speed in knots")
	if len(os.Args) < 2 {
		log.Println("Expected 'bestroute' subcommand")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "bestroute":
		brouteCmd.Parse(os.Args[2:])
		if *brouteSource == "" || *brouteTarget == "" {
			log.Println("Expected 'source' and 'target' flags")
			os.Exit(1)
		}
		duration, err := DirectRoute(*brouteSource, *brouteTarget, *brouteSpeed)
		if err != nil {
			log.Printf("Error calculating best route: %s", err)
			panic(err)
		}
		log.Printf("It should take %s to go from %s to %s at %f", duration, *brouteSource, *brouteTarget, *brouteSpeed)
	default:
		log.Println("Expected 'bestroute' subcommand")
		os.Exit(1)
	}
	// log.Println("NavComp Loaded")
	// var duration time.Duration
	// var speed float64 = 27
	// source := "magna"
	// dest := "303"
	// log.Printf("Calculating route from %s to %s with speed %f", source, dest, speed)
	// magna, err := NavComp.FindObject("magna")
	// if err != nil {
	// 	log.Printf("Error getting comp magna: %s", err)
	// 	panic(err)
	// }
	// log.Printf("%#v", magna)
	// epsilon303, err := NavComp.FindObject("303")
	// if err != nil {
	// 	log.Printf("Error getting comp 303: %s", err)
	// 	panic(err)
	// }
	// log.Printf("%#v", epsilon303)
	// timeToPoint, err := magna.TimeToObject(*epsilon303, speed)
	// if err != nil {
	// 	log.Printf("Error calculating time to point: %s", err)
	// 	panic(err)
	// }
	// log.Printf("It should take %f seconds to go from %s to %s at %f", timeToPoint, source, dest, speed)
	// val := int64(math.Ceil(timeToPoint * 1e9))
	// duration = time.Duration(val)
	// log.Printf("It should take %s to go from %s to %s at %f", duration, source, dest, speed)
}
