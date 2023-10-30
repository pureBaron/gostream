# GO-STREAM
##### processing go slices in a funcional way

## sample code:
```go
slc := []int{1, 19, 12, 6, 4, 29}
gostream.New(slc, gostream.SequentialType)
    .Filter(func(i int)bool {return i%2 == 0})
    .ForEach(func(i int) {fmt.Println(i)})
```