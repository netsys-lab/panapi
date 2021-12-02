package lua

import (
	//"io/ioutil"
	"log"
	"os"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type State struct {
	*lua.LState
	sync.Mutex
	*log.Logger
}

func NewState() *State {
	L := lua.NewState()
	//l := log.New(ioutil.Discard, "lua ", log.Ltime)
	l := log.Default()
	l.SetFlags(log.Ltime | log.Lshortfile)
	l.SetPrefix("lua ")
	return &State{L, sync.Mutex{}, l}
}

func (s *State) LoadScript(fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	if fn, err := s.Load(file, fname); err != nil {
		return err
	} else {
		s.Printf("loaded selector from file %s", fname)
		s.Push(fn)
		return s.PCall(0, lua.MultRet, nil)
	}
}
