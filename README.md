# Distributed Matrix Computation System

This project implements a distributed system for matrix computations using Go's RPC (Remote Procedure Call) framework. The system consists of a client, coordinator, and multiple worker nodes that collaborate to perform matrix operations.

## Project Structure

```
Blockchain_Asg1/
├── client/
│   └── client.go      # Client application for sending requests
├── common/
│   └── common.go      # Shared data structures and interfaces
├── coordinator/
│   └── coordinator.go # Coordinates tasks between client and workers
├── worker/
│   └── worker.go      # Performs the actual matrix computations
└── go.mod             # Go module definition
```

## Components

### Common Package

Defines shared data structures used across all components:
- `MatrixOp`: Represents a matrix operation request
- `Result`: Holds the operation result or error message

### Client

A command-line interface that allows users to:
1. Select matrix operations (Addition, Transpose, Multiplication)
2. Input matrix dimensions and values
3. Send operation requests to the coordinator
4. Display computation results

### Coordinator

Acts as an intermediary between clients and workers:
1. Receives operation requests from clients
2. Manages a pool of worker nodes
3. Assigns tasks to available workers
4. Returns results back to clients
5. Handles graceful shutdown of the system

### Workers

Perform the actual matrix computations:
1. Matrix addition
2. Matrix multiplication
3. Matrix transposition

## How to Run

1. Start the worker nodes:
   ```bash
   go run worker/worker.go
   ```

2. Start the coordinator:
   ```bash
   go run coordinator/coordinator.go
   ```

3. Run the client application:
   ```bash
   go run client/client.go
   ```

4. Follow the on-screen prompts to perform matrix operations.

## Supported Operations

1. **Addition**: Adds two matrices of the same dimensions  
2. **Transpose**: Transposes a matrix (flips rows and columns)  
3. **Multiplication**: Multiplies two matrices (if dimensions are compatible)  
4. **Shutdown**: Gracefully terminates the client  

## System Architecture

This system follows a coordinator-worker architecture:
- The client sends requests to the coordinator  
- The coordinator manages a pool of worker nodes  
- Workers perform computations in parallel  
- Results are returned through the coordinator back to the client  

Each component communicates using Go's RPC framework over TCP connections.

## Error Handling

The system includes error handling for:
- Invalid matrix dimensions  
- Unavailable workers  
- Connection failures  
- Invalid operations  

## Future Improvements

Potential enhancements to consider:
- Dynamic worker registration  
- Load balancing among workers  
- Fault tolerance mechanisms  
- Support for additional matrix operations  
- Performance optimizations for large matrices  
