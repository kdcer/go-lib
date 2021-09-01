package goftp

import (
	"errors"
	"fmt"
	"testing"

	"github.com/silenceper/wechat/v2/officialaccount/material"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/guuid"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/silenceper/wechat/v2/officialaccount/message"

	"github.com/gogf/gf/net/ghttp"

	"github.com/gogf/gf/frame/g"
	"github.com/kdcer/go-lib/lib/wechat"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

func Test_Example(t *testing.T) {
	wechatInit()

	res, err := wechatService.GetAccessToken()
	t.Log(res, err)

	// 直接将r交给wechat库处理
	var r *ghttp.Request // 模拟r
	wechatService.Serve(r)

	err = wechatService.Manager(r)
	t.Log(err)

	code, err := wechatService.Auth(r)
	t.Log(code, err)

	/*
	  Controller调用方式
	  type WechatController struct {
	  	*BaseController
	  }
	  func (c *WechatController) GetAccessToken(r *ghttp.Request) {
	  	res, err := wechatService.GetAccessToken(r)
	  	if err != nil {
	  		c.Fail(r, err.Error())
	  	}
	  	c.Ok(r, res)
	  }

	  func (c *WechatController) Serve(r *ghttp.Request) {
	  	// 直接将r交给wechat库处理
	  	wechatService.Serve(r)
	  }

	  func (c *WechatController) Manager(r *ghttp.Request) {
	  	err := wechatService.Manager(r)
	  	if err != nil {
	  		c.Fail(r, err.Error())
	  	}
	  	c.Ok(r, nil)
	  }

	  // Auth 授权登录 获取code返回给客户端url
	  func (c *WechatController) Auth(r *ghttp.Request) {
	  	code, err := wechatService.Auth(r)
	  	if err != nil {
	  		c.Error(r, err)
	  	}
	  	c.Ok(r, code)
	  }
	*/
}

func Test_Pay(t *testing.T) {
	wechatInit()
	var payService = wechat.NewPayService(&wechat.RedPacketConfig{
		WeixinPayKey:            g.Config().GetString("wechat.sendredpack.weixinPayKey"),
		WeixinMchID:             g.Config().GetString("wechat.sendredpack.weixinMchID"),
		WeixinAppID:             g.Config().GetString("wechat.appId"),
		WeixinClientCertPemPath: g.Config().GetString("wechat.sendredpack.weixinClientCertPemPath"),
		WeixinClientKeyPemPath:  g.Config().GetString("wechat.sendredpack.weixinClientKeyPemPath"),
		WeixinRootCaPath:        g.Config().GetString("wechat.sendredpack.weixinRootCaPath"),
		ClientIP:                g.Config().GetString("wechat.sendredpack.clientIP"),
	})
	redPacketRequest := &wechat.RedPacketRequest{
		ActName:     "红包活动",
		ReOpenid:    "user.OpenId",
		Remark:      "快来参与哦",
		SendName:    "上海坤鼎",
		TotalAmount: 100,
		Wishing:     "感谢您的参与！",
		SceneID:     "PRODUCT_2",
	}
	rsp, err := payService.SendRedPack(redPacketRequest)
	if err != nil {
		fmt.Println(rsp, err)
		return
	}

	if rsp.ReturnCode != "SUCCESS" || rsp.ResultCode != "SUCCESS" {
		t.Log(err)
	}
}

// 在boot.go中初始化wechat
func wechatInit() {
	wechat.InitWechat(&cache.RedisOpts{
		Host:        g.Config().GetString("redis.addr"),
		Password:    g.Config().GetString("redis.password"),
		Database:    0,
		MaxActive:   2000,
		MaxIdle:     10,
		IdleTimeout: 300,
	})
	wechat.InitOfficialAccount(&offConfig.Config{
		AppID:          g.Config().GetString("wechat.appId"),
		AppSecret:      g.Config().GetString("wechat.appSecret"),
		Token:          g.Config().GetString("wechat.token"),
		EncodingAESKey: g.Config().GetString("wechat.encodingAESKey"),
	})
	wechat.InitMiniProgram(&miniConfig.Config{
		AppID:     g.Config().GetString("wechat.appId"),
		AppSecret: g.Config().GetString("wechat.appSecret"),
	})
}

var wechatService *WechatService

type WechatService struct {
}

// GetAccessToken AccessToken
// 获取AccessToken，默认过期时间是7200秒，缓存设置过期时间为小于7200秒(目前是7200-1500)，使用时永远不会失效。刷新access_token时，在5分钟内，新老access_token都可用
func (s *WechatService) GetAccessToken() (res string, err error) {
	tk, err := wechat.OfficialAccount.GetAccessToken()
	return tk, err
}

// Serve 接口配置信息验证和被动接受用户发送给公众号的消息并处理和回复
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html
// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Receiving_standard_messages.html
func (s *WechatService) Serve(r *ghttp.Request) (err error) {
	server := wechat.OfficialAccount.GetServer(r.Request, r.Response.Writer)
	//关闭接口验证，则validate结果则一直返回true
	//server.SkipValidate(true)

	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content + "666")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err = server.Serve()
	if err != nil {
		glog.Error("Serve Error, err=%+v", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		glog.Error("Send Error, err=%+v", err)
		return
	}
	return err
}

// Manager 主动给用户发送消息
// 当用户和公众号产生特定动作的交互时,开发者可以在一段时间内调用客服接口 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html
func (s *WechatService) Manager(r *ghttp.Request) (err error) {
	manager := wechat.OfficialAccount.GetCustomerMessageManager()
	manager.Send(message.NewCustomerTextMessage("openid", "666"))
	return err
}

// Example Example
func (s *WechatService) Example(r *ghttp.Request) (err error) {

	// -----------------------------------小程序-----------------------------------------
	// 小程序登录
	res, err := wechat.MiniProgram.GetAuth().Code2Session("code")
	fmt.Println(res)
	// 小程序发送订阅消息
	wechat.MiniProgram.GetSubscribe().Send(&subscribe.Message{
		ToUser:     "openId",
		TemplateID: "templateId",
		Page:       "page",
		Data: map[string]*subscribe.DataItem{
			"thing2": {
				Value: "title",
				Color: "",
			}},
		MiniprogramState: "",
		Lang:             "",
	})
	// 获取小程序码，适用于需要的码数量极多的业务场景
	wechat.MiniProgram.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{
		Page:      "pages/index/index",
		Path:      "",
		Width:     0,
		Scene:     "id=1",
		AutoColor: false,
		LineColor: qrcode.Color{},
		IsHyaline: true,
	})
	// 小程序解密数据(用户信息/手机号信息)
	wechat.MiniProgram.GetEncryptor().Decrypt("sessionKey", "encryptedData", "iv")
	// -----------------------------------小程序-----------------------------------------

	// -----------------------------------公众号-----------------------------------------
	// 返回临时素材的下载地址供用户自己处理
	// URL 不可公开，因为含access_token 需要立即另存文件,AccessToken
	mediaURL, err := wechat.OfficialAccount.GetMaterial().GetMediaURL("mediaID")
	mediaBt := g.Client().GetBytes(mediaURL) // 获取图片byte
	fmt.Println(len(mediaBt))

	// GetWxConfig
	wechat.OfficialAccount.GetJs().GetConfig("url")

	// 授权两种模式
	// 1.snsapi_userinfo模式 可以调用sns/userinfo接口获取用户信息，但是授权后分享的菜单被隐藏
	// 2.snsapi_base模式可以分享，但是不能调用sns/userinfo接口获取用户信息，只能通过GetUserAccessToken获取openid去创建用户
	result, err := wechat.OfficialAccount.GetOauth().GetUserAccessToken("code")
	// 授权snsapi_userinfo获取用户授权后分享的菜单被隐藏，需要把auth接口把授权从snsapi_userinfo改成snsapi_base，然后这里就无法调用sns/userinfo接口，直接给定openid去创建用户
	userInfo, err := wechat.OfficialAccount.GetOauth().GetUserInfo(result.AccessToken, result.OpenID, "")
	fmt.Println(userInfo)
	// -----------------------------------公众号-----------------------------------------

	// -----------------------------------素材库-----------------------------------------
	// 上传永久素材，filename传文件地址
	wechat.OfficialAccount.GetMaterial().AddMaterial(material.MediaTypeImage, "./image/1.png")
	// -----------------------------------素材库-----------------------------------------
	return err
}

// Auth 授权登录 获取code返回给客户端url
func (s *WechatService) Auth(r *ghttp.Request) (code string, err error) {
	code = r.GetString("code")
	var scope string
	secret := r.GetString("secret")
	// 获取基础信息snsapi_base   获取用户信息，但是无法显示分享菜单snsapi_userinfo
	if secret == "userinfo" {
		scope = "snsapi_userinfo"
	} else if secret == "base" {
		scope = "snsapi_base"
	} else {
		return "", errors.New("err")
	}
	if len(code) == 0 {
		proto := "http://"
		if r.Request.TLS != nil {
			proto = "https://"
		}
		authURL := fmt.Sprintf("%s%s%s", proto, r.Request.Host, r.RequestURI)

		//redirectUrl, _ := wechat.OfficialAccount.GetOauth().GetRedirectURL(authURL, scope, guuid.New().String())
		//r.Response.RedirectTo(redirectUrl)
		// 第二种方式
		wechat.OfficialAccount.GetOauth().Redirect(r.Response.Writer, r.Request, authURL, scope, guuid.New().String())
		r.Response.Request.Exit() // gf框架需要手动结束
	}
	glog.Info(code)
	return code, nil
}
