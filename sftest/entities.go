package main

import "errors"

const (
	SfTokenUrl  = "https://test.salesforce.com/services/oauth2/token"
	SfCometdUrl = "/cometd/52.0"
	SfPublishUrl = "/services/data/v52.0/sobjects/TestJoomPro__e"

	SfClientId     = "3MVG95AcBeaB55lX.v7LGIE3EZIGm7NsaAb2Jwm2U.gfY6MCW9tcxKTThYXDmoM.O8wJh3wY0m5p1.UcRwrqN"
	SfClientSecret = "876E855BB42254E8D5C4B3C79991F7CC9D235538A60183D5356A6D971C8ACBA0"
	SfUsername     = "integrationuser@joom.com.joompro"
	SfPassword     = "***"
)

var (
	ErrSalesForce = errors.New("SalesForceError")
)

type TokenRes struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	InstanceUrl string `json:"instance_url"`
}

type HandshakeReq struct {
	Version                  string   `json:"version"`
	Channel                  string   `json:"channel"`
	SupportedConnectionTypes []string `json:"supportedConnectionTypes"`
}

type HandshakeResItem struct {
	ClientId   string `json:"clientId"`
	Successful bool   `json:"successful"`
}

type SubscribeReq struct {
	ClientId     string `json:"clientId"`
	Channel      string `json:"channel"`
	Subscription string `json:"subscription"`
}

type SubscribeResItem struct {
	Successful bool `json:"successful"`
}

type ConnectReq struct {
	ClientId       string `json:"clientId"`
	Channel        string `json:"channel"`
	ConnectionType string `json:"connectionType"`
}

type ConnectResItem struct {
	Channel    string `json:"channel"`
	Successful *bool  `json:"successful,omitempty"`
}

type PublishReq struct {
	TestText string `json:"TestText__c"`
}

type PublishRes struct {
	Id      string `json:"id"`
	Success bool   `json:"success"`
}

func NewHandshakeReq() *HandshakeReq {

	connectionTypes := make([]string, 1)
	connectionTypes[0] = "long-polling"

	return &HandshakeReq{
		Version: "1.0",
		Channel: "/meta/handshake",
		SupportedConnectionTypes: connectionTypes,
	}
}

func NewSubscribeReq(clientId string) *SubscribeReq {

	return &SubscribeReq{
		ClientId:     clientId,
		Channel:      "/meta/subscribe",
		Subscription: "/event/TestJoomProResponse__e",
	}
}

func NewConnectReq(clientId string) *ConnectReq {

	return &ConnectReq{
		ClientId:       clientId,
		Channel:        "/meta/connect",
		ConnectionType: "long-polling",
	}
}
