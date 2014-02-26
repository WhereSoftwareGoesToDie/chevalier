package chevalier

import (
	"math/rand"
	"time"
)

func randomAlphaByte() byte {
	b := rand.Intn(26) + 97
	return byte(b)
}

func genAlphaString(n int) string {
	data := make([]byte, n)
	for i, _ := range data {
		data[i] = randomAlphaByte()
	}
	return string(data[:n])
}

func genTestPattern() string {
	l := rand.Intn(20)
	return genAlphaString(l)
}

func genTestTags() []*SourceRequest_Tag {
	tags := make([]*SourceRequest_Tag, 0)
	tags = append(tags, NewSourceRequestTag("hostname", genTestPattern()))
	tags = append(tags, NewSourceRequestTag("metric", genTestPattern()))
	return tags
}

// GenTestSourceRequest returns a SourceRequest filled with random data.
func GenTestSourceRequest() *SourceRequest {
	rand.Seed(time.Now().UTC().UnixNano())
	tags := genTestTags()
	req := NewSourceRequest(tags)
	if rand.Intn(1) == 0 {
		startPage := rand.Int63n(10)
		pageLen := rand.Int63n(10)
		req.StartPage = &startPage
		req.SourcesPerPage = &pageLen
	}
	return req
}
