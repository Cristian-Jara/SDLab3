package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
	
)

type server struct {
	pb.UnimplementedChatServiceServer
}


type PlanetData struct{
	planet string
	X int32
	Y int32
	Z int32
}

func crearRegistro(path string, planet string) (error) {
	//Verifica que el archivo existe
	_, err := os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		PlanetsData = append(PlanetsData, PlanetData{planet, 0,0,0})
	}
	return nil
}

func crearLog(path string) (error) {
	//Verifica que el archivo existe
	_, err := os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func escribirArchivo(path string, message string) (error){ //Escribe al final del archivo
	b, errtxt := ioutil.ReadFile(path)

	if errtxt != nil {
		log.Fatal(errtxt)
		return errtxt
	}
	b = append(b, []byte(message)...)
	errtxt = ioutil.WriteFile(path, b, 0644)
	if errtxt != nil {
		log.Fatal(errtxt)
		return errtxt
	}
	return nil
}

func eliminarArchivo(path string){
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error eliminando archivo: %v\n", err)
	}
}

func SumarAlReloj(planet string)(int32,int32,int32){ //Uno para cada servidor en el campo correspondiente
	for idx,_ := range PlanetsData { 
		if PlanetsData[idx].planet == planet {
			PlanetsData[idx].X += 1     //S1
			//PlanetsData[idx].Y += 1     //S2
			//PlanetsData[idx].Z += 1     //S3
			return PlanetsData[idx].X,PlanetsData[idx].Y,PlanetsData[idx].Z
		}
	}
	return -1,-1,-1
}

func (s *server) AddCity(ctx context.Context, in *pb.ServerRequest) (*pb.ServerReply, error) {
	//aqui implementar la escritura del archivo de texto
	// Se debe crear el archivo del planeta si no existe y lo mismo con el log
	// Si existe hay que revisar si la ciudad existe, para ello se puede usar la estructura
	// Agregar cosa a la estructura también, manejar reloj de vector dependiendo del servidor
	var path = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"
	var path2 = "Servers/ServersData/Logs/"+ in.Planet +".txt"
	
	err := crearRegistro(path,in.Planet) 
	err2 := crearLog(path2)
	if err != nil {
		return &pb.ServerReply{ Status: "Error al crear el archivo del registro planetario" },err
	}
	if err2 != nil {
		return &pb.ServerReply{ Status: "Error al crear el archivo log" },err2
	}

	// Para planetregister
	message := in.Planet +" "+ in.City + " " + in.Value +"\n" // Agregar mensaje
	err = escribirArchivo(path, message)
	x, y, z := SumarAlReloj(in.Planet)
	//Agregar 1 al reloj dependiendo del servidor
	// Para log
	message = "AddCity " + in.Planet +" "+ in.City + " " + in.Value +"\n" // Agregar mensaje
	err2 = escribirArchivo(path2, message)

	if err != nil {
		return &pb.ServerReply{ Status: "Error al escribir en el archivo del registro planetario" },err
	}
	if err2 != nil {
		return &pb.ServerReply{ Status: "Error al escribir en el archivo log" },err2
	}

	return &pb.ServerReply{ Status: "OK" , X: x, Y: y, Z: z},nil
}

func (s *server) UpdateName(ctx context.Context, in *pb.ServerRequest) (*pb.ServerReply, error) {
	var path = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"
	var path2 = "Servers/ServersData/Logs/"+ in.Planet +".txt" 
	//Se asume que en este punto existen
	input, err := ioutil.ReadFile(path)
	if err != nil{
		log.Fatal(err)
	}	
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, in.City) {
			splitLine := strings.Split(string(line), " ")
			rebeldes := splitLine[2] //Se saca el número
			lines[i] = in.Planet + " " + in.Value + " " + rebeldes
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
		return &pb.ServerReply{ Status: "Error al escribir en el archivo de registro" },err
	}
	x, y, z := SumarAlReloj(in.Planet)
	message := "UpdateName " + in.Planet +" "+ in.City + " " + in.Value +"\n" // Agregar mensaje
	err2 := escribirArchivo(path2, message)
	if err2 != nil {
		return &pb.ServerReply{ Status: "Error al escribir en el archivo log" },err2
	}
	return &pb.ServerReply{ Status: "OK", X: x, Y: y, Z: z },nil
}

func (s *server) UpdateNumber(ctx context.Context, in *pb.ServerRequest) (*pb.ServerReply, error) {
	var path = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"
	var path2 = "Servers/ServersData/Logs/"+ in.Planet +".txt" 
	//Se asume que en este punto existen
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines { 
		if strings.Contains(line, in.City) {
			lines[i] = in.Planet + " " + in.City + " " + in.Value
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
		return &pb.ServerReply{ Status: "Error al escribir en el archivo de registro" },err
	}

	x, y, z := SumarAlReloj(in.Planet)
	message := "UpdateNumber " + in.Planet +" "+ in.City + " " + in.Value +"\n" // Agregar mensaje
	err2 := escribirArchivo(path2, message)
	if err2 != nil {
		return &pb.ServerReply{ Status: "Error al escribir en el archivo log" },err2
	}
	return &pb.ServerReply{ Status: "OK", X: x, Y: y, Z: z },nil
}

func (s *server) DeleteCity(ctx context.Context, in *pb.ServerRequest) (*pb.ServerReply, error) {
	var path = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"
	var path2 = "Servers/ServersData/Logs/"+ in.Planet +".txt" 
	//Se asume que en este punto existen
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var newlines []string
	for i, line := range lines { 
		if !strings.Contains(line, in.City) {
			newlines = append(newlines, lines[i])
		}
	}
	output := strings.Join(newlines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
		return &pb.ServerReply{ Status: "Error al escribir en el archivo de registro" },err
	}

	x, y, z := SumarAlReloj(in.Planet)
	message := "DeleteCity " + in.Planet +" "+ in.City +"\n" // Agregar mensaje
	err2 := escribirArchivo(path2, message)
	if err2 != nil {
		return &pb.ServerReply{ Status: "Error al escribir en el archivo log" },err2
	}

	return &pb.ServerReply{ Status: "OK", X: x, Y: y, Z: z },nil
}

func (s *server) GetNumberRebelds(ctx context.Context, in *pb.LeiaRequest)(*pb.LeiaReply,error){
	var path = "Servers/ServersData/PlanetRegisters/"+ in.Planet +".txt"
	var x,y,z = int32(-1), int32(-1), int32(-1)
	for idx,_ := range PlanetsData { 
		if PlanetsData[idx].planet == in.Planet {
			x,y,z = PlanetsData[idx].X,PlanetsData[idx].Y,PlanetsData[idx].Z
		}
	} // Buscar si existe
	if x == int32(-1){
		return &pb.LeiaReply{Status: "No se encontró el planeta y ciudad dados"}, nil
	}
	input, err := ioutil.ReadFile(path)
	if err != nil{
		log.Fatal(err)
		return &pb.LeiaReply{Status: "Error al leer el archivo"}, nil
	}
	lines := strings.Split(string(input), "\n")
	var quantity int32
	for _, line := range lines {
		if strings.Contains(line, in.City) {
			splitLine := strings.Split(string(line), " ")
			i, err := strconv.Atoi(splitLine[2])
			if err != nil{
				log.Fatal(err)
				return &pb.LeiaReply{Status: "Error obtener valor númerico de rebeldes"}, nil
			}
			quantity = int32(i) //Se saca el número  
		}
	}
	return &pb.LeiaReply{Status:"OK", Quantity: quantity, X: x, Y: y, Z: z},nil
}

/*func propagacion(){
	//Aquí iría el código de propagación, si tan solo tuviera uno
	//Nueva func para propagar info
	//Función que revise la data y con alguna elección arbitraria 
	// decida con que data quedarse
}*/

var PlanetsData []PlanetData

func main() {
	/*go func() {
		for{
			time.Sleep(120 * time.Second)
			propagacion()
		}
	}()*/
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
