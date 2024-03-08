package main

import (
	"fmt"
	"time"
)

type Route struct {
	Source, Target *AstralBody
	IsDirect       bool
	Stops          []*AstralBody
}

func DirectRoute(source, target string, speed float64) (time.Duration, error) {
	sourceObj, err := NavComp.FindObject(source)
	if err != nil {
		return 0, fmt.Errorf("error finding source object %s: %w", source, err)
	}
	targetObj, err := NavComp.FindObject(target)
	if err != nil {
		return 0, fmt.Errorf("error finding target object %s: %w", target, err)
	}
	tto, err := sourceObj.TimeToObject(*targetObj, speed)
	if err != nil {
		return 0, fmt.Errorf("error calculating time from %s to %s at %f: %w", source, target, speed, err)
	}
	return time.Duration(int64(tto * 1e9)), nil // Here we convert the seconds we get from TimtToObject to nanoseconds
}

func BestRoute(source, target *AstralBody, speed float64) (*Route, error) {
	directRouteTime, err := DirectRoute(source.Name, target.Name, speed)
	if err != nil {
		return nil, fmt.Errorf("error calculating direct route: %w", err)
	}

}
