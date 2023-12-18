package main

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// ms = sec / 1000
func countWriteDurationMP(client *redis.Client, amt int) {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(amt)

	for i := 0; i < amt; i++ {
		go func(i int) {
			defer wg.Done()
			err := client.Set("name"+string(i), "Leha", 0).Err()
			if err != nil {
				log.Panic(err)
			}
		}(i)
	}

	wg.Wait()
	end := time.Now()

	log.Printf("MPWrite: %v ms per op", float64(end.Sub(start).Milliseconds())/float64(amt))
}

func countWriteDuration(client *redis.Client, amt int) {
	start := time.Now()

	for i := 0; i < amt; i++ {
		err := client.Set("name"+string(i), "Leha", 0).Err()
		if err != nil {
			log.Panic(err)
		}
	}

	end := time.Now()

	log.Printf("Write: %v ms per op", float64(end.Sub(start).Milliseconds())/float64(amt))
}

func countReadDurationMP(client *redis.Client, amt int) {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(amt)

	for i := 0; i < amt; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := client.Get("name" + string(i)).Result()
			if err != nil {
				log.Panic(err)
			}
		}(i)
	}

	wg.Wait()
	end := time.Now()

	log.Printf("MPRead: %v ms per op", float64(end.Sub(start).Milliseconds())/float64(amt))
}

func countReadDuration(client *redis.Client, amt int) {
	start := time.Now()

	for i := 0; i < amt; i++ {
		_, err := client.Get("name" + string(i)).Result()
		if err != nil {
			log.Panic(err)
		}
	}

	end := time.Now()

	log.Printf("Read: %v ms per op", float64(end.Sub(start).Milliseconds())/float64(amt))
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:8080",
		Password: "",
		DB:       0,
	})

	countWriteDuration(client, 1e3)
	countWriteDurationMP(client, 1e4)
	countReadDuration(client, 1e3)
	countReadDurationMP(client, 1e4)
}
