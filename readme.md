# TAPS-API Project
This is the repository for the [TAPS-API Draft](https://www.ietf.org/archive/id/draft-ietf-taps-interface-10.html) example implementation. 

# Setup
GO 1.14+ is required.
To fix receive buffer size error run:
```
# sysctl -w net.core.rmem_max=2500000
```