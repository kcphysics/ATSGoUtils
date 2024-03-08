package main

import (
	"fmt"
	"math"
	"time"
)

type Route struct {
	Name           string
	Source, Target *AstralBody
	IsDirect       bool
	Distance       float64
	Stops          []*Route
}

func (r *Route) AverageCochranes() float64 {
	sCochranes := r.Source.Cochranes
	if sCochranes == 0 {
		sCochranes = AVG_COCHRANE_DENSITY
	}
	tCochranes := r.Target.Cochranes
	if tCochranes == 0 {
		tCochranes = AVG_COCHRANE_DENSITY
	}
	return (sCochranes + tCochranes) / 2
}

func (r *Route) TimeToExecute(speed float64) time.Duration {
	avgCochranes := r.AverageCochranes()
	velocity := math.Pow(speed, 3.33) * avgCochranes * LIGHTSPEED / PARSEC
	if r.IsDirect {
		rTime := (r.Distance / velocity) * 1e9
		return time.Duration(rTime)
	}
	timeToExecute := float64(0)
	for _, stop := range r.Stops {
		stopTime := stop.TimeToExecute(speed)
		timeToExecute += stopTime.Seconds()
	}
	rTime := (r.Distance / velocity) * 1e9
	return time.Duration(rTime)
}

func GetRouteName(source, target *AstralBody) string {
	return fmt.Sprintf("%s to %s", source.Name, target.Name)
}

func DirectRoute(sourceObj, targetObj *AstralBody) (Route, error) {
	distance := sourceObj.DistanceToObject(*targetObj)
	return Route{
		Name:     GetRouteName(sourceObj, targetObj),
		Source:   sourceObj,
		Target:   targetObj,
		IsDirect: true,
		Distance: distance,
	}, nil
}

func ShortestRouteToGates(astralObj *AstralBody) (*Route, error) {
	var targetRoute *Route
	shortestDistance := float64(0)
	for _, gate := range Gates {
		route, err := routeCache.GetRouteFromBodies(astralObj, gate)
		if err != nil {
			return nil, fmt.Errorf("error getting route from %s to %s: %w", astralObj.Name, gate.Name, err)
		}
		if shortestDistance == 0 || route.Distance < shortestDistance {
			shortestDistance = route.Distance
			targetRoute = route
		}
	}
	if targetRoute == nil {
		return nil, fmt.Errorf("unable to route from %s to any gates", astralObj.Name)
	}
	return targetRoute, nil
}

func BestRoute(source, target *AstralBody) (*Route, error) {
	firstRoute, err := routeCache.GetRouteFromBodies(source, target)
	if err != nil {
		return nil, fmt.Errorf("error getting route from %s to %s: %w", source.Name, target.Name, err)
	}
	if !firstRoute.IsDirect {
		// If this isn't a direct route, we know we have the best
		// in the cache and can return early
		return firstRoute, nil
	}
	firstLeg, err := ShortestRouteToGates(source)
	if err != nil {
		return nil, fmt.Errorf("error getting shortest route to gates from %s: %w", source.Name, err)
	}
	secondLeg, err := ShortestRouteToGates(target)
	if err != nil {
		return nil, fmt.Errorf("error getting shortest route to gates from %s: %w", target.Name, err)
	}
	if firstLeg.Distance+secondLeg.Distance < firstRoute.Distance {
		route := Route{
			Name:     GetRouteName(firstLeg.Source, secondLeg.Target),
			Source:   firstLeg.Source,
			Target:   secondLeg.Target,
			IsDirect: false,
			Distance: firstLeg.Distance + secondLeg.Distance,
			Stops: []*Route{
				firstLeg,
				secondLeg,
			},
		}
		routeCache.StoreRoute(route)
		return &route, nil
	}
	return firstRoute, nil
}
