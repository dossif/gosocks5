package logger

import "testing"

func TestAddField(t *testing.T) {
	lg, _ := NewLogger("info")
	lg2 := *lg
	lg.Lg.Info().Msgf("test log1")
	lg.AddField(map[string]string{"aaa": "bbb"})
	lg.Lg.Info().Msgf("test log2")
	lg.AddField(map[string]string{"ccc": "ddd"})
	lg.Lg.Info().Msgf("test log3")
	lg.AddField(map[string]string{"a1": "b1", "a2": "b2", "a3": "b3"})
	lg.Lg.Info().Msgf("test log3")
	lg2.Lg.Info().Msgf("test log4")

}
