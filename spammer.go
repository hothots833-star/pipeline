package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	in := make(chan interface{})
	for _, c := range cmds {
		out := make(chan interface{})
		go func(cmd cmd, in, out chan interface{}) {
			cmd(in, out)
			close(out)
		}(c, in, out)
		in = out
	}
	for range in {
	}
}

func SelectUsers(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	seen := make(map[uint64]bool)

	for val := range in {
		wg.Add(1)
		go func(val interface{}) {
			defer wg.Done()
			email, ok := val.(string)
			if !ok {
				return
			}
			user := GetUser(email)
			mu.Lock()
			if seen[user.ID] {
				mu.Unlock()
				return
			}
			seen[user.ID] = true
			mu.Unlock()
			out <- user
		}(val)
	}
	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	batch := make([]User, 0, GetMessagesMaxUsersBatch)

	flush := func(b []User) {
		wg.Add(1)
		go func(users []User) {
			defer wg.Done()
			ids, err := GetMessages(users...)
			if err != nil {
				return
			}
			for _, v := range ids {
				out <- v
			}
		}(b)
	}

	for val := range in {
		user, ok := val.(User)
		if !ok {
			continue
		}
		batch = append(batch, user)
		if len(batch) == GetMessagesMaxUsersBatch {
			flush(batch)
			batch = make([]User, 0, GetMessagesMaxUsersBatch)
		}
	}
	if len(batch) > 0 {
		flush(batch)
	}
	wg.Wait()
}

func CheckSpam(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	spamMaxAsyncRequests := make(chan uint8, HasSpamMaxAsyncRequests)
	for val := range in {
		wg.Add(1)
		go func(val interface{}) {
			spamMaxAsyncRequests <- 1
			defer func() {
				<-spamMaxAsyncRequests
				wg.Done()
			}()
			id, ok := val.(MsgID)
			if !ok {
				return
			}
			hasspam, err := HasSpam(id)
			if err != nil {
				return
			}
			out <- MsgData{ID: id, HasSpam: hasspam}
		}(val)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []MsgData
	for val := range in {
		msgData, ok := val.(MsgData)
		if !ok {
			continue
		}
		results = append(results, msgData)
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].HasSpam != results[j].HasSpam {
			return results[i].HasSpam
		}
		return results[i].ID < results[j].ID
	})
	for _, r := range results {
		out <- fmt.Sprintf("%v %v", r.HasSpam, r.ID)
	}
}
