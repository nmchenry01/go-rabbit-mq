package models

import "encoding/xml"

type signature struct {
	SignatureValue string
}

type headSignature struct {
	Signature signature `xml:"Signature"`
}

type appHdr struct {
	BizMsgIdr     string        `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.01 BizMsgIdr"`
	MsgDefIdr     string        `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.01 MsgDefIdr"`
	CreDt         string        `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.01 CreDt"`
	HeadSignature headSignature `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.01 Sgntr"`
}

type Pacs008Message struct {
	XMLName xml.Name `xml:"Message"`
	AppHdr  appHdr
}
