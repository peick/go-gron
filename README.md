# gron

A go library to make JSON greppable.

gron transforms JSON into discrete assignments to make it easier to grep for what you want and see the absolute 'path' to it.

## Examples

```go
import (
    "bufio"
    "fmt"
    "os"

    "github.com/peick/go-gron"
)

func main() {
    f, _ := os.Open("example.json")
    reader := bufio.NewReader(f)

    g := gron.New(reader)

    out := g.String()

    fmt.Println(out)
}
```
