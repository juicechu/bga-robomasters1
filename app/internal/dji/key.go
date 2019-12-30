package dji

import "fmt"

type Key int

const (
	KeyNone Key = iota
	KeyAirLinkConnection
	KeyGimbalResetPosition
	KeyGimbalAngleFrontYawRotation
	KeyGimbalAngleFrontPitchRotation
	KeysCount
	// TODO(bga): Add keys here as needed.
)

type DataType int

const (
	KeyDataTypeBool DataType = iota
	KeyDataTypeLong
	KeyDataTypeAbsoluteRotationParameter
	// TODO(bga): Add data types here as needed.
)

type AccessType int

const (
	KeyAccessTypeNone AccessType = 1 << iota
	KeyAccessTypeRead
	KeyAccessTypeWrite
	KeyAccessTypeAction
)

func (k Key) DataType() DataType {
	return keyAttributeMap[k].dataType
}

func (k Key) Value() uint32 {
	return keyAttributeMap[k].value
}

func (k Key) AccessType() AccessType {
	return keyAttributeMap[k].accessType
}

type keyAttributes struct {
	value      uint32
	dataType   DataType
	accessType AccessType
}

var (
	keyAttributeMap = map[Key]keyAttributes{
		KeyAirLinkConnection: keyAttributes{117440513, KeyDataTypeBool,
			KeyAccessTypeRead},
		KeyGimbalResetPosition: keyAttributes{67108870, KeyDataTypeLong,
			KeyAccessTypeAction},
		KeyGimbalAngleFrontYawRotation: keyAttributes{67108876,
			KeyDataTypeAbsoluteRotationParameter,
			KeyAccessTypeAction},
		KeyGimbalAngleFrontPitchRotation: keyAttributes{67108877,
			KeyDataTypeAbsoluteRotationParameter,
			KeyAccessTypeAction},
		// TODO(bga): Add other attributes here as needed. Needs to be kept in
		// 	sync with existing Keys.
	}

	keyByValueMap = map[int]Key{
		117440513: KeyAirLinkConnection,
		67108870:  KeyGimbalResetPosition,
		67108876:  KeyGimbalAngleFrontYawRotation,
		67108877:  KeyGimbalAngleFrontPitchRotation,
	}
)

func keyByValue(value int) Key {
	key, ok := keyByValueMap[value]
	if !ok {
		panic(fmt.Sprintf("Can't get key for value %d.", value))
	}

	return key
}
