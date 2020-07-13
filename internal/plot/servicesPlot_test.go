package plot

import (
	"fmt"
	"testing"
)

func TestGetPlotAlphaVantage(t *testing.T) {
	testSymbolReal := "IBM"
	testSymbolUnreal := "unrealSymbol"
	apiKey := "FUSAUHUZ3W0NG9V0"
	plotManagerTest := NewPlotManagerAlphaVantage(apiKey)

	t.Run(fmt.Sprintf("test real stock symbol:%s", testSymbolReal), func(t *testing.T) {
		plot, err := plotManagerTest.GetPlot(testSymbolReal)
		if err != nil {
			t.Error(err)
		}
		if len(plot) == 0 {
			t.Error("empty plot")
		}
	})

	t.Run(fmt.Sprintf("test unreal stock symbol:%s", testSymbolUnreal), func(t *testing.T) {
		plot, err := plotManagerTest.GetPlot(testSymbolUnreal)
		if err == nil {
			t.Error(err)
		}
		if len(plot) != 0 {
			t.Error("not empty plot")
		}
	})
}
