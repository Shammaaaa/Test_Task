package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"

	"github.com/shamil/Test_task/config"
	http5 "github.com/shamil/Test_task/internal/infrastructure/http"
	"github.com/shamil/Test_task/internal/infrastructure/usecase/api"
	"github.com/shamil/Test_task/internal/infrastructure/usecase/updater"
	"github.com/shamil/Test_task/internal/repository"
	"github.com/shamil/Test_task/internal/service"
	"github.com/shamil/Test_task/pkg/log"
)

func main() {

	// --config-file ./config.yml --listener 1
	application := cli.App{
		Name: "Api-Service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: true,
				Usage:    "YAML config filepath",
				EnvVars:  []string{"CONFIG_FILE"},
				FilePath: "/srv/secret/config_file",
			},
			&cli.StringFlag{
				Name:     "bind-address",
				Usage:    "IP и порт сервера, например: 0.0.0.0:3001",
				Required: false,
				Value:    "0.0.0.0:3003",
				EnvVars:  []string{"BIND_ADDRESS"},
			},
			&cli.StringFlag{
				Name:     "bind-socket",
				Usage:    "Путь к Unix сокет файлу",
				Required: false,
				Value:    "/tmp/api_service.sock",
				EnvVars:  []string{"BIND_SOCKET"},
			},
			&cli.IntFlag{
				Name:     "listener",
				Usage:    "Unix socket or TCP",
				Required: false,
				Value:    1,
				EnvVars:  []string{"LISTENER"},
			},
		},
		Action: Main,
		After: func(c *cli.Context) error {
			log.Info("stopped")
			return nil
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Error(err)
	}

}

func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		log.Info("app context is canceled, Api-Service is down!")
	}()

	cfg, err := config.New(ctx.String("config-file"))
	if err != nil {
		return err
	}

	apis, err := service.New(appContext, &service.Options{
		Database: &cfg.Database,
	})
	if err != nil {
		return err
	}

	defer func() {
		apis.Shutdown(func(err error) {
			log.Warning(err)
		})
		apis.Stacktrace()
	}()

	repo := repository.New(apis.Pool.Builder())
	useCase := updater.NewUpdaterUseCase(repo)

	go func() {
		useCase.Work(apis.Context())
	}()

	usecase := api.NewApiUseCase(repo)
	handler := http5.New(usecase)
	r := chi.NewRouter()
	handler.MountRoutes(r)
	r.Get("/alive", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive"))
	})
	if err := http.ListenAndServe(":3004", r); err != nil {
		log.Warningf("Failed to start server: %v", err)

	}
	return nil
}
