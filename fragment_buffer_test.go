package photon_spectator

import (
	"reflect"
	"testing"
)

func TestFragmentBuffer(t *testing.T) {
	fragmentOne := ReliableFragment{
		SequenceNumber: 7,
		FragmentNumber: 0,
		FragmentCount:  2,
		Data:           []byte{0xca},
	}

	fragmentTwo := ReliableFragment{
		SequenceNumber: 7,
		FragmentNumber: 1,
		FragmentCount:  2,
		Data:           []byte{0xfe},
	}

	buffer := NewFragmentBuffer()

	response := buffer.Offer(fragmentOne)

	if response != nil {
		t.Fail()
	}

	response = buffer.Offer(fragmentTwo)

	if response == nil {
		t.Fail()
	}

	if !reflect.DeepEqual((*response).Data, []byte{0xca, 0xfe}) {
		t.Fail()
	}

	if response.ReliableSequenceNumber != 7 {
		t.Errorf("response.ReliableSequenceNumber = %d, wanted 7", response.ReliableSequenceNumber)
	}
}
