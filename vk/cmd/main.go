package main

import (
	"fmt"
	"math/rand"
	"time"
	"vk/pkg/worker"
)

/*
TODO:
Реализовать примитивный worker-pool с возможностью динамически добавлять и
удалять воркеры. Входные данные (строки) поступают в канал, воркеры их
обрабатывают (например, выводят на экран номер воркера и сами данные). Задание
на базовые знания каналов и горутин.
*/

func GenerateRandomStrings(count int, minLen int, maxLen int) []string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		length := rand.Intn(maxLen-minLen+1) + minLen
		s := make([]rune, length)
		for j := range s {
			s[j] = letters[rand.Intn(len(letters))]
		}
		result = append(result, string(s))
	}
	return result
}

func main() {
	// example
	dictionary := GenerateRandomStrings(1000000, 5, 15)

	printer := func(id, msg string) {
		fmt.Printf("worker %s print: %s \n", id, msg)
	}

	wkPool := worker.NewWorkerPool(printer)

	wkPool.AddWorkers(2)

	go wkPool.Start()

	go wkPool.Submit(dictionary)

	time.Sleep(10 * time.Second)
	
	wkPool.Remove(0)

	time.Sleep(100 * time.Second)
	wkPool.WaitAndStop()

}
