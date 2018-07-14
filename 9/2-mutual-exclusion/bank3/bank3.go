// Package bank provides concurrency-safe bank with one account
package bank

import (
	"sync"
)

// Mutual exclusion is so useful that is supported directly by
// the `Mutex` type from the `sync` package.
// Nothing can Lock if mu is already Locked. Will await an Unlock
// before proceeding.
var (
	mu      sync.Mutex // guards balance
	balance int
)

// When you use a mutex, verify that both it and the variables
// it guards are not exported, whether they are package-level
// variables or the fields of a struct

// NOTE: defer runs even in occurences of `panic` (think `recover`). So Unlock in
// these situations may be appropriate.

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

// Deposit ...
func Deposit(amount int) {
	mu.Lock()
	deposit(amount) // section between a Lock and Unlock is the "critical section"
	mu.Unlock()
}

func deposit(amount int) { balance = balance + amount }

// Balance ...
func Balance() int {
	mu.Lock()
	defer mu.Unlock() // if branching code paths exist and unsure where to place Unlock, defer is to the rescue
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
