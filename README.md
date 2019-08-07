# gron

A go library to make JSON greppable inspired by https://github.com/tomnomnom/gron

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
    // g := gron.New(reader, gron.OriginalGronFormatter())

    out, _ := g.String()

    fmt.Println(out)
}
```
