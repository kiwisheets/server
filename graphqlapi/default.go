package graphqlapi

import (
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/emvi/hide"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kiwisheets/server"
	"github.com/kiwisheets/util"
	"github.com/maxtroughear/goenv"
	"github.com/maxtroughear/logrusextension"
	"github.com/maxtroughear/logrusnrhook"
	"github.com/maxtroughear/nrextension"
	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	NrApp      *newrelic.Application
	AppName    string
	Logger     *logrus.Entry
	gqlHandler *handler.Server
}

func (a *App) Shutdown() {
	if a.NrApp != nil {
		a.NrApp.Shutdown(30 * time.Second)
	}
}

func (a *App) Handler(es graphql.ExecutableSchema) *handler.Server {
	if a.gqlHandler == nil {
		gqlHandler := handler.New(es)
		a.injectExtensions(gqlHandler)
	}
	return a.gqlHandler
}

func (a *App) SetupServer(es graphql.ExecutableSchema, cfg *util.GqlConfig, db *gorm.DB) *gin.RouterGroup {
	return server.Setup(a.Handler(es), cfg, db)
}

func (a *App) injectExtensions(gqlHandler *handler.Server) {
	gqlHandler.Use(logrusextension.LogrusExtension{
		Logger: a.Logger,
	})
	gqlHandler.Use(nrextension.NrExtension{
		NrApp: a.NrApp,
	})
}

type env struct {
	appName      string
	nrLicenseKey string
	environment  string
	hashCfg      util.HashConfig
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
		logrus.AddHook(logrusnrhook.NewNrHook(env.appName, env.nrLicenseKey, true))

		var err error
		if app.NrApp, err = newrelic.NewApplication(
			newrelic.ConfigAppName(env.appName),
			newrelic.ConfigLicense(env.nrLicenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
			func(cfg *newrelic.Config) {
				cfg.ErrorCollector.RecordPanics = true
			},
			// newrelic.ConfigLogger(nrlogrus.StandardLogger()),
		); err != nil {
			logrus.Errorf("failed to start new relic agent %v", err)
		}
	}

	return app
}

func getEnv() env {
	godotenv.Load()
	return env{
		appName:      goenv.CanGet("APP_NAME", "unnamed"),
		nrLicenseKey: goenv.CanGet("NR_LICENSE_KEY", ""),
		environment:  goenv.MustGet("ENVIRONMENT"),
		hashCfg: util.HashConfig{
			Salt:      goenv.MustGetSecretFromEnv("HASH_SALT"),
			MinLength: goenv.CanGetInt32("HASH_MIN_LENGTH", 10),
		},
	}
}
