# mcongo
An in-memory cache written in Go.

## Example

    import (
        "fmt"
        "time"

        memcache "github.com/najork/mcongo/pkg/memcache"
    )

    func main() {
        cache := memcache.New()

        cache.Put("foo", "bar", time.Now().Add(time.Minute))

        val, valid, err := cache.Get("foo")
        if err != nil {
            fmt.Printf("err: %s", err.Error())
            return
        }
        if !valid {
            fmt.Printf("value is no longer valid")
            return
        }
        fmt.Printf("value: %s", val)
        return
    }
