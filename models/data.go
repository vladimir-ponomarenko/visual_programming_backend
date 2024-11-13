package models

import "time"

type GsmData struct {
	Registered    bool   `json:"registered"`
	Mcc           uint16 `json:"mcc"`
	Mnc           uint16 `json:"mnc"`
	Lac           int32  `json:"lac"`
	Cid           int32  `json:"cid"`
	Arfcn         uint16 `json:"arfcn"`
	Bsic          uint32 `json:"bsic"`
	Rssi          int32  `json:"rssi"`
	TimingAdvance uint32 `json:"timingAdvance"`
	BitErrorRate  uint32 `json:"bitErrorRate"`
}

type WcdmaData struct {
	Registered   bool   `json:"registered"`
	Mcc          uint16 `json:"mcc"`
	Mnc          uint16 `json:"mnc"`
	Lac          int32  `json:"lac"`
	Cid          int32  `json:"cid"`
	Rssi         int32  `json:"rssi"`
	Psc          int32  `json:"psc"`
	Uarfcn       uint16 `json:"uarfcn"`
	Rscp         int32  `json:"rscp"`
	Ecno         int32  `json:"ecno"`
	Level        uint8  `json:"level"`
	BitErrorRate uint32 `json:"bitErrorRate"`
}

type LteData struct {
	Type          string `json:"type"`
	Timestamp     int64  `json:"timestamp"`
	Registered    bool   `json:"registered"`
	Mcc           uint16 `json:"mcc"`
	Mnc           uint16 `json:"mnc"`
	Ci            uint32 `json:"ci"`
	Pci           uint16 `json:"pci"`
	Tac           uint16 `json:"tac"`
	Earfcn        uint32 `json:"earfcn"`
	Bandwidth     uint32 `json:"bandwidth"`
	Rsrp          int16  `json:"rsrp"`
	Rssi          int16  `json:"rssi"`
	Rsrq          int32  `json:"rsrq"`
	Rssnr         int32  `json:"rssnr"`
	Cqi           int32  `json:"cqi"`
	TimingAdvance int32  `json:"timingAdvance"`
	Level         int    `json:"level"`
	AsuLevel      int    `json:"asuLevel"`
}

type NRData struct {
	Type       string `json:"type"`
	Registered bool   `json:"registered"`
	Mcc        string `json:"mcc"`
	Mnc        string `json:"mnc"`
	Nci        int64  `json:"nci"`
	Nrarfcn    int32  `json:"nrarfcn"`
	Pci        int16  `json:"pci"`
	Tac        uint32 `json:"tac"`
	Bands      int    `json:"bands"`
}

type Message struct {
	Time      time.Time   `json:"time"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Altitude  float64     `json:"altitude"`
	Operator  string      `json:"operator"`
	Wcdma     []WcdmaData `json:"wcdma"`
	Gsm       []GsmData   `json:"gsm"`
	Lte       []LteData   `json:"lte"`
	Nr        []NRData    `json:"nr"`
}
