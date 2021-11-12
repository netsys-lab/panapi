import pandas as pd
import matplotlib.pyplot as plt

from time import sleep
from sys import argv

f = open(argv[1])
# prevent pandas to hickup after reading first line
metadata = f.readline()
print(pd.read_json(metadata))

fd = pd.read_json(f, lines=True)

d = dict()

for name in fd.name.unique():
    print("crunching %s" % name)
    tmp = fd[fd.name == str(name)][~fd.time.duplicated()].pivot(index="time", columns="name", values="data")[name]
    if name in ["transport:packet_received", "transport:packet_sent"]:
        for field in ["header", "raw", "frames"]:
            if field != "frames":
                d[name+":"+field] = pd.DataFrame([pd.Series(tmp[x].get(field, {})) for x in tmp.index], index=tmp.index)
            else:
                d[name+":"+field+":partial"] = pd.DataFrame([pd.Series(tmp[x].get(field, [None])[0]) for x in tmp.index], index=tmp.index)
            
    elif name == "recovery:packet_lost":
        for field in ["header", "trigger"]:
            d[name+":"+field] = pd.DataFrame([pd.Series(tmp[x].get(field, {})) for x in tmp.index], index=tmp.index)
    else:
        d[name] = pd.DataFrame([pd.Series(tmp[x]) for x in tmp.index], index=tmp.index)



for k, v in d.items():
    print("#"*80)
    print(k)
    print("-"*80)
    print(v)
    print("#"*80)


def filterDataFrame(frame, oldfields, newfields=None):
    if not newfields:
        newfields = oldfields
    return pd.DataFrame(dict([(newfield,frame[oldfield]) for newfield,oldfield in zip(newfields,oldfields)]), index=frame.index)

def fromTo(frame, field):
    last = 0
    f = 0
    t = 0
    l = []
    for now, p in zip(frame.index, frame[field]):
        if p > last+1:
            if f != t:
                l.append((f,t))
            f = now
        t = now
        last = p
    return l    
        

ts1 = d["recovery:metrics_updated"]
ts2 = d["transport:packet_received:frames:partial"]
ts3 = d["transport:packet_sent:frames:partial"]
ts4 = d["recovery:loss_timer_updated"]

#ts = pd.concat([filterDataFrame(ts1,["bytes_in_flight", "min_rtt", "smoothed_rtt", "latest_rtt"]),
ts = pd.concat([filterDataFrame(ts1,["smoothed_rtt", "latest_rtt"]),
                #filterDataFrame(ts2,["ack_delay"],["received_ack_delay"]),
                #filterDataFrame(ts3,["ack_delay"],["sent_ack_delay"]),
                #filterDataFrame(ts4,["delta"]),
                ])
               


#ts = pd.concat([ts1, ts2, ts3], ignore_index=True, sort=False, keys=["a", "b", "c"])

ts.plot(subplots=False,
        #y=["bytes_in_flight", "min_rtt", "smoothed_rtt", "latest_rtt", "ack_delay"],
        xticks=range(0, int(ts.index[-1]), int(ts.index[-1])//1000*100),
        ylim=(100,300),
        linestyle="solid",
        marker=".",
        )

loss = d.get("recovery:packet_lost:header")
if loss != None:
    for f, to in fromTo(loss, "packet_number"):
        plt.axvspan(f,to, facecolor="red", alpha=.2)

plt.show()
