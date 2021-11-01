package scion

import (
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"log"
	"sync"
)

type LuaSelector struct {
	mutex sync.Mutex
	L     *lua.LState
}

func NewLuaSelector(script string) (*LuaSelector, error) {
	//initialize Lua VM
	L := lua.NewState()
	//load script
	if err := L.DoFile(script); err != nil {
		return nil, err
	}
	log.Printf("loaded selector from file %s", script)
	selector := LuaSelector{L: L}
	return &selector, nil
}

func (s *LuaSelector) Path() *pan.Path {
	//log.Println("Path()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "selectpath" function from the Lua script
	//expect 1 return value
	s.L.CallByParam(lua.P{Fn: s.L.GetGlobal("selectpath"), NRet: 1})
	//convert top of the stack back to a UserData type
	lv := s.L.ToUserData(-1)
	//pop element from the stack
	s.L.Pop(1)
	//try casting the return value back to a *pan.Path
	if p, ok := lv.Value.(*pan.Path); ok {
		//		log.Printf("lua returned path %v", p)
		return p
	} else {
		//couldn't be casted
		panic("something went wrong with Lua")
	}
}

func (s *LuaSelector) SetPaths(addr pan.UDPAddr, paths []*pan.Path) {
	log.Println("SetPaths()")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//call the "setpaths" function in the Lua script
	//with two arguments
	//and don't expect a return value
	s.L.CallByParam(lua.P{
		Fn:   s.L.GetGlobal("setpaths"),
		NRet: 0}, luar.New(s.L, addr), luar.New(s.L, paths))
}

func (s *LuaSelector) OnPathDown(fp pan.PathFingerprint, pi pan.PathInterface) {
	log.Println("OnPathDown()")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//call the "onpathdown" function in the Lua script
	//with two arguments
	//and don't expect a return value
	s.L.CallByParam(lua.P{
		Fn:   s.L.GetGlobal("onpathdown"),
		NRet: 0}, luar.New(s.L, fp), luar.New(s.L, pi))

	log.Printf("OnPathDown called with fp %v and pi %v", fp, pi)
}

func (s *LuaSelector) Close() error {
	log.Println("Close called on LuaSelector")
	s.L.Close()
	return nil
}
