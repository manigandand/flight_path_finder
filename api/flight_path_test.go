package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlightPathCalculator(t *testing.T) {
	t.Run("test find flight path calculator - invalid routes", func(t *testing.T) {
		inputPaths := []flightRoute{}
		_, err := flightPathCalculator(inputPaths)
		assert.NotNil(t, err)
		assert.Equal(t, "flight routes required", err.Error())

		inputPaths = []flightRoute{
			{"SFO"},
		}
		_, err = flightPathCalculator(inputPaths)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid flight route [SFO]", err.Error())

		inputPaths = []flightRoute{
			{"SFO", "IND", "NYK"},
		}
		_, err = flightPathCalculator(inputPaths)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid flight route [SFO IND NYK]", err.Error())

		inputPaths = []flightRoute{
			{"ATL", "SFO"},
			{"SFO", "ATL"},
		}
		_, err = flightPathCalculator(inputPaths)
		assert.NotNil(t, err)
		assert.Equal(t, "can't able to identify the destination airport from the input", err.Error())
	})

	t.Run("test find flight path calculator", func(t *testing.T) {
		inputPaths := []flightRoute{
			{"SFO", "EWR"},
		}
		res, err := flightPathCalculator(inputPaths)
		assert.Nil(t, err)
		assert.Equal(t, []string{"SFO", "EWR"}, res.Result)
		assert.Equal(t, []string{"SFO", "EWR"}, res.FlightPath)

		inputPaths = []flightRoute{
			{"ATL", "EWR"},
			{"SFO", "ATL"},
		}
		res, err = flightPathCalculator(inputPaths)
		assert.Nil(t, err)
		assert.Equal(t, []string{"SFO", "EWR"}, res.Result)
		assert.Equal(t, []string{"SFO", "ATL", "EWR"}, res.FlightPath)

		inputPaths = []flightRoute{
			{"IND", "EWR"},
			{"SFO", "ATL"},
			{"GSO", "IND"},
			{"ATL", "GSO"},
		}
		res, err = flightPathCalculator(inputPaths)
		assert.Nil(t, err)
		assert.Equal(t, []string{"SFO", "EWR"}, res.Result)
		assert.Equal(t, []string{"SFO", "ATL", "GSO", "IND", "EWR"}, res.FlightPath)
	})
}
