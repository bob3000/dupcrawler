package fshash_test

import (
	"reflect"
	"testing"

	"github.com/bob3000/dupcrawler/fshash"
)

func assertVal(t *testing.T, got interface{}, want interface{}) {
	if got != want {
		t.Errorf("got '%s' want '%s'", got, want)
	}
}

func assertMapEqual(t *testing.T, got map[string][]string, want map[string][]string) {
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
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := map[string][]string{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/c/d/a.txt", "testdata/b/a.txt"},
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
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := map[string][]string{
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
			Verbose:     false,
		}
		got := fshash.ReadPath(args)
		want := map[string][]string{
			"68EedPgEqj52eHjjZa2A0OBFQm8=": []string{
				"testdata/a/a.txt", "testdata/c/d/a.txt"},
			"TTmnAadUv+7LU2O2oAq7FGX/+co=": []string{"testdata/c/c.txt"},
			"YICafv5TlzimYjDtEdqN2lXFRtg=": []string{"testdata/c/d/d.txt"},
		}
		assertMapEqual(t, got, want)
	})
}
