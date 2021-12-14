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
	"time"
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
	for _, line := range lines { 
		if !strings.Contains(line, in.City) {
			newlines = append(newlines, string(line))
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

func EmptyAll(){
	for _, value := range PlanetsData{
		path := "Servers/ServersData/Logs/"+ value.planet +".txt"
		err := ioutil.WriteFile(path, []byte(""), 0644)
		if err != nil {
			log.Fatal(err)
		}
		path2 := "Servers/ServersData/PlanetRegisters/"+ value.planet +".txt"
		err2 := ioutil.WriteFile(path2, []byte(""), 0644)
		if err2 != nil {
			log.Fatal(err)
		}
	}
}

func InfoToMessage() (*pb.PropagationReply){
	message := &pb.PropagationReply{Status: "OK", Planetsdata: []*pb.PlanetsData{}}
	for _, value := range PlanetsData{
		//Se revisa toda la data de los planetas
		path := "Servers/ServersData/PlanetRegisters/"+ value.planet +".txt"
		path2 := "Servers/ServersData/Logs/"+ value.planet +".txt"
		input, err := ioutil.ReadFile(path) //Archivo del planeta
		if err != nil {
			log.Fatalln(err)
		}
		input2, err := ioutil.ReadFile(path2) //Logs
		if err != nil {
			log.Fatalln(err)
		}
		planetData := &pb.PlanetsData{Planet: value.planet, X: value.X, Y: value.Y, Z: value.Z, Data: []*pb.Data{}, Logs: string(input2)}
		lines := strings.Split(string(input), "\n")
		for _, line := range lines {
			if line != "" {
				splitLine := strings.Split(string(line), " ")
				data := &pb.Data{City: splitLine[1], Value: splitLine[2]}
				planetData.Data = append(planetData.Data, data) //Agrega toda la data del planeta
			}
		}
		message.Planetsdata = append(message.Planetsdata, planetData) 
		//Agrega todos los datos de todos los planetas

	}
	return message
}



type DeletedCities struct{
	planet string
	cities []string
}

type CityValue struct{
	city string
	value string
}

type NameModifiedCities struct{
	planet string
	cities []CityValue
}

func Max(x,y int32) (int32){
	if x < y {
		return y
	}
	return x
}

func Fill(data []*pb.PlanetsData){
	for _, p := range data{
		// por cada registro planetario crea version nueva
		var path = "Servers/ServersData/PlanetRegisters/"+ p.Planet +".txt"
		var path2 = "Servers/ServersData/Logs/"+ p.Planet +".txt" 
		//crear registro nuevo
		crearRegistro(path,p.Planet) // Crea el archivo y añade Planetdata correspondiente
		crearLog(path2) //Crea el archivo log si no existe
		// se actualiza arreglo de planetas
		// revisa en que posición está y lo atualiza
		for idx, value := range PlanetsData {
			if p.Planet == value.planet {
				PlanetsData[idx].X = p.X
				PlanetsData[idx].Y = p.Y
				PlanetsData[idx].Z = p.Z
			}
		}
		// se escriben los datos de ciudades
		for _, c := range p.Data{
			text := p.Planet +  c.City + c.Value + "\n"
			escribirArchivo(path,text)
		}
	}
	return
}


func Merge(Sv2Data, Sv3Data []*pb.PlanetsData){
	Sv1Data := InfoToMessage().Planetsdata // Formatea la info para que este con el mismo formato que el los demás servidores
	//Asumir Sv1Data está bien y desde las acciones de los demás aplicar modificaciones
	flag := false
	var ANMC []NameModifiedCities //Conflictos de nombre
	var ADC []DeletedCities //Conflictos de eliminación
	for _, value := range Sv2Data {
		flag = false
		for idx, value2 := range Sv1Data {
			if value.Planet == value2.Planet {
				flag = true
				if value.X != value2.X || value.Y != value2.Y || value.Z != value2.Z {
					//Se debe revisar logs para posibles conflictos
					lines := strings.Split(string(value.Logs), "\n")
					lines2 := strings.Split(string(value2.Logs), "\n")
					var DeletedCitiesArray []string 
					var NameModifiedCitiesArray []CityValue 
					for _, line := range lines {
						if strings.Contains(line, "DeleteCity") {
							splitLine := strings.Split(string(line), " ")
							DeletedCitiesArray = append(DeletedCitiesArray, splitLine[2])
						}else if strings.Contains(line, "UpdateName") {
							splitLine := strings.Split(string(line), " ")
							var aux CityValue
							aux.city = splitLine[2]
							aux.value = splitLine[3]
							NameModifiedCitiesArray = append(NameModifiedCitiesArray, aux)
						}
					}
					for _, line := range lines2 {
						if strings.Contains(line, "DeleteCity") {
							splitLine := strings.Split(string(line), " ")
							DeletedCitiesArray = append(DeletedCitiesArray, splitLine[2])
						}else if strings.Contains(line, "UpdateName") {
							splitLine := strings.Split(string(line), " ")
							var aux CityValue
							aux.city = splitLine[2]
							aux.value = splitLine[3]
							NameModifiedCitiesArray = append(NameModifiedCitiesArray, aux)
						}
					}
					//Se guarda todos los posibles conflictos problematicos, se resolverán después
					if len(DeletedCitiesArray)>0 {
						var aux DeletedCities
						aux.planet = value.Planet
						aux.cities = DeletedCitiesArray
						ADC = append(ADC, aux)
					}
					if len(NameModifiedCitiesArray)>0 {
						var aux NameModifiedCities
						aux.planet = value.Planet
						aux.cities = NameModifiedCitiesArray
						ANMC = append(ANMC, aux)
					}
					flag2 := false
					for _, value3 := range value.Data {
						flag2 = false
						for idx2, value4 := range value2.Data {
							flag3 := false
							for _, value5 := range NameModifiedCitiesArray{
								if (value5.city == value3.City && value5.value == value4.City) || (value5.value == value3.City && value5.city == value4.City) {
									flag3 = true //Una de las dos fue modificada
								}
							}
							if value3.City == value4.City || flag3 { 
								flag2 = true
								i, err := strconv.Atoi(value3.Value)
								if err != nil{
									log.Fatal(err)
								}
								i2, err2 := strconv.Atoi(value4.Value)
								if err2 != nil{
									log.Fatal(err2)
								}
								if i > i2 {
									//Si el valor del sv2 es mayor me quedo con él
									Sv1Data[idx].Data[idx2].Value = value3.Value
									//Sin importar posibles conflictos me quedo con el dato mayor para cada dato
								} 
								break
							}
						}
						if !flag2 {
							//Agregar el elemento
							Sv1Data[idx].Data = append(Sv1Data[idx].Data, value3)
						}
					}
					x := Max(value.X, value2.X)
					y := Max(value.Y, value2.Y)
					z := Max(value.Z, value2.Z)
					Sv1Data[idx].X = x
					Sv1Data[idx].Y = y
					Sv1Data[idx].Z = z
				} //Si son distintos se hace algo, si no es un dato que no se ha modificado
			}
		}
		if !flag {
			//Se agrega ya que no está
			Sv1Data = append(Sv1Data, value)
		}
	}
	for _, value := range Sv3Data {
		flag = false
		for idx, value2 := range Sv1Data {
			if value.Planet == value2.Planet {
				flag = true
				if value.X != value2.X || value.Y != value2.Y || value.Z != value2.Z {
					//Se debe revisar logs para posibles conflictos
					lines := strings.Split(string(value.Logs), "\n")
					lines2 := strings.Split(string(value2.Logs), "\n")
					var DeletedCitiesArray []string 
					var NameModifiedCitiesArray []CityValue 
					for _, line := range lines {
						if strings.Contains(line, "DeleteCity") {
							splitLine := strings.Split(string(line), " ")
							DeletedCitiesArray = append(DeletedCitiesArray, splitLine[2])
						}else if strings.Contains(line, "UpdateName") {
							splitLine := strings.Split(string(line), " ")
							var aux CityValue
							aux.city = splitLine[2]
							aux.value = splitLine[3]
							NameModifiedCitiesArray = append(NameModifiedCitiesArray, aux)
						}
					}
					for _, line := range lines2 {
						if strings.Contains(line, "DeleteCity") {
							splitLine := strings.Split(string(line), " ")
							DeletedCitiesArray = append(DeletedCitiesArray, splitLine[2])
						}else if strings.Contains(line, "UpdateName") {
							splitLine := strings.Split(string(line), " ")
							var aux CityValue
							aux.city = splitLine[2]
							aux.value = splitLine[3]
							NameModifiedCitiesArray = append(NameModifiedCitiesArray, aux)
						}
					}
					//Se guarda todos los posibles conflictos problematicos, se resolverán después
					if len(DeletedCitiesArray)>0 {
						var aux DeletedCities
						aux.planet = value.Planet
						aux.cities = DeletedCitiesArray
						ADC = append(ADC, aux)
					}
					if len(NameModifiedCitiesArray)>0 {
						var aux NameModifiedCities
						aux.planet = value.Planet
						aux.cities = NameModifiedCitiesArray
						ANMC = append(ANMC, aux)
					}
					flag2 := false
					for _, value3 := range value.Data {
						flag2 = false
						for idx2, value4 := range value2.Data {
							flag3 := false
							for _, value5 := range NameModifiedCitiesArray{
								if (value5.city == value3.City && value5.value == value4.City) || (value5.value == value3.City && value5.city == value4.City) {
									flag3 = true //Una de las dos fue modificada
								}
							}
							if value3.City == value4.City || flag3 { 
								flag2 = true
								i, err := strconv.Atoi(value3.Value)
								if err != nil{
									log.Fatal(err)
								}
								i2, err2 := strconv.Atoi(value4.Value)
								if err2 != nil{
									log.Fatal(err2)
								}
								if i > i2 {
									//Si el valor del sv3 es mayor me quedo con él
									Sv1Data[idx].Data[idx2].Value = value3.Value
									//Sin importar posibles conflictos me quedo con el dato mayor para cada dato
								} 
								break
							}
						}
						if !flag2 {
							//Agregar el elemento
							Sv1Data[idx].Data = append(Sv1Data[idx].Data, value3)
						}
					}
					x := Max(value.X, value2.X)
					y := Max(value.Y, value2.Y)
					z := Max(value.Z, value2.Z)
					Sv1Data[idx].X = x
					Sv1Data[idx].Y = y
					Sv1Data[idx].Z = z
				} //Si son distintos se hace algo, si no es un dato que no se ha modificado
			}
		}
		if !flag {
			//Se agrega ya que no está
			Sv1Data = append(Sv1Data, value)
		}
	}
	
	for idx,value := range Sv1Data {
		// <- aqui definir el data por añadir de cada planeta
		for _,aux1 := range ADC{
			if value.Planet == aux1.planet {
				for idx2,dato := range value.Data{
					flag = true
					for _,city := range aux1.cities{
						if dato.City == city {
							flag = false
							Sv1Data[idx].Data = append(Sv1Data[idx].Data[:idx2], Sv1Data[idx].Data[idx2+1:]...)
							//ciudad encontrada para eliminar
							// Sv1Data[idx].Data[idx2] <- se debe eliminar o NO añadir
						}
					}
					if flag { //Vamo a buscarlo
						ciudad := ""
						for _,aux2 := range ANMC {
							if value.Planet == aux2.planet{
								for _,place := range aux2.cities{
									if place.city == dato.City { //busca posible nombre por actualizar de ciudad por eliminar 
										ciudad = place.value
										break
									} else if place.value == dato.City { //busca posible nombre viejo por eliminar
										ciudad = place.city
										break
									}
								}
							}
						}
						for _,city := range aux1.cities{
							if city == ciudad {
								Sv1Data[idx].Data = append(Sv1Data[idx].Data[:idx2], Sv1Data[idx].Data[idx2+1:]...)
								//ciudad encontrada para eliminar
								// Sv1Data[idx].Data[idx2] <- se debe eliminar, o NO añadir
							}
						}
					}
				}
			}			
		}
		for _,aux2 := range ANMC{ 
			if value.Planet == aux2.planet {
				for _,name := range aux2.cities{ //Recorre par (nombre viejo,nombre actual)
					for idx3,place := range value.Data { //Busca coincidencias para actualizar
						if name.city == place.City{
							Sv1Data[idx].Data[idx3].Value=name.value // <- hay que actualizar el nombre de la ciudad
						}
					}
				}
			}
		}
	}
	EmptyAll() //Vacía todo
	Fill(Sv1Data) //<----- copy paste de lo que hace el nacho en sv2    // ecole cuá 
	return
}

var PlanetsData []PlanetData

func main() {
	go func() {
		Server2 := "10.6.40.228:50058" // IP2
		Server3 := "10.6.40.229:50058" // IP3
		conn, err := grpc.Dial(Server2, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		conn2, err := grpc.Dial(Server3, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		Server2Service := pb.NewChatServiceClient(conn)
		Server3Service := pb.NewChatServiceClient(conn2)
		for{
			time.Sleep(120 * time.Second) //Se esperan dos minutos
			message := pb.Propagation{Status: "Oe, ya pasaron 2 min toca propagar"}
			response,err := Server2Service.PropagationRequest(context.Background(), &message)
			for err != nil{
				log.Fatalf("Error when calling PropagationRequest of server 2: %s",err)
				response,err = Server2Service.PropagationRequest(context.Background(), &message)
			}
			response2,err := Server3Service.PropagationRequest(context.Background(), &message)
			for err != nil{
				log.Fatalf("Error when calling PropagationRequest of server 3: %s",err)
				response2,err = Server3Service.PropagationRequest(context.Background(), &message)
			}

			Merge(response.Planetsdata, response2.Planetsdata) //Realiza el merge y guarda la data en los archivos y la estructura
			message2 := InfoToMessage() //Pasa la información guardada a mensaje
			response3,err := Server2Service.EventualConsistency(context.Background(), message2)
			for err != nil && response3.Status != "OK"{
				log.Fatalf("Error when calling EventualConsistency of server 2: %s",err)
				response3,err = Server2Service.EventualConsistency(context.Background(), message2)
			}
			response4,err := Server3Service.EventualConsistency(context.Background(), message2)
			for err != nil && response4.Status != "OK"{
				log.Fatalf("Error when calling EventualConsistency of server 3: %s",err)
				response4,err = Server3Service.EventualConsistency(context.Background(), message2)
			}
		}
	}()
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
