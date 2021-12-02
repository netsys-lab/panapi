package lua

import (
	"fmt"
	"log"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/rpc"
	"github.com/yuin/gopher-lua"
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
	*State
	state
	mod *lua.LTable
	d   time.Duration
}

//func NewLuaSelector(script string) (*LuaSelector, error) {
func NewSelector(state *State) rpc.ServerSelector {
	state.Lock()
	defer state.Unlock()

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
			state.Logger.Panic(s)
			return 0
		}
	}

	mod["log"] = func(L *lua.LState) int {
		s := ""
		for i := 1; i <= L.GetTop(); i++ {
			s += L.Get(i).String() + " "
		}
		state.Println(s)
		return 0
	}

	panapi := state.RegisterModule("panapi", mod).(*lua.LTable)

	s := &LuaSelector{state, new_state(), panapi, time.Second}

	go func(s *LuaSelector) {
		old := time.Now()
		for {
			time.Sleep(s.d)
			s.Lock()
			seconds := time.Since(old).Seconds()
			s.CallByParam(
				lua.P{
					Protect: true,
					Fn:      s.mod.RawGetString("periodic"),
					NRet:    0,
				},
				lua.LNumber(seconds),
			)
			old = time.Now()
			s.Unlock()
		}
	}(s)
	return s
}

func (s *LuaSelector) Initialize(local, remote pan.UDPAddr, paths []*pan.Path) {
	s.Printf("Initialize(%s,%s,[%d]pan.Path)", local, remote, len(paths))
	s.Lock()
	defer s.Unlock()

	//assume that setpaths is called with all the currently valid options
	//meaning that anything we already know can be flushed
	s.state.clear_addr(remote)
	lpaths := s.set_paths(remote, paths)

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.CallByParam(
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
		s.Fatal("Initialize", err)
	}

}

func (s *LuaSelector) Path(raddr pan.UDPAddr) *pan.Path {
	//log.Println("Path()")
	s.Lock()
	defer s.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.CallByParam(lua.P{
		Protect: true,
		Fn:      s.mod.RawGetString("selectpath"),
		NRet:    1}, lua.LString(raddr.String()))
	if err != nil {
		s.Fatal("Path", err)
	}
	lt := s.ToTable(-1)
	//pop element from the stack
	s.Pop(1)
	return s.state.get_pan_path(lt)
}

func (s *LuaSelector) PathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	//s.l.Println("PathDown()")
	s.Lock()
	defer s.Unlock()
	err := s.CallByParam(
		lua.P{
			Fn:      s.mod.RawGetString("pathdown"),
			NRet:    0,
			Protect: true,
		},
		lua.LString(raddr.String()),
		lua.LString(fp),
		new_lua_path_interface(pi),
	)
	s.Printf("PathDown called with fp %v and pi %v: %s", fp, pi, err)
	if err != nil {
		s.Fatal(err)
	}

}

func (s *LuaSelector) Refresh(remote pan.UDPAddr, paths []*pan.Path) {
	s.Println("Refresh()")
	s.Lock()
	defer s.Unlock()

	//assume that setpaths is called with all the currently valid options
	//meaning that anything we already know can be flushed
	s.state.clear_addr(remote)
	lpaths := s.state.set_paths(remote, paths)

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.CallByParam(
		lua.P{
			Protect: true,
			Fn:      s.mod.RawGetString("refresh"),
			NRet:    0,
		},
		lua.LString(remote.String()),
		lua_table_slice_to_table(lpaths),
	)
	if err != nil {
		s.Fatal("refresh", err)
	}
}

func (s *LuaSelector) Close(raddr pan.UDPAddr) error {
	s.Lock()
	defer s.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.CallByParam(
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
