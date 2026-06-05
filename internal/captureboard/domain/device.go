package domain

type DeviceType string

const (
	DeviceTypeVideo DeviceType = "video"
	DeviceTypeAudio DeviceType = "audio"
)

type Device struct {
	Name string
	Type DeviceType
}
