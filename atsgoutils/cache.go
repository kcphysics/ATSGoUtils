package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	routeCache *RouteCache
)

type RouteCache struct {
	Version   string           `json:"version"`
	RouteMap  map[string]Route `json:"routeMap"`
	NumHits   int              `json:"numHits"`
	NumMisses int              `json:"nuMisses"`
}

func (r *RouteCache) StoreRoute(route Route) {
	r.RouteMap[route.Name] = route
}

func (r *RouteCache) GetRouteFromStrings(source, target string) (*Route, error) {
	sourceObj, targetObj, err := NavComp.ResolveObjects(source, target)
	if err != nil {
		return nil, fmt.Errorf("error resolving objects: %w", err)
	}
	return r.GetRouteFromBodies(sourceObj, targetObj)
}

func (r *RouteCache) GetRouteFromBodies(source, target *AstralBody) (*Route, error) {
	rName := GetRouteName(source, target)
	route, ok := r.RouteMap[rName]
	if ok {
		r.NumHits = r.NumHits + 1
		return &route, nil
	}
	// Here we assume that if we can't find it, we're building a route from directs
	route, err := DirectRoute(source, target)
	if err != nil {
		return nil, fmt.Errorf("error getting direct route %s: %w", rName, err)
	}
	r.RouteMap[rName] = route
	r.NumMisses = r.NumMisses + 1
	return &route, nil
}

func (r RouteCache) WriteToFile(fname string) error {
	rbyte, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("unable to marshal cache: %w", err)
	}
	err = os.WriteFile(fname, rbyte, 0644)
	if err != nil {
		return fmt.Errorf("unable to write cache to %s: %w", fname, err)
	}
	return nil
}

func (r *RouteCache) ReadFromFile(fname string) error {
	rbyte, err := os.ReadFile(fname)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("File %s does not exist, loading new cache", fname)
			return nil
		} else {
			return fmt.Errorf("unable to read cache from %s: %w", fname, err)
		}
	}
	err = json.Unmarshal(rbyte, r)
	if err != nil {
		return fmt.Errorf("unable to unmarshal cache: %w", err)
	}
	return nil
}

func LoadCacheFromFile(fname string) (*RouteCache, error) {
	var r RouteCache
	err := r.ReadFromFile(fname)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
