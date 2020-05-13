package main

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"surface_attack/consts"
	"surface_attack/handlers/global/logger"
	"surface_attack/handlers/jsonHandler"
	"surface_attack/middleware"
	"surface_attack/providers/crudProvider"
	"surface_attack/webHandlers"
)

func main() {

	startLogger()
	loadEnvFile()

	crudProvider := crudProviderStart()
	defer crudProvider.EndConn()

	handleOnBootJsonData(crudProvider)
	httpServerStart(crudProvider)
}

func startLogger() {
	logger := logger.HandlerLogging{}
	logger.Init()
	log.Println("begin boot-up")
}

func loadEnvFile() {
	// loads values from .env into the system
	e := godotenv.Load()
	mainErrorHandler(e)
	log.Println(".env file loaded succesfully")
}

func crudProviderStart() *crudProvider.Implamintation {
	redis_addr, exists := os.LookupEnv(consts.REDIS_ADDR_ENV_KEY)
	mainExistHandler(exists, errors.New("value for key REDIS_ADDR does not exist"))

	crudProvider := crudProvider.Implamintation{}
	e := crudProvider.InitConn(redis_addr)
	mainErrorHandler(e)

	return &crudProvider
}

func handleOnBootJsonData(crudProvider *crudProvider.Implamintation) {
	loadJsonToCrudProvider, exists := os.LookupEnv(consts.ON_BOOT_JSON_LOAD)
	mainExistHandler(exists, errors.New("value for key ON_BOOT_JSON_LOAD does not exist"))

	b, e := strconv.ParseBool(loadJsonToCrudProvider)
	mainErrorHandler(e)

	// b true is relevant only for first worker that loades the map into the crud-provider.
	// any workers that we spin-up for scale purpuses and is working on same data will have
	//b = false.
	log.Println(consts.ON_BOOT_JSON_LOAD, " env param set to: ", b)
	if b {
		jsonHandler := jsonHandler.InputJsonHandler{
			CrudProvider: crudProvider,
		}
		e = jsonHandler.OnBoot(consts.DATA_FILE_LOCATION_NAME)
		mainErrorHandler(e)
	}
}

func httpServerStart(crudProvider *crudProvider.Implamintation) {
	//start http server
	timerWrapper := middleware.TimerWrapper{
		CrudProvider: crudProvider,
	}

	statsHandler := webHandlers.StatsHandler{
		CrudProvider: crudProvider,
	}

	attackHandler := webHandlers.AttackHandler{
		CrudProvider: crudProvider,
	}

	port, exists := os.LookupEnv(consts.PORT_ENV_KEY)
	mainExistHandler(exists, errors.New("value for key: "+consts.PORT_ENV_KEY+" does not exist"))

	http.Handle(consts.STATS_END_POINT, timerWrapper.Timer(http.HandlerFunc(statsHandler.Handle)))
	log.Println("stats end point available on: ", consts.STATS_END_POINT)
	http.Handle(consts.ATTACK_END_POINT, timerWrapper.Timer(http.HandlerFunc(attackHandler.Handle)))
	log.Println("attack end point available on: ", consts.ATTACK_END_POINT)
	log.Println("HTTP server is listening on port ", port)
	http.ListenAndServe(port, nil)
}

func mainErrorHandler(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}

func mainExistHandler(exists bool, e error) {
	if !exists {
		log.Println(e)
		os.Exit(1)
	}
}
