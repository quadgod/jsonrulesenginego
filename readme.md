# Path Resolver

Go lib which implements get value by path from Struct, Array, Slice or Map

Example:
```go
type Obj struct {
    FloatField    float64
}

obj := new(Obj)
obj.FloatField = 300.1

value, err := TryGetValueByPath("FloatField", obj)
// value == 300.1
```
