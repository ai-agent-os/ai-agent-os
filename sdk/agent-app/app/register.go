package app

import (
	"encoding/json"
	"errors"
	"github.com/ai-agent-os/ai-agent-os/pkg/logger"
	"github.com/ai-agent-os/ai-agent-os/sdk/agent-app/callback"
	"github.com/ai-agent-os/ai-agent-os/sdk/agent-app/response"
	"strings"
)

func routerKey(router string, method string) string {
	router = strings.Trim(router, "/")
	key := router + "." + method
	return key
}

func register(router string, method string, handleFunc HandleFunc, templater Templater) {
	if app == nil {
		initApp()
	}
	//router = strings.Trim(router, "/")
	//key := router + "." + method
	app.routerInfo[routerKey(router, method)] = &routerInfo{
		HandleFunc: handleFunc,
		Router:     router,
		Method:     method,
		Template:   templater,
	}
}

func GET(router string, handleFunc HandleFunc, templater Templater) {
	register(router, "GET", handleFunc, templater)
}

func POST(router string, handleFunc HandleFunc, templater Templater) {
	register(router, "POST", handleFunc, templater)
}
func PUT(router string, handleFunc HandleFunc, templater Templater) {
	register(router, "PUT", handleFunc, templater)
}

func DELETE(router string, handleFunc HandleFunc, templater Templater) {
	register(router, "DELETE", handleFunc, templater)
}

func initRouter(a *App) {
	//a.registerRouter(MethodPost, "/test/add", AddHandle, Temp)
	//a.registerRouter(MethodPost, "/test/get", GetHandle, Temp)
	a.registerRouter(MethodPost, "/_callback", a.CallbackRouter, &FormTemplate{})
	a.registerRouter(MethodGet, "/_callback", a.CallbackRouter, &FormTemplate{})
	a.registerRouter(MethodDelete, "/_callback", a.CallbackRouter, &FormTemplate{})
	a.registerRouter(MethodPut, "/_callback", a.CallbackRouter, &FormTemplate{})
}

type CallbackRouterReq struct {
	Type   string `json:"type" binding:"required" example:""`
	Method string `json:"method" binding:"required" example:""`
	Router string `json:"router" binding:"required" example:"/users/app/xxxx"`
	Body   []byte `json:"body" example:"eyJpZCI6MX0="`
}

func (a *App) CallbackRouter(ctx *Context, resp response.Response) error {
	var req CallbackRouterReq
	if err := json.Unmarshal(ctx.body, &req); err != nil {
		logger.Errorf(ctx, "CallbackRouter Unmarshal body:%s err: %v", ctx.body, err)
		return err
	}

	router, err := a.getRouter(req.Router, req.Method)
	if err != nil {
		return err
	}

	//callback只是代理路由，要重定向到真正的路由
	ctx.msg.Router = req.Router
	ctx.msg.Method = req.Method
	ctx.body = req.Body

	switch req.Type {
	case CallbackTypeOnTableAddRow:
		v, ok := router.Template.(*TableTemplate)
		if !ok {
			return errors.New("invalid type of TableTemplate")
		}
		var onTableReq callback.OnTableAddRowReq
		onTableResp, err := v.OnTableAddRow(ctx, &onTableReq)
		if err != nil {
			logger.Errorf(ctx, "callback onTableAddRow router:%s call error:%s", req.Type, err.Error())
			return err
		}
		err = resp.Form(onTableResp).Build()
		if err != nil {
			logger.Errorf(ctx, "callback onTableAddRow  router:%s Build error:%s", req.Type, err.Error())
			return err
		}
		logger.Infof(ctx, "CallbackRouter onTableAddRow success")
		return nil
	case CallbackTypeOnTableUpdateRow:
		v, ok := router.Template.(*TableTemplate)
		if !ok {
			return errors.New("invalid type of TableTemplate")
		}
		var onTableReq callback.OnTableUpdateRowReq
		err := json.Unmarshal(ctx.body, &onTableReq.Updates)
		if err != nil {
			return err
		}
		onTableResp, err := v.OnTableUpdateRow(ctx, &onTableReq)
		if err != nil {
			return err
		}
		err = resp.Form(onTableResp).Build()
		if err != nil {
			logger.Errorf(ctx, "callback OnTableUpdateRows router:%s error:%s", req.Type, err.Error())
			return err
		}
		logger.Infof(ctx, "CallbackRouter OnTableUpdateRows success")
		return nil
	case CallbackTypeOnTableDeleteRows:
		v, ok := router.Template.(*TableTemplate)
		if !ok {
			return errors.New("invalid type of TableTemplate")
		}
		var onTableReq callback.OnTableDeleteRowsReq
		err := json.Unmarshal(ctx.body, &onTableReq)
		if err != nil {
			return err
		}
		onTableResp, err := v.OnTableDeleteRows(ctx, &onTableReq)
		if err != nil {
			return err
		}
		err = resp.Form(onTableResp).Build()
		if err != nil {
			logger.Errorf(ctx, "callback OnTableDeleteRows router:%s error:%s", req.Type, err.Error())
			return err
		}
		logger.Infof(ctx, "CallbackRouter OnTableDeleteRows success")
		return nil
	}
	return nil

}
