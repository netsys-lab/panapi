package scion

import (
	"os"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsys-lab/panapi/rpc"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"

	"log"
	"sync"
)

func newLuaPath(path *pan.Path) *lua.LTable {
	t := lua.LTable{}
	if path != nil {
		t.RawSetString("Source", lua.LString(path.Source.String()))
		t.RawSetString("Destination", lua.LString(path.Destination.String()))
		t.RawSetString("Fingerprint", lua.LString(path.Fingerprint))
		t.RawSetString("Expiry", lua.LString(path.Expiry.String()))

		if path.Metadata != nil {
			meta := lua.LTable{}
			meta.RawSetString("MTU", lua.LString(path.Metadata.MTU))

			ifaces := lua.LTable{}
			for _, i := range path.Metadata.Interfaces {
				iface := lua.LTable{}
				iface.RawSetString("IA", lua.LString(i.IA.String()))
				iface.RawSetString("IfID", lua.LNumber(i.IfID))
				ifaces.Append(&iface)
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

func newTableFromTableSlice(s []*lua.LTable) *lua.LTable {
	res := lua.LTable{}
	for _, t := range s {
		res.Append(t)
	}
	return &res
}

type state struct {
	lpaths map[string]map[string]*lua.LTable
	ppaths map[string]map[string]*pan.Path
}

func (s state) get_pan_path(raddr pan.UDPAddr, fp string) *pan.Path {
	if tmp := s.ppaths[raddr.String()]; tmp != nil {
		return tmp[fp]
	}
	return nil
}

func (s state) clear_addr(addr pan.UDPAddr) {
	raddr := addr.String()
	s.lpaths[raddr] = map[string]*lua.LTable{}
	s.ppaths[raddr] = map[string]*pan.Path{}
}

func (s state) set_pan_paths(addr pan.UDPAddr, paths []*pan.Path) {
	raddr := addr.String()
	for _, p := range paths {
		fp := string(p.Fingerprint)
		s.ppaths[raddr][fp] = p
		s.lpaths[raddr][fp] = newLuaPath(p)
	}
}

func (s state) set_paths(addr pan.UDPAddr, ppaths []*pan.Path, lpaths []*lua.LTable) {
	if len(ppaths) != len(lpaths) {
		log.Panicf("len(ppaths) != len(lpaths) (%d != %d)", len(ppaths), len(lpaths))
	}
	raddr := addr.String()
	for i, ppath := range ppaths {
		fp_p := string(ppath.Fingerprint)
		fp_l := lpaths[i].RawGetString("Fingerprint").String()
		if fp_p != fp_l {
			log.Panicf("ppath.Fingerprint != lpaths.FingerPrint (%s != %s)", fp_p, fp_l)
		}
		s.ppaths[raddr][fp_p] = ppath
		s.lpaths[raddr][fp_l] = lpaths[i]
	}
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
			L: L,
			state: state{
				make(map[string]map[string]*lua.LTable),
				make(map[string]map[string]*pan.Path),
			},
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
	/*s.L.CallByParam(lua.P{
	Fn:   s.L.GetGlobal("setpaths"),
	NRet: 0}, luar.New(s.L, raddr), luar.New(s.L, paths))
	*/
	s.state.clear_addr(raddr)
	lpaths := make([]*lua.LTable, len(paths))
	for i, p := range paths {
		lpaths[i] = newLuaPath(p)
	}
	s.state.set_paths(raddr, paths, lpaths)
	s.L.CallByParam(
		lua.P{
			Fn:   s.L.GetGlobal("setpaths"),
			NRet: 0,
		},
		lua.LString(raddr.String()),
		newTableFromTableSlice(lpaths),
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
		NRet: 1}, luar.New(s.L, raddr))
	//convert top of the stack back to a UserData type
	lt := s.L.ToTable(-1)
	//lv := s.L.ToUserData(-1)
	//pop element from the stack
	s.L.Pop(1)
	//try casting the return value back to a *pan.Path
	/*if p, ok := lv.Value.(*pan.Path); ok {
		//		log.Printf("lua returned path %v", p)
		return p
	} else {
		//couldn't be casted
		panic("something went wrong with Lua")
	}*/
	tmp := s.state.ppaths[raddr.String()]
	if tmp == nil {
		panic("got nothing")
	}
	return tmp[lt.RawGetString("Fingerprint").String()]
}

func (s *LuaSelector) OnPathDown(raddr pan.UDPAddr, fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("OnPathDown()")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//call the "onpathdown" function in the Lua script
	//with two arguments
	//and don't expect a return value
	err := s.L.CallByParam(lua.P{
		Fn:   s.L.GetGlobal("onpathdown"),
		NRet: 0}, luar.New(s.L, raddr), luar.New(s.L, fp), luar.New(s.L, pi))

	log.Printf("OnPathDown called with fp %v and pi %v: %s", fp, pi, err)
}

func (s *LuaSelector) Close(raddr pan.UDPAddr) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	err := s.L.CallByParam(lua.P{
		Fn:   s.L.GetGlobal("close"),
		NRet: 1}, luar.New(s.L, raddr))

	log.Println("Close called on LuaSelector:", err)
	//s.L.Close()
	return err
}
