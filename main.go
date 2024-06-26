package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"
)

type User struct {
	Email string
	Name  string
}

var DataBase = []User{
	{Email: "alexander.davis@example.com", Name: "Alexander Davis"},
	{Email: "alexander.jackson@example.com", Name: "Alexander Jackson"},
	{Email: "avery.williams@example.com", Name: "Avery Williams"},
	{Email: "charlotte.smith@example.com", Name: "Charlotte Smith"},
	{Email: "daniel.miller@example.com", Name: "Daniel Miller"},
	{Email: "ella.smith@example.com", Name: "Ella Smith"},
	{Email: "jacob.white@example.com", Name: "Jacob White"},
	{Email: "james.martinez@example.com", Name: "James Martinez"},
	{Email: "james.miller@example.com", Name: "James Miller"},
	{Email: "jayden.jackson@example.com", Name: "Jayden Jackson"},
	{Email: "liam.robinson@example.com", Name: "Liam Robinson"},
	{Email: "mason.martin@example.com", Name: "Mason Martin"},
	{Email: "matthew.jackson@example.com", Name: "Matthew Jackson"},
	{Email: "mia.smith@example.com", Name: "Mia Smith"},
	{Email: "michael.white@example.com", Name: "Michael White"},
	{Email: "natalie.martin@example.com", Name: "Natalie Martin"},
	{Email: "sofia.garcia@example.com", Name: "Sofia Garcia"},
	{Email: "william.brown@example.com", Name: "William Brown"},
}

type Worker struct {
	users []User
	ch    chan *User
}

func NewWorker(users []User, ch chan *User) *Worker {
	return &Worker{users: users, ch: ch}
}

func (w *Worker) find(email string) {
	for _, u := range w.users {
		if strings.Contains(u.Email, email) {
			w.ch <- &u
		}
	}
}

func splitDb(numPerGroup int, db []User) [][]User {

	length := len(db)
	var numOfGroups int

	if numPerGroup > length {
		numOfGroups = length
	} else {
		numOfGroups = length / numPerGroup
		if length%numOfGroups > 0 {
			numOfGroups++
		}
	}
	result := make([][]User, numOfGroups)

	j := 0
	for i := 0; i < length; i += numPerGroup {
		end := i + numPerGroup
		if end > length {
			end = length
		}
		chunk := db[i:end]

		result[j] = append(result[j], chunk...)
		j++
	}
	return result
}

func main() {

	email := flag.String("email", "", "email address")
	flag.Parse()

	if email == nil || *email == "" {
		flag.PrintDefaults()
		return
	}

	result := make(chan *User)
	done := make(chan bool)
	wg := sync.WaitGroup{}

	chunks := splitDb(3, DataBase)

	wg.Add(len(chunks))
	for i := range chunks {

		go func(i int) {
			defer wg.Done()
			NewWorker(chunks[i], result).find(*email) // result on channel
		}(i)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

loop:
	for {
		select {
		case <-done:
			fmt.Println("Finished")
			break loop
		case found := <-result:
			fmt.Println("Found:", found)
		case <-time.After(5 * time.Second):
			fmt.Println("Timed out")
			break loop
		}
	}
}
