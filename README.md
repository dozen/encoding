# encoding

## base128
base128 using alphabet, number, hiragana.

### import
```
import "github.com/dozen/encoding/base128"
```

### encode
```

base128.StdEncoding.EncodeToString([]byte("Hello, world!")) //=>kZtけふ8ぢg7どぺmふQお
```

### decode
```
base128.StdEncoding.DecodeString("kZtけふ8ぢg7どぺmふQお") //=>Hello, world!
```

## base256
base256 using alphabet, number, hiragana, katakana, kanji.

### import
```
import "github.com/dozen/encoding/base256"
```

### encode
```
e := base256.NewEncoding(base256.StdEncoding)

e.EncodeToString([]byte("Hello, world!")) //=>かとははひsgずひふはでh
```

### decode
```
e := base256.NewEncoding(base256.StdEncoding)

fmt.Printf("%s\n", e.Decode("かとははひsgずひふはでh")) //=>Hello, world!
```
