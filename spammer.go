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

	wg := &sync.WaitGroup{}
	// запускаем все команды кроме последней
	for i := 0; i < len(cmds)-1; i++ {
		out := make(chan interface{})
		wg.Add(1)
		go func(c cmd, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			c(in, out)
		}(cmds[i], in, out)
		in = out
	}

	// последняя команда — отдельный канал
	lastOut := make(chan interface{})
	wg.Add(1)
	go func(c cmd, in, out chan interface{}) {
		defer wg.Done()
		defer close(out)
		c(in, out)
	}(cmds[len(cmds)-1], in, lastOut)

	// читаем всё из последнего канала
	go func() {
		for range lastOut {
			// пусто — дренируем
		}
	}()
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	var wg sync.WaitGroup
	seen := make(map[uint64]bool)
	var mtx sync.Mutex

	for v := range in {
		email := v.(string)
		wg.Add(1)
		go func(e string) {
			defer wg.Done()
			u := GetUser(e)

			mtx.Lock()
			have := seen[u.ID]
			if !have {
				seen[u.ID] = true
			}
			mtx.Unlock()
			if !have {
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
	res := make([]MsgData, 0)
	for v := range in {
		res = append(res, v.(MsgData))
	}
	// сначала HasSpam=true, потом false, везде по ID по возрастанию
	sort.Slice(res, func(i, j int) bool {
		if res[i].HasSpam && !res[j].HasSpam {
			return true
		}
		if !res[i].HasSpam && res[j].HasSpam {
			return false
		}
		return res[i].ID < res[j].ID
	})
	for _, r := range res {
		out <- fmt.Sprintf("%t %d", r.HasSpam, r.ID)
	}
}
