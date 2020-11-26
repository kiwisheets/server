//go:generate go-enum -f=$GOFILE --marshal

package server

// Health is an enumeration of health states
/*
ENUM(
Starting
Healthy
Unhealthy
Stopping
)
*/
type Health int

var health Health

func SetHealthStatus(h Health) {
	health = h
}

func GetHealthStatus() Health {
	return health
}
