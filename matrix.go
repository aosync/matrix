package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var mstick int = 0
var coe []byte = []byte{3, 15}

func Timer() {
	for {
		/*
		 * Accuracy of sleep seems to highly depend on the environment, but the length of sleep(x) seems to
		 * be consistent within the same environment so it can effectively be used as a time tick.
		 */
		time.Sleep(time.Nanosecond)
		mstick++
	}
}

func Cong(seed byte) byte {
	/* Here we take advantage of datatype overflow to avoid costy modulo, since they do the same thing in effect */
	return coe[0]*seed + coe[1]
}

func MultiplexTask(tasks int, mat []byte) {
	total := len(mat)
	var wg sync.WaitGroup
	for a := 0; a < tasks; a++ {
		wg.Add(1)
		go func(id int, wg *sync.WaitGroup) {
			lower := id * total / tasks
			for i, _ := range mat[lower : (id+1)*total/tasks] {
				mat[i+lower] = Cong(byte(mstick))
			}
			wg.Done()
		}(a, &wg)
	}
	wg.Wait()
}

func main() {
	fmt.Println("Begin")

	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(coe) /* Fill the congruential parameters (a and c) with better randomness for variety */

	num := 10000 * 10000
	mat := make([]byte, num)

	go Timer()
	start := time.Now()

	MultiplexTask(1, mat) /* Keep in mind that the timer technically should be counted as one thread too for it not to be smothered */

	fmt.Println("Took (real time): ", time.Now().Sub(start))
	// fmt.Println(mat) /* Print the matrix, it can take a long time */
}
