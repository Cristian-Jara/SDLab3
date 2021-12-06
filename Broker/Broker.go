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
func ChooseServer(X, Y, Z int)(string){ //Función que elige un servidor
	var ChoosenServer string
	rand.Seed(time.Now().UnixNano())
	if int(X) > int(Y) && int(X) > int(Z) { //[1,0,0] el último en editarse fue el primero
		ChoosenServer = ":50058"//"10.6.40.225:50058" // IP1
	} else if int(Y) > int(X) && int(Y) > int(Z) { //[0,1,0] el último en editarse el segundo
		ChoosenServer = ":50058"//"10.6.40.227:50058" // IP2	
	} else if int(Z) > int(Y) && int(Z) > int(X) { //[0,0,1] el último en editarse el tercero
		ChoosenServer = ":50058"//"10.6.40.229:50058" // IP3
	} else { // [1,1,1] no se sabe cual fue el último o no se ha editado, se elige al azar
		id := rand.Intn(3)
		if id == 0 {
			ChoosenServer = ":50058"//"10.6.40.225:50058" // IP1
		} else if id == 1 {
			ChoosenServer = ":50058"//"10.6.40.227:50058" // IP2
		} else {
			ChoosenServer = ":50058"//"10.6.40.229:50058" // IP3
		}
	}
	return ChoosenServer
}

func (s *Server) GetNumberRebelds(ctx context.Context, message *pb.LeiaRequest) (*pb.LeiaReply,error){
	log.Printf("Leia se ha conectado para saber el número de rebeldes")
	log.Printf("Los parámetros son:\nPlaneta: "+ message.Planet + "\nCiudad: "+message.City)
	ChoosenServer := ChooseServer(int(message.X),int(message.Y),int(message.Z)) 
	log.Printf("Server escogido es: "+ ChoosenServer)
	// Para cuando este funcional el servidor
	//conn, err := grpc.Dial(ChoosenServer, grpc.WithInsecure())
	//ServerService := pb.NewChatServiceClient(conn)
	//r, err := ServerService.GetNumberRebelds(context.Background(), &message)
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//return &pb.LeiaReply{Status:response.Status, Quantity: response.Quantity, X: response.X, Y: response.Y, Z: response.Z, Server: ChoosenServer},nil
	return &pb.LeiaReply{Status:"OK", Quantity: 10, X: 1, Y: 2, Z: 0, Server: ChoosenServer},nil
}

func (s *Server) GetServer(ctx context.Context, message *pb.BrokerRequest) (*pb.BrokerReply,error) {
	log.Printf("Un informante se ha conectado ... \nEntregando dirección de uno de los servidores Fulcrum ...\nDirección enviada")
	return &pb.BrokerReply{IP: ChooseServer(int(message.X),int(message.Y),int(message.Z))}, nil
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
