package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/lost-mountain/isard/broker"
	"github.com/lost-mountain/isard/configuration"
	"github.com/lost-mountain/isard/rpc"
	"github.com/lost-mountain/isard/rpc/api"
	"github.com/lost-mountain/isard/storage"
)

func main() {
	configPath := flag.String("config", "", "path to the configuration file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("the configuration file is missing, use the -config flag to set it")
		os.Exit(1)
	}

	config, err := configuration.Load(*configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		queue  broker.Broker
		bucket storage.Bucket
	)

	if config.GC != nil {
		b, err := broker.NewPubSubBroker(config.GC.Project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		queue = b

		s, err := storage.NewDatastore(config.GC.Project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bucket = s
	} else {
		queue = broker.NewChannelBroker()

		f, err := ioutil.TempFile("", "isard-")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		s, err := storage.NewBoltBucket(f.Name())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bucket = s
	}

	var server *grpc.Server
	if config.TLS != nil {
		creds, err := credentials.NewServerTLSFromFile(config.TLS.CertFile, config.TLS.KeyFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		server = grpc.NewServer()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.TCP.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proc := broker.NewDomainProcessor(bucket, queue, config.Domains)
	queue.Subscribe(proc)

	api := api.NewAPI(bucket, queue, config)
	rpc.RegisterAPIServer(server, api)

	if err := server.Serve(lis); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
