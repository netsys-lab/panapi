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
         print(formatting .. tostring(v))
      end
   end
end

-- global "path database variable",
-- not directly referenced from go
paths = {}
print("Hello Ranking SelectionServer")


function rankpaths(raddr)
   table.sort(
      paths[raddr],
      function(path_a, path_b)
         return path_a.Expiry < path_b.Expiry
      end
   )
   table.sort(
      paths[raddr],
      function(path_a, path_b)
         return #path_a.Metadata.Interfaces < #path_b.Metadata.Interfaces
      end
   )
end

-- gets called when a set of paths to addr is known
function panapi.initialize(laddr, raddr, ps)
   panapi.log("initialize", laddr, raddr)
   paths[raddr] = ps
   rankpaths(raddr)
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.selectpath(raddr)
   panapi.log("selectpath", raddr)
   if #paths[raddr] > 0 then
      return paths[raddr][1]
   end
end

-- gets called whenever a path disappears(?)
function panapi.pathdown(raddr, fp, pi)
   panapi.log("pathdown called with", raddr, fp, pi)
   for i,path in ipairs(paths[raddr]) do
      if path.Fingerprint == fp then
         panapi.log("found path at index " .. tostring(i) .. ", removing")
         table.remove(paths[raddr], i)
         break
      end
   end
end

function panapi.refresh(raddr, ps)
   panapi.log("refresh", raddr, ps)
   paths[raddr] = ps
   rankpaths(raddr)
   panapi.log("debug! exiting for introspection", #paths[raddr])
   os.exit(0)

end


function panapi.close(raddr)
   panapi.log("close", raddr)
   paths[raddr] = nil
end

function panapi.periodic(second)

end
