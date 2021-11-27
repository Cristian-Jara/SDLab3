package main

import(
	"log"
	"net"
	"context"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}

func (s *Server) GetNumberRebelds(ctx context.Context, message *pb.LeiaRequest) (*pb.LeiaReply,error){
	return &pb.LeiaReply{Status:"OK", Quantity: 10, X: 1, Y: 2, Z: 0},nil
}

const (
	puerto = ":50052"
	Server1Address = ""
	Server2Address = ""
	Server3Address = ""
)

func main(){
	lis, err := net.Listen("tcp", puerto)
	if err != nil{
		log.Fatalf("failed to listen: %v",err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s,&Server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil{
		log.Fatalf("Failed to servve: %v", err)
	}
}