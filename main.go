package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	var termsQuantity, threadsQnt int

	fmt.Print("Insira a quantidade de termos: ")
	_, _ = fmt.Scan(&termsQuantity)

	fmt.Print("Insira a quantidade de threads: ")
	_, _ = fmt.Scan(&threadsQnt)

	var numberOfExecutions = 5
	statistics := NewStatistics()
	calculator := NewCalculator(termsQuantity, threadsQnt)

	for i := 0; i < numberOfExecutions; i++ {
		statistics.AddResult(
			*calculator.Calculate(),
			)
	}

	fmt.Println("Valores de PI / Tempo de execução:  [")
	for _, r := range statistics.GetResult() {
		fmt.Println("    ", r)
	}
	fmt.Println("]\nMédia do tempo de execução (ms): ", statistics.GetAverageTimeInMilliseconds())
	fmt.Printf("Desvio padrão do tempo de execução: %.2f\n", statistics.GetStandardDeviation())
}


type Calculator struct {
	termsQnt int
	threadsQnt int
}

func NewCalculator (termsQnt, threadsQnt int) *Calculator {
	return &Calculator{
		termsQnt:   termsQnt,
		threadsQnt: threadsQnt,
	}
}

func (pc Calculator) run(c chan float64, initialN, termsQnt int) {
	var pi float64

	for n := initialN; n < initialN +termsQnt; n++ {
		pi += math.Pow(-1, float64(n)) / float64(2 * n + 1)
	}

	pi *= 4
	c <- pi
}

func (pc Calculator) Calculate() *Result {
	start := time.Now()

	termsQuantityByThread := pc.termsQnt / pc.threadsQnt
	initialNumber := 0

	var channels []chan float64
	for j := 0; j < pc.threadsQnt; j++ {
		channels = append(channels, make(chan float64))
		go pc.run(channels[j], initialNumber, termsQuantityByThread)

		initialNumber += termsQuantityByThread
	}

	var pi float64
	for _, c := range channels {
		pi += <-c
	}

	spentTime := time.Since(start)
	return &Result{pi, spentTime}
}

type Result struct {
	Pi  float64
	Duration time.Duration
}

type Statistics struct {
	result []Result
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

func (sc *Statistics) AddResult(set Result) {
	sc.result = append(sc.result, set)
}

func (sc Statistics) GetResult() []Result {
	return sc.result
}

func (sc Statistics) GetAverageTimeInMilliseconds() int64 {
	var sumTime int64
	for _, result := range sc.result {
		sumTime += result.Duration.Milliseconds()
	}
	return sumTime / int64(len(sc.result))
}

func (sc Statistics) GetStandardDeviation () float64 {
	var execTime float64
	for _, rs := range sc.result {
		aux := rs.Duration.Milliseconds() - sc.GetAverageTimeInMilliseconds()
		execTime += math.Pow(float64(aux), 2) / float64(len(sc.result))
	}
	return math.Sqrt(execTime)
}