// Package bank provides concurrency-safe bank with one account
package bank

import (
	"sync"
)

var (
	mu      sync.RWMutex // profitable when MOST operations are read. In this case, bank reads will occur more quickly.
	balance int
)

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

// Deposit ...
func Deposit(amount int) {
	mu.Lock()         // writers lock
	defer mu.Unlock() // writers unlock
	deposit(amount)
}

func deposit(amount int) { balance = balance + amount }

// Balance ...
func Balance() int {
	mu.RLock()         // readers lock
	defer mu.RUnlock() // if branching code paths exist and unsure where to place Unlock, defer is to the rescue
	return balance
}

// Withdraw ...
// If higher order than Deposit, it cannot run an additional mutex lock.
// In situations like this, it is best to split functions like Deposit into
// two.
func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	deposit(-amount)
	if balance < 0 {
		deposit(amount)
		return false
	}
	return true
}

// variables MUST BE confined to a single goroutine (or owner) at once to
// avoid data races
func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}
