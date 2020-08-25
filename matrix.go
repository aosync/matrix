package main

import (
	"fmt"
	"math/rand"
	"sync"
	"runtime"
	"time"
)

var coe []byte
var d byte
var e byte

const TestAmount = 100

func MultiplexTask(tasks int, mat []byte) {
	total := len(mat)
	var wg sync.WaitGroup
	for a := 0; a < tasks; a++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			submat := mat[id * total / tasks : (id+1)*total/tasks]
			for i, _ := range submat {
				submat[i] = d*byte(i) + e
			}
			wg.Done()
		}(a, &wg)
	}
	wg.Wait()
}

func main() {
	fmt.Println("Begin")

	rand.Seed(time.Now().UTC().UnixNano())
	coe = make([]byte, 2)
	rand.Read(coe) /* Fill the congruential parameters (a and c) with better randomness for variety */
	d = coe[0]
	e = coe[1]

	num := 10000 * 10000
	mat := make([]byte, num)

	var avr []int64
	for i := 1 ; i <= runtime.NumCPU() ; i++ {
		var lavr int64 = 0
		fmt.Printf("%d threads\n", i)
		for j := 0 ; j < TestAmount ; j++ {
			start := time.Now()

			MultiplexTask(i, mat) /* Keep in mind that the timer technically should be counted as one thread too for it not to be smothered if value passed is CPUNUM */
		
			telap := time.Now().Sub(start).Milliseconds()
			lavr += telap
			fmt.Printf("Cores: %d | Iteration: %d -> Took (real time): %d ms\n", i,  j, telap)
		}
		lavr /= TestAmount
		avr = append(avr, lavr)
	}
	fmt.Println("Testing finished!")
	for i, a := range avr {
		fmt.Printf("Average on %d cores: %d ms.\n", i+1, a)
	}
}
