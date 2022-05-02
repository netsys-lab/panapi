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

-- global "path database variable",
-- not directly referenced from go
paths = {}
print("Hello SelectionServer with stats")

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
   panapi.Log("Initialize", tprint(prefs), laddr, raddr)
   paths[raddr] = ps
   rankpaths(raddr)
end

function panapi.SetPreferences(prefs, laddr, raddr)
   panapi.Log("SetPreferences", tprint(prefs))
end

-- gets called for every packet
-- implementation needs to be efficient
function panapi.Path(laddr, raddr)
   panapi.Log("Path", laddr, raddr)
   if #paths[raddr] > 0 then
      return paths[raddr][1]
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
   panapi.Log("debug! exiting for introspection", #paths[raddr])
   os.exit(0)

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
