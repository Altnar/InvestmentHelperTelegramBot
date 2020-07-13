package news

import (
	"fmt"
	"testing"
)

func TestGetNewsYahoo(t *testing.T) {
	testSymbolReal := "IBM"
	testSymbolUnreal := "unrealSymbol"
	newsManagerTest := NewNewsManagerYahoo()
	t.Run(fmt.Sprintf("test real stock symbol:%s", testSymbolReal), func(t *testing.T) {
		news, err := newsManagerTest.GetNews(testSymbolReal)
		if err != nil {
			t.Error(err)
		}
		if len(news) == 0 {
			t.Error("empty news")
		}
	})

	t.Run(fmt.Sprintf("test unreal stock symbol:%s", testSymbolUnreal), func(t *testing.T) {
		news, err := newsManagerTest.GetNews(testSymbolUnreal)
		if err == nil {
			t.Error(err)
		}
		if len(news) != 0 {
			t.Error("not empty news")
		}
	})
}
