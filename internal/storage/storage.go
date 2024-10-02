package storage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	keyInfluxDB = "VRsOE9gM1OcQS5VFEwiR3pZERheEi4fXm909k_C7aaALGql509v3nNH4oWe5ziW70K0bXDG-P3Js-YxJfndHFA=="
)

type Storage struct{
	Client influxdb2.Client
}

func NewStorage(url string) *Storage{
	clientDB := influxdb2.NewClient(url, keyInfluxDB)
	return &Storage{
		Client: clientDB,
	}
}
