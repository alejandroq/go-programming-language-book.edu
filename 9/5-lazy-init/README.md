For expensive inititialization steps, consider using `sync.Once`. Conceptually the `sync.Once` contains a mutex and a boolean variable that records whether initialization has taken place. The mutex guards both the boolean and the data structure. The sole method, `Do`, accepts the initialization function as it's argument.

```go
var loadIconsOnce sync.Once
var icons map[string]image.Image

// concurrency-safe
func Icon(name string) image.Image {
    loadIconsOnce.Do(loadIcons) // deferred UNTIL use. Need only run once to guarentee icons map is filled

    // after 1 call, the mutex ensures that icons are in memory
    return icons[name]
}
```
