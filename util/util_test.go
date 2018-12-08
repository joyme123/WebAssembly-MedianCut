package util

import (
	"testing"
)

func TestRGB(t *testing.T) {

	if RGB(0, 0, 0) != 0 {
		t.Error("RGB should be 0")
	}

	if RGB(255, 255, 255) != 32767 {
		t.Error("RGB should be 32767")
	}

	if RGB(67, 119, 195) != 25032 {
		t.Error("RGB should be 25032")
	}
}

func TestRED(t *testing.T) {

	if RED(0) != 0 {
		t.Error("red should be 0")
	}

	if RED(25032) != 64 {
		t.Error("red should be 64")
	}

	if RED(32767) != 248 {
		t.Error("red should be 248")
	}
}

func TestGREEN(t *testing.T) {
	if GREEN(0) != 0 {
		t.Error("green should be 0")
	}

	if GREEN(25032) != 112 {
		t.Error("green should be 112")
	}

	if GREEN(32767) != 248 {
		t.Error("green should be 248")
	}
}

func TestBLUE(t *testing.T) {
	if GREEN(0) != 0 {
		t.Error("green should be 0")
	}

	if GREEN(25032) != 112 {
		t.Error("green should be 112")
	}

	if GREEN(32767) != 248 {
		t.Error("green should be 32767")
	}
}
