package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"fabric-sdk/internal/service"
)

type MultiChannelController struct{}

func (c *MultiChannelController) Init(r *ghttp.Request) {
	if err := service.MultiChannel.Init(r.Context()); err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "message": err.Error()})
		return
	}
	r.Response.WriteJson(g.Map{
		"code":     0,
		"message":  "MultiChannel service initialized",
		"channels": service.MultiChannel.ListChannels(),
	})
}

func (c *MultiChannelController) ListChannels(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"code":     0,
		"channels": service.MultiChannel.ListChannels(),
		"mapping": map[string]string{
			"mychannel":        "basic (asset transfer)",
			"access-channel":   "access_cc (door access control)",
			"billing-channel":  "billing_cc (utility billing)",
			"maintain-channel": "maintain_cc (equipment maintenance)",
			"esg-channel":      "esg_cc (ESG records)",
		},
	})
}

func (c *MultiChannelController) Invoke(r *ghttp.Request) {
	channelName := r.Get("channel").String()
	fn := r.Get("fn").String()
	argsJson := r.Get("args").String()

	if channelName == "" || fn == "" {
		r.Response.WriteJson(g.Map{"code": 400, "message": "channel and fn are required"})
		return
	}

	var args []string
	if argsJson != "" {
		json.Unmarshal([]byte(argsJson), &args)
	}

	result, err := service.MultiChannel.SubmitTransaction(r.Context(), channelName, fn, args...)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "message": fmt.Sprintf("invoke failed: %v", err)})
		return
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		data = string(result)
	}

	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "success",
		"channel": channelName,
		"fn":      fn,
		"data":    data,
	})
}

func (c *MultiChannelController) Query(r *ghttp.Request) {
	channelName := r.Get("channel").String()
	fn := r.Get("fn").String()
	argsJson := r.Get("args").String()

	if channelName == "" || fn == "" {
		r.Response.WriteJson(g.Map{"code": 400, "message": "channel and fn are required"})
		return
	}

	var args []string
	if argsJson != "" {
		json.Unmarshal([]byte(argsJson), &args)
	}

	result, err := service.MultiChannel.EvaluateTransaction(r.Context(), channelName, fn, args...)
	if err != nil {
		r.Response.WriteJson(g.Map{"code": 500, "message": fmt.Sprintf("query failed: %v", err)})
		return
	}

	var data interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		data = string(result)
	}

	r.Response.WriteJson(g.Map{
		"code":    0,
		"message": "success",
		"channel": channelName,
		"fn":      fn,
		"data":    data,
	})
}
