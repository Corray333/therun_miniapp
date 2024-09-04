package app

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Corray333/therun_miniapp/docs"
	"github.com/Corray333/therun_miniapp/internal/config"
	"github.com/Corray333/therun_miniapp/internal/domains/user"
	"github.com/Corray333/therun_miniapp/internal/files"
	"github.com/Corray333/therun_miniapp/internal/storage"
	"github.com/Corray333/therun_miniapp/internal/telegram"
	"github.com/Corray333/therun_miniapp/pkg/server/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	server http.Server
}

func New() *App {
	config.MustInit("../.env")

	router := chi.NewMux()
	router.Use(logger.NewLoggerMiddleware())

	// TODO: get allowed origins, headers and methods from cfg
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования предзапроса (в секундах)
	}))

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)

	// TODO: add timeouts
	server := http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}

	telegramClient := telegram.NewClient(os.Getenv("BOT_TOKEN"))
	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	fileManager := files.New()

	user.NewUserController(router, telegramClient.GetBot(), store, fileManager)

	return &App{
		server: server,
	}
}

func (app *App) Run() {

	slog.Error(app.server.ListenAndServe().Error())
}
