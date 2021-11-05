
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

local pathstructure = {
   Source = "",
   Destination = "",
   Fingerprint = "",
   Expiry = "",
   Metadata = {
      Interfaces = {
         {
            IA = "",
            IfID = 0,
         },
      },
      MTU = 0,
      Latency = { 0 },
      Bandwidth = { 0 },
      Geo = {
         {
            Latitude = 0,
            Longitude = 0,
            Address = "",
         },
      },
      LinkType = { 0 },
      InternalHops = { 0 },
      Notes = {
         "",
         },
      },
}

-- parse userdata into lua-native table structure
function gopath2luapath(userdata, structure)
   if type(structure) == "string" then
      return tostring(userdata)
   end
   if type(structure) == "number" then
      return tonumber(userdata)
   end
   if type(structure) == "boolean" then
      return tostring(userdata) == "true"
   end
   if type(structure) == "table" then
      local t = {}
      if structure[1] ~= nil then
         -- structure is "array"
         for index in userdata() do
            table.insert(t, gopath2luapath(userdata[index], structure[1]))
         end
      else
         for key, substructure in pairs(structure) do
            t[key] = gopath2luapath(userdata[key], substructure)
         end
      end
      return t
   end
   return "error, something went wrong"
end


-- global "path database variable",
-- not directly referenced from go
paths = {}
ranking = {}
print("Hello PathRanking")


function rankpaths()
   table.sort(
      ranking,
      function(a, b)
         return #paths[a].Metadata.Interfaces < #paths[b].Metadata.Interfaces
         --return paths[a].SumBandwidth > paths[b].SumBandwidth
      end
   )
   --[[table.sort(
      ranking,
      function(a, b)
         return paths[a].Expiry < paths[b].Expiry
         --return paths[a].SumLatency < paths[b].SumLatency
      end
      )]]
end


-- gets called when a set of paths to addr is known
function setpaths(addr, ps)
   --tprint(gopath2luapath(ps, pathstructure))
   for index in ps() do
      path = ps[index]
      luapath = gopath2luapath(path, pathstructure)
      
      -- sum over bandwidth
      local bw = 0
      for _,mbw in ipairs(luapath.Metadata.Bandwidth) do
         bw = bw+mbw
      end
      luapath.SumBandwidth = bw

      -- sum over latency
      local ls = 0
      for _,l in ipairs(luapath.Metadata.Latency) do
         ls = ls+l
      end
      luapath.SumLatency = ls
      
      paths[path] = luapath
      -- debug print
      tprint(paths[path])

      table.insert(ranking, path)
   end
   print(string.format("lua output: setpaths called with %s and %d paths", addr, #ranking))
   rankpaths()
end

-- gets called for every packet
-- implementation needs to be efficient
function selectpath()
   if #ranking > 0 then
      path = ranking[1]
      --print("lua output: selecting path with soonest expiry", paths[path].Expiry)
      --tprint(paths[path])
      return path
   else
      print("lua output: couldn't select a path")
   end
end

-- gets called whenever a path disappears(?)
-- TODO tidy this up
function onpathdown(fp, pi)
   print(string.format("lua output: onpathdown called with fp %s and pi %s",
                       tostring(fp),
                       tostring(pi)))
   if paths[pi] == nil then
      print("path is not known in database")
      for ud,path in pairs(paths) do
         if path.Fingerprint == tostring(fp) then
            print("found path via fingerprint")
            paths[ud] = nil
            for i,candidate in ipairs(ranking) do
               if candidate == path then
                  table.remove(ranking, i)
               end
            end
         end
      end
   else
      print("path is known in database, deleting")
      paths[pi] = nil
      for i,candidate in ipairs(ranking) do
         if candidate == pi then
            table.remove(ranking, i)
            break
         end
      end
   end
   print("deleted path, exiting for introspection")
   os.exit(0)
   
   -- todo, remove path from "database"
end
