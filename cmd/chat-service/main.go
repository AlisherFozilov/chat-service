package main

import (
	"flag"
	"github.com/AlisherFozilov/adisher/pkg/di"
	"github.com/AlisherFozilov/chat-service/cmd/chat-service/app"
	"github.com/AlisherFozilov/chat-service/pkg/services/messaging"
	user "github.com/AlisherFozilov/db-storage/pkg/api"
	"github.com/AlisherFozilov/mymux/pkg/exactmux"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"

	"net"
	"net/http"
)

var (
	host     = flag.String("host", "0.0.0.0", "Server host")
	port     = flag.String("port", "8888", "Server port")
	grpcHost = flag.String("grpcHost", "0.0.0.0", "grpc host")
	grpcPort = flag.String("grpcPort", "7777", "grpc port")
	dbSvcURL = flag.String("filesvc", "http://localhost:9999/api/files", "file service url")
)

func main() {
	flag.Parse()
	addr := net.JoinHostPort(*host, *port)

	start(addr)
}

func start(addr string) {

	container := di.NewContainer()
	container.Provide(
		exactmux.NewExactMux,
		func() *websocket.Upgrader {
			return &websocket.Upgrader{
				ReadBufferSize:  0,
				WriteBufferSize: 0,
			}
		},
		messaging.NewConnectorService,
		messaging.NewService,
		func() messaging.RemoteURL { return messaging.RemoteURL(*dbSvcURL) },
		func() user.StorageClient {
			conn, err := grpc.Dial(net.JoinHostPort(*grpcHost, *grpcPort), grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			return user.NewStorageClient(conn)
		},
		//func() *dbconnector.ConnectorService {
		//	return dbconnector.NewConnectorService(*dbSvcURL)
		//},
		messaging.NewService,
		app.NewServer,
	)

	container.Start()

	server := &app.Server{}
	container.Component(&server)
	panic(http.ListenAndServe(addr, server))
}
