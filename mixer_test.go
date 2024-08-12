package beep_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Crystal303/beep"
	"github.com/Crystal303/beep/internal/testtools"
)

func TestMixer_MixesSamples(t *testing.T) {
	epsilon := 0.000001

	s1, data1 := testtools.RandomDataStreamer(200)
	s2, data2 := testtools.RandomDataStreamer(200)

	m := beep.Mixer{}
	m.Add(s1)
	m.Add(s2)

	samples := testtools.CollectNum(100, &m)
	for i, s := range samples {
		wantL := data1[i][0] + data2[i][0]
		wantR := data1[i][1] + data2[i][1]

		if s[0] < wantL-epsilon || s[0] > wantL+epsilon {
			t.Fatalf("unexpected value for mixed samples; expected: %f, got: %f", wantL, s[0])
		}
		if s[1] < wantR-epsilon || s[1] > wantR+epsilon {
			t.Fatalf("unexpected value for mixed samples; expected: %f, got: %f", wantR, s[1])
		}
	}

	s3, data3 := testtools.RandomDataStreamer(100)
	m.Add(s3)

	samples = testtools.CollectNum(100, &m)
	for i, s := range samples {
		wantL := data1[100+i][0] + data2[100+i][0] + data3[i][0]
		wantR := data1[100+i][1] + data2[100+i][1] + data3[i][1]

		if s[0] < wantL-epsilon || s[0] > wantL+epsilon {
			t.Fatalf("unexpected value for mixed samples; expected: %f, got: %f", wantL, s[0])
		}
		if s[1] < wantR-epsilon || s[1] > wantR+epsilon {
			t.Fatalf("unexpected value for mixed samples; expected: %f, got: %f", wantR, s[1])
		}
	}
}

func TestMixer_DrainedStreamersAreRemoved(t *testing.T) {
	s1, _ := testtools.RandomDataStreamer(50)
	s2, _ := testtools.RandomDataStreamer(60)

	m := beep.Mixer{}
	m.Add(s1)
	m.Add(s2)

	// Drain s1 but not so far it returns false.
	samples := testtools.CollectNum(50, &m)
	assert.Len(t, samples, 50)
	assert.Equal(t, 2, m.Len())

	// Fully drain s1.
	// Drain s2 but not so far it returns false.
	samples = testtools.CollectNum(10, &m)
	assert.Len(t, samples, 10)
	assert.Equal(t, 1, m.Len())

	// Fully drain s2.
	samples = testtools.CollectNum(10, &m)
	assert.Len(t, samples, 10)
	assert.Equal(t, 0, m.Len())
}

func TestMixer_PlaysSilenceWhenNoStreamersProduceSamples(t *testing.T) {
	m := beep.Mixer{}

	// Test silence before streamers are added.
	samples := testtools.CollectNum(10, &m)
	assert.Len(t, samples, 10)
	assert.Equal(t, make([][2]float64, 10), samples)

	// Test silence after streamer is partly drained.
	s, _ := testtools.RandomDataStreamer(50)
	m.Add(s)

	samples = testtools.CollectNum(100, &m)
	assert.Len(t, samples, 100)
	assert.Equal(t, 1, m.Len())
	assert.Equal(t, make([][2]float64, 50), samples[50:])

	// Test silence when streamer is fully drained.
	samples = testtools.CollectNum(10, &m)
	assert.Len(t, samples, 10)
	assert.Equal(t, 0, m.Len())
	assert.Equal(t, make([][2]float64, 10), samples)

	// Test silence after streamer was fully drained.
	samples = testtools.CollectNum(10, &m)
	assert.Len(t, samples, 10)
	assert.Equal(t, make([][2]float64, 10), samples)
}

func BenchmarkMixer_MultipleStreams(b *testing.B) {
	s1, _ := testtools.RandomDataStreamer(b.N)
	s2, _ := testtools.RandomDataStreamer(b.N)

	m := beep.Mixer{}
	m.Add(s1)
	m.Add(s2)

	b.StartTimer()

	testtools.CollectNum(b.N, &m)
}

func BenchmarkMixer_OneStream(b *testing.B) {
	s, _ := testtools.RandomDataStreamer(b.N)

	m := beep.Mixer{}
	m.Add(s)

	b.StartTimer()
	testtools.CollectNum(b.N, &m)
}

func BenchmarkMixer_Silence(b *testing.B) {
	m := beep.Mixer{}
	// Don't add any streamers

	b.StartTimer()
	testtools.CollectNum(b.N, &m)
}
