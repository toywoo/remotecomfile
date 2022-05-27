
package main

import (
	handler "RemoteComfile/handler"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type HttpSrv struct {
	Port      int
	server    *http.Server
	isStarted bool
	mtx       *sync.Mutex
	mux       *http.ServeMux
	handler   *handler.Handler
}

func NewHttpSrv(port int) *HttpSrv {
	return &HttpSrv{
		Port:      port,
		server:    nil,
		isStarted: false,
		mtx:       &sync.Mutex{},
		mux:       http.NewServeMux(),
		handler:   handler.NewHandler(),
	}
}

func (srv *HttpSrv) startServer() error {
	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	if srv.isStarted {
		return errors.New("server is already started")
	}

	handler := handler.NewHandler()

	srv.mux.HandleFunc("/", handler.PathNav)
	srv.mux.Handle("/Client/", http.StripPrefix("/Client", http.FileServer(http.Dir("./Client"))))

	srv.isStarted = true

	addr := fmt.Sprintf(":%v", srv.Port)
	srv.server = &http.Server{Addr: addr, Handler: srv.mux}

	go func() {
		log.Println("Start Server localhost:3002")
		fmt.Print(">> ")
		err := srv.server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Printf("Server closed under request: %v\n", err)
			} else {
				log.Fatalf("Server closed unexpect: %v\n", err)
			}
			srv.isStarted = false
		}

	}()

	return nil
}

func (srv *HttpSrv) shutdownServer() error {
	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	if !srv.isStarted || srv.server == nil {
		return errors.New("server is not started")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stop := make(chan bool)
	log.Println("Close Server...")
	go func() {
		_ = srv.server.Shutdown(ctx)
		stop <- true
	}()

	select {
	case <-ctx.Done():
		log.Printf("Timeout: %v", ctx.Err())
		break
	case <-stop:
		log.Println("Finished")
	}

	log.Println("Closed Server!")

	return nil
}

func interruptExit(srv *HttpSrv) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM,
	)

	for {
		sig := <-signalChan
		if sig != nil {
			log.Println("Kill Interrupt!")
			if srv.isStarted {
				log.Println("Server will Turn off")
				srv.shutdownServer()
			}
			break
		}
	}

	os.Exit(0)
}

func shellScript(srv *HttpSrv) {
	var command string
	fmt.Print(">> ")
	for {
		fmt.Scanf("%s\n", &command)

		command = strings.ToLower(command)

		switch command {
		case "start":
			srv.startServer()

		case "shutdown":
			if srv.isStarted {
				err := srv.shutdownServer()
				if err != nil {
					log.Printf("Server Shutdown failed: %v\n", err)
				} else {
					return
				}
			}

		case "exit":
			return

		case "help":
			fmt.Println("start: start server\n shutdown: if you turn on server, this program will close server or \n exit: turn off this program")

		default:
			fmt.Println("You entered an invalid command. If you want to know what commands are available, type help.")
		}

		fmt.Print(">> ")
	}
}

func main() {

	srv := NewHttpSrv(3002)

	go interruptExit(srv)

	shellScript(srv)
}