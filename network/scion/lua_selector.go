package scion

import (
	"os"

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
}

//func NewLuaSelector(script string) (*LuaSelector, error) {
func NewLuaSelector(script string) (rpc.ServerSelector, error) {
	//load script
	file, err := os.Open(script)
	if err != nil {
		return nil, err
	}
	//initialize Lua VM
	L := lua.NewState()
	if fn, err := L.Load(file, script); err != nil {
		return nil, err
	} else {
		log.Printf("loaded selector from file %s", script)
		L.Push(fn)
		err = L.PCall(0, lua.MultRet, nil)
		return &LuaSelector{
			L:     L,
			state: new_state(),
		}, err
	}
}

func (s *LuaSelector) SetPaths(raddr pan.UDPAddr, paths []*pan.Path) {
	log.Println("SetPaths()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	s.state.clear_addr(raddr)
	lpaths := s.state.set_paths(raddr, paths)
	s.L.CallByParam(
		lua.P{
			Fn:   s.L.GetGlobal("setpaths"),
			NRet: 0,
		},
		lua.LString(raddr.String()),
		lua_table_slice_to_table(lpaths),
	)
}

func (s *LuaSelector) Path(raddr pan.UDPAddr) *pan.Path {
	//log.Println("Path()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	s.L.CallByParam(lua.P{
		Fn:   s.L.GetGlobal("selectpath"),
		NRet: 1}, lua.LString(raddr.String()))
	lt := s.L.ToTable(-1)
	//pop element from the stack
	s.L.Pop(1)
	return s.state.get_pan_path(lt)
}

func (s *LuaSelector) OnPathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("OnPathDown()")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//call the "onpathdown" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.L.CallByParam(
		lua.P{
			Fn:   s.L.GetGlobal("onpathdown"),
			NRet: 0,
		},
		lua.LString(raddr.String()),
		lua.LString(fp),
		new_lua_path_interface(pi),
	)

	log.Printf("OnPathDown called with fp %v and pi %v: %s", fp, pi, err)
}

func (s *LuaSelector) Close(raddr pan.UDPAddr) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.L.CallByParam(
		lua.P{Fn: s.L.GetGlobal("close"),
			NRet: 1},
		lua.LString(raddr.String()))

	log.Println("Close called on LuaSelector:", err)
	//s.L.Close()
	return err
}
