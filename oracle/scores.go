package oracle

import (
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

type pathScore float64
type scoringServiceName string
type pathFingerprint pan.PathFingerprint
type subscribedScores map[string][]scoringServiceName

type scoreRequest struct {
	SubQueries subscribedScores `json:"scorings"`
}

type reponseEntry struct {
	Fingerprint pathFingerprint    `json:"fingerprint"`
	Scores      map[string]float64 `json:"scores"`
}

type scoreResponse map[string][]reponseEntry

func ScoresToLuaTable(scores *scoreResponse) *lua.LTable {
	if scores == nil {
		return &lua.LTable{}
	}

	t := lua.LTable{}

	for dst, entries := range *scores {
		dstT := &lua.LTable{}
		for _, e := range entries {
			val := &lua.LTable{}
			val.RawSetString("Fingerprint", lua.LString(e.Fingerprint))
			stT := &lua.LTable{}
			for k, v := range e.Scores {
				stT.RawSetString(k, lua.LNumber(v))
			}
			val.RawSetString("Stats", stT)
			dstT.Append(val)
		}
		t.RawSetString(dst, dstT)
	}
	return &t
}

func InjectPathRefs(oracleEntries *lua.LTable, pathRefs []*lua.LTable) {
	oe := make(map[string]*lua.LTable)
	for _, pRef := range pathRefs {
		fp := pRef.RawGetString("Fingerprint").String()
		oe[fp] = pRef
	}

	for _, pRef := range pathRefs {
		pathDst := pRef.RawGetString("Destination").String()
		fp := pRef.RawGetString("Fingerprint").String()
		dstAS := strings.Split(pathDst, ",")[0]

		if oracleEntry := GetPathEntry(oracleEntries, fp, dstAS); oracleEntry != nil {
			oracleEntry.RawSetString("PathRef", pRef)
		}
	}
}

func GetPathEntry(scoreTable *lua.LTable, fp, dst string) *lua.LTable {
	sDst := scoreTable.RawGetString(dst)
	if sDst == lua.LNil {
		return nil
	}

	var res *lua.LTable
	sDst.(*lua.LTable).ForEach(func(_ lua.LValue, v lua.LValue) {
		if e, ok := v.(*lua.LTable); ok {
			if e.RawGetString("Fingerprint").String() == fp {
				res = e
			}
		}
	})

	return res
}
