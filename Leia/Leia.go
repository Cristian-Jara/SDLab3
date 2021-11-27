package main

import (
	"log"
	"fmt"
	"context"
	"strconv"
	"google.golang.org/grpc"
	pb "github.com/Cristian-Jara/SDLab3.git/proto"
)

const (
	LocalIP = ""
	BrokerIP = "10.6.40.227"
	Puerto = ":50052"
)

var (
	Connected = true
	input string
	planet string
	city string
)

func main(){
	conn,err := grpc.Dial(fmt.Sprint(LocalIP,Puerto), grpc.WithInsecure())
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
			message := pb.LeiaRequest{Planet: planet,City: city}
			response,err := serviceClient.GetNumberRebelds(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetNumberRebelds: %s",err)
			}
			if response.Status == "OK"{
				log.Printf("La cantidad de rebeldes que contiene el planeta y ciudad dado es de: " + strconv.Itoa(int(response.Quantity)))
				//Agregarse el reloj y el dato, también revisar el dato con la memoria
			}else {
				log.Printf("Los valores ingresados no fueron encontrados")
			}
		}
	}
}