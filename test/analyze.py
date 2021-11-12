import pandas as pd

from sys import argv

f = open(argv[1])
# prevent pandas to hickup after reading first line
metadata = f.readline()
print(pd.read_json(metadata))

fd = pd.read_json(f, lines=True)

d = dict()

for name in fd.name.unique():
    tmp = fd[fd.name == str(name)][~fd.time.duplicated()].pivot(index="time", columns="name", values="data")[name]
    d[name] = pd.DataFrame([pd.Series(tmp[x]) for x in tmp.index], index=tmp.index)

for k, v in d.items():
    print("#"*80)
    print(k)
    print("-"*80)
    print(v)
    print("#"*80)
