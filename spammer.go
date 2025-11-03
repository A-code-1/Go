package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	if len(cmds) == 0 {
		return
	}
	// стартовый канал
	in := make(chan interface{})
	out := in

	wg := &sync.WaitGroup{}
	for _, c := range cmds {
		out = make(chan interface{})
		wg.Add(1)

		go func(c cmd, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			c(in, out)
		}(c, in, out)

		in = out
	}

	// дренируем последний канал
	go func() {
		for range out {
		}
	}()
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	seen := make(map[uint64]struct{})

	for v := range in {
		email := v.(string)
		wg.Add(1)
		go func(e string) {
			defer wg.Done()
			u := GetUser(e)
			mu.Lock()
			_, exists := seen[u.ID]
			if !exists {
				seen[u.ID] = struct{}{}
			}
			mu.Unlock()
			if !exists {
				out <- u
			}
		}(email)
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	var wg sync.WaitGroup
	// канал батчей
	batches := make(chan []User, 10)
	// собирающий
	go func() {
		defer close(batches)
		b := make([]User, 0, GetMessagesMaxUsersBatch)
		for v := range in {
			u := v.(User)
			b = append(b, u)
			if len(b) == GetMessagesMaxUsersBatch {
				// отправляем копию (чтобы безопасно параллелить)
				tmp := make([]User, len(b))
				copy(tmp, b)
				batches <- tmp
				b = make([]User, 0, GetMessagesMaxUsersBatch)
			}
		}
		if len(b) > 0 {
			tmp := make([]User, len(b))
			copy(tmp, b)
			batches <- tmp
		}
	}()
	// обрабатываем батчи
	for batch := range batches {
		wg.Add(1)
		go func(users []User) {
			defer wg.Done()
			msgs, err := GetMessages(users...)
			if err != nil {
				return
			}
			for _, m := range msgs {
				out <- m
			}
		}(batch)
	}
	wg.Wait()
}

// ограничиваем число одновременных вызовов HasSpam
func CheckSpam(in, out chan interface{}) {
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	var wg sync.WaitGroup

	for v := range in {
		id := v.(MsgID)
		wg.Add(1)
		go func(mid MsgID) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			ok, err := HasSpam(mid)
			if err == nil {
				out <- MsgData{ID: mid, HasSpam: ok}
			}
		}(id)
	}
	wg.Wait()
}

// аккумулируем, сортируем и отправляем строки
func CombineResults(in, out chan interface{}) {
	var res []MsgData

	for v := range in {
		res = append(res, v.(MsgData))
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].HasSpam != res[j].HasSpam {
			return res[i].HasSpam
		}
		return res[i].ID < res[j].ID
	})

	for _, r := range res {
		out <- fmt.Sprintf("%t %d", r.HasSpam, r.ID)
	}
}
