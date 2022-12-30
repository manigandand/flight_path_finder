package api

import (
	"fmt"
	"net/http"

	"github.com/manigandand/adk/api"
	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

type flightPathReq struct {
	Paths []flightRoute `json:"paths"`
}

type flightPathResponse struct {
	Result     []string `json:"result"`
	FlightPath []string `json:"flight_path"`
}

func flightPathCalculatorHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	var payload flightPathReq
	if err := api.Decode(r, &payload); err.NotNil() {
		return err
	}
	resp, err := flightPathCalculator(payload.Paths)
	if err.NotNil() {
		return err
	}

	respond.OK(w, resp)
	return nil
}

type flightRoute []string

// IsValid checks if the flight route is valid or not
// valid path should have only 2 airport codes
func (f flightRoute) IsValid() bool {
	return len(f) == 2
}

// flightPathCalculator calculates the user source and destination routes
// based on his flight journey histories.
func flightPathCalculator(inputPaths []flightRoute) (*flightPathResponse, *errors.AppError) {
	routesData, err := parseInputFlightRoutes(inputPaths)
	if err.NotNil() {
		return nil, err
	}

	for _, source := range routesData.sourceAirpots {
		routesData.tempPaths = []string{
			source,
		}
		foundPath := findRoutes(routesData, source, 1)
		if foundPath {
			break
		}
		// reset tempPaths
		routesData.tempPaths = make([]string, 0)
	}
	if len(routesData.tempPaths) == 0 {
		return nil, errors.InternalServer("can't able to find routes")
	}
	// fmt.Printf("Source: %s --> Destination: %s\n", routesData.tempPaths[0],
	// 	routesData.tempPaths[len(routesData.tempPaths)-1])
	// fmt.Println(strings.Repeat("----", 20))

	return &flightPathResponse{
		Result: []string{
			routesData.tempPaths[0],
			routesData.tempPaths[len(routesData.tempPaths)-1],
		},
		FlightPath: routesData.tempPaths,
	}, nil
}

func findRoutes(routesData *flightRoutesData, sourceAirport string, sourceCount int) bool {
	destionation := routesData.airportData[sourceAirport]
	// fmt.Printf("Current Source: %s --> Destination: %s\n", sourceAirport, destionation)
	if sourceAirport == "" || destionation == "" {
		// log.Fatal("Panic")
		return false
	}

	if destionation == routesData.destinationAirport &&
		sourceCount != len(routesData.sourceAirpots) {
		// fmt.Println("reached destination but not from the valid source, Continuing")
		return false
	}

	routesData.tempPaths = append(routesData.tempPaths, destionation)
	if destionation == routesData.destinationAirport &&
		sourceCount == len(routesData.sourceAirpots) {
		// fmt.Println("FOUND VALID ROUTE")
		return true
	}

	// untill we reach destination and correct path
	return findRoutes(routesData, destionation, sourceCount+1)
}

type flightRoutesData struct {
	airportData        map[string]string
	sourceAirpots      []string
	destinationAirport string

	// stores the routes for the source airport
	tempPaths []string
}

func parseInputFlightRoutes(inputPaths []flightRoute) (*flightRoutesData, *errors.AppError) {
	var (
		// airPortData contains each airport's next destinations
		// NOTE: currently we consider the airport has only one destinations
		// airport1 -> airport2, airport2 -> airport3, etc...
		airportData = make(map[string]string)

		// all the source airports
		sourceAirpots = make([]string, 0)

		destinationAirport string
	)

	if len(inputPaths) == 0 {
		return nil, errors.BadRequest("flight routes required")
	}

	for _, route := range inputPaths {
		if !route.IsValid() {
			return nil, errors.BadRequest(fmt.Sprintf("invalid flight route %+v", route))
		}

		// fmt.Printf("Source: %s --> Destination: %s\n", route[0], route[1])
		airportData[route[0]] = route[1]
		sourceAirpots = append(sourceAirpots, route[0])
	}

	// find the destinations airport.
	// NOTE: the assumption from the input data is, the destination airport is
	// not a another-source airport from the list.
	for _, destination := range airportData {
		if _, ok := airportData[destination]; !ok {
			destinationAirport = destination
			break
		}
	}
	if destinationAirport == "" {
		return nil, errors.BadRequest("can't able to identify the destination airport from the input")
	}
	// fmt.Printf("destinationAirport: %s \n", destinationAirport)

	return &flightRoutesData{
		airportData:        airportData,
		sourceAirpots:      sourceAirpots,
		destinationAirport: destinationAirport,
	}, nil
}
