# HUPX Hungarian Power exchange daily data
Get data from HUPX 
![Data](hupxdata.png)

To install this package
```
go get github.com/mishop/hupxapi 
```

## How to use?
```
import (
	"github.com/mishop/hupxapi"
)
```

Retrieve data
```
hupxapi.GetHUPX("2006-01-02")
```

## Full example daily data
```
package main

import (
	"fmt"
	"time"

	"github.com/mishop/hupxapi"
)

func main() {
	// download the target HTML document
	currentTime := time.Now()
	data := hupxapi.GetHUPX(currentTime.Format("2006-01-02"))
	fmt.Println("Baseload price", data["Baseload price"])
}
```

## Make influxdb client for colect data

```
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/mishop/hupxapi"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://localhost:8086", "my-token")
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("my-org", "my-bucket")

	// Generiranje trenutnog vremena
	now := time.Now()
	data := hupxapi.GetHUPX(now.Format("2006-01-02"))

	for key, value := range data {
		// Create point using fluent style
		p := influxdb2.NewPointWithMeasurement("stat").
			AddTag("key", key).
			AddField("value", value).
			SetTime(now)
		writeAPI.WritePoint(context.Background(), p)
	}

	client.Close()

	fmt.Println("Data writed in InfluxDB.")
}
```