package evlistener

import (
	"reflect"
	"testing"
)

type MockWriter struct {
	data []string
}

func (mw *MockWriter) Write(b []byte) (n int, err error) {
	mw.data = append(mw.data, string(b))
	return len(b), nil
}

func TestObserver(t *testing.T) {
	evManager := EventManager{observers: map[EventType][]Observer{}}
	fileWatcher := Watcher{eventManager: &evManager}

	store := MockWriter{}
	logger := LoggingListener{store: &store}
	service := MockWriter{}
	alerts := AlertListener{service: &service}

	expectedDataInStore := make([]string, 0)
	expectedDataInService := make([]string, 0)

	eqDataOfWriter := func(t *testing.T) {
		t.Helper()

		if !reflect.DeepEqual(store.data, expectedDataInStore) {
			t.Errorf("store expected:%v, got:%v", expectedDataInStore, store.data)
		}
		if !reflect.DeepEqual(service.data, expectedDataInService) {
			t.Errorf("service expected:%v, got:%v", expectedDataInService, store.data)
		}
	}

	t.Run("Attach", func(t *testing.T) {
		evManager.Attach(EventTypeCreate, &logger)
		evManager.Attach(EventTypeRemove, &logger)
		evManager.Attach(EventTypeModify, &logger)
		evManager.Attach(EventTypeCreate, &alerts)
		evManager.Attach(EventTypeModify, &alerts)

		if len(evManager.observers[EventTypeCreate]) != 2 {
			t.Errorf("expected 2 listeners for EventTypeCreate, got:%d", len(evManager.observers[EventTypeCreate]))
		}
		if len(evManager.observers[EventTypeModify]) != 2 {
			t.Errorf("expected 2 listeners for EventTypeModify, got:%d", len(evManager.observers[EventTypeModify]))
		}
		if len(evManager.observers[EventTypeRemove]) != 1 {
			t.Errorf("expected 1 listeners for EventTypeRemove, got:%d", len(evManager.observers[EventTypeRemove]))
		}
	})

	t.Run("OnCreate", func(t *testing.T) {
		fi := FileInfo{Name: "testName"}
		fileWatcher.OnCreate(fi)

		expectedDataInStore = append(expectedDataInStore, "create|testName")
		expectedDataInService = append(expectedDataInService, "create|testName")

		eqDataOfWriter(t)
	})
	t.Run("OnCreate 2", func(t *testing.T) {
		fi := FileInfo{Name: "testName2"}
		fileWatcher.OnCreate(fi)

		expectedDataInStore = append(expectedDataInStore, "create|testName2")
		expectedDataInService = append(expectedDataInService, "create|testName2")

		eqDataOfWriter(t)
	})

	t.Run("OnRemove", func(t *testing.T) {
		fi := FileInfo{Name: "testName"}
		fileWatcher.OnRemove(fi)

		expectedDataInStore = append(expectedDataInStore, "remove|testName")

		eqDataOfWriter(t)
	})

	t.Run("OnModify", func(t *testing.T) {
		fi := FileInfo{Name: "testName"}
		fileWatcher.OnModify(fi)

		expectedDataInStore = append(expectedDataInStore, "modify|testName")
		expectedDataInService = append(expectedDataInService, "modify|testName")

		eqDataOfWriter(t)
	})

	t.Run("Detach", func(t *testing.T) {
		evManager.Detach(EventTypeCreate, &logger)

		if len(evManager.observers[EventTypeCreate]) != 1 {
			t.Errorf("expected 1 listeners for EventTypeCreate, got:%d", len(evManager.observers[EventTypeCreate]))
		}
		if len(evManager.observers[EventTypeModify]) != 2 {
			t.Errorf("expected 2 listeners for EventTypeModify, got:%d", len(evManager.observers[EventTypeModify]))
		}
		if len(evManager.observers[EventTypeRemove]) != 1 {
			t.Errorf("expected 1 listeners for EventTypeRemove, got:%d", len(evManager.observers[EventTypeRemove]))
		}
	})

	t.Run("OnCreate after detach", func(t *testing.T) {
		fi := FileInfo{Name: "testName3"}
		fileWatcher.OnCreate(fi)

		expectedDataInService = append(expectedDataInService, "create|testName3")

		eqDataOfWriter(t)
	})
}
