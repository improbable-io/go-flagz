package flagz

import (
	"testing"
	flag "github.com/spf13/pflag"

	"github.com/stretchr/testify/assert"
	"time"
	"fmt"
)


func TestDynDuration_SetAndGet(t *testing.T) {
	set := flag.NewFlagSet("foobar", flag.ContinueOnError)
	dynFlag := DynDuration(set, "some_duration_1", 5 * time.Second, "Use it or lose it")
	assert.Equal(t, 5 * time.Second, dynFlag.Get(), "value must be default after create")
	err := set.Set("some_duration_1", "10h")
	assert.NoError(t, err, "setting value must succeed")
	assert.Equal(t, 10 * time.Hour, dynFlag.Get(), "value must be set after update")
}

func TestDynDuration_IsMarkedDynamic(t *testing.T) {
	set := flag.NewFlagSet("foobar", flag.ContinueOnError)
	DynDuration(set, "some_duration_1", 5* time.Minute, "Use it or lose it")
	assert.True(t, IsFlagDynamic(set.Lookup("some_duration_1")))
}

func TestDynDuration_FiresValidators(t *testing.T) {
	set := flag.NewFlagSet("foobar", flag.ContinueOnError)
	validator := func (x time.Duration) error {
		if x > 1 * time.Hour {
			return fmt.Errorf("too long")
		}
		return nil
	}
	DynDuration(set, "some_duration_1", 5 * time.Second, "Use it or lose it").WithValidator(validator)

	assert.NoError(t, set.Set("some_duration_1", "50m"), "no error from validator when in range")
	assert.Error(t, set.Set("some_duration_1", "2h"), "error from validator when value out of range")
}

func Benchmark_Duration_Dyn_Get(b *testing.B) {
	set := flag.NewFlagSet("foobar", flag.ContinueOnError)
	value := DynDuration(set, "some_duration_1", 5 * time.Second, "Use it or lose it")
	set.Set("some_duration_1", "10s")
	for i := 0; i < b.N; i++ {
		value.Get()
	}
}

func Benchmark_Duration_Normal_get(b *testing.B) {
	set := flag.NewFlagSet("foobar", flag.ContinueOnError)
	valPtr := set.Duration("some_duration_1", 5 * time.Second, "Use it or lose it")
	set.Set("some_duration_1", "10s")
	for i := 0; i < b.N; i++ {
		x := *valPtr
		x = x
	}
}