package web

import (
	"community-sig/config"
	"community-sig/constants"
	"community-sig/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	server *http.Server
)

//启动web服务
func Run() {
	router := gin.New()
	// 使用跨域中间件
	router.Use(cors())

	v1 := router.Group("v1")
	bindRouterV1(v1)

	server = &http.Server{
		Handler: router,
	}
	router.Run(config.GetConf().Base.WebAddr) // listen and serve on 0.0.0.0:port

}

//服务退出
func Shutdown() {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}
}

// 处理跨域请求,支持options访问
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//				允许跨域设置																										可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //	跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

// v1版本鉴权处理
func filterUserAuthV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("user_id")
		token := c.GetHeader("token")
		//TODO 验证token是否有效
		if userId == "" || token == "" {
			c.JSON(200,
				model.ResponseCode{
					Code: constants.ResponseCodeAuthError,
					Desc: constants.ResponseCodeAuthError.String(),
				})
			return
		}
		// 处理请求
		c.Next()
	}
}

//读取body数据,并检查错误
func readBody(c *gin.Context, isCheckBodyIsNull bool) ([]byte, error) {
	body, err := c.GetRawData()
	if isCheckBodyIsNull && len(body) == 0 {
		c.JSON(200,
			model.ResponseCode{
				Code: constants.ResponseCodeBodyIsNull,
				Desc: constants.ResponseCodeBodyIsNull.String(),
			})
	}
	return body, err
}

/**
解析请求参数
jsonBinary
	josn二进制数据
value
	json解析结构体
context
	请求上下文
isCheckParamError
	解析参数异常验证
*/
func parseParamV1(jsonBinary []byte, value interface{}, context *gin.Context, isCheckParamError bool) error {
	err := json.Unmarshal(jsonBinary, value)
	if isCheckParamError && err != nil {
		context.JSON(200,
			model.ResponseCode{
				Code: constants.ResponseCodeJsonParsError,
				Desc: constants.ResponseCodeJsonParsError.String(),
			})
	}
	return err
}

//v1版本URL与方法绑定
func bindRouterV1(v1 *gin.RouterGroup) {
	//v1.Use(filterUserAuthV1()) //进行鉴权
	v1.GET("/wsmsg", wsHandler) //注册websocket
}
