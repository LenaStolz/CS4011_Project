// 18224296 Thomas Yong
// 17210577 Lena Stolz


// ############################ Matrix Multiplication using Go Channels ######################################################
package main

import (
	"math/rand"
	"sync"
	"time"
	"fmt"
)

type matrix struct {
	rowNo, colNo int
}

// This constant declares the size of the matrixes that will be formed at a later stage
const size = 120

func main() {
	// First we create a constant that holds the number of goroutines we want to use
	const noOfRoutines = 3
	// Next we create a buffered channel of type matrix (perfomance reasons)
	pairs := make(chan matrix, 10)
	// Here we create a variable of type waitGroup (to insure program is fully done at later stage)
	var wait sync.WaitGroup
	// Now we create a 2D array of size "size" to store the results of our calculation in
	var result [size][size]int


	// Here we create two matrixes of size "size" each and fill them with random values
	var a [size][size]int
	var b [size][size]int
	for i := 0; i < size; i++ {
		for j:= 0; j < size; j++ {
			a[i][j] = rand.Intn(10)
			b[i][j] = rand.Intn(10)
		}
	}

	// Here we initialize the 'start' variable with time.Now() to compute
	// the elapsed time at a later stage (this was done for comparison purposes
	// and to see if own logic was correct)
	var start = time.Now()

	// Next we set the number of goroutines we want to wait for by using the previously 
	wait.Add(noOfRoutines)

	// Here we start the goroutines and calculate the results of the matrix multiplication
	for i := 0; i < noOfRoutines; i++ {
        go MatrixMulti(pairs, &a, &b, &result, &wait)
    }
	
	// In the next step we send the results of the calculation to 'pairs'
	for i := 0; i < size; i++ {
        for j := 0; j < size; j++ {
            pairs <- matrix{rowNo: i, colNo: j}
        }
    }

	// At this stage we wait for all the goroutines to finish before going to the
	// next step
	close(pairs)
	wait.Wait()

	// Here we ask for the time that has passed since the calculation was done
	// We then print the result for comparison purposes
	passed := time.Since(start)
	fmt.Println("Time of the calculation is ", passed)

	// This loop prints out the result matrix
	// CAUTION : it is unproductive to use for big matrixes 
	// (e.g. size < 20 )	
	/*for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			print(result[i][j])
			print(" ")
		}
		println()
	}*/
}

// Function that multiplies the two matrixes (here arrays) 
// It take the channel
func MatrixMulti(pairs chan matrix, a, b, result *[size][size]int, wait *sync.WaitGroup){
	for {
		// Here 'pair' and 'ok' recieve from 'pairs' channel which is of type matrix
		pair, ok := <- pairs
		// If 'ok' is false it means that there are no values to recieve and therefore
		// the channel will be closed
		if !ok {
			break
		}
		// After we checked for that we intialize the corresponding field in the 'result' matrix
		result[pair.rowNo][pair.colNo] = 0
		// Here we multiply the two values at [pair.rowNo][i] of matrix 'a'  with the [i][pair.colNo]
		// of the matrix 'b' and store the result in the 'result' matrix in the the corrensponding
		// field
		for i:= 0; i < size; i++ {
			result[pair.rowNo][pair.colNo] += a[pair.rowNo][i] * b[i][pair.colNo]
		}
	}
	// Last but not least here we signal the end of our goroutine
	wait.Done()
}
