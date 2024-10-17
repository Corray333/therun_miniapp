package types

type Element string

const (
	ElementDesert = Element("desert")
	ElementTrack  = Element("track")
	ElementCity   = Element("city")
)

var Elements = []Element{
	ElementDesert,
	ElementTrack,
	ElementCity,
}

type Round struct {
	ID      int     `json:"id"`
	EndTime int64   `json:"endTime"`
	Element Element `json:"element"`
}
