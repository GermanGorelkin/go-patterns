package functional_options

import "testing"

func TestNewClient(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := NewClient(":").(*client)

		if c.address != ":" {
			t.Errorf("address expected:%s, got:%s", ":", c.address)
		}
		if c.timeout != 0 {
			t.Errorf("timeout expected:%d, got:%d", 0, c.timeout)
		}
		if c.retries != 0 {
			t.Errorf("retries expected:%d, got:%d", 0, c.retries)
		}
		if c.isCheatMode != false {
			t.Errorf("isCheatMode expected:%t, got:%t", false, c.isCheatMode)
		}
	})
	t.Run("with WithTimeout and WithRetries", func(t *testing.T) {
		c := NewClient(":", WithRetries(5), WithTimeout(10)).(*client)

		if c.address != ":" {
			t.Errorf("address expected:%s, got:%s", ":", c.address)
		}
		if c.timeout != 10 {
			t.Errorf("timeout expected:%d, got:%d", 10, c.timeout)
		}
		if c.retries != 5 {
			t.Errorf("retries expected:%d, got:%d", 5, c.retries)
		}
		if c.isCheatMode != false {
			t.Errorf("isCheatMode expected:%t, got:%t", false, c.isCheatMode)
		}
	})
	t.Run("with WithTimeout and SetCheatMode", func(t *testing.T) {
		c := NewClient(":", WithRetries(5), SetCheatMode()).(*client)

		if c.address != ":" {
			t.Errorf("address expected:%s, got:%s", ":", c.address)
		}
		if c.timeout != 0 {
			t.Errorf("timeout expected:%d, got:%d", 0, c.timeout)
		}
		if c.retries != 5 {
			t.Errorf("retries expected:%d, got:%d", 5, c.retries)
		}
		if c.isCheatMode != true {
			t.Errorf("isCheatMode expected:%t, got:%t", true, c.isCheatMode)
		}
	})
}
