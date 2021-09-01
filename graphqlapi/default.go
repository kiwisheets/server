package graphqlapi

import (
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/emvi/hide"
	"github.com/joho/godotenv"
	"github.com/kiwisheets/server"
	"github.com/kiwisheets/util"
	"github.com/maxtroughear/goenv"
	"github.com/maxtroughear/logrusextension"
	"github.com/sirupsen/logrus"
)

type App struct {
	AppName    string
	Logger     *logrus.Entry
	gqlHandler *handler.Server
}

func (a *App) Shutdown() {}

func (a *App) Handler(es graphql.ExecutableSchema) *handler.Server {
	if a.gqlHandler == nil {
		a.gqlHandler = handler.New(es)
		a.injectExtensions(a.gqlHandler)
	}
	return a.gqlHandler
}

func (a *App) SetupServer(es graphql.ExecutableSchema, cfg *util.GqlConfig) *server.Server {
	return server.Setup(a.Handler(es), cfg)
}

func (a *App) injectExtensions(gqlHandler *handler.Server) {
	gqlHandler.Use(logrusextension.LogrusExtension{
		Logger: a.Logger,
	})
}

type env struct {
	appName     string
	environment string
	logLevel    string
	hashCfg     util.HashConfig
}

func NewDefault() App {
	env := getEnv()

	hide.UseHash(hide.NewHashID(env.hashCfg.Salt, env.hashCfg.MinLength))

	logrus.SetOutput(os.Stdout)
	logLevel, err := logrus.ParseLevel(env.logLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
		logrus.Errorf("Failed to parse LOG_LEVEL environment variable: %v", err)
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	hostname, _ := os.Hostname()

	app := App{
		AppName: env.appName,
		Logger: logrus.WithFields(logrus.Fields{
			"service":  env.appName,
			"env":      env.environment,
			"hostname": hostname,
		}),
	}

	return app
}

func getEnv() env {
	godotenv.Load()
	return env{
		appName:     goenv.CanGet("APP_NAME", "unnamed"),
		environment: goenv.MustGet("ENVIRONMENT"),
		logLevel:    goenv.CanGet("LOG_LEVEL", "info"),
		hashCfg: util.HashConfig{
			Salt:      goenv.MustGetSecretFromEnv("HASH_SALT"),
			MinLength: goenv.CanGetInt("HASH_MIN_LENGTH", 10),
		},
	}
}
