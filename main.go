package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"gobot/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Env      string
	BotToken string
	Dsn      string
}
type MyHandler struct {
	config *Config
	bot    bot_init.UpgradeBot
}
type Recipient interface {
	getUsr() string
}

func getUsr() string {
	return "ASDASD"
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &Config{}
	_, err := toml.DecodeFile(*configPath, cfg)

	if err != nil {
		log.Fatalf("Ошибка декодирования файла конфигов %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Ошибка подключения к БД %v", err)
	}

	upgradeBot := bot_init.UpgradeBot{
		Bot:   bot_init.InitBot(cfg.BotToken),
		Users: &models.UserModel{Db: db},
	}
	handler := &MyHandler{
		config: cfg,
		bot:    upgradeBot,
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
	upgradeBot.Bot.Handle("/start", upgradeBot.StartHandler)
	upgradeBot.Bot.Start()

}
func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	x := r.Form.Get("text")
	log.Printf(x)
	//Я сделал вывод в консоль пост запроса с нашим сообщением, осталось вывести её в переменную и разослать всем с помощью метода telebot.bot.Send()
}
