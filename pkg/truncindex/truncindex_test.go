package truncindex

import (
	"math/rand"
	"testing"

	"github.com/dotcloud/docker/utils"
)

// Test the behavior of TruncIndex, an index for querying IDs from a non-conflicting prefix.
func TestTruncIndex(t *testing.T) {
	ids := []string{}
	index := NewTruncIndex(ids)
	// Get on an empty index
	if _, err := index.Get("foobar"); err == nil {
		t.Fatal("Get on an empty index should return an error")
	}

	// Spaces should be illegal in an id
	if err := index.Add("I have a space"); err == nil {
		t.Fatalf("Adding an id with ' ' should return an error")
	}

	id := "99b36c2c326ccc11e726eee6ee78a0baf166ef96"
	// Add an id
	if err := index.Add(id); err != nil {
		t.Fatal(err)
	}

	// Add an empty id (should fail)
	if err := index.Add(""); err == nil {
		t.Fatalf("Adding an empty id should return an error")
	}

	// Get a non-existing id
	assertIndexGet(t, index, "abracadabra", "", true)
	// Get an empty id
	assertIndexGet(t, index, "", "", true)
	// Get the exact id
	assertIndexGet(t, index, id, id, false)
	// The first letter should match
	assertIndexGet(t, index, id[:1], id, false)
	// The first half should match
	assertIndexGet(t, index, id[:len(id)/2], id, false)
	// The second half should NOT match
	assertIndexGet(t, index, id[len(id)/2:], "", true)

	id2 := id[:6] + "blabla"
	// Add an id
	if err := index.Add(id2); err != nil {
		t.Fatal(err)
	}
	// Both exact IDs should work
	assertIndexGet(t, index, id, id, false)
	assertIndexGet(t, index, id2, id2, false)

	// 6 characters or less should conflict
	assertIndexGet(t, index, id[:6], "", true)
	assertIndexGet(t, index, id[:4], "", true)
	assertIndexGet(t, index, id[:1], "", true)

	// 7 characters should NOT conflict
	assertIndexGet(t, index, id[:7], id, false)
	assertIndexGet(t, index, id2[:7], id2, false)

	// Deleting a non-existing id should return an error
	if err := index.Delete("non-existing"); err == nil {
		t.Fatalf("Deleting a non-existing id should return an error")
	}

	// Deleting an empty id should return an error
	if err := index.Delete(""); err == nil {
		t.Fatal("Deleting an empty id should return an error")
	}

	// Deleting id2 should remove conflicts
	if err := index.Delete(id2); err != nil {
		t.Fatal(err)
	}
	// id2 should no longer work
	assertIndexGet(t, index, id2, "", true)
	assertIndexGet(t, index, id2[:7], "", true)
	assertIndexGet(t, index, id2[:11], "", true)

	// conflicts between id and id2 should be gone
	assertIndexGet(t, index, id[:6], id, false)
	assertIndexGet(t, index, id[:4], id, false)
	assertIndexGet(t, index, id[:1], id, false)

	// non-conflicting substrings should still not conflict
	assertIndexGet(t, index, id[:7], id, false)
	assertIndexGet(t, index, id[:15], id, false)
	assertIndexGet(t, index, id, id, false)
}

func assertIndexGet(t *testing.T, index *TruncIndex, input, expectedResult string, expectError bool) {
	if result, err := index.Get(input); err != nil && !expectError {
		t.Fatalf("Unexpected error getting '%s': %s", input, err)
	} else if err == nil && expectError {
		t.Fatalf("Getting '%s' should return an error, not '%s'", input, result)
	} else if result != expectedResult {
		t.Fatalf("Getting '%s' returned '%s' instead of '%s'", input, result, expectedResult)
	}
}

func BenchmarkTruncIndexAdd100(b *testing.B) {
	var testSet []string
	for i := 0; i < 100; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexAdd250(b *testing.B) {
	var testSet []string
	for i := 0; i < 250; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexAdd500(b *testing.B) {
	var testSet []string
	for i := 0; i < 500; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexGet100(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 100; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	index := NewTruncIndex([]string{})
	for _, id := range testSet {
		if err := index.Add(id); err != nil {
			b.Fatal(err)
		}
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}

func BenchmarkTruncIndexGet250(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 250; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	index := NewTruncIndex([]string{})
	for _, id := range testSet {
		if err := index.Add(id); err != nil {
			b.Fatal(err)
		}
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}

func BenchmarkTruncIndexGet500(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	index := NewTruncIndex([]string{})
	for _, id := range testSet {
		if err := index.Add(id); err != nil {
			b.Fatal(err)
		}
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}

func BenchmarkTruncIndexDelete100(b *testing.B) {
	var testSet []string
	for i := 0; i < 100; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		b.StartTimer()
		for _, id := range testSet {
			if err := index.Delete(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexDelete250(b *testing.B) {
	var testSet []string
	for i := 0; i < 250; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		b.StartTimer()
		for _, id := range testSet {
			if err := index.Delete(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexDelete500(b *testing.B) {
	var testSet []string
	for i := 0; i < 500; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		b.StartTimer()
		for _, id := range testSet {
			if err := index.Delete(id); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkTruncIndexNew100(b *testing.B) {
	var testSet []string
	for i := 0; i < 100; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTruncIndex(testSet)
	}
}

func BenchmarkTruncIndexNew250(b *testing.B) {
	var testSet []string
	for i := 0; i < 250; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTruncIndex(testSet)
	}
}

func BenchmarkTruncIndexNew500(b *testing.B) {
	var testSet []string
	for i := 0; i < 500; i++ {
		testSet = append(testSet, utils.GenerateRandomID())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewTruncIndex(testSet)
	}
}

func BenchmarkTruncIndexAddGet100(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ {
		id := utils.GenerateRandomID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}

func BenchmarkTruncIndexAddGet250(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ {
		id := utils.GenerateRandomID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}

func BenchmarkTruncIndexAddGet500(b *testing.B) {
	var testSet []string
	var testKeys []string
	for i := 0; i < 500; i++ {
		id := utils.GenerateRandomID()
		testSet = append(testSet, id)
		l := rand.Intn(12) + 12
		testKeys = append(testKeys, id[:l])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := NewTruncIndex([]string{})
		for _, id := range testSet {
			if err := index.Add(id); err != nil {
				b.Fatal(err)
			}
		}
		for _, id := range testKeys {
			if res, err := index.Get(id); err != nil {
				b.Fatal(res, err)
			}
		}
	}
}
