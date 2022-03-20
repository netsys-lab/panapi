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

-- global "scores database variable", used to reference scores fetched from a path oracle
-- not directly referenced from go
scores = {}

print("Loading SelectionServer - selecting the path with highest oracle score for a connection")

function rankpaths(raddr)
   table.sort(
           scores[GetAsOfAddr(raddr)],
           function(path_a, path_b)
              return path_a.Scores.bandwidth_v2 < path_b.Scores.bandwidth_v2
           end
   )
end

-- gets called when a set of paths to addr is known
function panapi.Initialize(laddr, raddr, ps, sc)
   panapi.Log("Initialize - abc", laddr, raddr)
   paths[raddr] = ps
   scores = sc
   rankpaths(raddr)
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   local ras = GetAsOfAddr(raddr)
   if #scores[ras] > 0 then
      panapi.Log("using path chosen by using dynmaic path metadata: ", scores[ras][1].Fingerprint)
      return scores[ras][1].PathRef
   end

   if #paths[raddr] > 0 then
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
