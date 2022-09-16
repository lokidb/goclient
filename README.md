# goclient
LokiDB golang client sdk

---

## Example
```golang
package main

import (
	"fmt"
	"time"

	"github.com/lokidb/goclient"
)

func main() {
	addresss := []goclient.NodeAddress{
		{Host: "127.0.0.1", Port: 50051},
		{Host: "127.0.0.1", Port: 50052},
		{Host: "127.0.0.1", Port: 50053},
		{Host: "127.0.0.1", Port: 50054},
		{Host: "127.0.0.1", Port: 50055},
	}

	c := goclient.New(addresss, time.Minute)
	defer c.Close()

	c.Set("a", "A")

	val, _ := c.Get("a")
	fmt.Println(val)

	deleted, _ := c.Del("a")
	fmt.Println(deleted)

	keys, _ := c.Keys()
	fmt.Println(keys)

	c.Flush()
}
```

## API
| Method | Input                  | Output                          |
|--------|------------------------|---------------------------------|
| Get    | key (str)              | value (str), error              |
| Set    | key (str), value (str) | error                           |
| Del    | key(str)               | error                           |
| Keys   |                        | list of keys (list[str]), error |
| Flush  |                        | error                           |