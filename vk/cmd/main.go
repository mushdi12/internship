package main

import (
	"fmt"
	"vk/internal/utils"
	"vk/internal/worker"
)

/*
TODO:
Реализовать примитивный worker-pool с возможностью динамически добавлять и
удалять воркеры. Входные данные (строки) поступают в канал, воркеры их
обрабатывают (например, выводят на экран номер воркера и сами данные). Задание
на базовые знания каналов и горутин.
*/

func main() {
	// example
	dictionary := utils.GenerateRandomStrings(1000000, 5, 15)

	printer := func(id, msg string) {
		fmt.Printf("worker %s print: %s \n", id, msg)
	}

	wkPool := worker.NewWorkerPool(printer)

	wkPool.AddWorkers(10)

	go wkPool.Start()

	wkPool.Submit(dictionary)

	wkPool.WaitAndStop()

}

//time.Sleep(1 * time.Second)
//
//wkPool.AddWorkers(2)
//
//wkPool.Remove(0)
//
//wkPool.Remove(2)
