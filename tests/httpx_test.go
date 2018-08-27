package tests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.xiagaogao.com/coffee/boot/errors"
	"git.xiagaogao.com/coffee/httpx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/check.v1"
)

var _ = check.Suite(&MySuite{})

func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct {
	dir     string // 测试用的临时目录
	f       string // 测试用的临时文件
	logger  *zap.Logger
	service httpx.Service
	errorService errors.Service
}

func (s *MySuite) TearDownSuite(c *check.C) {
	defer s.logger.Sync()
	s.service.Shutdown()
}

func (s *MySuite) SetUpSuite(c *check.C) {
	s.errorService = errors.NewService("test")
	logger, _ := zap.NewDevelopment()
	s.logger = logger
	service, err := httpx.NewService(&httpx.Config{}, logger)
	if err != nil {
		c.Error(err)
		c.FailNow()
	}
	s.service = service
	service.Start(func() {
		logger.Debug("服务器关闭")
	})
}

func (s *MySuite) TestRequest(c *check.C) {
	routerGroup := s.service.NewRouterGroup("/test")
	v1RouterGroup := routerGroup.Group("/v1")
	v1RouterGroup.GET("/hi", func(c *gin.Context) {
		c.String(http.StatusOK,"hi: %s","coffee")
	})
	req,_ := http.NewRequest(http.MethodGet,"/test/v1/hi",nil)
	rr :=httptest.NewRecorder()
	s.service.GetGinEngine().ServeHTTP(rr,req)
	c.Assert(rr.Code,check.Equals,http.StatusOK)
	c.Assert(rr.Body.String(),check.Equals,fmt.Sprintf("hi: %s","coffee"))
}

func (s *MySuite) TestURLParamRequest(c *check.C) {
	routerGroup := s.service.NewRouterGroup("/test")
	v1RouterGroup := routerGroup.Group("/v1")
	v1RouterGroup.GET("/:name", func(c *gin.Context) {
		c.String(http.StatusOK,"hi: %s",c.Param("name"))
	})
	req,_ := http.NewRequest(http.MethodGet,"/test/v1/coffee",nil)
	rr :=httptest.NewRecorder()
	s.service.GetGinEngine().ServeHTTP(rr,req)
	c.Assert(rr.Code,check.Equals,http.StatusOK)
	c.Assert(rr.Body.String(),check.Equals,fmt.Sprintf("hi: %s","coffee"))
}

func (s *MySuite) TestURLQueryRequest(c *check.C) {
	routerGroup := s.service.NewRouterGroup("/test")
	v1RouterGroup := routerGroup.Group("/v1")
	now := time.Now().Format("2006-01-02 15:04:05.999999999")
	v1RouterGroup.GET("/:name", func(c *gin.Context) {
		c.String(http.StatusOK,"hi: %s\n%s",c.Param("name"),c.Query("time"))
	})
	param := "coffee"
	req,_ := http.NewRequest(http.MethodGet,fmt.Sprintf("/test/v1/%s?time=%s",param,now),nil)
	rr :=httptest.NewRecorder()
	s.service.GetGinEngine().ServeHTTP(rr,req)
	c.Assert(rr.Code,check.Equals,http.StatusOK)
	c.Assert(rr.Body.String(),check.Equals,fmt.Sprintf("hi: %s\n%s",param,now))
}

func (s *MySuite) TestJSONRequest(c *check.C) {
	routerGroup := s.service.NewRouterGroup("/test")
	v1RouterGroup := routerGroup.Group("/v1")
	now := time.Now().Format("2006-01-02 15:04:05.999999999")
	v1RouterGroup.POST("/:name", func(c *gin.Context) {
		c.String(http.StatusOK,"hi: %s\n%s",c.Param("name"),c.Query("time"))
	})
	param := "coffee"
	req,_ := http.NewRequest(http.MethodGet,fmt.Sprintf("/test/v1/%s?time=%s",param,now),nil)
	rr :=httptest.NewRecorder()
	s.service.GetGinEngine().ServeHTTP(rr,req)
	c.Assert(rr.Code,check.Equals,http.StatusOK)
	c.Assert(rr.Body.String(),check.Equals,fmt.Sprintf("hi: %s\n%s",param,now))
}
