package mutex

import (
	"testing"
	"fmt"
	"sync"
)

func TestStartInstance(t *testing.T) {
	singleton := GetInstance()
	singleton2 := GetInstance()

	n := 5000

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			singleton.AddOne()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			singleton2.AddOne()
			wg.Done()
		}()
	}

	fmt.Printf("Before loop, current count is %d\n", singleton.GetCount())

	wg.Wait()

	fmt.Printf("Current count is %d\n", singleton.GetCount())

	currentCount1 := singleton.GetCount()
	currentCount2 := singleton2.GetCount()
	if currentCount1 != currentCount2 {
		t.Errorf("Counts not match\nCurrentCount1=%d\nCurrentCount2=%d", currentCount1, currentCount2)
	}

	if currentCount1 != n*2 {
		t.Errorf("Counts not match\nCurrentCount1=%d\nN*2=%d", currentCount1, n*2)
	}
}
