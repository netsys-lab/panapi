-- Print contents of `tbl`, with indentation.
-- `indent` sets the initial level of indentation.
function tprint (tbl, indent)
   if not indent then indent = 0 end
   s = ""
   for k, v in pairs(tbl) do
      formatting = string.rep("  ", indent) .. k .. ": "
      if type(v) == "table" then
         --print(formatting)
         s = s ..  formatting .. "\n" .. tprint(v, indent+1)
      else
         s = s .. formatting .. tostring(v) .. "\n"
      end
   end
   return s
end

function GetAsOfAddr(raddr)
   local addrParts = {}
   for str in string.gmatch(raddr, "([^,]+)") do
      table.insert(addrParts, str)
   end
   return addrParts[1]
end

-- global "path database variable",
-- not directly referenced from go
paths = {}

math.randomseed(os.time())
print("Loading SelectionServer - selecting a path with the lowest amount of hops for a connection")

function rankpaths(raddr)
   panapi.Log("Ranking paths for ", raddr)

   -- sort by amount of hops
   table.sort(
           paths[raddr],
           function(a, b)
              return #paths[a].Metadata.Interfaces < #paths[b].Metadata.Interfaces
           end
   )

   if #paths[raddr] == 1 then
      return
   end

   -- shuffle all paths with the lowest amount of hops
   local index_last_path_with_min_hops
   for i = 1, #paths[raddr], 1 do
      if #paths[raddr][i].Interfaces > #paths[raddr][1].Interfaces then
         index_last_path_with_min_hops = i - 1
         break
      end
   end

   for i = index_last_path_with_min_hops, 2, -1 do
      local j = math.random(i)
      paths[raddr][i], paths[raddr][j] = paths[raddr][j], paths[raddr][i]
   end
end

-- gets called when a set of paths to addr is known
function panapi.Initialize(laddr, raddr, ps, sc)
   panapi.Log("Initialize - abc", laddr, raddr)
   paths[raddr] = ps
   rankpaths(raddr)
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   if #paths[raddr] > 0 then
      panapi.Log("Path for ", raddr, " ", paths[raddr][1].Fingerprint)
      return paths[raddr][1]
   end
end

-- gets called whenever a path disappears(?)
function panapi.PathDown(laddr, raddr, fp, pi)
   panapi.Log("PathDown called with", laddr, raddr, fp, pi)
end

function panapi.Refresh(laddr, raddr, ps)
   panapi.Log("Refresh", raddr, ps)
   paths[raddr] = ps
   rankpaths(raddr)
end


function panapi.Close(laddr, raddr)
   panapi.Log("Close", laddr, raddr)
   paths[raddr] = nil
end

function panapi.Periodic(seconds)
   panapi.Log("Periodic", seconds)
end



function stats.TracerForConnection(id, p, odcid)
   panapi.Log("id:", id, "perspective", p, "odcid", odcid)
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
function stats.UpdatedMetrics(laddr, raddr, rttStats, cwnd, bytesInFlight, packetsInFlight)
   panapi.Log("UpdatedMetrics", cwnd, bytesInFlight, packetsInFlight)
   panapi.Log("\n", tprint(rttStats, 1))
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
