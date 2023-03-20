package huaweidns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jicjoy/huaweidns/core"
)

const APIHOST = "https://dns.%s.myhuaweicloud.com/%s"

type DnsClient struct {
	Signer     *core.Signer
	client     *http.Client
	entryPoint string
	ResCode    int
	ZoneID     string
	Method     string
	Body       io.Reader
}

type BuilderApi struct {
	config WithConfig
}
type WithConfig struct {
	AccKeyID     string
	AccKeySecret string
	Version      string
	RegionID     string
	ZoneName     string
}

func (b *BuilderApi) Build() (*DnsClient, error) {

	point := fmt.Sprintf(APIHOST, b.config.RegionID, b.config.Version)

	//res["resp"] = &[]RecordTag{}
	dns := &DnsClient{
		Signer: &core.Signer{
			Key:    b.config.AccKeyID,
			Secret: b.config.AccKeySecret,
		},
		client:     &http.Client{},
		entryPoint: point,
		//response:   res,
	}

	if len(b.config.ZoneName) > 0 {
		dns.GetZoneList(context.TODO(), b.config.ZoneName, false)
	}

	return dns, nil
}
func (b *BuilderApi) WithConfig(config *WithConfig) *BuilderApi {
	if len(config.Version) == 0 {
		config.Version = "v2"
	}
	if len(config.RegionID) == 0 {
		config.RegionID = "cn-east-2"
	}
	b.config = *config

	return b

}

func (c *DnsClient) GetZoneList(ctx context.Context, zoneName string, isAll bool) (Zones, error) {
	url := fmt.Sprintf("%s/%s", c.entryPoint, "zones")

	if len(zoneName) > 0 {
		url = fmt.Sprintf("%s?name=%s", url, zoneName)
	}

	response := &Zones{}
	err := c.ApiRequest(ctx, url, response)
	if err != nil {

		return *response, err
	}

	if isAll {
		return *response, nil
	}
	//fmt.Println(response)
	for _, v := range response.Response {

		if strings.Trim(v.Name, ".") == strings.Trim(zoneName, ".") {
			c.ZoneID = v.ID
			break
		}
	}

	return *response, nil
}

func (c *DnsClient) GetRecordLists(ctx context.Context, name string, dType string) (Recordsets, error) {
	var url string
	if len(name) == 0 {
		url = fmt.Sprintf("%s/zones/%s/recordsets", c.entryPoint, c.ZoneID)
	} else {
		url = fmt.Sprintf("%s/zones/%s/recordsets?name=%s&type=%s", c.entryPoint, c.ZoneID, name, dType)
	}

	response := &Recordsets{}
	err := c.ApiRequest(ctx, url, response)
	return *response, err
}

func (c *DnsClient) CreateRecord(ctx context.Context, rec *RecordTag) (RecordTag, error) {

	url := fmt.Sprintf("%s/zones/%s/recordsets", c.entryPoint, c.ZoneID)

	c.Method = "POST"
	recStr, err := json.Marshal(rec)
	if err != nil {

		return RecordTag{}, err
	}
	c.Body = strings.NewReader(string(recStr))
	response := &RecordTag{}

	err = c.ApiRequest(ctx, url, response)
	return *response, err

}

func (c *DnsClient) UpdateRecord(ctx context.Context, rec *RecordTag) (RecordTag, error) {
	url := fmt.Sprintf("%s/zones/%s/recordsets/%s", c.entryPoint, c.ZoneID, rec.ID)
	c.Method = "PUT"

	response := &RecordTag{}
	recStr, err := json.Marshal(rec)
	if err != nil {

		return *response, err
	}
	c.Body = strings.NewReader(string(recStr))
	err = c.ApiRequest(ctx, url, response)
	return *response, err

}
func (c *DnsClient) DeleteRecord(ctx context.Context, recId string) (Recrod, error) {
	url := fmt.Sprintf("%s/zones/%s/recordsets/%s", c.entryPoint, c.ZoneID, recId)
	c.Method = "DELETE"

	response := &Recrod{}
	err := c.ApiRequest(ctx, url, response)
	return *response, err
}
func (c *DnsClient) ApiRequest(ctx context.Context, url string, response interface{}) error {

	if len(c.Method) == 0 {
		c.Method = "GET"
	}
	r, err := http.NewRequestWithContext(ctx, c.Method, url, c.Body)

	if err != nil {
		return err
	}
	r.Header.Add("content-type", "application/json; charset=utf-8")
	r.Header.Add("x-stage", "RELEASE")
	r.Header.Add("User-Agent", "huaweicloud-usdk-go/3.0")
	c.Signer.Sign(r)
	fmt.Printf("api: %s, method:  %s,body: %+v", url, c.Method, c.Body)

	res, err := c.client.Do(r)
	if err != nil {

		return err
	}
	defer res.Body.Close()
	var buf []byte
	buf, err = io.ReadAll(res.Body)

	strBody := string(buf)
	if err != nil {
		return err
	}

	//fmt.Println(strBody)
	err = json.Unmarshal([]byte(strBody), &response)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Printf("body:%+v:  res:%+v", strBody, res)
		return fmt.Errorf("huawei DNS status: HTTP %d: URI %s", res.StatusCode, url)
	}

	return err
}
