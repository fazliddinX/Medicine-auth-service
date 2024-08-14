package servers

import (
	"auth-service/api"
	"auth-service/api/handler"
	"auth-service/generated/users"
	"auth-service/pkg/config"
	"auth-service/service"
	"auth-service/storage"
	"auth-service/storage/postgres"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

type Server interface {
	RunGinServer()
	RunGRPCServer()
}

func NewServer(log *slog.Logger, config2 config.Config) Server {
	return &server{log: log, config: config2}
}

type server struct {
	log    *slog.Logger
	config config.Config
}

func (s *server) RunGinServer() {
	db, err := postgres.ConnectPostgres(s.config)
	if err != nil {
		s.log.Error("Error connecting to database", "error", err)
		log.Fatal(err)
	}

	storage1 := storage.NewAuthRepo(db)

	service1 := service.NewAuthService(s.log, storage1)

	handler1 := handler.NewAuthHandler(s.log, service1)

	router := api.NewRouter(handler1)

	router.InitRouter()

	log.Fatal(router.RunRouter(s.config))
}

func (s *server) RunGRPCServer() {
	db, err := postgres.ConnectPostgres(s.config)
	if err != nil {
		s.log.Error("Error connecting to database", "error", err)
		log.Fatal(err)
	}
	a := s.config.HOST + s.config.GRPC_SERVER_PORT
	log.Println(a)
	log.Println(a)
	log.Println(a)

	listner, err := net.Listen("tcp", a)
	if err != nil {
		s.log.Error("Error listening gRPC server", "error", err)
		log.Fatal(err)
	}
	server := grpc.NewServer()

	storage2 := storage.NewUserRepo(db)

	userStorage := service.NewUserService(s.log, storage2)

	users.RegisterUserServiceServer(server, userStorage)

	s.log.Info("gRPC server listening on " + s.config.GRPC_SERVER_PORT)
	log.Println("gRPC server listening on -->" + a)

	log.Fatal(server.Serve(listner))
}
