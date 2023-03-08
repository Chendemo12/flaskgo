package main

import (
	"context"
	"fmt"
	"github.com/Chendemo12/flaskgo"
	"time"
)

// Configuration 配置文件类
type Configuration struct {
	HTTP struct {
		Host string `json:"host" yaml:"host"` // API host
		Port string `json:"port" yaml:"port"` // API port
	}
}

// ServiceContext 全局服务依赖
type ServiceContext struct {
	Conf *Configuration
}

func (c *ServiceContext) Config() any { return c.Conf }

type TunnelWorkParams struct {
	flaskgo.BaseModel
	ModType     string      `json:"mod_type"`
	FecRate     string      `json:"fec_rate"`
	FecType     string      `json:"fec_type"`
	TunnelNo    int         `json:"tunnel_no" oneof:"0 1" binding:"required"`
	IfFrequency int         `json:"if_frequency" description:"中频频点"`
	SymbolRate  int         `json:"symbol_rate" description:"符号速率"`
	FreqOffset  int         `json:"freq_offset"`
	PositionCu  PositionGeo `json:"position_cu"`
	PositionRcu PositionGeo `json:"position_rcu"`
	PositionSat PositionGeo `json:"position_sat"`
	Power       float32     `json:"power"`
	Reset       bool        `json:"reset"`
}

type PositionGeo struct {
	flaskgo.BaseModel
	Longi float32 `json:"longi" binding:"required"`
	Lati  float32 `json:"lati" binding:"required"`
}

// ReturnLinkInfo 反向链路参数，仅当网管代理配置的参数与此匹配时才转发小站消息到NCC
type ReturnLinkInfo struct {
	flaskgo.BaseModel
	ModType     string            `json:"mod_type"`
	FecRate     string            `json:"fec_rate"`
	ForwardLink []ForwardLinkInfo `json:"forward_link" description:"前向链路"`
	IfFrequency int               `json:"if_frequency" description:"中频频点"`
	SymbolRate  int               `json:"symbol_rate" description:"符号速率"`
}

func (m ReturnLinkInfo) SchemaDesc() string {
	return "反向链路参数，仅当网管代理配置的参数与此匹配时才转发小站消息到NCC"
}

type ForwardLinkInfo struct {
	flaskgo.BaseModel
	IfFrequency int `json:"if_frequency" description:"中频频点"`
	SymbolRate  int `json:"symbol_rate" description:"符号速率"`
}

func (m ForwardLinkInfo) SchemaDesc() string { return "前向链路参数" }

type SimpleForm struct {
	flaskgo.BaseModel
	Name     string `json:"name" description:"姓名" validate:"required"`
	Age      int    `json:"age" description:"年龄" default:"23" gte:"50" validate:"required"`
	Content  any    `json:"content" description:"任意内容"`
	Contents []any  `json:"contents" description:"任意一些内容"`
}

func (s *SimpleForm) SchemaDesc() string { return "简单的表单" }

type Step struct {
	Click string `json:"click"`
}

type Action struct {
	OneStep  Step     `json:"one_step"`
	TwoSteps []Step   `json:"two_steps"`
	Next     []string `json:"next"`
}

type ExampleForm struct {
	flaskgo.BaseModel
	Name    string    `json:"name"`
	Action  *Action   `json:"action"`
	Actions []*Action `json:"actions"`
}

func makeTunnelWork(s *flaskgo.Context) *flaskgo.Response {
	p := TunnelWorkParams{}
	if err := s.ShouldBindJSON(&p); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 200) // 休眠200ms,模拟设置硬件时长
	return s.OKResponse(p.TunnelNo)
}

func setNccReturnLinks(s *flaskgo.Context) *flaskgo.Response {
	p := make([]ReturnLinkInfo, 0)
	if err := s.ShouldBindJSON(&p); err != nil {
		return err
	}
	return s.OKResponse(s)
}

func getSimpleFrom(s *flaskgo.Context) *flaskgo.Response {
	s.Logger().Info("query fields: ", s.QueryFields)
	s.Logger().Info("path fields: ", s.PathFields)

	form := SimpleForm{}
	resp := s.ShouldBindJSON(&form)
	if resp != nil {
		return resp
	}

	return s.OKResponse(form)
}

func getExampleForm(s *flaskgo.Context) *flaskgo.Response {

	time.Sleep(10 * time.Second)
	//return s.OKResponse(&ExampleForm{Name: s.PathFields["name"]})

	return s.OKResponse(SimpleForm{
		Name:     s.PathFields["name"],
		Age:      0,
		Content:  nil,
		Contents: nil,
	})
}

func makeRouter() *flaskgo.Router {
	router := flaskgo.APIRouter("/api/device", []string{"Tunnel"})
	{
		router.POST(
			"/simple/:name/:age?",
			&SimpleForm{}, &SimpleForm{},
			"提交一个个人信息表单",
			getSimpleFrom,
		)

		router.GET("/form/:name", &ExampleForm{}, "获得一个随机表单", getExampleForm)

		router.POST("/tunnel/:no", &TunnelWorkParams{}, flaskgo.Int, "设置通道工作参数", makeTunnelWork).
			SetDescription("设置通道的工作参数，表单内部的`tunnel_no`必须与路径参数保持一致")

		router.POST(
			"/ncc/return_links",
			flaskgo.L(&ReturnLinkInfo{}), &ReturnLinkInfo{}, "设置NCC反向链路参数",
			setNccReturnLinks,
		)
	}
	return router
}

// Clock 定时任务
type Clock struct {
	flaskgo.CronJobFunc
}

func (c *Clock) String() string          { return "Clock" }
func (c *Clock) Interval() time.Duration { return time.Second * 5 }

func (c *Clock) Do(ctx context.Context) error {
	fmt.Println("current second:", time.Now().Second())
	return nil
}

func ExampleFlaskGo_App() {
	conf := &Configuration{}
	conf.HTTP.Host = "0.0.0.0"
	conf.HTTP.Port = "8088"
	ctx := &ServiceContext{Conf: conf}

	app := flaskgo.NewFlaskGo("FlaskGo Example", "0.2.1", true, ctx)
	app.DisableMultipleProcess().
		EnableDumpPID().
		DisableRequestValidate().
		SetShutdownTimeout(5).
		IncludeRouter(makeRouter()).
		SetDescription("一个简单的FlaskGo应用程序,在启动app之前首先需要创建并替换ServiceContext,最后调用Run来运行程序").
		AddCronjob(&Clock{})

	app.OnEvent("startup", func() { app.Service().Logger().Info("current pid: ", app.PID()) })
	app.OnEvent("startup", func() { app.Service().Logger().Info("startup event: 1") })
	app.OnEvent("startup", func() { app.Service().Logger().Info("startup event: 2") })
	app.OnEvent("shutdown", func() { app.Service().Logger().Info("shutdown event: 1") })
	app.OnEvent("shutdown", func() { app.Service().Logger().Info("shutdown event: 2") })

	app.Run(conf.HTTP.Host, conf.HTTP.Port) // 阻塞运行
}

// -----------------------------------------------------------------

func main() {
	ExampleFlaskGo_App()
}
