// Package bank provides concurrency-safe bank with one account
package bank

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

// Deposit ...
func Deposit(amount int) { deposits <- amount }

// Balance ...
func Balance() int { return <-balances }

// Variables MUST BE confined to a single goroutine (or owner) at once to
// avoid data races
// If difficult, practice "serial confinemnet". I.E. a Cake must be passed
// on, but the baker never touches it again.
// If impossible, practice "mutual exclusion".
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
