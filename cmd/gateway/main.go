package main

import (
	"log"
	"net"

	"github.com/AntonTyurin87/Recon_Com_protoc/gen/go/gateway"
	lib "github.com/AntonTyurin87/Recon_Com_protoc/gen/go/librarian"
	tg_bot_lib "github.com/AntonTyurin87/Recon_Com_protoc/gen/go/tg_bot_librarian"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gatewayPort       = ":50051"
	TGBotLibrarianAdr = "localhost:50052"
	LibrarianAdr      = "localhost:50053"
)

type gatewayServer struct {
	gateway.UnimplementedGatewayServer
	tgBotLibrarianClient tg_bot_lib.TG_Bot_LibrarianClient //Service1
	librarianClient      lib.LibrarianClient               // Service2
}

func NewGatewayServer() (*gatewayServer, error) {
	// Подключаемся к TG_Bot_Librarian
	tgBotLibrarianConnect, err := grpc.NewClient(TGBotLibrarianAdr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Подключаемся к Librarian
	librarianConnect, err := grpc.NewClient(LibrarianAdr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &gatewayServer{
		tgBotLibrarianClient: tg_bot_lib.NewTG_Bot_LibrarianClient(tgBotLibrarianConnect),
		librarianClient:      lib.NewLibrarianClient(librarianConnect),
	}, nil
}

//
//func (s *gatewayServer) ProcessWithTGBotLibrarian(ctx context.Context, in *gateway.GatewayRequest) (*gateway.GatewayResponse, error) {
//	log.Printf("Gateway: Processing with Service1 - %v", in.GetInput())
//
//	resp, err := s.tgBotLibrarianClient.SendMessage(ctx, &tg_bot_lib.SendMessageRequest{})
//	if err != nil {
//		return nil, err
//	}
//
//	return &pb.GatewayResponse{
//		Result:      resp.GetResult(),
//		ServiceUsed: "service1",
//		Timestamp:   time.Now().Format(time.RFC3339),
//	}, nil
//}
//
//func (s *gatewayServer) ProcessWithService2(ctx context.Context, in *pb.GatewayRequest) (*pb.GatewayResponse, error) {
//	log.Printf("Gateway: Processing with Service2 - %v", in.GetInput())
//
//	resp, err := s.librarianClient.TransformString(ctx, &pb.StringRequest{Input: in.GetInput()})
//	if err != nil {
//		return nil, err
//	}
//
//	return &pb.GatewayResponse{
//		Result:      resp.GetResult(),
//		ServiceUsed: "service2",
//		Timestamp:   time.Now().Format(time.RFC3339),
//	}, nil
//}
//
//func (s *gatewayServer) ProcessWithBoth(ctx context.Context, in *pb.GatewayRequest) (*pb.GatewayResponse, error) {
//	log.Printf("Gateway: Processing with both services - %v", in.GetInput())
//
//	// Обрабатываем через Service1
//	resp1, err := s.tgBotLibrarianClient.ProcessString(ctx, &pb.StringRequest{Input: in.GetInput()})
//	if err != nil {
//		return nil, err
//	}
//
//	// Обрабатываем через Service2
//	resp2, err := s.librarianClient.TransformString(ctx, &pb.StringRequest{Input: in.GetInput()})
//	if err != nil {
//		return nil, err
//	}
//
//	result := resp1.GetResult() + " | " + resp2.GetResult()
//
//	return &pb.GatewayResponse{
//		Result:      result,
//		ServiceUsed: "both",
//		Timestamp:   time.Now().Format(time.RFC3339),
//	}, nil
//}

func main() {
	// Создаем newGatewayServer сервер
	newGatewayServer, err := NewGatewayServer()
	if err != nil {
		log.Fatalf("Failed to create newGatewayServer server: %v", err)
	}

	// Запускаем newGatewayServer сервер
	lis, err := net.Listen("tcp", gatewayPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gateway.RegisterGatewayServer(s, newGatewayServer)

	log.Printf("Gateway gRPC server listening at %v", lis.Addr())
	log.Printf("TGBotLibrarian at: %s", TGBotLibrarianAdr)
	log.Printf("Librarian at: %s", LibrarianAdr)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
