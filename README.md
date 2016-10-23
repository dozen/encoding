# encoding

## base128
base128 using alphabet, number, hiragana.

### import

```
import "github.com/dozen/encoding/base128"
```

### encode

```
e := base128.NewEncoding(base128.StdEncoding)

e.EncodeToString([]byte("Hello, world!"))
//kZtけふ8ぢg7どぺmふQお
```

### decode

```
en := base128.NewEncoding(base128.StdEncoding)

fmt.Printf("%s\n", en.Decode("kZtけふ8ぢg7どぺmふQお"))
```
