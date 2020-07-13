package news

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

//Структура News содержит заголовок новости и ссылку на источник с полным текстом
type News struct {
	Headline string
	Link     string
}

//интерфейс менеджера новостей, реализующие его струтуры должны иметь метод получающий символ финансового актива
//и возвращать список новостей в виде списка экземпляров структуры News
//(Tesla - название компании, TSLA - символ акций (финансового актива) этой компании на рынке)
type NewsManager interface {
	GetNews(string) ([]News, error) // принимает символ финансового актива, возвращать список новостей в виде списка экземпляров структуры News
}

//Реализация интерфейса NewsManager, отвечает за получение новостей с сайта https://finance.yahoo.com
type NewsManagerYahoo struct {
}

//Конструктор для структуры NewsManagerYahoo
func NewNewsManagerYahoo() NewsManager {
	newsManager := NewsManagerYahoo{}
	return newsManager
}

//Метод структуры NewsManagerYahoo, принимает символ финансового актива, возвращает список экземпляров структуры News
func (newsManager NewsManagerYahoo) GetNews(symbol string) ([]News, error) {
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/news?p=%s", symbol, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	news := []News{}
	html.Find("a").Each(func(i int, selection *goquery.Selection) {
		isBox := selection.Children().HasClass("StretchedBox")
		if isBox {
			headline := selection.Text()
			link, _ := selection.Attr("href")
			if !strings.Contains(link, "http") {
				link = "https://finance.yahoo.com" + link
			}
			news = append(news, News{headline, link})
		}
	})
	if len(news) == 0 {
		err = errors.New("emptyNews")
		return nil, err
	}
	return news, err
}
