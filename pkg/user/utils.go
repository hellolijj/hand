package user

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitlab.alibaba-inc.com/zhoushua.ljj/hand/pkg/dd"
	"io/ioutil"
	"net/http"
	"net/url"
	//dd "server/common/dd"
	//jwt "server/common/jwt"
	//models "server/models"
	//table_user "server/models/table/user"
	"time"
	"github.com/prometheus/common/log"
)

type DdLogin struct {
}

type dDLoginReq struct {
	TmpAuthCode string `json:"tmp_auth_code"` //临时授权码
}
type RetDDUserInfo struct {
	Nick                 string `json:"nick"`                     //钉钉昵称
	Unionid              string `json:"unionid"`                  //Unionid
	DingId               string `json:"dingId"`                   //DingId
	Openid               string `json:","openid"`                 //Openid
	MainOrgAuthHighLevel bool   `json:"main_org_auth_high_level"` //MainOrgAuthHighLevel
}

type DDResp struct {
	Errcode  int64         `json:"errcode"`   //错误码
	Errmsg   string        `json:"errmsg"`    //错误信息
	UserInfo RetDDUserInfo `json:"user_info"` //用户信息
}

type RetddModule struct {
	Id   int64  `json:"id"`
	Name string `json:"name"` //模块名字
	Op   int64  `json:"op"`   //模块权限定义值
}
type dDLoginResp struct {
	StatusCode int    `json:"status_code"` //状态码
	StatusMsg  string `json:"status_msg"`  //状态信息
	UserId     int64  `json:"user_id"`     //用户id
	Name       string `json:"name"`        //用户名称
	RoleId     int64  `json:"role_id"`     //用户所属角色id
	RoleName   string `json:"role_name"`   //用户所属角色名称

	Phone     string        `json:"phone"`
	Email     string        `json:"email"`
	Avatar    string        `json:"avatar"`
	JobNumber string        `json:"job_number"`
	Token     string        `json:"token"`   //登录token
	Modules   []RetddModule `json:"modules"` //模块定义
}

// @Summary  钉钉登录接口
// @Description 无
// @Tags user
// @Accept  json
// @Produce json
// @Param 请求体 body user.dDLoginReq true "请求体"
// @Success 200 {object} user.dDLoginResp "返回体"
// @Router /algorithm_platform_api/v1/user/dd_login [post]

func DdLoginHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var resp dDLoginResp
	str_code := query.Get("code")
	timestamp := time.Now().UnixNano() / 1e6
	strTimeStamp := fmt.Sprintf("%d", timestamp)

	appKey := dd.DdConf.AppKey // 读取配置文件 appid
	appSecret := dd.DdConf.AppSecret
	signature := ComputeHmacSha256(strTimeStamp, appSecret) //签名
	signature = url.QueryEscape(signature)
	//post请求提交json数据
	var ddreq dDLoginReq
	ddreq.TmpAuthCode = str_code
	ba, _ := json.Marshal(ddreq)
	targetUrl := fmt.Sprintf("%s?accessKey=%s&timestamp=%d&signature=%s", dd.DdConf.DDServerAddress, appKey, timestamp, signature)
	log.Infof("target url %v", targetUrl)

	tr := &http.Transport{
		//"把从服务器传过来的非叶子证书，添加到中间证书的池中，使用设置的根证书和中间证书对叶子证书进行验证。"
		//TLSClientConfig: &tls.Config{RootCAs: pool},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //InsecureSkipVerify用来控制客户端是否证书和服务器主机名。如果设置为true,//
		//则不会校验证书以及证书中的主机名和服务器主机名是否一致。
	}
	client := &http.Client{Transport: tr}
	resp_dingding, err := client.Post(targetUrl, "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		log.Error(err)
		resp.StatusCode = 1001
		resp.StatusMsg = "请求参数错误"
		responseObject(w, resp, http.StatusOK)
		return
	}

	defer resp_dingding.Body.Close()
	body, err := ioutil.ReadAll(resp_dingding.Body)
	if err != nil {
		log.Error(err)
		resp.StatusCode = 1001
		resp.StatusMsg = "请求参数错误"
		responseObject(w, resp, http.StatusOK)
		return
	}

	log.Infof("到这里    钉钉登录=user信息==========================:", string(body))
	log.Infof("下面是我平台系统的逻辑  ==========================:")


	var ddResp DDResp
	//解析json结构体
	json.Unmarshal([]byte(body), &ddResp)


	//先查找钉钉用户表,用Unionid查找,找到返回信息,找不到插入信息

	// visit url https://blog.csdn.net/qq_33878858/article/details/106429741
}


//钉钉签名
func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func responseObject(w http.ResponseWriter, obj interface{}, statusCode int) {
	resp := &HttpRespResult{
		Code: statusCode,
		Data: obj,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type HttpRespResult struct {
	Code       int               `json:"code"`
	Message    string            `json:"message"`
	Status     string            `json:"status"`
	Data       interface{}       `json:"data,omitempty"`
}

func responseErrorMessage(w http.ResponseWriter, msg string, statusCode int) {
	resp := &HttpRespResult{
		Code:    statusCode,
		Message: msg,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}