// Copyright (C) 2025 - Damien Dejean <dam.dejean@gmail.com>

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type plug struct {
	url string
}

type relay struct {
	IsOn           bool    `json:"ison"`
	HasTimer       bool    `json:"has_timer"`
	TimerStartedAt int     `json:"timer_started_at"`
	Duration       float32 `json:"timer_duration"`
	Remaining      float32 `json:"timer_remaining"`
	Overpower      bool    `json:"overpower"`
	Source         string  `json:"source"`
}

func newPlug(ip string) *plug {
	return &plug{
		url: fmt.Sprintf("http://%s/relay/0", ip),
	}
}

func (p *plug) isOn() (bool, error) {
	res, err := http.Get(p.url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	var r relay
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return false, err
	}
	return r.IsOn, nil
}

func (p *plug) turnOn(on bool) error {
	req, err := http.NewRequest(http.MethodGet, p.url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	if on {
		q.Add("turn", "on")
	} else {
		q.Add("turn", "off")
	}
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var r relay
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	if r.IsOn != on {
		status := "off"
		if on {
			status = "on"
		}
		return fmt.Errorf("failed to turn the relay %s", status)
	}
	return nil
}
