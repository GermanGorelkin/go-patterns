package classic

import (
	"testing"
)

func TestObserver(t *testing.T) {
	subject := ConcreteSubject{}

	obs1 := ConcreteObserver{subject: &subject, state: "obs1"}
	obs2 := ConcreteObserver{subject: &subject, state: "obs2"}

	t.Run("Attach", func(t *testing.T) {
		subject.Attach(&obs1)
		subject.Attach(&obs2)

		if len(subject.observers) != 2 {
			t.Errorf("expected len: 2, got:%d", len(subject.observers))
		}
	})

	t.Run("Notify 2 obs", func(t *testing.T) {
		state := "new test"

		subject.SetState(state)

		for _, o := range subject.observers {
			obs := o.(*ConcreteObserver)
			if obs.state != state {
				t.Errorf("expected state:%v, got:%v", state, obs.state)
			}
		}
	})

	t.Run("Detach", func(t *testing.T) {
		subject.Detach(&obs1)

		if len(subject.observers) != 1 {
			t.Errorf("expected len: 1, got:%d", len(subject.observers))
		}
	})

	t.Run("Notify after detach", func(t *testing.T) {
		state := "new test2"

		subject.SetState(state)

		for _, o := range subject.observers {
			obs := o.(*ConcreteObserver)
			if obs.state != state {
				t.Errorf("expected state:%v, got:%v", state, obs.state)
			}
		}
	})
}
