package domain

type DeviceKind string

const (
	DeviceKindVideo DeviceKind = "video"
	DeviceKindAudio DeviceKind = "audio"
)

type Device struct {
	Name string
	Kind DeviceKind
}

type DeviceCatalog struct {
	Videos []Device
	Audios []Device
}

func (c DeviceCatalog) HasSelectableDevices() bool {
	return len(c.Videos) > 0 && len(c.Audios) > 0
}

type Selection struct {
	Video string
	Audio string
}
