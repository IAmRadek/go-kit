package envconfig_test

import (
	"testing"
	"time"

	"github.com/IAmRadek/go-kit/envconfig"
	"github.com/google/go-cmp/cmp"
)

type Config struct {
	NotPopulated string            `env:"-"`
	A            string            `env:"A"`
	B            int               `env:"B"`
	C            int8              `env:"C"`
	D            int16             `env:"D"`
	E            int32             `env:"E"`
	F            int64             `env:"F"`
	G            uint              `env:"G"`
	H            uint8             `env:"H"`
	I            uint16            `env:"I"`
	J            uint32            `env:"J"`
	K            uint64            `env:"K"`
	L            float32           `env:"L"`
	M            float64           `env:"M"`
	N            bool              `env:"N"`
	O            [2]string         `env:"O"`
	P            []string          `env:"P"`
	Q            map[string]string `env:"Q"`

	Default string `env:"DEFAULT" default:"12345"`

	Sub SubConfig `prefix:"SUB"`

	CustomTextUnmarshaler CustomTextUnmarshaler    `prefix:"CUSTOM"`
	Duration              time.Duration            `env:"DURATION"`
	SDur                  []time.Duration          `env:"SDUR"`
	MDur                  map[string]time.Duration `env:"MDUR"`
}

type SubConfig struct {
	A string `env:"AA"`
}

type CustomTextUnmarshaler struct {
	Value string
}

func (c *CustomTextUnmarshaler) UnmarshalText(text []byte) error {
	c.Value = "***" + string(text) + "***"
	return nil
}

func TestRead(t *testing.T) {
	le := func(key string) (string, bool) {
		switch key {
		case "-":
			t.Fatalf("- should not be searched for.")
		case "A":
			return "hello", true
		case "B", "C", "D", "E", "F", "G", "H", "I", "J", "K":
			return "42", true
		case "L", "M":
			return "42.42", true
		case "N":
			return "true", true
		case "O", "P":
			return "hello,world", true
		case "Q":
			return "key1=value1,key2=value2", true
		case "SUB_AA":
			return "sub", true
		case "CUSTOM":
			return "custom", true
		case "DURATION":
			return "1h", true
		case "SDUR":
			return "1h,2h,3h", true
		case "MDUR":
			return "key1=1h,key2=2h,key3=3h", true
		}

		return "", false
	}
	_ = le

	var cfg Config
	if err := envconfig.Read(&cfg, le); err != nil {
		t.Error(err)
	}

	diff := cmp.Diff(cfg, Config{
		A: "hello",
		B: 42,
		C: 42,
		D: 42,
		E: 42,
		F: 42,
		G: 42,
		H: 42,
		I: 42,
		J: 42,
		K: 42,
		L: 42.42,
		M: 42.42,
		N: true,
		O: [2]string{"hello", "world"},
		P: []string{"hello", "world"},
		Q: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},

		Default: "12345",

		Sub: SubConfig{
			A: "sub",
		},
		CustomTextUnmarshaler: CustomTextUnmarshaler{
			Value: "***custom***",
		},
		Duration: time.Hour,
		SDur: []time.Duration{
			time.Hour,
			2 * time.Hour,
			3 * time.Hour,
		},
		MDur: map[string]time.Duration{
			"key1": time.Hour,
			"key2": 2 * time.Hour,
			"key3": 3 * time.Hour,
		},
	})
	if diff != "" {
		t.Error("expected equal: got\n", diff)
	}
}
