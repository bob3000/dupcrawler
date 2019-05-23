package fshash_test

import (
	"reflect"
	"testing"

	"github.com/bob3000/dupcrawler/fshash"
)

func assertMapEqual(t *testing.T, got fshash.Map, want fshash.Map) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%s' want '%s'", got, want)
	}

}

func TestReadPath(t *testing.T) {
	t.Run("Read path", func(t *testing.T) {
		args := fshash.ReadPathArgs{
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    0,
			Excludes:    []string{},
			Parallel:    false,
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := fshash.Map{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/b/a.txt", "testdata/c/d/a.txt", },
			"LSuQi2CqnyqcwBxqNkbkyumeibI=": []string{"testdata/b/b.txt"},
			"TTmnAadUv+7LU2O2oAq7FGX/+co=": []string{"testdata/c/c.txt"},
			"YICafv5TlzimYjDtEdqN2lXFRtg=": []string{"testdata/c/d/d.txt"},
		}
		assertMapEqual(t, got, want)
	})

	t.Run("Read path in parallel", func(t *testing.T) {
		args := fshash.ReadPathArgs{
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    0,
			Excludes:    []string{},
			Parallel:    true,
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := fshash.Map{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/b/a.txt", "testdata/c/d/a.txt", },
			"LSuQi2CqnyqcwBxqNkbkyumeibI=": []string{"testdata/b/b.txt"},
			"TTmnAadUv+7LU2O2oAq7FGX/+co=": []string{"testdata/c/c.txt"},
			"YICafv5TlzimYjDtEdqN2lXFRtg=": []string{"testdata/c/d/d.txt"},
		}
		assertMapEqual(t, got, want)
	})

	t.Run("Read path with max depth", func(t *testing.T) {
		args := fshash.ReadPathArgs{
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    3,
			Excludes:    []string{},
			Parallel:    false,
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := fshash.Map{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/b/a.txt"},
			"LSuQi2CqnyqcwBxqNkbkyumeibI=": []string{"testdata/b/b.txt"},
			"TTmnAadUv+7LU2O2oAq7FGX/+co=": []string{"testdata/c/c.txt"},
		}
		assertMapEqual(t, got, want)
	})

	t.Run("Read path with excludes", func(t *testing.T) {
		args := fshash.ReadPathArgs{
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    0,
			Excludes:    []string{"testdata/b"},
			Parallel:    false,
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := fshash.Map{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/c/d/a.txt"},
			"TTmnAadUv+7LU2O2oAq7FGX/+co=": []string{"testdata/c/c.txt"},
			"YICafv5TlzimYjDtEdqN2lXFRtg=": []string{"testdata/c/d/d.txt"},
		}
		assertMapEqual(t, got, want)
	})
}

func BenchmarkReadPath(b *testing.B) {
	b.Run("benchmark sequential", func(b *testing.B) {
		args := fshash.ReadPathArgs{
			// you might want to change FPath to a larger data set for real
			// benchmarking
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    0,
			Excludes:    []string{},
			Parallel:    false,
			Verbose:     false,
		}
		fshash.ReadPath(args)
	})

	b.Run("benchmark parallel", func(b *testing.B) {
		args := fshash.ReadPathArgs{
			// you might want to change FPath to a larger data set for real
			// benchmarking
			FPath:       "testdata",
			FollowLinks: false,
			CurDepth:    0,
			MaxDepth:    0,
			Excludes:    []string{},
			Parallel:    true,
			Verbose:     false,
		}
		fshash.ReadPath(args)
	})
}