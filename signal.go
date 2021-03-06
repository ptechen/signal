package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var Signal = make(chan int, 0)
var ExitSignal int32

func Add()  {
	atomic.AddInt32(&ExitSignal, 1)
}

func Sub()  {
	atomic.AddInt32(&ExitSignal, -1)
}

func main() {

	go func() {
		Add()
		defer Sub()
	forTag:
		for {
			select {
			case <-Signal:
				log.Println("End")
				break forTag
			}
		}
	}()
	TestChan := make(chan int, 100000)
	for i :=0; i < 100000; i ++ {
		TestChan <- i
	}
	time.Sleep(time.Second)
	log.Println("ok")
	go func() {
		Add()
		defer Sub()
	forTag:
		for {
			select {
			case <- TestChan:
			case <-Signal:
				if len(TestChan) == 0 {
					break forTag
				}
			}
		}
	}()

	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Printf("server get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			close(Signal)
			for ExitSignal > 0 {
				time.Sleep(time.Second)
			}
			log.Println("server quit !!!")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
