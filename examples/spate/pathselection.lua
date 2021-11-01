-- global "path database variable",
-- not directly referenced from go
paths = {}
len = 0
calls = 0
print("Hello SPATE")

-- gets called when a set of paths to addr is known
function setpaths(addr, ps)
   print("lua output: setpaths called with", addr, "and", ps)
   for index in ps() do
      path = ps[index]
      --printpath(path)
      table.insert(paths, path)
      len = len + 1
  end
end

-- gets called for every packet
-- implementation needs to be efficient
function selectpath()
   calls = calls + 1
   if #paths > 0 then
      -- simply return the first path we have in the "database"
      i = math.random(len)
      print(string.format("lua output: %dth call to selectpath, returning: %dth path", calls, i))
      return paths[i]
   end
end

-- gets called whenever a path disappears(?)
function onpathdown(fp, pi)
   print(string.format("lua output: onpathdown called with fp %s and pi %s",
                       tostring(fp),
                       tostring(pi)))
   -- todo, remove path from "database"
end


-- just a helper function to print out everything we know
-- about a path
function printpath(path)
   print(string.format([[lua output: got path %s with
   source %s
   destination %s
   forwardingpath %s
   metadata %s
   fingerprint %s
   expiry %s]],
   tostring(path),
   tostring(path.Source),
   tostring(path.Destination),
   tostring(path.ForwardingPath),
   tostring(path.Metadata),
   tostring(path.Fingerprint),
   tostring(path.Expiry)))
end

