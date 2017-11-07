// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// guid implements the guid type, which is used to identify about anything
// in UEFI. This package is used in reading DXE dependency sections. The
// reason we can't use the google/uuid package is that EFI GUIDS are little
// endian.

package guid

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	glen     = 16
	gExample = "01234567-89AB-CDEF-0123-456789ABCDEF"
	gTextLen = len(gExample)
	gFmt     = "%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X"
)

var (
	fields = [...]int{4, 2, 2, 2, 1, 1, 1, 1, 1, 1}
)

type GUID [glen]byte

func reverse(b []byte) {
	for i := 0; i < len(b)/2; i++ {
		other := len(b) - i - 1
		b[other], b[i] = b[i], b[other]
	}
}

func (guid *GUID) Parse(guidString string) error {
	// remove all hyphens to make it easier to parse.
	strippedGUID := strings.Replace(guidString, "-", "", -1)
	decoded, err := hex.DecodeString(strippedGUID)
	if err != nil {
		return fmt.Errorf("guid string not correct, need string of the format \n%v\n, got \n%v\n",
			gExample, guidString)
	}

	if len(decoded) != glen {
		return fmt.Errorf("guid string has incorrect length, need string of the format \n%v\n, got \n%v\n",
			gExample, guidString)
	}

	// copy into guid now that we've swapped the endianness
	copy(guid[:], decoded[:])
	return nil
}

func (guid *GUID) String() string {
	return fmt.Sprintf(gFmt,
		guid[0], guid[1], guid[2], guid[3],
		guid[4], guid[5],
		guid[6], guid[7],
		guid[8], guid[9],
		guid[10], guid[11], guid[12], guid[13], guid[14], guid[15])
}

func (guid *GUID) MarshalBinary() ([]byte, error) {
	var data []byte
	data = make([]byte, glen, glen)

	copy(data[:], guid[:])
	guidIndex := 0
	for _, fieldlen := range fields {
		reverse(data[guidIndex : guidIndex+fieldlen])
		guidIndex += fieldlen
	}

	return data, nil
}

func (guid *GUID) UnmarshalBinary(data []byte) error {
	if size := len(data); size != glen {
		return fmt.Errorf("invalid GUID, expected %d bytes, got %d",
			glen, size)
	}

	copy(guid[:], data[:])

	// god I hate EFI.
	guidIndex := 0
	for _, fieldlen := range fields {
		reverse(guid[guidIndex : guidIndex+fieldlen])
		guidIndex += fieldlen
	}

	return nil
}
