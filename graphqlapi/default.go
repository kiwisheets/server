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
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (a *App) SetupServer(es graphql.ExecutableSchema, cfg *util.GqlConfig, db *gorm.DB) *server.Server {
	return server.Setup(a.Handler(es), cfg, db)
}

func (a *App) injectExtensions(gqlHandler *handler.Server) {
	gqlHandler.Use(logrusextension.LogrusExtension{
		Logger: a.Logger,
	})
}

type env struct {
	appName     string
	environment string
	hashCfg     util.HashConfig
}

func NewDefault() App {
	env := getEnv()

	hide.UseHash(hide.NewHashID(env.hashCfg.Salt, env.hashCfg.MinLength))

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	hostname, _ := os.Hostname()

	app := App{
		AppName: env.appName,
		Logger: logrus.WithFields(logrus.Fields{
			"service":  env.appName,
			"env":      env.environment,
			"hostname": hostname,
		}),
	}

	// gin.SetMode(gin.ReleaseMode)

	if env.environment == "production" {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(nrlogrusplugin.ContextFormatter{})
	}

	return app
}

func getEnv() env {
	godotenv.Load()
	return env{
		appName:     goenv.CanGet("APP_NAME", "unnamed"),
		environment: goenv.MustGet("ENVIRONMENT"),
		hashCfg: util.HashConfig{
			Salt:      goenv.MustGetSecretFromEnv("HASH_SALT"),
			MinLength: goenv.CanGetInt32("HASH_MIN_LENGTH", 10),
		},
	}
}
