print("HELLO FROM THE SCAVENGER SCRIPT")

-- global variable to track "time"
tick = 0

function panapi.Periodic(seconds)
   --panapi.Log("Tick", seconds)
   tick = tick + seconds
end

-- global "path database variable",
-- not directly referenced from go
paths = {}


-- global table to keep track of connections
conns = {}

fingerprints = {}


function nextScavengePath(raddr)
   table.sort(
      fingerprints[raddr],
      function(fp_a, fp_b)
         return paths[raddr][fp_a].LatestTick < paths[raddr][fp_b].LatestTick
      end
   )
   return paths[raddr][fingerprints[raddr][1]].Path
end


function nextBestPath(raddr)
   table.sort(
      fingerprints[raddr],
      function(fp_a, fp_b)
         return (paths[raddr][fp_a].RTTStats.LatestRTT or 0) < (paths[raddr][fp_b].RTTStats.LatestRTT or 0)
      end
   )
   return paths[raddr][fingerprints[raddr][1]].Path
end





-- gets called when a set of paths to addr is known
function panapi.Initialize(prefs, laddr, raddr, ps)
   panapi.Log("New connection [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)
   paths[raddr] = paths[raddr] or {}
   fingerprints[raddr] = {}
   for _, path in ipairs(ps) do
      table.insert(fingerprints[raddr], path.Fingerprint)
      paths[raddr][path.Fingerprint] = {
         Path = path,
         LatestTick = tick,
         RTTStats = {},
      }
   end
   conns[raddr] = conns[raddr] or {}
   conns[raddr][laddr] = conns[raddr][laddr] or {}
   if prefs ~= nil then
      conns[raddr][laddr].Preferences = prefs
   end
end

function panapi.SetPreferences(prefs, laddr, raddr)
   panapi.Log("Update Preferences [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)
   conns[raddr][laddr].Preferences = prefs
--   panapi.Log("(Will return " .. tprint(panapi.Path(laddr, raddr).Fingerprint) .. ")")
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   --
   if next(paths[raddr]) ~= nil then
      local conn = conns[raddr][laddr]
      if conn.Preferences["ConnCapacityProfile"] == "Scavenger" then
         local p = nextScavengePath(raddr)
         conn.LastPath = p
         panapi.Log("Scavenger Path", p.Fingerprint)
         return p
      else
         local p = nextBestPath(raddr)
         conn.LastPath = p
         panapi.Log("Best Path", p.Fingerprint)
         return p
      end
   end
end

-- gets called whenever a path disappears(?)
function panapi.PathDown(laddr, raddr, fp, pi)
   panapi.Log("PathDown called with", laddr, raddr, fp, pi)
   paths[raddr][fp] = nil
end

function panapi.Refresh(laddr, raddr, ps)
   panapi.Log("Refresh", raddr, ps)
   panapi.Initialize(nil, laddr, raddr, ps)
end


function panapi.Close(laddr, raddr)
   panapi.Log("Close", laddr, raddr)
   conns[raddr][laddr] = nil
   --paths[raddr] = nil
end


function stats.UpdatedMetrics(laddr, raddr, rttStats, cwnd, bytesInFlight, packetsInFlight)
   --panapi.Log("UpdatedMetrics", cwnd, bytesInFlight, packetsInFlight)
   local lastpath = conns[raddr][laddr].LastPath
   if lastpath == nil then return end
   local path = paths[raddr][lastpath.Fingerprint]
   path.LatestTick = tick
   path.RTTStats = rttStats
--   panapi.Log("\n", tprint(rttStats, 1))
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

-- CURRENTLY UNUSED
function stats.TracerForConnection(id, p, odcid)
   --panapi.Log("id:", id, "perspective", p, "odcid", odcid)
end


function stats.StartedConnection(laddr, raddr, srcid, dstid)
   
end
function stats.NegotiatedVersion(laddr, raddr)

end
function stats.ClosedConnection(laddr, raddr)

end
function stats.SentTransportParameters(laddr, raddr)

end
function stats.ReceivedTransportParameters(laddr, raddr)

end
function stats.RestoredTransportParameters(laddr, raddr)

end
function stats.SentPacket(laddr, raddr)

end
function stats.ReceivedVersionNegotiationPacket(laddr, raddr)

end
function stats.ReceivedRetry(laddr, raddr)

end
function stats.ReceivedPacket(laddr, raddr)

end
function stats.BufferedPacket(laddr, raddr)

end
function stats.DroppedPacket(laddr, raddr)

end
function stats.AcknowledgedPacket(laddr, raddr)

end
function stats.LostPacket(laddr, raddr)

end
function stats.UpdatedCongestionState(laddr, raddr)

end
function stats.UpdatedPTOCount(laddr, raddr)

end
function stats.UpdatedKeyFromTLS(laddr, raddr)

end
function stats.UpdatedKey(laddr, raddr)

end
function stats.DroppedEncryptionLevel(laddr, raddr)

end
function stats.DroppedKey(laddr, raddr)

end
function stats.SetLossTimer(laddr, raddr)

end
function stats.LossTimerExpired(laddr, raddr)

end
function stats.LossTimerCanceled(laddr, raddr)

end
function stats.Debug(laddr, raddr)

end

