
-- Print contents of `tbl`, with indentation.
-- `indent` sets the initial level of indentation.
function tprint (tbl, indent)
  if not indent then indent = 0 end
  for k, v in pairs(tbl) do
    formatting = string.rep("  ", indent) .. k .. ": "
    if type(v) == "table" then
      print(formatting)
      tprint(v, indent+1)
    else
      print(formatting .. v)
    end
  end
end

-- global "path database variable",
-- not directly referenced from go
paths = {}
print("Hello Simple SelectionServer")

-- gets called when a set of paths to addr is known
function setpaths(addr, ps)
   print("setpath", tostring(addr))
   tprint(ps)
   paths = ps
end

-- gets called for every packet
-- implementation needs to be efficient
function selectpath(addr)
   print("selectpath", tostring(addr))
   return paths[1]
end

-- gets called whenever a path disappears(?)
function onpathdown(addr, fp, pi)
   print(string.format("lua output: onpathdown called with fp %s and pi %s",
                       tostring(fp),
                       tostring(pi)))
end

function close(addr)
   print("close", tostring(addr))
end
