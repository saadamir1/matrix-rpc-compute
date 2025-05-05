package common

// MatrixOp represents a matrix operation request
type MatrixOp struct {
	Operation string
	MatrixA   [][]int
	MatrixB   [][]int
}

// Result holds the matrix operation result
type Result struct {
	Matrix [][]int
	Error  string
}