package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Value any
	//Expiration int64
	Expiration time.Time
}
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// Constructor function for Cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// implementation of set method
func (c *Cache) Set(key string, value any, ttl time.Duration) { //Store a value in the cache with a key, and make it expire after some time.//
	c.mu.Lock()
	defer c.mu.Unlock()
	//expiration := time.Now().Add(ttl).Unix()
	expiration := time.Now().Add(ttl)
	c.items[key] = CacheItem{Value: value, Expiration: expiration}
}

// implementation of get method
func (c *Cache) Get(key string) (any, bool) { //Try to fetch a value—but only if it exists and hasn’t expired.//
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	//if !found || (item.Expiration > 0 && time.Now().Unix() > item.Expiration) { //This item had a TTL, and it has now expired//
	//if !found OR expired//

	if !found || (item.Expiration.IsZero() || time.Now().After(item.Expiration)) { //This item had a TTL, and it has now expired//
		return nil, false //Either it doesn’t exist, or it is no longer valid//
	}
	return item.Value, true //return item.Value, true//
}

// Implemementation of automatic cleanup and DEFINE AN ANONYMOUS FUNCTION AND CALL IT IMMEDIATELY but in a goroutine
// Run a background process that periodically removes expired items from the cache//

//func (c *Cache) StartCleanup(interval time.Duration) {

func (c *Cache) StartCleanup(ctx context.Context, interval time.Duration) {

	go func() { //Start this process in the background, independently of the main flow
		//ticker := time.NewTicker(interval) //Ticker is part of the cache’s internal lifecycle—it periodically cleans expired entries.
		//Sleep is only used in main to simulate time passing so we can observe TTL behavior.
		//  In production, we wouldn’t rely on Sleep; the system would run continuously and expiration would happen naturally.
		
		// for range ticker.C {

		// 	c.mu.Lock() //lock for writing as we are about to modify the map, so nobody can read and write now
		// 	//now := time.Now().Unix()         //capture the current time

		// 	time.Now()
		// 	for key, item := range c.items { //Check every stored entry and decide if it should stay or go

		// 		//if item.Expiration > 0 && now > item.Expiration { //This item had a TTL, and it has now expired
		// 		if item.Expiration.IsZero() || time.Now().After(item.Expiration) { //This item had a TTL, and it has now expired
		// 			delete(c.items, key)
		// 		}
		// 	}
        
		ticker := time.NewTicker(interval)
		defer ticker.Stop() //ensuring ticker stops when the goroutine exit.

        for{
			select { //either do work when ticker ticks or exit when context is cancelled
			case <-ticker.C:
				c.cleanup() //call a separate method to handle the cleanup logic 
			case <-ctx.Done():  //ctx.Done already returns a channel that gets closed when context is cancelled.
				return //exit the goroutine when the context is cancelled. GOROUTINE HAS A LIFECYCLE NOW.
			}
		}
	}()
}
func(c *Cache) cleanup(){
			c.mu.Lock()
			defer c.mu.Unlock() //ensure the lock is released even if an error occurs

			now := time.Now()
			for key,item := range c.items{
				if !item.Expiration.IsZero() && now.After(item.Expiration){ //IsZero() means no expiration is set.
					delete(c.items, key)
				}
			}
		}
        
	//The interval passed to StartCleanup is just a configuration value. The ticker is what actually
	// uses that interval to trigger periodic execution. Without ticker (or sleep), the cleanup
	// loop would either not run periodically or run continuously. So ticker is not redundant—it’s
	// the mechanism that enforces the interval


// main func
func main() {
	//cache := &Cache{items: make(map[string]CacheItem)} //creating a Cache instance and allocating memory for the map

	cache := NewCache() //“Encapsulated object initialization to ensure the internal map is always properly
	// allocated, improving safety and abstraction

	//WithCancel gives us a handle(cancel) to stop the worker.
	ctx, cancel := context.WithCancel((context.Background()))
	defer cancel()//ensure that the cleaup is stopped when main exits.

	go cache.StartCleanup(ctx ,10 * time.Second) // a goroutine starts and every 10 seconds → expired items are removed
	fmt.Println("Setting value...")
	cache.Set("username", "message_dev", 5*time.Second) //store this value but only for 5 secs
	fmt.Println("Immediately fetching...")
	value, found := cache.Get("username")
	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found!")
	}

	time.Sleep(11 * time.Second) // Wait for TTL to expire
	value, found = cache.Get("username")
	if found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Not found!") // Should print "Not found!"
	}
}
