package server

import (
	_ "prueba_tecnica/api/docs"
	"prueba_tecnica/api/endpoints"
	"prueba_tecnica/api/repository"
	"prueba_tecnica/api/service"
	transports "prueba_tecnica/api/transports/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Server struct {
	router  *gin.Engine
	grpcSrv *grpc.Server
	client  *mongo.Client
	logger  logrus.FieldLogger
}

func NewServer(client *mongo.Client, logger logrus.FieldLogger) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		client: client,
		logger: logger,
	}
}

func (s *Server) Run() {
	eventRepo := repository.NewMongoEventRepository(s.client, s.logger)
	eventService := service.NewEventService(eventRepo, s.logger)
	eventEndpoints := endpoints.NewEventEndpoints(eventService)

	transports.NewEventRouter(s.router, eventEndpoints, s.logger)

	s.setupSwagger()

	s.router.Run(":8080")
}

func (s *Server) setupSwagger() {
	// Configuración de Swagger
	url := ginSwagger.URL("/swagger/doc.json") // La URL del archivo generado
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
