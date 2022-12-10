package main

import (
	"testing"
)

func TestConvertKWhToWZero(t *testing.T) {
	var kWh float64 = 0
	var expected float64 = 0
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh != %8fW", kWh, expected)
	}
}

func TestConvertKWhToWsUnit(t *testing.T) {
	var kWh float64 = 1
	var expected float64 = 1 * 1000 * 3600
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh == %8fW != %8fW", kWh, expected, actual)
	}
}

func TestConvertKWhToWKilo(t *testing.T) {
	var kWh float64 = 1000
	var expected float64 = 1000 * 1000 * 3600
	actual := ScaleKWhToWs(kWh)
	if actual != expected {
		t.Fatalf("Incorrect conversion: %8fkWh == %8fW != %8fW", kWh, expected, actual)
	}
}
