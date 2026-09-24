package changelog

import (
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"go.l0nax.org/typact"
)

func TestEntryMarshalOmitsNoneAuthor(t *testing.T) {
	none := Entry{ChangeTypeID: "bug_fix", Title: "Fix the thing", Author: typact.None[string]()}
	some := Entry{ChangeTypeID: "bug_fix", Title: "Fix the thing", Author: typact.Some("emanuel")}

	for name, e := range map[string]Entry{"none": none, "some": some} {
		raw, err := toml.Marshal(e)
		if err != nil {
			t.Fatalf("%s: marshal: %v", name, err)
		}

		t.Logf("%s author ->\n%s", name, raw)

		var back Entry
		if err := toml.Unmarshal(raw, &back); err != nil {
			t.Fatalf("%s: unmarshal: %v", name, err)
		}

		if back.Title != e.Title || back.ChangeTypeID != e.ChangeTypeID {
			t.Errorf("%s: round-trip mismatch: %+v", name, back)
		}

		if back.Author.IsSome() != e.Author.IsSome() {
			t.Errorf("%s: author presence changed: got %v want %v",
				name, back.Author.IsSome(), e.Author.IsSome())
		}
	}
}

func TestReleaseInfoRoundTrip(t *testing.T) {
	// the committed ReleaseInfo files use second precision
	in := ReleaseInfo{
		Version:      "v2.0.0-rc.1",
		ReleaseDate:  time.Date(2024, 12, 4, 11, 48, 21, 0, time.FixedZone("CET", 3600)),
		IsPreRelease: true,
	}

	raw, err := toml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	t.Logf("marshaled ->\n%s", raw)

	var back ReleaseInfo
	if err := toml.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !back.ReleaseDate.Equal(in.ReleaseDate) {
		t.Errorf("date changed: got %v want %v", back.ReleaseDate, in.ReleaseDate)
	}

	if back.Version != in.Version || back.IsPreRelease != in.IsPreRelease {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}

// the 16 committed ReleaseInfo files must still parse
func TestParseCommittedReleaseInfo(t *testing.T) {
	const committed = `version = "v2.0.0-rc.1"
date = 2024-12-04T11:48:21+01:00
pre_release = true
`

	var got ReleaseInfo
	if err := toml.Unmarshal([]byte(committed), &got); err != nil {
		t.Fatalf("committed ReleaseInfo no longer parses: %v", err)
	}

	if got.Version != "v2.0.0-rc.1" || !got.IsPreRelease {
		t.Errorf("unexpected: %+v", got)
	}

	t.Logf("parsed: %+v", got)
}
