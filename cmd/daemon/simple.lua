-- global "path database variable",
-- not directly referenced from go
paths = {}
print("Hello Simple SelectionServer")

-- gets called when a set of paths to addr is known
function panapi.initialize(laddr, raddr, ps)
   panapi.log("initialize")
   paths = ps
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.selectpath(addr)
   panapi.log("selectpath", addr)
   return paths[1]
end

-- gets called whenever a path disappears(?)
function panapi.pathdown(addr, fp, pi)
   panapi.log("pathdown called with", fp, pi)
end

function panapi.refresh(addr, ps)
   panapi.log("refresh", addr, ps)
   paths = ps
end

function panapi.close(addr)
   panapi.log("close", addr)
end



