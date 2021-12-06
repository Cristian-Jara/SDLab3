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
	value string
	PlanetsList []PlanetInfo
)

func main(){
	conn,err := grpc.Dial(fmt.Sprint(LocalIP,Puerto), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s",err)
	}
	serviceClient := pb.NewChatServiceClient(conn)
	for Connected != false{
		log.Printf("Ingresa alguno de los siguientes números para realizar una acción:\n1: Añadir un nuevo registro de ciudad a un planeta (AddCity)\n2: Actualizar el nombre de una ciudad (UpdateName)\n3: Actualizar el valor de un registro (UpdateNumber)\n4: Eliminar un registro (DeleteCity)\n5: Salir")
		fmt.Scanln(&input)
		if input == "5"{
			break
		}
		log.Printf("Ingresa el nombre del planeta:")
		fmt.Scanln(&planet)
		log.Printf("Ingresa el nombre de la ciudad:")
		fmt.Scanln(&city)
		x:= int32(-1)
		y:= int32(-1)
		z:= int32(-1)
		Exist = false
		// Revisar si la planeta y ciudad existe y enviar el reloj si es así
		for idx,_ := range PlanetsList{
			if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
				x = PlanetsList[idx].X
				y = PlanetsList[idx].Y
				z = PlanetsList[idx].Z
				Exist = true
			}
		}
		if input == "1"{
			log.Printf("Ingresa la cantidad de rebeldes:")
			fmt.Scanln(&value)
			if value != ""{
				_, err := strconv.Atoi(value);
				for err != nil {
					log.Printf("Debes ingresar un valor númerico o dejar vacio")
					fmt.Scanln(&value)
					if value == ""{
						value = "0"
						break
					}
					_, err = strconv.Atoi(value);
				}
			}else {
				value = "0"
			}
			message := pb.BrokerRequest{X: x,Y: y,Z: z} 
			//Si se envia -1,-1,-1 quiere decir que aún no se tienen datos de él
			response,err := serviceClient.GetServer(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetServer: %s",err)
			}
			conn,err := grpc.Dial(response.IP, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %s",err)
			}
			ServerClient := pb.NewChatServiceClient(conn)
			message2 := pb.ServerRequest{Planet: planet, City: city, Value: value,X: x, Y:y, Z: z}
			response2, err := ServerClient.AddCity(context.Background(), &message2)
			if response2.Status == "OK"{
				/*for idx,_ := range PlanetsList{
					if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
						Exist = true
						if PlanetsList[idx].X <= response2.X && PlanetsList[idx].Y <= response2.Y && PlanetsList[idx].Z <= response2.Z {
							PlanetsList[idx].X = response2.X
							PlanetsList[idx].Y = response2.Y
							PlanetsList[idx].Z = response2.Z
							PlanetsList[idx].lastserver = response.IP
						} 
						break
					} 
				}*/
				PlanetsList = append(PlanetsList, PlanetInfo{planet, city, response2.X,response2.Y,response2.Z, response.IP})
			}else {
				log.Printf("No se pudo realizar la acción ingresada " + response2.Status)
			}
		}else if input == "2"{
			log.Printf("Ingresa el nuevo nombre de la ciudad:")
			fmt.Scanln(&value)
			message := pb.BrokerRequest{X: x,Y: y,Z: z} 
			//Si se envia -1,-1,-1 quiere decir que aún no se tienen datos de él
			response,err := serviceClient.GetServer(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetServer: %s",err)
			}
			conn,err := grpc.Dial(response.IP, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %s",err)
			}
			ServerClient := pb.NewChatServiceClient(conn)
			message2 := pb.ServerRequest{Planet: planet, City: city, Value: value,X: x, Y:y, Z: z}
			response2, err := ServerClient.UpdateName(context.Background(), &message2)
			if response2.Status == "OK"{
				for idx,_ := range PlanetsList{
					if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
						Exist = true
						if PlanetsList[idx].X <= response2.X && PlanetsList[idx].Y <= response2.Y && PlanetsList[idx].Z <= response2.Z {
							PlanetsList[idx].X = response2.X
							PlanetsList[idx].Y = response2.Y
							PlanetsList[idx].Z = response2.Z
							PlanetsList[idx].lastserver = response.IP
						} 
						break
					} 
				}
				if Exist == false{
					PlanetsList = append(PlanetsList, PlanetInfo{planet, city, response2.X,response2.Y,response2.Z, response.IP})
				}
			}else {
				log.Printf("No se pudo realizar la acción ingresada " + response2.Status)
			}
		}else if input == "3"{
			log.Printf("Ingresa la nueva cantidad de rebeldes:")
			fmt.Scanln(&value)
			_, err := strconv.Atoi(value);
			for err != nil {
				log.Printf("Debes ingresar un valor númerico")
				fmt.Scanln(&value)
				_, err = strconv.Atoi(value);
			}
			message := pb.BrokerRequest{X: x,Y: y,Z: z} 
			//Si se envia -1,-1,-1 quiere decir que aún no se tienen datos de él
			response,err := serviceClient.GetServer(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetServer: %s",err)
			}
			conn,err := grpc.Dial(response.IP, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %s",err)
			}
			ServerClient := pb.NewChatServiceClient(conn)
			message2 := pb.ServerRequest{Planet: planet, City: city, Value: value,X: x, Y:y, Z: z}
			response2, err := ServerClient.UpdateNumber(context.Background(), &message2)
			if response2.Status == "OK"{
				for idx,_ := range PlanetsList{
					if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
						Exist = true
						if PlanetsList[idx].X <= response2.X && PlanetsList[idx].Y <= response2.Y && PlanetsList[idx].Z <= response2.Z {
							PlanetsList[idx].X = response2.X
							PlanetsList[idx].Y = response2.Y
							PlanetsList[idx].Z = response2.Z
							PlanetsList[idx].lastserver = response.IP
						} 
						break
					} 
				}
				if Exist == false{
					PlanetsList = append(PlanetsList, PlanetInfo{planet, city, response2.X,response2.Y,response2.Z, response.IP})
				}
			}else {
				log.Printf("No se pudo realizar la acción ingresada " + response2.Status)
			}
		}else if input == "4"{
			message := pb.BrokerRequest{X: x,Y: y,Z: z} 
			//Si se envia -1,-1,-1 quiere decir que aún no se tienen datos de él
			response,err := serviceClient.GetServer(context.Background(), &message)
			if err != nil{
				log.Fatalf("Error when calling GetServer: %s",err)
			}
			conn,err := grpc.Dial(response.IP, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %s",err)
			}
			value := ""
			ServerClient := pb.NewChatServiceClient(conn)
			message2 := pb.ServerRequest{Planet: planet, City: city, Value: value,X: x, Y:y, Z: z}
			response2, err := ServerClient.DeleteCity(context.Background(), &message2)
			if response2.Status == "OK"{
				for idx,_ := range PlanetsList{
					if planet == PlanetsList[idx].planet && city == PlanetsList[idx].city {
						Exist = true
						if PlanetsList[idx].X <= response2.X && PlanetsList[idx].Y <= response2.Y && PlanetsList[idx].Z <= response2.Z {
							PlanetsList[idx].X = response2.X
							PlanetsList[idx].Y = response2.Y
							PlanetsList[idx].Z = response2.Z
							PlanetsList[idx].lastserver = response.IP
						} 
						break
					} 
				}
				if Exist == false{
					PlanetsList = append(PlanetsList, PlanetInfo{planet, city, response2.X,response2.Y,response2.Z, response.IP})
				}
			}else {
				log.Printf("No se pudo realizar la acción ingresada " + response2.Status)
			}
		}
	}
}