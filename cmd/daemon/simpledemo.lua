print("HELLO FROM SIMPLE LUA DE(A)MO(N) SCRIPT")

-- global "path database variable",
-- not directly referenced from go
paths = {}

-- global table to keep track of connections
conns = {}

-- global table to keep track of connection identites
connids = {}

-- helper function
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
function panapi.Initialize(prefs, laddr, raddr, ps)
   panapi.Log("New connection [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)
   paths[raddr] = ps
   conns[raddr] = conns[raddr] or {}
   conns[raddr][laddr] = prefs
   rankpaths(raddr)
end

function panapi.SetPreferences(prefs, laddr, raddr)
   panapi.Log("Update Preferences [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)
   conns[raddr][laddr] = prefs
   
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   --panapi.Log("Path", laddr, raddr)
   if #paths[raddr] > 0 then
      if conns[raddr][laddr]["ConnCapacityProfile"] == "Scavenger" then
         p = paths[raddr][#paths[raddr]]
         return p
      else
         return paths[raddr][1]
      end
   end
end

-- gets called whenever a path disappears(?)
function panapi.PathDown(laddr, raddr, fp, pi)
   panapi.Log("PathDown called with", laddr, raddr, fp, pi)
   for i,path in ipairs(paths[raddr]) do
      if path.Fingerprint == fp then
         panapi.Log("found path at index " .. tostring(i) .. ", removing")
         table.remove(paths[raddr], i)
         break
      end
   end
end

function panapi.Refresh(laddr, raddr, ps)
   panapi.Log("Refresh", raddr, ps)
   paths[raddr] = ps
   rankpaths(raddr)
end


function panapi.Close(laddr, raddr)
   panapi.Log("Close", laddr, raddr)
   conns[raddr][laddr] = nil
   --paths[raddr] = nil
end


function panapi.Periodic(seconds)
   --panapi.Log("Tick", seconds)
end

-- HELPER FUNCTIONS ---
-- 
-- Print contents of `tbl`, with indentation.
-- `indent` sets the initial level of indentation.
function tprint (tbl, indent)
   if not indent then indent = 0 end
   if type(tbl) == "table" then
      local s = ""
      for k, v in pairs(tbl) do
         formatting = string.rep("  ", indent) .. k .. ": "
         if type(v) == "table" then
            --print(formatting)
            s = s ..  formatting .. "\n" .. tprint(v, indent+1)
         else
            s = s .. formatting .. tprint(v) .. "\n"
         end
      end
      return s
   else
      return tostring(tbl)
   end
end


