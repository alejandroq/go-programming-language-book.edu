// Package bank provides concurrency-safe bank with one account
package bank

var (
	sema    = make(chan struct{}, 1) // a binary semaphore guarding balance
	balance int
)

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

// Deposit ...
func Deposit(amount int) {
	sema <- struct{}{} // acquire token
	balance = balance + amount
	<-sema
}

// Balance ...
func Balance() int {
	sema <- struct{}{} // acquire token
	b := balance
	<-sema // release token
	return b
}

// variables MUST BE confined to a single goroutine (or owner) at once to
// avoid data races
func teller() {
	var balance int // balance is confied to teller goroutine
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
