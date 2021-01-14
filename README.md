# ColorfulClouds API

[彩云天气 API](https://open.caiyunapp.com/%E5%BD%A9%E4%BA%91%E5%A4%A9%E6%B0%94_API_%E4%B8%80%E8%A7%88%E8%A1%A8)
> 个人用户, `2个` 可用接口
> - 实况天气
> - 通用预报

## Installation

```shell
go get github.com/electricbubble/colorful-clouds-api
```

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	caiYunApi "github.com/electricbubble/colorful-clouds-api"
	"log"
	"os"
)

func main() {
	cyToken := os.Getenv("Token")
	cyLongitudeAndLatitude := os.Getenv("LongitudeAndLatitude")
	cyAPi := caiYunApi.NewColorfulCloudsApi(cyToken, cyLongitudeAndLatitude)

	realtime, err := cyAPi.Realtime()
	if err != nil {
		log.Fatal(err)
	}

	buffer := bytes.NewBufferString("🌈天气预报员\n")
	buffer.WriteString(fmt.Sprintf("当前温度: %.0f˚\n", realtime.Result.Realtime.Temperature))
	buffer.WriteString(fmt.Sprintf("体感温度: %.0f˚\n", realtime.Result.Realtime.ApparentTemperature))
	// buffer.WriteString(fmt.Sprintf("主要天气现象: %s\n", realtime.Result.Realtime.Skycon))
	buffer.WriteString(fmt.Sprintf("能见度: %.1fkm\n", realtime.Result.Realtime.Visibility))
	buffer.WriteString(fmt.Sprintf("风速: %.1fkm/hr\n", realtime.Result.Realtime.Wind.Speed))
	buffer.WriteString(fmt.Sprintf("空气质量: %s\n", realtime.Result.Realtime.AirQuality.Description.Chn))
	buffer.WriteString(fmt.Sprintf("紫外线指数: %s\n", realtime.Result.Realtime.LifeIndex.Ultraviolet.Desc))
	buffer.WriteString(fmt.Sprintf("舒适度指数: %s\n", realtime.Result.Realtime.LifeIndex.Comfort.Desc))

	weather, err := cyAPi.Weather(caiYunApi.Alert(true), caiYunApi.DailySteps(1), caiYunApi.HourlySteps(24))
	if err == nil {
		buffer.WriteString(weather.Result.ForecastKeypoint + "\n")
		astro := weather.Result.Daily.Astro
		if len(astro) > 0 {
			buffer.WriteString(fmt.Sprintf("日出时间: %s\n", astro[0].Sunrise))
			buffer.WriteString(fmt.Sprintf("日落时间: %s\n", astro[0].Sunset))
		}
	}

	fmt.Println(buffer.String())
}

```
