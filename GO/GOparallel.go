//GOcparallel.go

package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	wg        sync.WaitGroup
	num_procs int
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: make parallel [matrizN] [seed]")
	}
	args := os.Args
	N, err1 := strconv.Atoi(args[1])
	if err1 != nil {
		fmt.Println(err1)
	}
	seed, err2 := strconv.Atoi(args[2])
	if err2 != nil {
		fmt.Println(err2)
	}

	// A * X  = B
	B := make([]float64, N)
	X := make([]float64, N)
	A := make([][]float64, N)
	for i := range A {
		A[i] = make([]float64, N)
	}

	num_procs = runtime.NumCPU()

	fmt.Printf("Matriz dimension size: %d.\nSeed: %d.\nNum_Procs: %d.\n", N, seed, num_procs)

	r := rand.New(rand.NewSource(int64(seed)))
	// Initialize A, B and X
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			A[j][i] = r.Float64() / 32768.0
		}
		B[i] = r.Float64() / 32768.0
		X[i] = 0.0
	}

	// Print inputs
	printIn(N, A, B, X)

	//Gauss elimination
	start := time.Now()
	wg.Add(num_procs)
	for i := 0; i < num_procs; i++ {
		go gauss(N, X, B, A, &wg)
	}
	wg.Wait()
	end := time.Now()

	//Print result and time
	printOut(N, X)
	fmt.Printf("Time taken: \n" + fmt.Sprint(end.Sub(start)))

}

func printIn(N int, A [][]float64, B, X []float64) {
	fmt.Printf("\n")
	fmt.Println("--A--")
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("[ %v ]", A[i][j])
		}
		fmt.Printf("\n")
	}
	fmt.Println("--B--")
	for i := 0; i < N; i++ {
		fmt.Printf("[ %v ]", B[i])
	}
	fmt.Printf("\n")
	fmt.Println("--X--")
	for i := 0; i < N; i++ {
		fmt.Printf("[ %v ]", X[i])
	}
	fmt.Printf("\n")

}

func printOut(N int, X []float64) {
	fmt.Printf("\n")
	fmt.Printf("--Answer--\n")
	for i := 0; i < N; i++ {
		fmt.Printf("[ %v ]", X[i])
	}
	fmt.Printf("\n")
}

func gauss(N int, X, B []float64, A [][]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	var multiplier float64
	var norm, col, row int

	for norm = 0; norm < N-1; norm++ {
		for row = norm + 1; row < N; row++ {
			multiplier = A[row][norm] / A[norm][norm]
			for col = norm; col < N; col++ {
				A[row][col] -= A[norm][col] * multiplier
			}
			B[row] -= B[norm] * multiplier
		}
	}

	for row = N - 1; row >= 0; row-- {
		X[row] = B[row]
		for col = N - 1; col > row; col-- {
			X[row] -= A[row][col] * X[col]
		}
		X[row] /= A[row][row]
	}
}
