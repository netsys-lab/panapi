-- global "path database variable",
-- not directly referenced from go
paths = {}
print("Hello Simple SelectionServer")

-- gets called when a set of paths to addr is known
function panapi.setpaths(addr, ps)
   panapi.log("setpath", addr)
   tprint(ps)
   paths = ps
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.selectpath(addr)
   panapi.log("selectpath", addr)
   return paths[1]
end

-- gets called whenever a path disappears(?)
function panapi.onpathdown(addr, fp, pi)
   panapi.log(string.format("lua output: onpathdown called with fp %s and pi %s", fp, pi))
end

function panapi.close(addr)
   panapi.log("close", addr)
end



