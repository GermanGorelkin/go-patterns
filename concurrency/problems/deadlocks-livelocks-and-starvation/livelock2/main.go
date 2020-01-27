package main

import (
	"fmt"
	"sync"
	"time"
)

type spoon struct {
	owner *diner
}

func (s spoon) use() {
	fmt.Printf("%s has eaten!\n", s.owner.name)
}

type diner struct {
	name     string
	isHungry bool
}

func (d *diner) eatWith(sp *spoon, spouse *diner) {
	for d.isHungry {

		// 1
		if sp.owner != d {
			time.Sleep(1 * time.Second)
			continue
		}

		// 2
		if spouse.isHungry {
			fmt.Printf("%s: You eat first my darling %s!\n", d.name, spouse.name)
			sp.owner = spouse
			continue
		}

		// 3
		sp.use()
		d.isHungry = false
		fmt.Printf("%s: I'm stuffed? my darling %s!\n", d.name, spouse.name)
		sp.owner = spouse
	}
}

func main() {
	husband := &diner{
		name:     "Bob",
		isHungry: true,
	}
	wifi := &diner{
		name:     "Alice",
		isHungry: true,
	}

	sp := &spoon{owner: husband}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		husband.eatWith(sp, wifi)
		wg.Done()
	}()

	go func() {
		wifi.eatWith(sp, husband)
		wg.Done()
	}()

	wg.Wait()
}
