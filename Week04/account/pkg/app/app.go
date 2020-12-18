package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr             string `yaml:"addr"`
	ReflectionEnable bool   `yaml:"reflectionEnable"`
}

func ReadConfig(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

type AccountApp struct {
	Conf         *Config
	ConfRaw      []byte
	GrpcServer   *grpc.Server
	SignalChan   chan os.Signal
	ServeErrChan chan error
}

func NewAccountApp(confRaw []byte) (*AccountApp, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(confRaw, cfg)
	if err != nil {
		return nil, err
	}
	return &AccountApp{
		Conf:       cfg,
		ConfRaw:    confRaw,
		GrpcServer: grpc.NewServer(),
	}, nil
}

func (app *AccountApp) Start() error {
	ctx := context.Background()
	group, _ := errgroup.WithContext(ctx)

	if app.Conf.ReflectionEnable {
		reflection.Register(app.GrpcServer)
	}

	app.SignalChan = make(chan os.Signal, 10)
	app.ServeErrChan = make(chan error, 1)
	group.Go(func() error {
		return app.listenStopSignal()
	})
	group.Go(func() error {
		err := app.serve()
		select {
		default:
			app.ServeErrChan <- err
		case <-app.ServeErrChan:
		}
		return err
	})

	err := group.Wait()
	if err != nil {
		return err
	}
	return nil
}
func (app *AccountApp) serve() error {
	lis, err := net.Listen("tcp", app.Conf.Addr)
	if err != nil {
		return err
	}
	fmt.Printf("service serve at %s\n", app.Conf.Addr)
	return app.GrpcServer.Serve(lis)
}
func (app *AccountApp) listenStopSignal() error {
	signal.Notify(app.SignalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-app.SignalChan:
		fmt.Println("receive close signal!")
	case err := <-app.ServeErrChan:
		fmt.Printf("receive server close! %+v\n", err)
	}
	signal.Stop(app.SignalChan)
	close(app.ServeErrChan)
	app.GrpcServer.GracefulStop()
	return nil
}
