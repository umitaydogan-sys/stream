package rtmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// AMF0 type markers
const (
	AMF0Number    = 0x00
	AMF0Boolean   = 0x01
	AMF0String    = 0x02
	AMF0Object    = 0x03
	AMF0Null      = 0x05
	AMF0Undefined = 0x06
	AMF0ECMAArray = 0x08
	AMF0EndObject = 0x09
)

// AMFValue represents a decoded AMF value
type AMFValue struct {
	Type   byte
	Str    string
	Num    float64
	Bool   bool
	Obj    map[string]AMFValue
	Array  map[string]AMFValue
}

// ReadAMF0 reads a sequence of AMF0 values from the reader
func ReadAMF0(r io.Reader) ([]AMFValue, error) {
	var values []AMFValue
	for {
		val, err := readAMF0Value(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return values, err
		}
		values = append(values, val)
	}
	return values, nil
}

func readAMF0Value(r io.Reader) (AMFValue, error) {
	marker := make([]byte, 1)
	if _, err := io.ReadFull(r, marker); err != nil {
		return AMFValue{}, err
	}

	switch marker[0] {
	case AMF0Number:
		return readAMF0Number(r)
	case AMF0Boolean:
		return readAMF0Bool(r)
	case AMF0String:
		return readAMF0String(r)
	case AMF0Object:
		return readAMF0Object(r)
	case AMF0Null:
		return AMFValue{Type: AMF0Null}, nil
	case AMF0Undefined:
		return AMFValue{Type: AMF0Undefined}, nil
	case AMF0ECMAArray:
		return readAMF0ECMAArray(r)
	default:
		return AMFValue{}, fmt.Errorf("unknown AMF0 type: 0x%02x", marker[0])
	}
}

func readAMF0Number(r io.Reader) (AMFValue, error) {
	b := make([]byte, 8)
	if _, err := io.ReadFull(r, b); err != nil {
		return AMFValue{}, err
	}
	bits := binary.BigEndian.Uint64(b)
	return AMFValue{Type: AMF0Number, Num: math.Float64frombits(bits)}, nil
}

func readAMF0Bool(r io.Reader) (AMFValue, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return AMFValue{}, err
	}
	return AMFValue{Type: AMF0Boolean, Bool: b[0] != 0}, nil
}

func readAMF0String(r io.Reader) (AMFValue, error) {
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return AMFValue{}, err
	}
	strLen := binary.BigEndian.Uint16(lenBuf)
	str := make([]byte, strLen)
	if _, err := io.ReadFull(r, str); err != nil {
		return AMFValue{}, err
	}
	return AMFValue{Type: AMF0String, Str: string(str)}, nil
}

func readAMF0Object(r io.Reader) (AMFValue, error) {
	obj := make(map[string]AMFValue)
	for {
		// Read property name
		lenBuf := make([]byte, 2)
		if _, err := io.ReadFull(r, lenBuf); err != nil {
			return AMFValue{}, err
		}
		nameLen := binary.BigEndian.Uint16(lenBuf)

		if nameLen == 0 {
			// Check for end marker
			end := make([]byte, 1)
			if _, err := io.ReadFull(r, end); err != nil {
				return AMFValue{}, err
			}
			if end[0] == AMF0EndObject {
				break
			}
			// Not end marker, this is an error
			return AMFValue{}, fmt.Errorf("expected end of object marker")
		}

		name := make([]byte, nameLen)
		if _, err := io.ReadFull(r, name); err != nil {
			return AMFValue{}, err
		}

		val, err := readAMF0Value(r)
		if err != nil {
			return AMFValue{}, err
		}
		obj[string(name)] = val
	}
	return AMFValue{Type: AMF0Object, Obj: obj}, nil
}

func readAMF0ECMAArray(r io.Reader) (AMFValue, error) {
	// Read array count (4 bytes, but we don't really use it)
	countBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, countBuf); err != nil {
		return AMFValue{}, err
	}

	// Read like object
	obj, err := readAMF0Object(r)
	if err != nil {
		return AMFValue{}, err
	}
	obj.Type = AMF0ECMAArray
	obj.Array = obj.Obj
	return obj, nil
}

// ─── AMF0 Write Functions ─────────────────────────────────

// WriteAMF0String writes an AMF0 string
func WriteAMF0String(w io.Writer, s string) error {
	if _, err := w.Write([]byte{AMF0String}); err != nil {
		return err
	}
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(len(s)))
	if _, err := w.Write(b); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

// WriteAMF0Number writes an AMF0 number
func WriteAMF0Number(w io.Writer, n float64) error {
	if _, err := w.Write([]byte{AMF0Number}); err != nil {
		return err
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(n))
	_, err := w.Write(b)
	return err
}

// WriteAMF0Bool writes an AMF0 boolean
func WriteAMF0Bool(w io.Writer, val bool) error {
	if _, err := w.Write([]byte{AMF0Boolean}); err != nil {
		return err
	}
	b := byte(0)
	if val {
		b = 1
	}
	_, err := w.Write([]byte{b})
	return err
}

// WriteAMF0Null writes an AMF0 null
func WriteAMF0Null(w io.Writer) error {
	_, err := w.Write([]byte{AMF0Null})
	return err
}

// WriteAMF0ObjectStart writes object start marker
func WriteAMF0ObjectStart(w io.Writer) error {
	_, err := w.Write([]byte{AMF0Object})
	return err
}

// WriteAMF0Property writes an object property name
func WriteAMF0Property(w io.Writer, name string) error {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(len(name)))
	if _, err := w.Write(b); err != nil {
		return err
	}
	_, err := w.Write([]byte(name))
	return err
}

// WriteAMF0ObjectEnd writes object end marker
func WriteAMF0ObjectEnd(w io.Writer) error {
	_, err := w.Write([]byte{0x00, 0x00, AMF0EndObject})
	return err
}
