package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/parnurzeal/gorequest"

	"UniqueRecruitmentBackend/pkg"
)

type Client struct {
	opts *Opts
	cli  *gorequest.SuperAgent
}

type Opts struct {
	Addr string
}

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewClient(opts *Opts) (*Client, error) {
	if opts.Addr == "" {
		return nil, errors.New("addr is empty")
	}
	cli := gorequest.New()
	return &Client{
		opts: opts,
		cli:  cli,
	}, nil
}

func parseBaseResp(errs []error, resp BaseResponse) error {
	if len(errs) > 0 {
		if resp.Code == 0 {
			return errs[0]
		} else {
			return fmt.Errorf("%v: operation failed: code: %d, msg: %s", errs[0], resp.Code, resp.Msg)
		}
	}
	if resp.Code != 0 {
		return fmt.Errorf("operation failed: code: %d, msg： %s", resp.Code, resp.Msg)
	}
	return nil
}

var cookie = &http.Cookie{
	Name:   "SSO_SESSION",
	Value:  "unique_web_admin",
	Path:   "/",
	Domain: "hustunique.com",
}

func (c *Client) Ping() error {
	uri := fmt.Sprintf("%s/ping", c.opts.Addr)

	type resp struct {
		BaseResponse
	}
	var r resp
	_, _, errs := c.cli.Get(uri).EndStruct(&r)

	return parseBaseResp(errs, r.BaseResponse)
}

func (c *Client) CreateRecruitment(opts *pkg.CreateRecOpts) (*pkg.Recruitment, error) {
	uri := fmt.Sprintf("%s/recruitments/", c.opts.Addr)

	type resp struct {
		BaseResponse
		Data *pkg.Recruitment
	}
	var r resp
	re, _, errs := c.cli.Post(uri).AddCookies([]*http.Cookie{cookie}).Send(*opts).EndStruct(&r)
	log.Print(re)
	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) ListRecruitment(rid string) (*pkg.Recruitment, error) {
	uri := fmt.Sprintf("%s/recruitments/%s", c.opts.Addr, rid)

	type resp struct {
		BaseResponse
		Data *pkg.Recruitment
	}
	var r resp
	_, _, errs := c.cli.Get(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) GetLastestRecruitment() (*pkg.Recruitment, error) {
	uri := fmt.Sprintf("%s/recruitments/pending", c.opts.Addr)

	type resp struct {
		BaseResponse
		Data *pkg.Recruitment
	}
	var r resp
	_, _, errs := c.cli.Get(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) CreateApplication(opts *pkg.CreateAppOpts) (*pkg.Application, error) {
	uri := fmt.Sprintf("%s/applications/", c.opts.Addr)

	type resp struct {
		BaseResponse
		Data *pkg.Application
	}
	var r resp
	_, _, errs := c.cli.Post(uri).AddCookies([]*http.Cookie{cookie}).Send(*opts).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) ListApplication(aid string) ([]*pkg.Application, error) {
	uri := fmt.Sprintf("%s/applications/%s", c.opts.Addr, aid)

	type resp struct {
		BaseResponse
		Data []*pkg.Application
	}
	var r resp
	_, _, errs := c.cli.Get(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) CreateInterview(opts []pkg.CreateInterviewOpts, rid string, group pkg.Group) error {
	uri := fmt.Sprintf("%s/recruitments/%s/interviews/%s", c.opts.Addr, rid, group)

	type resp struct {
		BaseResponse
		Data *pkg.Interview
	}
	var r resp
	_, _, errs := c.cli.Post(uri).AddCookies([]*http.Cookie{cookie}).Send(opts).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteInterviews(opts []pkg.DeleteInterviewOpts, rid string, group pkg.Group) error {
	uri := fmt.Sprintf("%s/recruitments/%s/interviews/%s", c.opts.Addr, rid, group)

	type resp struct {
		BaseResponse
		Data *pkg.Interview
	}
	var r resp
	_, _, errs := c.cli.Delete(uri).AddCookies([]*http.Cookie{cookie}).Send(opts).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return err
	}
	return nil
}

func (c *Client) ListInterviews(rid string, group pkg.Group) ([]pkg.Interview, error) {
	uri := fmt.Sprintf("%s/recruitments/%s/interviews/%s", c.opts.Addr, rid, group)

	type resp struct {
		BaseResponse
		Data []pkg.Interview
	}
	var r resp
	_, _, errs := c.cli.Get(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) SetApplicationInterviewTime(aid string, interviewType pkg.GroupOrTeam) error {
	uri := fmt.Sprintf("%s/applications/%s/interviews/%s", c.opts.Addr, aid, interviewType)
	type resp struct {
		BaseResponse
		Data []pkg.Interview
	}
	var r resp
	_, _, errs := c.cli.Put(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return err
	}
	return nil
}

func (c *Client) SelectInterviewSlots(aid string, interviewType pkg.GroupOrTeam) error {
	uri := fmt.Sprintf("%s/applications/%s/slots/%s", c.opts.Addr, aid, interviewType)
	type resp struct {
		BaseResponse
		Data []pkg.Interview
	}
	var r resp
	_, _, errs := c.cli.Put(uri).AddCookies([]*http.Cookie{cookie}).EndStruct(&r)

	if err := parseBaseResp(errs, r.BaseResponse); err != nil {
		return err
	}
	return nil
}
