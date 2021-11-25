package scion

import (
	"fmt"
	"os"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/rpc"
	"github.com/yuin/gopher-lua"

	"log"
	"sync"
)

func new_lua_path_interface(intf pan.PathInterface) *lua.LTable {
	iface := lua.LTable{}
	iface.RawSetString("IA", lua.LString(intf.IA.String()))
	iface.RawSetString("IfID", lua.LNumber(intf.IfID))
	return &iface

}

func new_lua_path(path *pan.Path) *lua.LTable {
	t := lua.LTable{}
	if path != nil {
		t.RawSetString("Source", lua.LString(path.Source.String()))
		t.RawSetString("Destination", lua.LString(path.Destination.String()))
		t.RawSetString("Fingerprint", lua.LString(path.Fingerprint))
		t.RawSetString("Expiry", lua.LString(path.Expiry.String()))

		if path.Metadata != nil {
			meta := lua.LTable{}
			meta.RawSetString("MTU", lua.LNumber(path.Metadata.MTU))

			ifaces := lua.LTable{}
			for _, i := range path.Metadata.Interfaces {
				ifaces.Append(new_lua_path_interface(i))
			}
			meta.RawSetString("Interfaces", &ifaces)

			latencies := lua.LTable{}
			for _, l := range path.Metadata.Latency {
				latencies.Append(lua.LNumber(l))
			}
			meta.RawSetString("Latency", &latencies)

			bandwidths := lua.LTable{}
			for _, b := range path.Metadata.Bandwidth {
				bandwidths.Append(lua.LNumber(b))
			}
			meta.RawSetString("Bandwidth", &bandwidths)

			linktypes := lua.LTable{}
			for _, l := range path.Metadata.LinkType {
				linktypes.Append(lua.LNumber(l))
			}
			meta.RawSetString("LinkType", &linktypes)

			internalhops := lua.LTable{}
			for _, h := range path.Metadata.InternalHops {
				internalhops.Append(lua.LNumber(h))
			}
			meta.RawSetString("InternalHops", &internalhops)

			notes := lua.LTable{}
			for _, n := range path.Metadata.Notes {
				notes.Append(lua.LString(n))
			}
			meta.RawSetString("Notes", &notes)

			geo := lua.LTable{}
			for _, g := range path.Metadata.Geo {
				pos := lua.LTable{}
				pos.RawSetString("Latitude", lua.LNumber(g.Latitude))
				pos.RawSetString("Longitude", lua.LNumber(g.Longitude))
				pos.RawSetString("Address", lua.LString(g.Address))
				geo.Append(&pos)
			}
			meta.RawSetString("Geo", &geo)

			t.RawSetString("Metadata", &meta)
		}
	}
	return &t
}

func lua_table_slice_to_table(s []*lua.LTable) *lua.LTable {
	res := lua.LTable{}
	for _, t := range s {
		res.Append(t)
	}
	return &res
}

// help to translate lua to pan pointers back and forth
type state struct {
	lpaths map[string]map[string]*lua.LTable
	ppaths map[*lua.LTable]*pan.Path
}

func new_state() state {
	return state{
		make(map[string]map[string]*lua.LTable),
		make(map[*lua.LTable]*pan.Path),
	}
}

func (s state) get_pan_path(lpath *lua.LTable) *pan.Path {
	return s.ppaths[lpath]
}

func (s state) clear_addr(addr pan.UDPAddr) {
	raddr := addr.String()
	for _, lt := range s.lpaths[raddr] {
		s.ppaths[lt] = nil
	}
	s.lpaths[raddr] = map[string]*lua.LTable{}
}

func (s state) set_paths(addr pan.UDPAddr, ppaths []*pan.Path) (lpaths []*lua.LTable) {
	raddr := addr.String()
	lpaths = make([]*lua.LTable, len(ppaths))
	for i, ppath := range ppaths {
		lpath := new_lua_path(ppath)
		s.lpaths[raddr][string(ppath.Fingerprint)] = lpath
		s.ppaths[lpath] = ppath
		lpaths[i] = lpath
	}
	return
}

type LuaSelector struct {
	mutex sync.Mutex
	L     *lua.LState
	state state
	l     *log.Logger
	mod   *lua.LTable
	d     time.Duration
}

//func NewLuaSelector(script string) (*LuaSelector, error) {
func NewLuaSelector(script string) (rpc.ServerSelector, error) {
	//load script
	file, err := os.Open(script)
	if err != nil {
		return nil, err
	}
	l := log.Default() //.New(os.Stderr, "lua", log.Ltime|log.Lshortfile)
	l.SetFlags(log.Ltime)
	l.SetPrefix("lua ")

	//initialize Lua VM
	L := lua.NewState()
	mod := map[string]lua.LGFunction{}
	for _, fn := range []string{
		"initialize",
		"selectpath",
		"pathdown",
		"refresh",
		"close",
		"periodic",
	} {
		s := fmt.Sprintf("function %s not implemented in script", fn)
		mod[fn] = func(L *lua.LState) int {
			l.Panic(s)
			return 0
		}
	}

	mod["log"] = func(L *lua.LState) int {
		s := ""
		for i := 1; i <= L.GetTop(); i++ {
			s += L.Get(i).String() + " "
		}
		l.Println(s)
		return 0
	}

	panapi := L.RegisterModule("panapi", mod).(*lua.LTable)

	if fn, err := L.Load(file, script); err != nil {
		return nil, err
	} else {
		l.Printf("loaded selector from file %s", script)
		L.Push(fn)
		err = L.PCall(0, lua.MultRet, nil)
		if err != nil {
			return nil, err
		}
		s := &LuaSelector{
			L:     L,
			state: new_state(),
			l:     l,
			mod:   panapi,
			d:     time.Second,
		}

		go func(s *LuaSelector) {
			old := time.Now()
			for {
				time.Sleep(s.d)
				s.mutex.Lock()
				seconds := time.Since(old).Seconds()
				s.L.CallByParam(
					lua.P{
						Protect: true,
						Fn:      s.mod.RawGetString("periodic"),
						NRet:    0,
					},
					lua.LNumber(seconds),
				)
				old = time.Now()
				s.mutex.Unlock()
			}
		}(s)
		return s, err
	}
}

func (s *LuaSelector) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	s.l.Printf("Initialize(%s,%s,[%d]pan.Path)", local, remote, len(paths))
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//assume that setpaths is called with all the currently valid options
	//meaning that anything we already know can be flushed
	s.state.clear_addr(remote)
	lpaths := s.state.set_paths(remote, paths)

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.L.CallByParam(
		lua.P{
			Protect: true,
			Fn:      s.mod.RawGetString("initialize"),
			NRet:    0,
		},
		lua.LString(local.String()),
		lua.LString(remote.String()),
		lua_table_slice_to_table(lpaths),
	)
	if err != nil {
		s.l.Fatal("Initialize", err)
	}

}

func (s *LuaSelector) Path(raddr pan.UDPAddr) *pan.Path {
	//log.Println("Path()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.L.CallByParam(lua.P{
		Protect: true,
		Fn:      s.mod.RawGetString("selectpath"),
		NRet:    1}, lua.LString(raddr.String()))
	if err != nil {
		s.l.Fatal("Path", err)
	}
	lt := s.L.ToTable(-1)
	//pop element from the stack
	s.L.Pop(1)
	return s.state.get_pan_path(lt)
}

func (s *LuaSelector) PathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	//s.l.Println("PathDown()")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.L.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("pathdown"),
			NRet:    0,
			Protect: true,
		},
		lua.LString(raddr.String()),
		lua.LString(fp),
		new_lua_path_interface(pi),
	)
	s.l.Printf("PathDown called with fp %v and pi %v: %s", fp, pi, err)
	if err != nil {
		s.l.Fatal(err)
	}

}

func (s *LuaSelector) Refresh(remote pan.UDPAddr, paths []*pan.Path) {
	s.l.Println("Refresh()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//assume that setpaths is called with all the currently valid options
	//meaning that anything we already know can be flushed
	s.state.clear_addr(remote)
	lpaths := s.state.set_paths(remote, paths)

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.L.CallByParam(
		lua.P{
			Protect: true,
			Fn:      s.mod.RawGetString("refresh"),
			NRet:    0,
		},
		lua.LString(remote.String()),
		lua_table_slice_to_table(lpaths),
	)
	if err != nil {
		s.l.Fatal("refresh", err)
	}
}

func (s *LuaSelector) Close(raddr pan.UDPAddr) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.L.CallByParam(
		lua.P{
			Protect: true,
			Fn:      s.mod.RawGetString("close"),
			NRet:    1,
		},
		lua.LString(raddr.String()))

	log.Println("Close called on LuaSelector:", err)
	//s.L.Close()
	return err
}
