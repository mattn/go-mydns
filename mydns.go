package mydns

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

type MX struct {
	Name     string
	Priority int
}

type Record struct {
	Hostname   string
	Type       string
	Content    string
	DelegateID string
}

type DomainInfo struct {
	Domain string
	MX     []MX
	Record []Record
}

type Client struct {
	browser *browser.Browser
}

func NewClient() *Client {
	return &Client{
		browser: surf.NewBrowser(),
	}
}

func (di *DomainInfo) ToValues() url.Values {
	params := url.Values{}
	params.Set("DNSINFO[domainname]", di.Domain)
	for i, mx := range di.MX {
		params.Set(fmt.Sprintf("DNSINFO[mx][%d]", i), mx.Name)
		params.Set(fmt.Sprintf("DNSINFO[prio][%d]", i), fmt.Sprint(mx.Priority))
	}
	for i, record := range di.Record {
		params.Set(fmt.Sprintf("DNSINFO[hostname][%d]", i), record.Hostname)
		params.Set(fmt.Sprintf("DNSINFO[type][%d]", i), record.Type)
		params.Set(fmt.Sprintf("DNSINFO[content][%d]", i), record.Content)
		params.Set(fmt.Sprintf("DNSINFO[delegateid][%d]", i), record.DelegateID)
	}
	return params
}

func setString(form browser.Submittable, key string, valueref *string) bool {
	v, err := form.Value(key)
	if err != nil {
		return false
	}
	*valueref = v
	return true
}

func setInt(form browser.Submittable, key string, valueref *int) bool {
	v, err := form.Value(key)
	if err != nil {
		return false
	}
	_, err = fmt.Sscanf(v, "%d", valueref)
	if err != nil {
		return false
	}
	return true
}

func findForm(forms []browser.Submittable, key string) browser.Submittable {
	for _, form := range forms {
		if form != nil {
			if key == "" {
				return form
			}
			_, err := form.Value(key)
			if err == nil {
				return form
			}
		}
	}
	return nil
}

func (c *Client) Login(masterid, masterpwd string) error {
	err := c.browser.Open("https://www.mydns.jp/")
	if err != nil {
		return err
	}
	err = c.browser.PostForm(
		"https://www.mydns.jp/",
		url.Values{
			"MENU":      []string{"100"},
			"masterid":  []string{masterid},
			"masterpwd": []string{masterpwd},
		},
	)
	if err != nil {
		return err
	}
	form := findForm(c.browser.Forms(), "masterid")
	if form != nil {
		return errors.New("invalid masterid or password")
	}
	return nil
}

func (c *Client) FetchRecords() (*DomainInfo, error) {
	err := c.browser.Open("https://www.mydns.jp/?MENU=300")
	if err != nil {
		return nil, err
	}
	form := findForm(c.browser.Forms(), "DNSINFO[domainname]")
	if form == nil {
		return nil, errors.New("cannot fetch records")
	}

	var info DomainInfo
	setString(form, "DNSINFO[domainname]", &info.Domain)
	for i := 0; i < 8; i++ {
		var mx MX
		setString(form, fmt.Sprintf("DNSINFO[mx][%d]", i), &mx.Name)
		setInt(form, fmt.Sprintf("DNSINFO[prio][%d]", i), &mx.Priority)
		if mx.Priority <= 0 {
			mx.Priority = 10
		}
		info.MX = append(info.MX, mx)
	}

	for i := 0; i < 16; i++ {
		var record Record
		setString(form, fmt.Sprintf("DNSINFO[hostname][%d]", i), &record.Hostname)
		setString(form, fmt.Sprintf("DNSINFO[type][%d]", i), &record.Type)
		setString(form, fmt.Sprintf("DNSINFO[content][%d]", i), &record.Content)
		setString(form, fmt.Sprintf("DNSINFO[delegateid][%d]", i), &record.DelegateID)
		info.Record = append(info.Record, record)
	}
	return &info, nil
}

func (c *Client) UpdateRecords(di *DomainInfo) error {
	form := findForm(c.browser.Forms(), "DNSINFO[domainname]")
	if form == nil {
		return errors.New("cannot update records")
	}
	form.Set("DNSINFO[domainname]", di.Domain)
	for i, mx := range di.MX {
		form.Set(fmt.Sprintf("DNSINFO[mx][%d]", i), mx.Name)
		form.Set(fmt.Sprintf("DNSINFO[prio][%d]", i), fmt.Sprint(mx.Priority))
	}
	for i, record := range di.Record {
		form.Set(fmt.Sprintf("DNSINFO[hostname][%d]", i), record.Hostname)
		form.Set(fmt.Sprintf("DNSINFO[type][%d]", i), record.Type)
		form.Set(fmt.Sprintf("DNSINFO[content][%d]", i), record.Content)
		form.Set(fmt.Sprintf("DNSINFO[delegateid][%d]", i), record.DelegateID)
	}
	err := form.Submit()
	if err != nil {
		return err
	}
	form = findForm(c.browser.Forms(), "JOB")
	if form == nil {
		return errors.New("cannot update records")
	}
	return form.Submit()
}
