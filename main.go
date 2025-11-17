package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	task := flag.Int("task", 1, "номер задачи (1-4)")
	workerCount := flag.Int("workers", 4, "количество воркеров для задачи 2")
	flag.Parse()

	switch *task {
	case 1:
		runTask1()
	case 2:
		runTask2(*workerCount)
	case 3:
		runTask3()
	case 4:
		runTask4()
	default:
		fmt.Println("Неизвестная задача, используйте -task=1..4")
	}
}

///////////////////////
// ЗАДАНИЕ 1
///////////////////////

func squareWorker(n int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	out <- n * n
}

func runTask1() {
	nums := []int{2, 4, 6, 8, 10}

	out := make(chan int, len(nums))
	var wg sync.WaitGroup

	for _, n := range nums {
		wg.Add(1)
		go squareWorker(n, out, &wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	sum := 0
	for sq := range out {
		sum += sq
	}

	fmt.Println("Сумма квадратов:", sum)
}

///////////////////////
// ЗАДАНИЕ 2
///////////////////////

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: контекст отменён, выходим\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("worker %d: канал закрыт, выходим\n", id)
				return
			}
			fmt.Printf("worker %d: получил %d\n", id, job)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func runTask2(workerCount int) {
	// Контекст, отменяемый по Ctrl+C / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jobs := make(chan int)

	var wg sync.WaitGroup

	// Стартуем воркеров
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg)
	}

	// Продюсер: постоянно пишет данные в канал до отмены контекста
	go func() {
		defer close(jobs)
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("producer: контекст отменён, прекращаем запись в канал")
				return
			case jobs <- i:
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	fmt.Println("Задание 2: запущено", workerCount, "воркеров. Нажмите Ctrl+C для завершения.")

	<-ctx.Done()
	fmt.Println("Получен сигнал завершения, ожидаем остановки воркеров...")
	wg.Wait()
	fmt.Println("Все воркеры завершены, выходим.")
}

///////////////////////
// ЗАДАНИЕ 3
///////////////////////

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]int),
	}
}

func (s *SafeMap) Set(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) Get(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SafeMap) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

func runTask3() {
	sm := NewSafeMap()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			keyPrefix := fmt.Sprintf("worker-%d", id)
			for j := 0; j < 100; j++ {
				sm.Set(fmt.Sprintf("%s-%d", keyPrefix, j), j)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Задание 3: всего элементов в карте:", sm.Len())
}

///////////////////////
// ЗАДАНИЕ 4
///////////////////////

func gen(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func multiplyByTwo(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for x := range in {
			out <- x * 2
		}
	}()
	return out
}

func runTask4() {
	nums := []int{1, 2, 3, 4, 5}

	ch1 := gen(nums)
	ch2 := multiplyByTwo(ch1)

	fmt.Println("Задание 4: результаты конвейера x*2:")
	for v := range ch2 {
		fmt.Println(v)
	}
}
