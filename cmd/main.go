package main

import (
	"InvestmentHelperTelegramBot/internal/news"
	"InvestmentHelperTelegramBot/internal/plot"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	plotik "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
)

type TelegramMesageHandler struct {
	NewsManager news.NewsManager
	PlotManager plot.PlotManager
}

func NewTelegramMessageHandler(newsManager news.NewsManager, plotManager plot.PlotManager) TelegramMesageHandler {
	return TelegramMesageHandler{newsManager, plotManager}
}

func (handler TelegramMesageHandler) getPlotPNG(symbol string) error {
	plotData, err := handler.PlotManager.GetPlot(symbol)
	if err != nil {
		return err
	}

	p, err := plotik.New()
	if err != nil {
		return err
	}

	p.Title.Text = fmt.Sprintf("%s Plot", symbol)
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Price"

	err = plotutil.AddLinePoints(p, "First", getPlotPoints(plotData))
	if err != nil {
		return err
	}

	if err := p.Save(16*vg.Inch, 8*vg.Inch, "points.png"); err != nil {
		return err
	}
	return nil
}

func getPlotPoints(plotData []plot.Candle)  plotter.XYs{
	pts := make(plotter.XYs, len(plotData))
	for i, _ := range pts {
		pts[i].Y = plotData[i].Close
		pts[i].X = float64(i)
	}
	return pts
}

func (handler TelegramMesageHandler) getTextNews(symbol string, newsCount int) ([]string, error){
	result, err := handler.NewsManager.GetNews(symbol)
	if err != nil{
		return nil, err
	}
	if len(result) < newsCount{
		return nil, errors.New("few News")
	}
	resultTexts := []string{}
	for _, value := range(result[:newsCount]){
		resultText := fmt.Sprintf("Headline:%s\nLink:%s", value.Headline, value.Link)
		resultTexts = append(resultTexts, resultText)
	}
	return resultTexts, nil
}

func (handler TelegramMesageHandler) handleUpdate (bot *tgbotapi.BotAPI, update tgbotapi.Update){

	symbol := update.Message.Text
	if update.Message == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pls send only text")
		bot.Send(msg)
		return
	}

	msgTexts, err :=  handler.getTextNews(symbol, 10)
	if err != nil{
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "problem with Symbol")
		bot.Send(msg)
		return
	}

	for _, msgText := range(msgTexts){
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		bot.Send(msg)
	}

	err = handler.getPlotPNG(symbol)
	if err != nil{
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "problem with Plot")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "points.png")
	bot.Send(msg)

}

func main() {
	botApi, ok := os.LookupEnv("BotApi")
	if !ok {
		log.Panic("No BotApi")
	}
	alphaApi, ok := os.LookupEnv("AlphaApiKey")
	if !ok {
		log.Panic("No AlphaApi")
	}

	newsManager := news.NewNewsManagerYahoo()
	plotManager := plot.NewPlotManagerAlphaVantage(alphaApi)
	mesageHandler := NewTelegramMessageHandler(newsManager, plotManager)

	bot, err := tgbotapi.NewBotAPI(botApi)
	if err != nil{
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		mesageHandler.handleUpdate(bot, update)
	}
}