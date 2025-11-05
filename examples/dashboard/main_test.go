package main

import "testing"

func TestSnapshot(t *testing.T) {
	m := newModel()
	m.width = 80
	m.height = 24
	m.addCamera()
	out := m.View()
	t.Logf("\n%s", out)
}
