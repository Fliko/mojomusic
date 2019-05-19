package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "github.com/Fliko/mojoMusic/mojoroutes"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.RouteGuide_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("vim-go")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln("fuck it aint listnin")
	}

	s := grpc.NewServer()
	greetpb.RegisterRouteGuideServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Fuck it aint servin")
	}
}
