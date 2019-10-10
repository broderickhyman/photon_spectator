package photon_spectator

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	NilType               = 42
	DictionaryType        = 68
	StringSliceType       = 97
	Int8Type              = 98
	Custom                = 99
	DoubleType            = 100
	EventDateType         = 101
	Float32Type           = 102
	Hashtable             = 104
	Int32Type             = 105
	Int16Type             = 107
	Int64Type             = 108
	Int32SliceType        = 110
	BooleanType           = 111
	OperationResponseType = 112
	OperationRequestType  = 113
	StringType            = 115
	Int8SliceType         = 120
	SliceType             = 121
	ObjectSliceType       = 122
)

type ReliableMessageParamaters map[uint8]interface{}

// Converts the parameters of a reliable message into a hash suitable for use in
// hashmap.
func DecodeReliableMessage(msg ReliableMessage) ReliableMessageParamaters {
	buf := bytes.NewBuffer(msg.Data)
	params := make(map[uint8]interface{})

	for i := 0; i < int(msg.ParamaterCount); i++ {
		var paramID uint8
		var paramType uint8

		binary.Read(buf, binary.BigEndian, &paramID)
		binary.Read(buf, binary.BigEndian, &paramType)

		paramsKey := paramID
		params[paramsKey] = decodeType(buf, paramType)
	}

	return params
}

func decodeType(buf *bytes.Buffer, paramType uint8) interface{} {
	switch paramType {
	case NilType, 0:
		// Do nothing
		return nil
	case Int8Type:
		return decodeInt8Type(buf)
	case Float32Type:
		return decodeFloat32Type(buf)
	case Int32Type:
		return decodeInt32Type(buf)
	case Int16Type, 7:
		return decodeInt16Type(buf)
	case Int64Type:
		return decodeInt64Type(buf)
	case StringType:
		return decodeStringType(buf)
	case BooleanType:
		result, err := decodeBooleanType(buf)

		if err != nil {
			return fmt.Sprintf("ERROR - Boolean - %v", err.Error())
		} else {
			return result
		}
	case Int8SliceType:
		result, err := decodeSliceInt8Type(buf)
		if err != nil {
			return fmt.Sprintf("ERROR - Slice Int8 - %v", err.Error())
		} else {
			return result
		}
	case SliceType:
		array, err := decodeSlice(buf)
		if err != nil {
			return fmt.Sprintf("ERROR - Slice - %v", err.Error())
		} else {
			return array
		}
	case DictionaryType:
		dict, err := decodeDictionaryType(buf)
		if err != nil {
			return fmt.Sprintf("ERROR - Dictionary - %v", err.Error())
		} else {
			return dict
		}
	}
	return fmt.Sprintf("ERROR - Invalid type of %v", paramType)
}

func decodeSlice(buf *bytes.Buffer) (interface{}, error) {
	var length uint16
	var sliceType uint8

	binary.Read(buf, binary.BigEndian, &length)
	binary.Read(buf, binary.BigEndian, &sliceType)

	switch sliceType {
	case Float32Type:
		array := make([]float32, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeFloat32Type(buf)
		}

		return array, nil
	case Int32Type:
		array := make([]int32, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeInt32Type(buf)
		}

		return array, nil
	case Int16Type:
		array := make([]int16, length)

		for j := 0; j < int(length); j++ {
			var temp int16
			binary.Read(buf, binary.BigEndian, &temp)
			array[j] = temp
		}

		return array, nil
	case Int64Type:
		array := make([]int64, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeInt64Type(buf)
		}

		return array, nil
	case StringType:
		array := make([]string, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeStringType(buf)
		}

		return array, nil
	case BooleanType:
		array := make([]bool, length)

		for j := 0; j < int(length); j++ {
			result, err := decodeBooleanType(buf)

			if err != nil {
				return array, err
			}

			array[j] = result
		}

		return array, nil
	case Int8SliceType:
		array := make([][]int8, length)

		for j := 0; j < int(length); j++ {
			result, err := decodeSliceInt8Type(buf)
			if err != nil {
				return nil, err
			}
			array[j] = result
		}

		return array, nil
	case SliceType:
		array := make([]interface{}, length)

		for j := 0; j < int(length); j++ {
			subArray, error := decodeSlice(buf)

			if error != nil {
				return nil, error
			}

			array[j] = subArray
		}

		return array, nil
	default:
		return nil, fmt.Errorf("Invalid slice type of %d", sliceType)
	}
}

func decodeInt8Type(buf *bytes.Buffer) (temp int8) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeFloat32Type(buf *bytes.Buffer) (temp float32) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt16Type(buf *bytes.Buffer) (temp int16) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt32Type(buf *bytes.Buffer) (temp int32) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt64Type(buf *bytes.Buffer) (temp int64) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeStringType(buf *bytes.Buffer) string {
	var length uint16

	binary.Read(buf, binary.BigEndian, &length)

	strBytes := make([]byte, length)
	buf.Read(strBytes)

	return string(strBytes[:])
}

func decodeBooleanType(buf *bytes.Buffer) (bool, error) {
	var value uint8

	binary.Read(buf, binary.BigEndian, &value)

	if value == 0 {
		return false, nil
	} else if value == 1 {
		return true, nil
	} else {
		return false, fmt.Errorf("Invalid value for boolean of %d", value)
	}

}

func decodeSliceInt8Type(buf *bytes.Buffer) ([]int8, error) {
	var length uint32

	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	array := make([]int8, length)

	for j := 0; j < int(length); j++ {
		var temp int8
		err := binary.Read(buf, binary.BigEndian, &temp)
		if err != nil {
			return nil, err
		}
		array[j] = temp
	}

	return array, nil
}

func decodeDictionaryType(buf *bytes.Buffer) (map[interface{}]interface{}, error) {
	var keyTypeCode uint8
	var valueTypeCode uint8
	var dictionarySize uint16

	err := binary.Read(buf, binary.BigEndian, &keyTypeCode)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.BigEndian, &valueTypeCode)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.BigEndian, &dictionarySize)
	if err != nil {
		return nil, err
	}

	// TODO: The map[interface{}]interface{} may not actually work in real use-cases
	dictionary := make(map[interface{}]interface{})
	for i := uint16(0); i < dictionarySize; i++ {
		// TODO: We may need to read another byte for either key or value if they equal 0 or 42 in order to determine actual type
		key := decodeType(buf, keyTypeCode)
		value := decodeType(buf, valueTypeCode)
		dictionary[key] = value
	}

	return dictionary, nil
}
