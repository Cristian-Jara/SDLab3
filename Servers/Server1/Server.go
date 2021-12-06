package main

import (
	//"context"
	"fmt"
	//"io/ioutil"
	"log"
	"net"
	"os"
	//"strconv"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
	
)

type server struct {
	pb.UnimplementedChatServiceServer
}

type Data struct{
	city string
	value string
}

type PlanetData struct{
	planet string
	data []Data
	RegisterPath string
	LogPath string
	X int32
	Y int32
	Z int32
}

/*func AppendData(planet string, city string, value string){
	data = Data{city,value}
	for _,info := range PlanetData { 
		if info.planet == planet { //Se podría revisar si existe o no, o se hace antes
			info.data = append(info.data, data)
		}
	}
}

func UpdtName(planet string, city string, value string){
	for idx,_ := range PlanetData { 
		if PlanetData[idx].planet == planet { //Se podría revisar si existe o no, o se hace antes
			for i,_ := PlanetData[idx].data {
				if PlanetData[idx].data[i].city == city {
					PlanetData[idx].data[i].city = value
					break
				}
			}
			break
		}
	}
}

func UpdtNumber(planet string, city string, value string){
	for idx,_ := range PlanetData { 
		if PlanetData[idx].planet == planet { //Se podría revisar si existe o no, o se hace antes
			for i,_ := PlanetData[idx].data {
				if PlanetData[idx].data[i].city == city {
					PlanetData[idx].data[i].value = value
					break
				}
			}
			break
		}
	}
}

func ReadData(player string)([]string){
	for _,info := range PlayersData {
		if info.player == player {
			return info.paths
		}
	}
	return nil
}

func (s *server) AddCity(ctx context.Context, in *pb.ServerRequest) (*pb.ServerReply, error) {
	//aqui implementar la escritura del archivo de texto
	// Se debe crear el archivo del planeta si no existe y lo mismo con el log
	// Si existe hay que revisar si la ciudad existe, para ello se puede usar la estructura
	// Agregar cosa a la estructura también, manejar reloj de vector dependiendo del servidor
	var path = "Servers/ServersData/Logs/"+ in.Planet +".txt"
	var path2 = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"

	//AppendData(in.Player,path) // Agregar la info a PlanetData
	//Verifica que el archivo existe
	var _, err = os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return &pb.ServerReply{ Status: "Error al crear el archivo" },err
		}
		defer file.Close()
	}

	// añadir al texto
	b, errtxt := ioutil.ReadFile(path)

	if errtxt != nil {
		log.Fatal(errtxt)
	}
	message := "\n"
	b = append(b, []byte(message)...)
	errtxt = ioutil.WriteFile(path, b, 0644)
	if errtxt != nil {
		log.Fatal(errtxt)
		return &pb.ServerReply{ Status: "Error al escribir en el archivo" },err
	}
	return &pb.ServerReply{ Status: "OK" },nil
}

func (s *server) GetPlayerInfo(ctx context.Context, in *pb.PlayerInfo) (*pb.PlayerInfo, error) {
	paths := ReadData(in.Message)
	if paths == nil{
		return &pb.PlayerInfo{Message: ""}, nil
	}
	message := ""
	for _, path := range paths {
		//Leer el archivo y chantar todo
		b, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		message += string(b)
	}
	return &pb.PlayerInfo{Message: message}, nil
}*/
var PlanetsData []PlanetData

func main() {
	var path = "Servers/ServersData"
	var _, err = os.Stat(path)
	if !os.IsNotExist(err){
		os.RemoveAll(path) //Si existe se borra para no guardar datos anteriores
	}
	if _,err := os.Stat(path); os.IsNotExist(err){
		err = os.Mkdir(path, 0755)
		if err != nil{
			log.Fatalf("Failed to create the directory: %v",err)
		}
	}
	var path2 = path + "/PlanetRegisters" //Carpeta para guardar los planetas
	if _,err := os.Stat(path2); os.IsNotExist(err){
		err = os.Mkdir(path2, 0755)
		if err != nil{
			log.Fatalf("Failed to create the directory: %v",err)
		}
	}

	var path3 = path + "/Logs" //Carpeta para guardar los Logs
	if _,err := os.Stat(path3); os.IsNotExist(err){
		err = os.Mkdir(path3, 0755)
		if err != nil{
			log.Fatalf("Failed to create the directory: %v",err)
		}
	}

	listner, err := net.Listen("tcp", ":50058")

	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	FulcrumService := grpc.NewServer()
	pb.RegisterChatServiceServer(FulcrumService, &server{})
	fmt.Println("Servidor disponible en el puerto 50058")
	if err = FulcrumService.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
