package app

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Corray333/therun_miniapp/docs"
	"github.com/Corray333/therun_miniapp/internal/config"
	"github.com/Corray333/therun_miniapp/internal/domains/battle"
	"github.com/Corray333/therun_miniapp/internal/domains/car"
	"github.com/Corray333/therun_miniapp/internal/domains/cases"
	"github.com/Corray333/therun_miniapp/internal/domains/city"
	"github.com/Corray333/therun_miniapp/internal/domains/farming"
	"github.com/Corray333/therun_miniapp/internal/domains/round"
	"github.com/Corray333/therun_miniapp/internal/domains/task"
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

type controller interface {
	Build()
	Run()
}

type App struct {
	server      *http.Server
	controllers []controller
}

func (app *App) AddController(c controller) {
	app.controllers = append(app.controllers, c)
}

func New() *App {
	config.MustInit("../.env")

	app := &App{}

	router := chi.NewMux()
	router.Use(logger.NewLoggerMiddleware())

	// TODO: get allowed origins, headers and methods from cfg
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования предзапроса (в секундах)
	}))

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)

	// TODO: add timeouts
	server := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}

	app.server = server

	telegramClient := telegram.NewClient(os.Getenv("BOT_TOKEN"))
	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	fileManager := files.New()

	userController := user.NewUserController(router, telegramClient, store, fileManager)
	app.AddController(userController)

	farmingController := farming.NewFarmingController(router, store, userController.GetService())
	app.AddController(farmingController)

	taskController := task.NewTaskController(router, store, telegramClient)
	app.AddController(taskController)

	roundController := round.NewRoundController(router, store)
	app.AddController(roundController)

	battleController := battle.NewBattleController(router, store, userController.GetService(), roundController.GetService())
	app.AddController(battleController)

	casesController := cases.NewCasesController(router, store, userController.GetService())
	app.AddController(casesController)

	cityController := city.NewCityController(router, store, userController.GetRepository())
	app.AddController(cityController)

	carsControlles := car.NewCarController(router, store, userController.GetService(), roundController.GetService())
	app.AddController(carsControlles)

	return app
}

func (app *App) Init() *App {
	for _, c := range app.controllers {
		c.Build()
	}
	return app
}

func (app *App) Run() {
	slog.Info("Server started at " + app.server.Addr)
	for _, c := range app.controllers {
		go c.Run()
	}
	slog.Error(app.server.ListenAndServe().Error())
}
