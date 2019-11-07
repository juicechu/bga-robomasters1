package internal

type Attribute struct {
	Creator string `xml:"creator"`

	// "YYYY/MM/DD".
	CreationDate string `xml:"creation_date"`

	// "MM/DD/YYYY HH:MM:SS TZ"
	ModifyTime string `xml:"modify_time"`

	// "00.00.0000".
	FirmwareVersionDependency string `xml:"firmware_version_dependency"`

	Title string `xml:"title"`

	// Most likely can be anything.
	Guid string `xml:"guid"`

	// "python" or "scratch".
	CodeType string `xml:"code_type"`

	// Currently empty.
	AppMinVersion string `xml:"app_min_version"`

	// Currently not present in resulting XML.
	AppMaxVersion string `xml:"app_max_version,omitempty"`

	// Currently not checked, but we know how to compute.
	Sign string `xml:"sign"`
}
