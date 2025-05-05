package main

import (
	"log"
	"net"
	"net/rpc"
	"os"

	"Blockchain_Asg1/common" // Import the common package
)

// Worker handles matrix computations
type Worker struct{}

// Compute performs the requested matrix operation
func (w *Worker) Compute(op common.MatrixOp, result *common.Result) error {
	switch op.Operation {
	case "add":
		result.Matrix = addMatrices(op.MatrixA, op.MatrixB)
	case "multiply":
		result.Matrix = multiplyMatrices(op.MatrixA, op.MatrixB)
	case "transpose":
		result.Matrix = transposeMatrix(op.MatrixA)
	default:
		result.Error = "Invalid operation"
	}
	return nil
}

func addMatrices(a, b [][]int) [][]int {
	rows, cols := len(a), len(a[0])
	result := make([][]int, rows)
	for i := range result {
		result[i] = make([]int, cols)
		for j := range result[i] {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result
}

func multiplyMatrices(a, b [][]int) [][]int {
	rows, cols, common := len(a), len(b[0]), len(a[0])
	result := make([][]int, rows)
	for i := range result {
		result[i] = make([]int, cols)
		for j := range result[i] {
			sum := 0
			for k := 0; k < common; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

func transposeMatrix(a [][]int) [][]int {
	rows, cols := len(a), len(a[0])
	result := make([][]int, cols)
	for i := range result {
		result[i] = make([]int, rows)
		for j := range result[i] {
			result[i][j] = a[j][i]
		}
	}
	return result
}

// Start the worker and listen for shutdown signals (optional)
func StartWorker(address string, shutdown chan struct{}) {
	worker := new(Worker)
	rpc.Register(worker)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Worker failed to start:", err)
	}
	defer listener.Close()

	log.Println("Worker is listening on", address)

	// Listen for shutdown signal
	go func() {
		<-shutdown // Wait for shutdown signal
		log.Println("Worker shutting down.")
		listener.Close() // Gracefully shut down the worker listener
		os.Exit(0)       // Exit worker
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
	// Create a shutdown channel to stop workers when required.
	shutdown := make(chan struct{})

	go StartWorker("localhost:9001", shutdown) // Worker 1
	go StartWorker("localhost:9002", shutdown) // Worker 2
	StartWorker("localhost:9003", shutdown)    // Worker 3 (last one runs in the foreground)
}
