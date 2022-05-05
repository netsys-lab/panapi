print("HELLO FROM LUA DE(A)MO(N) SCRIPT")

-- global "path database variable",
-- not directly referenced from go
paths = {}

-- global table to keep function call statistics
calls = {
   cur = {},
   old = {},
}

-- global table to keep track of connections
conns = {}

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
   panapi.Log("New connection between", laddr, "and", raddr .. ", Profile:", prefs.ConnCapacityProfile)
   calls.cur.Initialize = (calls.cur.Initialize or 0) + 1
   paths[raddr] = ps
   conns[raddr] = conns[raddr] or {}
   conns[raddr][laddr] = prefs
   rankpaths(raddr)
end

function panapi.SetPreferences(prefs, laddr, raddr)
   calls.cur.SetPreferences = (calls.cur.SetPreferences or 0) + 1
   panapi.Log("SetPreferences, Profile:", prefs,ConnCapacityProfile)
   conns[raddr][laddr] = prefs
   
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   calls.cur.Path = (calls.cur.Path or 0) + 1
   --panapi.Log("Path", laddr, raddr)
   if #paths[raddr] > 0 then
      panapi.Log(tprint(conns))
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
   calls.cur.PathDown = (calls.cur.PathDown or 0) + 1
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
   calls.cur.Refresh = (calls.cur.Refresh or 0) + 1
   panapi.Log("Refresh", raddr, ps)
   paths[raddr] = ps
   rankpaths(raddr)
   panapi.Log("debug! exiting for introspection", #paths[raddr])
   os.exit(0)

end


function panapi.Close(laddr, raddr)
   calls.cur.Close = (calls.cur.Close or 0) + 1
   panapi.Log("Close", laddr, raddr)
   --paths[raddr] = nil
end

function diffcallstats()
   local d = {}
   for k,v in pairs(calls.cur) do
      local old = calls.old[k]
      if v ~= old then
         d[k] = v - (old or 0)
      end
   end
   calls.old = copy(calls.cur)
   panapi.Log("function calls since last tick:\n" .. tprint(d, 2))
end

function abscallstats()
   local d = {}
   for k,v in pairs(calls.cur) do
      if v ~= 0 then
         d[k] = v
      end
   end
   panapi.Log("function calls:\n" .. tprint(d, 2))
end


function panapi.Periodic(seconds)
   calls.cur.Periodic = (calls.cur.Periodic or 0) + 1
   --abscallstats()
end


--[[function stats.TracerForConnection(id, p, odcid)
   calls.cur.TracerForConnection = (calls.cur.TracerForConnection or 0) + 1
   --panapi.Log("id:", id, "perspective", p, "odcid", odcid)
end
function stats.StartedConnection(laddr, raddr, srcid, dstid)
   calls.cur.StartedConnection = (calls.cur.StartedConnection or 0) + 1
   
end
function stats.NegotiatedVersion(laddr, raddr)
   calls.cur.NegotiatedVersion = (calls.cur.NegotiatedVersion or 0) + 1

end
function stats.ClosedConnection(laddr, raddr)
   calls.cur.ClosedConnection = (calls.cur.ClosedConnection or 0) + 1

end
function stats.SentTransportParameters(laddr, raddr)
   calls.cur.SentTransportParameters = (calls.cur.SentTransportParameters or 0) + 1

end
function stats.ReceivedTransportParameters(laddr, raddr)
   calls.cur.ReceivedTransportParameters = (calls.cur.ReceivedTransportParameters or 0) + 1

end
function stats.RestoredTransportParameters(laddr, raddr)
   calls.cur.RestoredTransportParameters = (calls.cur.RestoredTransportParameters or 0) + 1

end
function stats.SentPacket(laddr, raddr)
   calls.cur.SentPacket = (calls.cur.SentPacket or 0) + 1

end
function stats.ReceivedVersionNegotiationPacket(laddr, raddr)
   calls.cur.ReceivedVersionNegotiationPacket = (calls.cur.ReceivedVersionNegotiationPacket or 0) + 1

end
function stats.ReceivedRetry(laddr, raddr)
   calls.cur.ReceivedRetry = (calls.cur.ReceivedRetry or 0) + 1

end
function stats.ReceivedPacket(laddr, raddr)
   calls.cur.ReceivedPacket = (calls.cur.ReceivedPacket or 0) + 1

end
function stats.BufferedPacket(laddr, raddr)
   calls.cur.BufferedPacket = (calls.cur.BufferedPacket or 0) + 1

end
function stats.DroppedPacket(laddr, raddr)
   calls.cur.DroppedPacket = (calls.cur.DroppedPacket or 0) + 1

end
function stats.UpdatedMetrics(laddr, raddr, rttStats, cwnd, bytesInFlight, packetsInFlight)
   calls.cur.UpdatedMetrics = (calls.cur.UpdatedMetrics or 0) + 1
   --panapi.Log("UpdatedMetrics", cwnd, bytesInFlight, packetsInFlight)
   --panapi.Log("\n", tprint(rttStats, 1))
end
function stats.AcknowledgedPacket(laddr, raddr)
   calls.cur.AcknowledgedPacket = (calls.cur.AcknowledgedPacket or 0) + 1

end
function stats.LostPacket(laddr, raddr)
   calls.cur.LostPacket = (calls.cur.LostPacket or 0) + 1

end
function stats.UpdatedCongestionState(laddr, raddr)
   calls.cur.UpdatedCongestionState = (calls.cur.UpdatedCongestionState or 0) + 1

end
function stats.UpdatedPTOCount(laddr, raddr)
   calls.cur.UpdatedPTOCount = (calls.cur.UpdatedPTOCount or 0) + 1

end
function stats.UpdatedKeyFromTLS(laddr, raddr)
   calls.cur.UpdatedKeyFromTLS = (calls.cur.UpdatedKeyFromTLS or 0) + 1

end
function stats.UpdatedKey(laddr, raddr)
   calls.cur.UpdatedKey = (calls.cur.UpdatedKey or 0) + 1

end
function stats.DroppedEncryptionLevel(laddr, raddr)
   calls.cur.DroppedEncryptionLevel = (calls.cur.DroppedEncryptionLevel or 0) + 1

end
function stats.DroppedKey(laddr, raddr)
   calls.cur.DroppedKey = (calls.cur.DroppedKey or 0) + 1

end
function stats.SetLossTimer(laddr, raddr)
   calls.cur.SetLossTimer = (calls.cur.SetLossTimer or 0) + 1

end
function stats.LossTimerExpired(laddr, raddr)
   calls.cur.LossTimerExpired = (calls.cur.LossTimerExpired or 0) + 1

end
function stats.LossTimerCanceled(laddr, raddr)
   calls.cur.LossTimerCanceled = (calls.cur.LossTimerCanceled or 0) + 1

end
function stats.Debug(laddr, raddr)
   calls.cur.Debug = (calls.cur.Debug or 0) + 1
   end
]]--

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

-- recursively perform a deep copy of a table
function copy(thing)
   if type(thing) == "table" then
      local r = {}
      for k,v in pairs(thing) do
         r[k] = copy(v)
      end
      return r
   else
      return thing
   end
end

-- return a new thing containing everything about thing1 that is different from thing2
function diff(thing1, thing2)
   if type(thing1) == "table" then
      local thing = {}
      if type(thing2) == "table" then
         local thing = {}
         for k,v in pairs(thing1) do
            thing[k] = diff(thing1[k], thing2[k])
         end
         return thing
      else
         return copy(thing1)
      end
   else
      if thing1 == thing2 then
         return nil
      else
         return thing1
      end
   end
end

