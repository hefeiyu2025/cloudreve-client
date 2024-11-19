package cloudreve

import (
	"github.com/imroc/req/v3"
	"net/http"
	"strconv"
	"time"
)

// 注意参考 https://github.com/cloudreve/Cloudreve.git
var defaultUa = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

type SessionRefresh func(session string)

type CloudreveClient struct {
	cloudreveUrl     string
	cloudreveSession string
	sessionClient    *req.Client
	defaultClient    *req.Client
	sessionRefresh   SessionRefresh
}

func NewClient(cloudreveUrl, cloudreveSession string) *CloudreveClient {
	client := &CloudreveClient{
		cloudreveUrl:  cloudreveUrl,
		sessionClient: initSessionClient(cloudreveUrl),
		defaultClient: initDefaultClient(),
	}
	client.refreshSession(cloudreveSession)
	return client
}

func NewClientWithRefresh(cloudreveUrl, cloudreveSession string, sessionRefresh SessionRefresh) *CloudreveClient {
	client := &CloudreveClient{
		cloudreveUrl:   cloudreveUrl,
		sessionClient:  initSessionClient(cloudreveUrl),
		defaultClient:  initDefaultClient(),
		sessionRefresh: sessionRefresh,
	}
	client.refreshSession(cloudreveSession)
	_, err := client.Config()
	if err != nil {
		panic(err)
	}
	return client
}

func NewClientWithLogin(cloudreveUrl, username, password string) *CloudreveClient {
	// TODO 登录获取session
	return &CloudreveClient{
		cloudreveUrl: cloudreveUrl,
	}
}

func (c *CloudreveClient) refreshSession(cloudreveSession string) *req.Client {
	c.cloudreveSession = cloudreveSession
	if c.sessionRefresh != nil {
		c.sessionRefresh(cloudreveSession)
	}
	return c.sessionClient.SetCommonCookies(&http.Cookie{Name: "cloudreve-session", Value: cloudreveSession})
}

func initSessionClient(cloudreveUrl string) *req.Client {
	sessionClient := req.C().SetCommonHeader("User-Agent", defaultUa).
		SetCommonHeader("Accept", "application/json, text/plain, */*").
		SetTimeout(30 * time.Minute).SetBaseURL(cloudreveUrl + "/api/v3")
	return sessionClient
}

func initDefaultClient() *req.Client {
	defaultClient := req.C().SetCommonHeader("User-Agent", defaultUa)
	defaultClient.GetTransport().
		WrapRoundTripFunc(func(rt http.RoundTripper) req.HttpRoundTripFunc {
			return func(req *http.Request) (resp *http.Response, err error) {
				// 由于内容长度部分是由后台计算的，所以这里需要手动设置,http默认会过滤掉header.reqWriteExcludeHeader
				if req.ContentLength <= 0 {
					if req.Header.Get("Content-Length") != "" {
						req.ContentLength, _ = strconv.ParseInt(req.Header.Get("Content-Length"), 10, 64)
					}
				}
				return rt.RoundTrip(req)
			}
		})
	return defaultClient
}
