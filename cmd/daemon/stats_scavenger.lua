print("HELLO FROM THE SCAVENGER SCRIPT")

-- global variable to track "time"
tick = 0

function panapi.Periodic(seconds)
   --panapi.Log("Tick", seconds)
   tick = tick + 1
end

raddr2fps = {}

fp2path = {}

fp2tick = {}

fp2rtt = {}

laddr2prefs = {}
laddr2fp = {}

-- global table to keep track of connections
conns = {}

--fingerprints = {}


function nextScavengePath(raddr)
   local fingerprints = {}
   for _, fp in ipairs(raddr2fps[raddr]) do
      table.insert(fingerprints, fp)
   end
   table.sort(
      fingerprints,
      function(fp_a, fp_b)
         return fp2tick[fp_a] < fp2tick[fp_b]
      end
   )
   panapi.Log("Chosen tick:", fp2tick[fingerprints[1]])
   return fp2path[fingerprints[1]]
end


function nextBestPath(raddr)
   local fingerprints = {}
   for _, fp in ipairs(raddr2fps[raddr]) do
      table.insert(fingerprints, fp)
   end
    table.sort(
      fingerprints,
      function(fp_a, fp_b)
         return (fp2rtt[fp_a] or 1000) < (fp2rtt[fp_b] or 0)
      end
    )
    panapi.Log("Chosen tick:", fp2tick[fingerprints[1]])
   return fp2path[fingerprints[1]]
end



-- gets called when a set of paths to addr is known
function panapi.Initialize(prefs, laddr, raddr, ps)
   panapi.Log("New connection [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)
   raddr2fps[raddr] = raddr2fps[raddr] or {}
   for _, path in ipairs(ps) do
      local fp = path.Fingerprint
      fp2path[fp] = path
      fp2tick[fp] = tick
      table.insert(raddr2fps[raddr], fp)
   end
   panapi.SetPreferences(prefs, laddr, raddr)
end

function panapi.SetPreferences(prefs, laddr, raddr)
   panapi.Log("Update Preferences [" .. laddr, "|", raddr .. "] Profile:", prefs.ConnCapacityProfile)

   if prefs ~= nil then
      laddr2prefs[laddr] = laddr2prefs[laddr] or {}
      laddr2prefs[laddr] = prefs
   end
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   if #raddr2fps[raddr] > 0 then
      local p = nil
      if laddr2prefs[laddr]["ConnCapacityProfile"] == "Scavenger" then
         p = nextScavengePath(raddr)
         panapi.Log("Scavenger Path", p.Fingerprint)
      else
         p = nextBestPath(raddr)
         panapi.Log("Best Path", p.Fingerprint)
      end
      laddr2fp[laddr] = p.Fingerprint
      return p
   end
end

-- gets called whenever a path disappears(?)
function panapi.PathDown(laddr, raddr, fp, pi)
   panapi.Log("PathDown called with", laddr, raddr, fp, pi)
   fp2path[fp] = nil
   fp2rtt[fp] = nil
   fp2tick[fp] = nil
end

function panapi.Refresh(laddr, raddr, ps)
   panapi.Log("Refresh", raddr, ps)
   panapi.Initialize(nil, laddr, raddr, ps)
end


function panapi.Close(laddr, raddr)
   panapi.Log("Close", laddr, raddr)
   raddr2fps[raddr] = nil
   laddr2fp[laddr] = nil
   laddr2prefs[laddr] = nil
end


function stats.UpdatedMetrics(laddr, raddr, rttStats, cwnd, bytesInFlight, packetsInFlight)
   --panapi.Log("UpdatedMetrics", cwnd, bytesInFlight, packetsInFlight)
   local fp = laddr2fp[laddr]
   if fp == nil then return end
   fp2tick[fp] = tick
   fp2rtt[fp] = rttStats.LatestRTT
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

