package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type SalesForceSyncEventSender struct {
	client http.Client

	accessToken     string
	accessTokenType string
	clientId        string
	instanceUrl     string

	channel chan string
}

func NewSalesForceSyncEventSender() *SalesForceSyncEventSender {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	sender := &SalesForceSyncEventSender{
		client: http.Client{
			Jar: jar,
		},

		channel: make(chan string, 10),
	}

	err = sender.start()
	if err != nil {
		panic(err)
	}

	return sender
}

func (s *SalesForceSyncEventSender) getCometdUrl() string {
	return s.instanceUrl + SfCometdUrl
}

func (s *SalesForceSyncEventSender) getPublishUrl() string {
	return s.instanceUrl + SfPublishUrl
}

func (s *SalesForceSyncEventSender) doRequest(url string, v interface{}) ([]byte, error) {

	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	logReq(buf)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.accessTokenType+" "+s.accessToken)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	logRes(body)

	return body, nil
}

func (s *SalesForceSyncEventSender) obtainAccessToken() error {

	res, err := s.client.PostForm(SfTokenUrl, url.Values{
		"grant_type":    {"password"},
		"client_id":     {SfClientId},
		"client_secret": {SfClientSecret},
		"username":      {SfUsername},
		"password":      {SfPassword},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	logRes(body)

	tokenRes := &TokenRes{}
	err = json.Unmarshal(body, tokenRes)
	if err != nil {
		return err
	}

	s.accessToken, s.accessTokenType, s.instanceUrl =
		tokenRes.AccessToken, tokenRes.TokenType, tokenRes.InstanceUrl
	return nil
}

func (s *SalesForceSyncEventSender) performHandshake() error {

	body, err := s.doRequest(s.getCometdUrl(), NewHandshakeReq())

	var handshakeRes []*HandshakeResItem
	err = json.Unmarshal(body, &handshakeRes)
	if err != nil {
		return err
	}

	if !handshakeRes[0].Successful {
		return ErrSalesForce
	}

	s.clientId = handshakeRes[0].ClientId
	return nil
}

func (s *SalesForceSyncEventSender) subscribeToTestJoomProResponse() error {

	body, err := s.doRequest(s.getCometdUrl(), NewSubscribeReq(s.clientId))

	var subscribeRes []*SubscribeResItem
	err = json.Unmarshal(body, &subscribeRes)
	if err != nil {
		return err
	}

	if !subscribeRes[0].Successful {
		return ErrSalesForce
	}

	return nil
}

func (s *SalesForceSyncEventSender) connect() error {

	body, err := s.doRequest(s.getCometdUrl(), NewConnectReq(s.clientId))
	if err != nil {
		return err
	}

	var connectRes []*ConnectResItem
	err = json.Unmarshal(body, &connectRes)
	if err != nil {
		return err
	}

	for _, v := range connectRes {
		if v.Channel == "/event/TestJoomProResponse__e" {
			s.channel <- "TestJoomProResponse"
		}
	}

	return nil
}

func (s *SalesForceSyncEventSender) start() error {

	if err := s.obtainAccessToken(); err != nil {
		fmt.Println("obtainAccessToken FAILED")
		return err
	}
	fmt.Println("obtainAccessToken OK")

	if err := s.performHandshake(); err != nil {
		fmt.Println("performHandshake FAILED")
		return err
	}
	fmt.Println("performHandshake OK")

	if err := s.subscribeToTestJoomProResponse(); err != nil {
		fmt.Println("subscribeToTestJoomProResponse FAILED")
		return err
	}
	fmt.Println("subscribeToTestJoomProResponse OK")

	go func() {
		for {
			if err := s.connect(); err != nil {
				fmt.Println("connect FAILED")
			}
		}
	}()

	return nil
}

func (s *SalesForceSyncEventSender) SendJoomProTestEventSync() (string, error) {

	body, err := s.doRequest(s.getPublishUrl(), &PublishReq{
		TestText: "JoomProTestEvent: " + time.Now().String(),
	})
	if err != nil {
		return "", err
	}

	publishRes := PublishRes{}
	err = json.Unmarshal(body, &publishRes)
	if err != nil {
		return "", err
	}

	if !publishRes.Success {
		return "", ErrSalesForce
	}

	res := <- s.channel
	return res, nil
}
