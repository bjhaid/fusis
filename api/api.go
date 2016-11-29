package api

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/luizbafilho/fusis/types"
	"github.com/toorop/gin-logrus"
)

// ApiService ...
type ApiService struct {
	*gin.Engine
	balancer Balancer
	env      string
}

type Balancer interface {
	GetServices() []types.Service
	AddService(*types.Service) error
	GetService(string) (*types.Service, error)
	DeleteService(string) error

	AddDestination(*types.Service, *types.Destination) error
	GetDestination(string) (*types.Destination, error)
	GetDestinations(svc *types.Service) []types.Destination
	DeleteDestination(*types.Destination) error

	AddCheck(check types.CheckSpec) error
	DeleteCheck(check types.CheckSpec) error

	IsLeader() bool
}

//NewAPI ...
func NewAPI(balancer Balancer) ApiService {
	gin.SetMode(gin.ReleaseMode)
	as := ApiService{
		Engine:   gin.New(),
		balancer: balancer,
		env:      getEnv(),
	}

	as.registerLoggerAndRecovery()
	as.registerRoutes()
	return as
}

func (as ApiService) registerLoggerAndRecovery() {
	as.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
}

func (as ApiService) registerRoutes() {
	as.GET("/services", as.serviceList)
	as.GET("/services/:service_name", as.serviceGet)
	as.POST("/services", as.serviceCreate)
	as.DELETE("/services/:service_name", as.serviceDelete)
	as.POST("/services/:service_name/destinations", as.destinationCreate)
	as.DELETE("/services/:service_name/destinations/:destination_name", as.destinationDelete)
	as.POST("/services/:service_name/check", as.checkCreate)
	as.DELETE("/services/:service_name/check", as.checkDelete)
}

// Serve starts the api.
// Binds IP using HOST env or 0.0.0.0
// Binds to port using PORT env or 8000
func (as ApiService) Serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	address := host + ":" + port

	log.Infof("Listening on %s", address)
	as.Run(address)
}

func getEnv() string {
	env := os.Getenv("FUSIS_ENV")
	if env == "" {
		env = "development"
	}
	return env
}
