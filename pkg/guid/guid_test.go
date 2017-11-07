package guid

import (
	"bytes"
	"strings"
	"testing"
)

func TestBinary(t *testing.T) {
	var g GUID
	var testGUID = [...]byte{
		0x67, 0x45, 0x23, 0x01,
		0xAB, 0x89,
		0xEF, 0xCD,
		0x23, 0x01,
		0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

	if err := g.UnmarshalBinary(testGUID[:]); err != nil {
		t.Errorf("Unable to unmarshal the binary: %v", err)
	}

	data, err := g.MarshalBinary()

	if err != nil {
		t.Errorf("Unable to marshall the binary: %v", err)
	}
	if !bytes.Equal(testGUID[:], data) {
		t.Errorf("Bytes are not equal after unmarshal and marshal, original: %X, new: %X", testGUID, data)
	}
}

func TestString(t *testing.T) {
	var g GUID
	testGUIDString := "01234567-89AB-CDEF-0123-456789ABCDEF"

	if err := g.Parse(testGUIDString); err != nil {
		t.Errorf("Unable to parse the guid: %v", err)
	}

	newString := g.String()
	if strings.Compare(newString, testGUIDString) != 0 {
		t.Errorf("Strings are not equal: original %v, new %v", testGUIDString, newString)
	}
}
