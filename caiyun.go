package colorful_clouds_api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type colorfulCloudsApi struct {
	version string // API 版本
	lang    string // 语言选项
	unit    string // 单位制选项

	realtimeUrl string // 实况天气接口 url

	weather weather
}

type weather struct {
	weatherUrl  string // 通用预报接口 url
	hourlySteps int    // 小时步长选项
	dailySteps  int    // 天步长选项
	alert       bool   // 预警信息选项
}

func NewColorfulCloudsApi(token string, longitudeAndLatitude string, opts ...Option) ColorfulCloudsApi {
	cy := new(colorfulCloudsApi)
	for _, opt := range opts {
		opt(cy)
	}
	if cy.version == "" {
		Version("v2.5")(cy)
	}
	if cy.lang == "" {
		LangZhCN()(cy)
	}
	if cy.unit == "" {
		UnitMetric()(cy)
	}
	cy.realtimeUrl =
		fmt.Sprintf(RealtimeUrl, cy.version, token, longitudeAndLatitude) +
			"?lang=" + cy.lang +
			"&unit=" + cy.unit

	cy.weather.weatherUrl =
		fmt.Sprintf(WeatherUrl, cy.version, token, longitudeAndLatitude) +
			"?lang=" + cy.lang +
			"&unit=" + cy.unit

	return cy
}

func (cy colorfulCloudsApi) Realtime() (reply RealtimeReply, err error) {
	var rawResp []byte
	if rawResp, err = executeGet(cy.realtimeUrl); err != nil {
		return RealtimeReply{}, err
	}
	if err = json.Unmarshal(rawResp, &reply); err != nil {
		return RealtimeReply{}, fmt.Errorf("unexpected response: %w\nraw response: %s", err, rawResp)
	}

	return
}

func (cy colorfulCloudsApi) Weather(opts ...WeatherOption) (reply WeatherReply, err error) {
	for _, opt := range opts {
		opt(&cy)
	}

	rawUrl := cy.weather.weatherUrl
	var tmpURL *url.URL
	if tmpURL, err = url.Parse(rawUrl); err != nil {
		return WeatherReply{}, err
	}

	query := tmpURL.Query()
	if cy.weather.hourlySteps != 0 {
		query.Set("hourlysteps", strconv.Itoa(cy.weather.hourlySteps))
	}
	if cy.weather.dailySteps != 0 {
		query.Set("dailysteps", strconv.Itoa(cy.weather.dailySteps))
	}
	query.Set("alert", strconv.FormatBool(cy.weather.alert))
	tmpURL.RawQuery = query.Encode()

	var rawResp []byte
	if rawResp, err = executeGet(tmpURL.String()); err != nil {
		return WeatherReply{}, err
	}
	if err = json.Unmarshal(rawResp, &reply); err != nil {
		return WeatherReply{}, fmt.Errorf("unexpected response: %w\nraw response: %s", err, rawResp)
	}

	return
}

// Option 通用设置
type Option func(cy *colorfulCloudsApi)

func Version(version string) Option {
	return func(cy *colorfulCloudsApi) {
		cy.version = version
	}
}

// LangZhCN 简体中文
func LangZhCN() Option {
	return func(cy *colorfulCloudsApi) {
		cy.lang = "zh_CN"
	}
}

// LangZhTW 繁体中文
func LangZhTW() Option {
	return func(cy *colorfulCloudsApi) {
		cy.lang = "zh_TW"
	}
}

// LangEnUS 美式英语
func LangEnUS() Option {
	return func(cy *colorfulCloudsApi) {
		cy.lang = "en_US"
	}
}

// LangEnGB 英式英语
func LangEnGB() Option {
	return func(cy *colorfulCloudsApi) {
		cy.lang = "en_GB"
	}
}

// LangJa 日语
func LangJa() Option {
	return func(cy *colorfulCloudsApi) {
		cy.lang = "ja"
	}
}

// UnitMetric 公制 `metric`
func UnitMetric() Option {
	return func(cy *colorfulCloudsApi) {
		cy.unit = "metric"
	}
}

// UnitImperial 英制 `imperial`
func UnitImperial() Option {
	return func(cy *colorfulCloudsApi) {
		cy.unit = "imperial"
	}
}

// UnitSI 科学单位制 `SI`
func UnitSI() Option {
	return func(cy *colorfulCloudsApi) {
		cy.unit = "SI"
	}
}

// WeatherOption 通用预报接口设置
type WeatherOption func(cy *colorfulCloudsApi)

// HourlySteps 小时步长选项: 可选, 缺省值是 `48`，选择范围 `1 ~ 360`
func HourlySteps(steps ...int) WeatherOption {
	if len(steps) == 0 {
		steps = []int{48}
	}
	if steps[0] < 1 {
		steps[0] = 1
	}
	if steps[0] > 360 {
		steps[0] = 360
	}
	return func(cy *colorfulCloudsApi) {
		cy.weather.hourlySteps = steps[0]
	}
}

// DailySteps 天步长选项: 可选, 缺省值是 `5`，选择范围 `1 ~ 15`
func DailySteps(steps ...int) WeatherOption {
	if len(steps) == 0 {
		steps = []int{5}
	}
	if steps[0] < 1 {
		steps[0] = 1
	}
	if steps[0] > 15 {
		steps[0] = 15
	}
	return func(cy *colorfulCloudsApi) {
		cy.weather.dailySteps = steps[0]
	}
}

// Alert 预警信息选项: 可选, 缺省值是 `false`
func Alert(steps ...bool) WeatherOption {
	if len(steps) == 0 {
		steps = []bool{false}
	}
	return func(cy *colorfulCloudsApi) {
		cy.weather.alert = steps[0]
	}
}
