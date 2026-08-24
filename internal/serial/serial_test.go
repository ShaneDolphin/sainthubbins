package serial

import "testing"

func TestSerialNew(t *testing.T) {
	s := New("/dev/ttyUSB0")
	if s.Port != "/dev/ttyUSB0" {
		t.Fatalf("serial port")
	}
	if n, _ := s.Write([]byte{1,2,3}); n != 3 {
		t.Fatalf("write")
	}
}
