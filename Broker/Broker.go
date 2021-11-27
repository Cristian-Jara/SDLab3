package main

import(
	"log"
	"net"
	"context"
	"math/rand"
	"time"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedChatServiceServer
}
func ChooseServer()(string){
	var ChoosenServer string
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(3)
	if id == 0 {
		ChoosenServer = ":50058"//"10.6.40.225:50058" // IP1
	} else if id == 1 {
		ChoosenServer = ":50058"//"10.6.40.227:50058" // IP2
	} else {
		ChoosenServer = ":50058"//"10.6.40.229:50058" // IP3
	}
	return ChoosenServer
}

func (s *Server) GetNumberRebelds(ctx context.Context, message *pb.LeiaRequest) (*pb.LeiaReply,error){
	log.Printf("Leia se ha conectado para saber el número de rebeldes")
	log.Printf("Los parámetros son:\nPlaneta: "+ message.Planet + "\nCiudad: "+message.City)
	ChoosenServer := ChooseServer()
	log.Printf("Server escogido es: "+ChoosenServer)
	//conn, err := grpc.Dial(ChoosenServer, grpc.WithInsecure())
	//ServerService := pb.NewChatServiceClient(conn)
	//r, err := ServerService.GetNumberRebelds(context.Background(), &message)
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//return &pb.LeiaReply{Status:response.Status, Quantity: response.Quantity, X: response.X, Y: response.Y, Z: response.Z, Server: ChoosenServer},nil
	return &pb.LeiaReply{Status:"OK", Quantity: 10, X: 1, Y: 2, Z: 0, Server: ChoosenServer},nil
}

const (
	puerto = ":50052"
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