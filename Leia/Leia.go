package main

import (
	"log"
	"fmt"
	"context"
	"strconv"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
)

type PlanetInfo struct{
	planet string
	city string
	rquantity int32
	X int32
	Y int32
	Z int32
	lastserver string
}


const (
	LocalIP = ""
	BrokerIP = "10.6.40.227"
	Puerto = ":50052"
)

var (
	Connected = true
	Exist = false
	input string
	planet string
	city string
	PlanetsList []PlanetInfo
)

func main(){
	conn,err := grpc.Dial(fmt.Sprint(BrokerIP,Puerto), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s",err)
	}


	serviceClient := pb.NewChatServiceClient(conn)
	for Connected != false{
		log.Printf("Presiona 1 para obtener el número de rebeldes y 2 para salir")
		fmt.Scanln(&input)
		if input == "2"{
			break
		}else if input == "1"{
			log.Printf("Ingresa el planeta a buscar")
			fmt.Scanln(&planet)
			log.Printf("Ingresa la ciudad a buscar")
			fmt.Scanln(&city)
			x:= int32(-1)
			y:= int32(-1)
			z:= int32(-1)
			// Revisar si la planeta y ciudad existe en los ya revisados y enviar el reloj si es así
			for idx,_ := range PlanetsList{
				if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
					x = PlanetsList[idx].X
					y = PlanetsList[idx].Y
					z = PlanetsList[idx].Z
				}
			}
			message := pb.LeiaRequest{Planet: planet,City: city,X: x,Y: y,Z: z} 
			//Si se envia -1,-1,-1 quiere decir que aún no se tienen datos de él
			response,err := serviceClient.GetNumberRebelds(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetNumberRebelds: %s",err)
			}
			if response.Status == "OK"{
				for idx,_ := range PlanetsList{
					if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
						Exist = true
						if PlanetsList[idx].X <= response.X && PlanetsList[idx].Y <= response.Y && PlanetsList[idx].Z <= response.Z {
							PlanetsList[idx].rquantity = response.Quantity
							PlanetsList[idx].X = response.X
							PlanetsList[idx].Y = response.Y
							PlanetsList[idx].Z = response.Z
							PlanetsList[idx].lastserver = response.Server
						} 
						// Se asume que x,y,z vendrán bien en el sentido que 
						// aumentará alguno de los vectores al menos para la actualización
						// Tipo 0,0,0 ->1,0,0 -> 1,1,0 -> 1,2,0
						// Se ignoran posibles errores como 1,0,0 -> 0,1,0 
						// en el peor de los casos simplemente asume que la nueva opción 
						// tiene algún problema así que se queda con la opción anterior
						// ya que no es posible distinguir cual opción es la más nueva 
						log.Printf("La cantidad de rebeldes que contiene el planeta y ciudad dado es de: " + strconv.Itoa(int(PlanetsList[idx].rquantity)))
						break
					} 
				}
				if Exist == false {
					PlanetsList = append(PlanetsList, PlanetInfo{planet, city, response.Quantity, response.X,response.Y,response.Z, response.Server})
					log.Printf("La cantidad de rebeldes que contiene el planeta y ciudad dado es de: " + strconv.Itoa(int(response.Quantity)))
				}
				Exist = false
			}else {
				log.Printf("Los valores ingresados no fueron encontrados")
			}
		}
	}
}