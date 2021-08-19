module github.com/kiwisheets/server

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/abice/go-enum v0.3.4
	github.com/emvi/hide v1.1.2
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/kiwisheets/auth v0.0.9
	github.com/kiwisheets/util v0.0.13
	github.com/maxtroughear/goenv v0.0.4
	github.com/maxtroughear/logrusextension v0.0.1
	github.com/maxtroughear/logrusnrhook v0.0.2
	github.com/maxtroughear/nrextension v0.0.1
	github.com/newrelic/go-agent/v3 v3.14.1
	github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin v1.0.1
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	gorm.io/gorm v1.21.13
)
