package colorful_clouds_api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// RealtimeUrl 实况天气接口 url
var RealtimeUrl = "https://api.caiyunapp.com/%s/%s/%s/realtime.json"

// WeatherUrl 通用预报接口 url
var WeatherUrl = "https://api.caiyunapp.com/%s/%s/%s/weather.json"

type ColorfulCloudsApi interface {
	Realtime() (reply RealtimeReply, err error)
	Weather(opts ...WeatherOption) (reply WeatherReply, err error)
}

var HTTPClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	},
}

type RealtimeReply struct {
	Status     string    `json:"status"`
	APIVersion string    `json:"api_version"` // 版本号
	APIStatus  string    `json:"api_status"`  // 版本号状态
	Lang       string    `json:"lang"`        // 语言
	Unit       string    `json:"unit"`        // 单位制
	Location   []float64 `json:"location"`    // 经纬度
	ServerTime int       `json:"server_time"` // 服务器本次返回的 `UTC` 时间戳
	TZShift    int       `json:"tzshift"`     // 时区的偏移秒数，如东八区就是 28800 秒
	Timezone   string    `json:"timezone"`
	Result     struct {
		Realtime struct {
			Status              string  `json:"status"`               // 实况模块返回状态
			Temperature         float64 `json:"temperature"`          // 温度
			ApparentTemperature float64 `json:"apparent_temperature"` // 体感温度
			Pressure            float64 `json:"pressure"`             // 气压
			Humidity            float64 `json:"humidity"`             // 相对湿度
			Cloudrate           float64 `json:"cloudrate"`            // 云量
			Skycon              string  `json:"skycon"`               // 主要天气现象
			Visibility          float64 `json:"visibility"`           // 能见度
			Dswrf               float64 `json:"dswrf"`                // 向下短波辐射通量
			Wind                struct {
				Direction float64 `json:"direction"` // 风向，单位是度。正北方向为0度，顺时针增加到360度。
				Speed     float64 `json:"speed"`     // 风速，米制下是公里每小时
			} `json:"wind"`
			Precipitation struct {
				Nearest struct {
					Status    string  `json:"status"`
					Distance  float64 `json:"distance"`  // 最近的降水带距离
					Intensity float64 `json:"intensity"` // 最近的降水带降水强度（单位为雷达降水强度）
				} `json:"nearest"`
				Local struct {
					Status     string  `json:"status"`
					Intensity  float64 `json:"intensity"`  // 本地降水强度（单位为雷达降水强度）
					Datasource string  `json:"datasource"` // 本地降水观测的数据源（radar，GFS）
				} `json:"local"`
			} `json:"precipitation"`
			AirQuality struct {
				Pm25 float64 `json:"pm25"` // pm25，质量浓度值
				Pm10 float64 `json:"pm10"` // pm10，质量浓度值
				O3   float64 `json:"o3"`   // 臭氧，质量浓度值
				No2  float64 `json:"no2"`  // 二氧化氮，质量浓度值
				So2  float64 `json:"so2"`  // 二氧化硫，质量浓度值
				Co   float64 `json:"co"`   // 一氧化碳，质量浓度值
				Aqi  struct {
					Chn float64 `json:"chn"` // AQI（中国标准）
					Usa float64 `json:"usa"` // AQI（美国标准）
				} `json:"aqi"`
				Description struct {
					Chn string `json:"chn"`
					Usa string `json:"usa"`
				} `json:"description"`
			} `json:"air_quality"`
			LifeIndex struct {
				Ultraviolet struct { // 紫外线指数及其自然语言描述
					Index float64 `json:"index"`
					Desc  string  `json:"desc"`
				} `json:"ultraviolet"`
				Comfort struct { // 舒适度指数及其自然语言描述
					Index int    `json:"index"`
					Desc  string `json:"desc"`
				} `json:"comfort"`
			} `json:"life_index"`
		} `json:"realtime"`
		Primary int `json:"primary"`
	} `json:"result"`
}

type WeatherReply struct {
	Status     string    `json:"status"`
	APIVersion string    `json:"api_version"`
	APIStatus  string    `json:"api_status"`
	Lang       string    `json:"lang"`
	Unit       string    `json:"unit"`
	TZShift    int       `json:"tzshift"`
	Timezone   string    `json:"timezone"`
	ServerTime int       `json:"server_time"`
	Location   []float64 `json:"location"`
	Result     struct {
		Realtime struct {
			Status      string  `json:"status"`
			Temperature float64 `json:"temperature"`
			Humidity    float64 `json:"humidity"`
			Cloudrate   float64 `json:"cloudrate"`
			Skycon      string  `json:"skycon"`
			Visibility  float64 `json:"visibility"`
			Dswrf       float64 `json:"dswrf"`
			Wind        struct {
				Speed     float64 `json:"speed"`
				Direction float64 `json:"direction"`
			} `json:"wind"`
			Pressure            float64 `json:"pressure"`
			ApparentTemperature float64 `json:"apparent_temperature"`
			Precipitation       struct {
				Local struct {
					Status     string  `json:"status"`
					Datasource string  `json:"datasource"`
					Intensity  float64 `json:"intensity"`
				} `json:"local"`
				Nearest struct {
					Status    string  `json:"status"`
					Distance  float64 `json:"distance"`
					Intensity float64 `json:"intensity"`
				} `json:"nearest"`
			} `json:"precipitation"`
			AirQuality struct {
				Pm25 float64 `json:"pm25"`
				Pm10 float64 `json:"pm10"`
				O3   float64 `json:"o3"`
				So2  float64 `json:"so2"`
				No2  float64 `json:"no2"`
				Co   float64 `json:"co"`
				Aqi  struct {
					Chn float64 `json:"chn"`
					Usa float64 `json:"usa"`
				} `json:"aqi"`
				Description struct {
					Chn string `json:"chn"`
					Usa string `json:"usa"`
				} `json:"description"`
			} `json:"air_quality"`
			LifeIndex struct {
				Ultraviolet struct {
					Index float64 `json:"index"`
					Desc  string  `json:"desc"`
				} `json:"ultraviolet"`
				Comfort struct {
					Index int    `json:"index"`
					Desc  string `json:"desc"`
				} `json:"comfort"`
			} `json:"life_index"`
		} `json:"realtime"`
		Minutely struct {
			Status          string    `json:"status"`
			Datasource      string    `json:"datasource"`
			Precipitation2H []float64 `json:"precipitation_2h"`
			Precipitation   []float64 `json:"precipitation"`
			Probability     []float64 `json:"probability"`
			Description     string    `json:"description"`
		} `json:"minutely"`
		Hourly struct {
			Status        string `json:"status"`
			Description   string `json:"description"`
			Precipitation []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"precipitation"`
			Temperature []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"temperature"`
			Wind []struct {
				Datetime  string  `json:"datetime"`
				Speed     float64 `json:"speed"`
				Direction float64 `json:"direction"`
			} `json:"wind"`
			Humidity []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"humidity"`
			Cloudrate []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"cloudrate"`
			Skycon []struct {
				Datetime string `json:"datetime"`
				Value    string `json:"value"`
			} `json:"skycon"`
			Pressure []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"pressure"`
			Visibility []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"visibility"`
			Dswrf []struct {
				Datetime string  `json:"datetime"`
				Value    float64 `json:"value"`
			} `json:"dswrf"`
			AirQuality struct {
				Aqi []struct {
					Datetime string `json:"datetime"`
					Value    struct {
						Chn int `json:"chn"`
						Usa int `json:"usa"`
					} `json:"value"`
				} `json:"aqi"`
				Pm25 []struct {
					Datetime string `json:"datetime"`
					Value    int    `json:"value"`
				} `json:"pm25"`
			} `json:"air_quality"`
		} `json:"hourly"`
		Daily struct {
			Status string `json:"status"`
			Astro  []struct {
				Date    string `json:"date"`
				Sunrise struct {
					Time string `json:"time"`
				} `json:"sunrise"`
				Sunset struct {
					Time string `json:"time"`
				} `json:"sunset"`
			} `json:"astro"`
			Precipitation []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"precipitation"`
			Temperature []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"temperature"`
			Wind []struct {
				Date string `json:"date"`
				Max  struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"max"`
				Min struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"min"`
				Avg struct {
					Speed     float64 `json:"speed"`
					Direction float64 `json:"direction"`
				} `json:"avg"`
			} `json:"wind"`
			Humidity []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"humidity"`
			Cloudrate []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"cloudrate"`
			Pressure []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"pressure"`
			Visibility []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"visibility"`
			Dswrf []struct {
				Date string  `json:"date"`
				Max  float64 `json:"max"`
				Min  float64 `json:"min"`
				Avg  float64 `json:"avg"`
			} `json:"dswrf"`
			AirQuality struct {
				Aqi []struct {
					Date string `json:"date"`
					Max  struct {
						Chn int `json:"chn"`
						Usa int `json:"usa"`
					} `json:"max"`
					Avg struct {
						Chn float64 `json:"chn"`
						Usa float64 `json:"usa"`
					} `json:"avg"`
					Min struct {
						Chn int `json:"chn"`
						Usa int `json:"usa"`
					} `json:"min"`
				} `json:"aqi"`
				Pm25 []struct {
					Date string  `json:"date"`
					Max  int     `json:"max"`
					Avg  float64 `json:"avg"`
					Min  int     `json:"min"`
				} `json:"pm25"`
			} `json:"air_quality"`
			Skycon []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon"`
			Skycon08H20H []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon_08h_20h"`
			Skycon20H32H []struct {
				Date  string `json:"date"`
				Value string `json:"value"`
			} `json:"skycon_20h_32h"`
			LifeIndex struct {
				Ultraviolet []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"ultraviolet"`
				CarWashing []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"carWashing"`
				Dressing []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"dressing"`
				Comfort []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"comfort"`
				ColdRisk []struct {
					Date  string `json:"date"`
					Index string `json:"index"`
					Desc  string `json:"desc"`
				} `json:"coldRisk"`
			} `json:"life_index"`
		} `json:"daily"`
		Primary          int    `json:"primary"`
		ForecastKeypoint string `json:"forecast_keypoint"`
	} `json:"result"`
}

func executeGet(rawUrl string) (rawResp []byte, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, rawUrl, nil); err != nil {
		return nil, err
	}
	var resp *http.Response
	if resp, err = HTTPClient.Do(req); err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	var reply = new(struct {
		Status string `json:"status"`
	})
	rawResp, err = ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, fmt.Errorf("unexpected response: %w\nraw response: %s", err, rawResp)
	}
	if reply.Status != "ok" {
		return nil, fmt.Errorf("unknown Status: %s", rawResp)
	}
	return
}
