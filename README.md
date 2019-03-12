# Abstract JSON


# TODO

- Global
- [ ] array support
- [ ] object support
- Node
- [x] key => *string
- [x] ‌Value() => Source()
- [ ] add func Value() interface{}
- Functions 
- [ ] func JsonPath(data [] byte, path string, clone bool) ([]*Node, error) 
- [ ] func (n *Node) JsonPath(path string) ([]*Node, error)
- Shugar
- [ ] ‌Node: IsArray, Is Numeric,...
- buffer
- [ ] ‌const: coma
- [ ] add tests
- [ ] func scan(b byte, skip bool) error
- [ ] func skip(...) error
- errors
- [ ] expected error: `wrong symbol '%s' expected %s, on %d`
- [ ] add buffer in error: detect column and line from index
- [ ] ‌error*(b *buffer) error
- [ ] fix iota use
- future
- [ ] use io.Reader instead of []byte
