package main

import (
	"InvestmentHelperTelegramBot/internal/loggerbot"
	"InvestmentHelperTelegramBot/internal/news"
	"InvestmentHelperTelegramBot/internal/plot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	plotik "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type TelegramMesageHandler struct {
	NewsManager news.NewsManager
	PlotManager plot.PlotManager
	LoggerBot   loggerbot.LoggerBot
}

func NewTelegramMessageHandler(newsManager news.NewsManager, plotManager plot.PlotManager, loggerBot loggerbot.LoggerBot) TelegramMesageHandler {
	return TelegramMesageHandler{newsManager, plotManager, loggerBot}
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

func getPlotPoints(plotData []plot.Candle) plotter.XYs {
	pts := make(plotter.XYs, len(plotData))
	for i, _ := range pts {
		pts[i].Y = plotData[i].Close
		pts[i].X = float64(i)
	}
	return pts
}

func (handler TelegramMesageHandler) getTextNews(symbol string, newsCount int) ([]string, error) {
	result, err := handler.NewsManager.GetNews(symbol)
	if err != nil {
		return nil, err
	}
	if len(result) < newsCount {
		return nil, errors.New("few News")
	}
	resultTexts := []string{}
	for _, value := range result[:newsCount] {
		resultText := fmt.Sprintf("Headline:%s\nLink:%s", value.Headline, value.Link)
		resultTexts = append(resultTexts, resultText)
	}
	return resultTexts, nil
}

func (handler TelegramMesageHandler) handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	symbol := update.Message.Text
	if update.Message == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pls send only text")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		return errors.New("user message is not text")
	}

	msgTexts, err := handler.getTextNews(symbol, 10)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "problem with Symbol")
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return errors.New("user symbol is wrong")
	}

	for _, msgText := range msgTexts {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	}

	err = handler.getPlotPNG(symbol)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "problem with Plot")
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return errors.New("problem with sending plot")
	}
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "points.png")
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	botAPI, ok := os.LookupEnv("BotApi")
	if !ok {
		log.Panic("No BotApi")
	}
	alphaAPI, ok := os.LookupEnv("AlphaApiKey")
	if !ok {
		log.Panic("No AlphaApi")
	}
	loggerBotAPI, ok := os.LookupEnv("LoggerBotAPI")
	if !ok {
		log.Panic("No LoggerBotAPI")
	}
	loggerChatIDString, ok := os.LookupEnv("LoggerChatId")
	if !ok {
		log.Panic("No LoggerChatId")
	}
	loggerChatID, err := strconv.Atoi(loggerChatIDString)
	if err != nil {
		log.Panic("wrong LoggetChatId")
	}

	newsManager := news.NewNewsManagerYahoo()
	plotManager := plot.NewPlotManagerAlphaVantage(alphaAPI)
	logger := loggerbot.NewLoggerBot(loggerBotAPI, loggerChatID)
	mesageHandler := NewTelegramMessageHandler(newsManager, plotManager, logger)

	bot, err := tgbotapi.NewBotAPI(botAPI)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		err := mesageHandler.handleUpdate(bot, update)
		if err != nil {
			log.Println(err)
		}
		logMessage := fmt.Sprintf("%s\n%s %s\n%s",
			update.Message.From.UserName,
			update.Message.From.FirstName,
			update.Message.From.LastName,
			update.Message.Text)
		err = mesageHandler.LoggerBot.SendLog(logMessage)
		if err != nil {
			log.Println(err)
		}
	}
}
