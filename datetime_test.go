package httpexpect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_FailedChain(t *testing.T) {
	chain := newMockChain(t)
	chain.setFailed()

	tm := time.Unix(0, 0)
	value := newDateTime(chain, tm)
	value.chain.assertFailed(t)

	value.Alias("foo")

	value.IsEqual(tm)
	value.NotEqual(tm)
	value.InRange(tm, tm)
	value.NotInRange(tm, tm)
	value.InList(tm, tm)
	value.NotInList(tm, tm)
	value.Gt(tm)
	value.Ge(tm)
	value.Lt(tm)
	value.Le(tm)

	value.Zone().chain.assertFailed(t)
	value.Year().chain.assertFailed(t)
	value.Month().chain.assertFailed(t)
	value.Day().chain.assertFailed(t)
	value.WeekDay().chain.assertFailed(t)
	value.YearDay().chain.assertFailed(t)
	value.Hour().chain.assertFailed(t)
	value.Minute().chain.assertFailed(t)
	value.Second().chain.assertFailed(t)
	value.Nanosecond().chain.assertFailed(t)

	value.AsUTC().chain.assertFailed(t)
	value.AsLocal().chain.assertFailed(t)
}

func TestDateTime_Constructors(t *testing.T) {
	time := time.Unix(0, 1234)

	t.Run("reporter", func(t *testing.T) {
		reporter := newMockReporter(t)
		value := NewDateTime(reporter, time)
		value.IsEqual(time)
		value.chain.assertNotFailed(t)
	})

	t.Run("config", func(t *testing.T) {
		reporter := newMockReporter(t)
		value := NewDateTimeC(Config{
			Reporter: reporter,
		}, time)
		value.IsEqual(time)
		value.chain.assertNotFailed(t)
	})

	t.Run("chain", func(t *testing.T) {
		chain := newMockChain(t)
		value := newDateTime(chain, time)
		assert.NotSame(t, value.chain, chain)
		assert.Equal(t, value.chain.context.Path, chain.context.Path)
	})
}

func TestDateTime_Alias(t *testing.T) {
	reporter := newMockReporter(t)

	value1 := NewDateTime(reporter, time.Unix(0, 1234))
	assert.Equal(t, []string{"DateTime()"}, value1.chain.context.Path)
	assert.Equal(t, []string{"DateTime()"}, value1.chain.context.AliasedPath)

	value2 := value1.Alias("foo")
	assert.Equal(t, []string{"DateTime()"}, value2.chain.context.Path)
	assert.Equal(t, []string{"foo"}, value2.chain.context.AliasedPath)
}

func TestDateTime_Getters(t *testing.T) {
	reporter := newMockReporter(t)

	parsedTime, _ := time.Parse(time.UnixDate, "FRI Dec 30 15:04:05 IST 2022")

	value := NewDateTime(reporter, parsedTime)

	value.chain.assertNotFailed(t)

	value.Zone().chain.assertNotFailed(t)
	value.Year().chain.assertNotFailed(t)
	value.Month().chain.assertNotFailed(t)
	value.Day().chain.assertNotFailed(t)
	value.WeekDay().chain.assertNotFailed(t)
	value.YearDay().chain.assertNotFailed(t)
	value.Hour().chain.assertNotFailed(t)
	value.Minute().chain.assertNotFailed(t)
	value.Second().chain.assertNotFailed(t)
	value.Nanosecond().chain.assertNotFailed(t)
	value.AsUTC().chain.assertNotFailed(t)
	value.AsLocal().chain.assertNotFailed(t)

	expectedTime := parsedTime
	expectedZone, _ := expectedTime.Zone()
	assert.Equal(t, expectedZone, value.Zone().Raw())
	assert.Equal(t, float64(expectedTime.Year()), value.Year().Raw())
	assert.Equal(t, float64(expectedTime.Month()), value.Month().Raw())
	assert.Equal(t, float64(expectedTime.Day()), value.Day().Raw())
	assert.Equal(t, float64(expectedTime.Weekday()), value.WeekDay().Raw())
	assert.Equal(t, float64(expectedTime.YearDay()), value.YearDay().Raw())
	assert.Equal(t, float64(expectedTime.Hour()), value.Hour().Raw())
	assert.Equal(t, float64(expectedTime.Minute()), value.Minute().Raw())
	assert.Equal(t, float64(expectedTime.Second()), value.Second().Raw())
	assert.Equal(t, float64(expectedTime.Nanosecond()), value.Nanosecond().Raw())
}

func TestDateTime_Equal(t *testing.T) {
	reporter := newMockReporter(t)

	value := NewDateTime(reporter, time.Unix(0, 1234))

	assert.True(t, time.Unix(0, 1234).Equal(value.Raw()))

	value.IsEqual(time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.IsEqual(time.Unix(0, 4321))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotEqual(time.Unix(0, 4321))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotEqual(time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()
}

func TestDateTime_Greater(t *testing.T) {
	reporter := newMockReporter(t)

	value := NewDateTime(reporter, time.Unix(0, 1234))

	value.Gt(time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Gt(time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.Ge(time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Ge(time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Ge(time.Unix(0, 1234+1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()
}

func TestDateTime_Lesser(t *testing.T) {
	reporter := newMockReporter(t)

	value := NewDateTime(reporter, time.Unix(0, 1234))

	value.Lt(time.Unix(0, 1234+1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Lt(time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.Le(time.Unix(0, 1234+1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Le(time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.Le(time.Unix(0, 1234-1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()
}

func TestDateTime_InRange(t *testing.T) {
	reporter := newMockReporter(t)

	value := NewDateTime(reporter, time.Unix(0, 1234))

	value.InRange(time.Unix(0, 1234), time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234), time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InRange(time.Unix(0, 1234-1), time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234-1), time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InRange(time.Unix(0, 1234), time.Unix(0, 1234+1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234), time.Unix(0, 1234+1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InRange(time.Unix(0, 1234+1), time.Unix(0, 1234+2))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234+1), time.Unix(0, 1234+2))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.InRange(time.Unix(0, 1234-2), time.Unix(0, 1234-1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234-2), time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.InRange(time.Unix(0, 1234+1), time.Unix(0, 1234-1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInRange(time.Unix(0, 1234+1), time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()
}

func TestDateTime_InList(t *testing.T) {
	reporter := newMockReporter(t)

	value := NewDateTime(reporter, time.Unix(0, 1234))

	value.InList()
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInList()
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234), time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234), time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234-1), time.Unix(0, 1234))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234-1), time.Unix(0, 1234))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234), time.Unix(0, 1234+1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234), time.Unix(0, 1234+1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234+1), time.Unix(0, 1234+2))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234+1), time.Unix(0, 1234+2))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234-2), time.Unix(0, 1234-1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234-2), time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()

	value.InList(time.Unix(0, 1234+1), time.Unix(0, 1234-1))
	value.chain.assertFailed(t)
	value.chain.clearFailed()

	value.NotInList(time.Unix(0, 1234+1), time.Unix(0, 1234-1))
	value.chain.assertNotFailed(t)
	value.chain.clearFailed()
}
