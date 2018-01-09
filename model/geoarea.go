package model


type GeoArea struct {
	Latitude  float64 `db:"lat",json:"latitude"`
	Longitude float64 `db:"lon",json:"longitude"`
	Radius	  float64 `db:"radius",json:"radius"`
}
