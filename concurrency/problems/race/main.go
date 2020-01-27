package main

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex

func main() {
	//runDeposit()
	//runTransfer1()

	//for i:=0;i<100;i++ {
	//	raceConditionBasic()
	//}

	//runTransfer2_2()

	//raceConditionBasic2_sync()
	//raceConditionBasic2_lc()
}

func raceConditionBasic2_lc() {
	x := 0
	for i := 0; i < 1000; i++ {
		go func() {
			x++
		}()
		go func() {
			y := x
			if y%2 == 0 {
				time.Sleep(1 * time.Millisecond)
				fmt.Println(y)
			}
		}()
	}
	time.Sleep(1 * time.Second)
}

func raceConditionBasic2_sync() {
	var mu sync.Mutex
	x := 0
	for i := 0; i < 1000; i++ {
		go func() {
			mu.Lock()
			x++
			mu.Unlock()
		}()
		go func() {
			mu.Lock()
			if x%2 == 0 {
				time.Sleep(1 * time.Millisecond)
				fmt.Println(x)
			}
			mu.Unlock()
		}()
	}
	time.Sleep(1 * time.Second)
}

func raceConditionBasic2() {
	x := 0
	for {
		go func() {
			x++
		}()
		go func() {
			if x%2 == 0 {
				time.Sleep(1 * time.Millisecond)
				fmt.Println(x)
			}
		}()
	}
}

func raceConditionBasic() {
	go func() {
		fmt.Printf("A->")
	}()

	go func() {
		fmt.Printf("B")
	}()

	fmt.Println()
	time.Sleep(100 * time.Millisecond)
}

type account struct {
	balance int
}

func runDeposit() {
	acc := account{balance: 0}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(n int) {
			deposit(&acc, 1)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("balance=%d\n", acc.balance)
}

func deposit(acc *account, amount int) {
	acc.balance += amount
}

func runTransfer1() {
	accFrom := account{balance: 1000}
	accTo := account{balance: 0}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(n int) {
			err := transfer1(&accFrom, &accTo, 1)
			if err != nil {
				fmt.Printf("error for n=%d\n", n)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("accFrom.balance=%d\naccTo.balance=%d\n", accFrom.balance, accTo.balance)
}

func transfer1(accFrom, accTo *account, amount int) error {
	if accFrom.balance < amount {
		return fmt.Errorf("accFrom.balance<amount")
	}
	accTo.balance += amount
	accFrom.balance -= amount
	return nil
}

func runTransfer2() {
	accFrom := account{balance: 1000}
	accTo := account{balance: 0}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(n int) {
			err := transfer2(&accFrom, &accTo, 1)
			if err != nil {
				fmt.Printf("error for n=%d\n", n)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("accFrom.balance=%d\naccTo.balance=%d\n", accFrom.balance, accTo.balance)
}

func runTransfer2_2() {
	accFrom := account{balance: 1000}
	accTo := account{balance: 0}
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(n int) {
			err := transfer2(&accFrom, &accTo, 1)
			if err != nil {
				fmt.Printf("error for n=%d\n", n)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("accFrom.balance=%d\naccTo.balance=%d\n", accFrom.balance, accTo.balance)
}

func transfer2(accFrom, accTo *account, amount int) error {
	mu.Lock()
	bal := accFrom.balance
	mu.Unlock()

	if bal < amount {
		return fmt.Errorf("accFrom.balance<amount")
	}
	mu.Lock()
	accTo.balance += amount
	mu.Unlock()

	mu.Lock()
	accFrom.balance -= amount
	mu.Unlock()

	return nil
}
