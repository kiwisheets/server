package server

func checkIfDev(environment string) bool {
	return environment == "dev" || environment == "development"
}
