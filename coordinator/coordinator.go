package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"

	"Blockchain_Asg1/common" // Corrected import path
)

type Worker struct {
	Address string
	Busy    bool
}

type Coordinator struct {
	mu        sync.Mutex
	workers   []Worker
	queue     []common.MatrixOp
	shutdown  chan struct{}
	workerWg  sync.WaitGroup
}

func (c *Coordinator) ProcessRequest(op common.MatrixOp, result *common.Result) error {
	c.mu.Lock()
	c.queue = append(c.queue, op)
	c.mu.Unlock()

	return c.assignTask(&op, result)
}

func (c *Coordinator) assignTask(op *common.MatrixOp, result *common.Result) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var worker *Worker
	for i := range c.workers {
		if !c.workers[i].Busy {
			worker = &c.workers[i]
			break
		}
	}

	if worker == nil {
		result.Error = "No available workers"
		return nil
	}

	worker.Busy = true
	client, err := rpc.Dial("tcp", worker.Address)
	if err != nil {
		worker.Busy = false
		result.Error = "Worker connection failed"
		return nil
	}
	defer client.Close()

	err = client.Call("Worker.Compute", *op, result)
	worker.Busy = false
	if err != nil {
		result.Error = "Worker computation failed"
		return nil
	}

	return nil
}

// Shutdown the coordinator gracefully
func (c *Coordinator) Shutdown(request string, reply *string) error {
	c.shutdown <- struct{}{}  // Trigger shutdown
	*reply = "Coordinator shutting down."
	return nil
}

func StartCoordinator(port string) {
	coordinator := &Coordinator{
		workers: []Worker{
			{Address: "localhost:9001", Busy: false},
			{Address: "localhost:9002", Busy: false},
			{Address: "localhost:9003", Busy: false},
		},
		shutdown: make(chan struct{}),
	}

	rpc.Register(coordinator)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting coordinator:", err)
	}
	defer listener.Close()

	log.Println("Coordinator is running on port", port)

	// Handle shutdown signal
	go func() {
		<-coordinator.shutdown // Wait for shutdown signal
		for i := range coordinator.workers {
			client, err := rpc.Dial("tcp", coordinator.workers[i].Address)
			if err != nil {
				log.Println("Error connecting to worker:", err)
				continue
			}
			client.Close() // Close the worker connections
		}
		log.Println("Coordinator shutdown complete.")
		listener.Close() // Shutdown listener
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	log.Println("Coordinator starting...")
	StartCoordinator("8080")
}
