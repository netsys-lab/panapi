# Traffic Control SCION Topology

The script [tctopo.sh](tctopo.sh) allows to setup traffic control (currently
bandwidth only) for SCION topologies generated with [scionproto/scion topology](https://github.com/scionproto/scion/tree/master/topology).

## Usage

1. Setup a SCION development environment

    see: [SCION Documentation - Setting up the development environment](https://docs.scion.org/en/latest/build/setup.html)

2. Install Requirements

    Following packages/applications are required before running tctopo.sh:

    - tcconfig, see: [tcconfig installation](https://tcconfig.readthedocs.io/en/latest/pages/introduction/index.html#installation)
    - jq, see: [jq installation](https://wiki.ubuntuusers.de/jq/#Installation)
    - yq, see: [yq installation](https://github.com/mikefarah/yq#install)

3. Adapt your topology file

    Adapt the links in your topology files with the following keys:

    | Key  | Description                     |Sample Value |
    |------|---------------------------------|-------------|
    | tcBw | bidirectional maximum Bandwidth |10 Mbps      |

    E.g. the topology file below, will create a SCION topology with one parent/core
    AS and two child ASes. The links interconnecting the ASes will be bandwidth
    limited to a throughput of 35 Mbps and respectively 20 Mbps.

    ```yaml
    ---
    ASes:
    "1-ff00:0:110":
        core: true
        voting: true
        authoritative: true
        issuing: true
    "1-ff00:0:111":
        cert_issuer: 1-ff00:0:110
    "1-ff00:0:112":
        cert_issuer: 1-ff00:0:110
    links:
    - {a: "1-ff00:0:110#1", b: "1-ff00:0:111#1", linkAtoB: CHILD, tcBw: 35Mbps}
    - {a: "1-ff00:0:110#2", b: "1-ff00:0:112#1", linkAtoB: CHILD, tcBw: 20Mbps}
    ```

4. Run SCION with your topology

    see: [SCION Documentation - Running SCION locally](https://docs.scion.org/en/latest/build/setup.html#running-scion-locally)

    Pass your customized topology file with:

    ```bash
    user@host:~/scion$ ./scion.sh topology -c $PATH_OF_YOUR_TOPO_FILE
    ```

5. Run tctopo.sh

    To apply the custom maximum bandwidth to your local SCION topology run:

    ```bash
    user@host:~/inter-domain-testbed$ ./path/to/tctopo.sh $PATH_OF_YOUR_TOPO_FILE $GEN_DIR $NET_INTERFACE
    ```

    Replace the variables:
      - _$PATH_OF_YOUR_TOPO_FILE_ with the path of the SCION topology file used
        in step four
      - _$GEN_DIR_ with the directory of the generated SCION configuration files
        (usually the subdirectory _gen_ of your local SCION repository)
      - _$NET_INTERFACE_ with the network interface to apply the tcconfig to (defaults to lo)

    Attention: Running this script will reset all tc rules applied to
    _$NET_INTERFACE_ before generating the new set of rules.
