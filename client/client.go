package main

import (
	"fmt"
	"log"
	"net/rpc"

	"Blockchain_Asg1/common" // Import the common package
)

// sendRequest sends a matrix operation request to the coordinator
func sendRequest(coordinatorAddr string, op common.MatrixOp) {
	client, err := rpc.Dial("tcp", coordinatorAddr)
	if err != nil {
		log.Fatal("Connection to coordinator failed:", err)
	}
	defer client.Close()

	var result common.Result
	err = client.Call("Coordinator.ProcessRequest", op, &result)
	if err != nil {
		log.Fatal("RPC call failed:", err)
	}

	if result.Error != "" {
		fmt.Println("Error:", result.Error)
	} else {
		fmt.Println("Result:")
		for _, row := range result.Matrix {
			fmt.Println(row)
		}
	}
}

// Prompt and read integer inputs safely
func readInt(prompt string) int {
	var input int
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal("Invalid input:", err)
	}
	return input
}

// Prompt and read 2D matrix inputs
func readMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			fmt.Printf("Enter element for position [%d][%d]: ", i+1, j+1)
			_, err := fmt.Scanln(&matrix[i][j])
			if err != nil {
				log.Fatal("Invalid input:", err)
			}
		}
	}
	return matrix
}

// Check if matrix multiplication is possible
func validateMultiplicationDimensions(cols1, rows2 int) bool {
	return cols1 == rows2
}

func main() {
	coordinatorAddr := "localhost:8080"

	var op common.MatrixOp
	var operation string

	// Loop for continuously displaying the operation menu until "Shutdown" is selected
	for {
		// Provide an interface for the user to choose an operation
		fmt.Println("\nSelect an operation:")
		fmt.Println("1. Addition")
		fmt.Println("2. Transpose")
		fmt.Println("3. Multiplication")
		fmt.Println("4. Shutdown")

		// Prompt for user input
		fmt.Print("Enter the operation number (1/2/3/4): ")
		_, err := fmt.Scanln(&operation)
		if err != nil {
			log.Fatal("Invalid input:", err)
		}

		// Check operation choice and handle input
		switch operation {
		case "1": // Addition
			rows1 := readInt("Enter the number of rows for the first matrix: ")
			cols1 := readInt("Enter the number of columns for the first matrix: ")
			matrixA := readMatrix(rows1, cols1)

			// Automatically set the second matrix dimensions to match the first matrix
			matrixB := make([][]int, rows1)
			for i := 0; i < rows1; i++ {
				matrixB[i] = make([]int, cols1)
			}

			// Ask the user for the elements of the second matrix
			fmt.Println("Enter the elements for the second matrix:")
			matrixB = readMatrix(rows1, cols1)

			op.Operation = "add"
			op.MatrixA = matrixA
			op.MatrixB = matrixB

		case "2": // Transpose
			rows := readInt("Enter the number of rows for the matrix: ")
			cols := readInt("Enter the number of columns for the matrix: ")
			matrixA := readMatrix(rows, cols)

			op.Operation = "transpose"
			op.MatrixA = matrixA

		case "3": // Multiplication
			rows1 := readInt("Enter the number of rows for the first matrix: ")
			cols1 := readInt("Enter the number of columns for the first matrix: ")
			matrixA := readMatrix(rows1, cols1)

			rows2 := readInt("Enter the number of rows for the second matrix: ")
			cols2 := readInt("Enter the number of columns for the second matrix: ")

			// Validate multiplication dimensions
			if !validateMultiplicationDimensions(cols1, rows2) {
				fmt.Println("Error: Matrix dimensions are not compatible for multiplication.")
				continue // Go back to operation menu
			}

			matrixB := readMatrix(rows2, cols2)

			op.Operation = "multiply"
			op.MatrixA = matrixA
			op.MatrixB = matrixB

		case "4": // Shutdown
			fmt.Println("Shutting down client...")
			return // Exit the client gracefully

		default:
			log.Fatal("Invalid operation selected.")
		}

		// Send the operation to the coordinator
		sendRequest(coordinatorAddr, op)
	}
}
