package base

import "testing"

func TestGetInstance(t *testing.T) {
	var counter1 Singleton
	counter1 = GetInstance()
	if counter1 == nil {
		t.Fatalf("expected pointer to Singleton after calling GetInstance(), not nil")
	}

	currentCount := counter1.AddOne()
	if currentCount != 1 {
		t.Errorf("After calling for the first time to count, the count must be 1 but it is %d\n", currentCount)
	}

	var counter2, expectedCounter Singleton
	expectedCounter = counter1
	counter2 = GetInstance()
	if counter2 != expectedCounter {
		t.Error("Expected same instance in counter2 but it got a different instance")
	}

	currentCount = counter2.AddOne()
	if currentCount != 2 {
		t.Errorf("After calling 'AddOne' using the second counter, the current count must be 2 but was %d\n", currentCount)
	}
}
