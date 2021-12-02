# Lua Selector

* Follows pan.Selector interface

## Lua API

Lua scripts needs to implement the following functions:

```
panapi.initialize(laddr, raddr, paths)
panapi.selectpath(raddr)
panapi.pathdown(raddr, fingerprint, pathinterface)
panapi.refresh(raddr, paths)
panapi.close(raddr)
panapi.periodic(delta)
```

Lua scripts can call the following functions from the panapi module:
```
panapi.log(...)
```

