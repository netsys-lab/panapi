#!/usr/bin/env bash

src_topo_file="$1"
gen_dir=${2:-gen}
net_interface=${3:-lo}

function as_topo_file() {
  local as=$1
  local dir_postfix=$(echo "$as" | cut -d '-' -f2 | tr ":" _)
  echo "$gen_dir/AS${dir_postfix}/topology.json"
}

function tc_links_from_topo_file() {
  yq ".links[] | select(.tcBw != null)" $1
}

function underlay_public() {
      local fromAS=$1
      local toAS=$2
      local interface=$3
      local as_topofile=$(as_topo_file "$fromAS")
      jq -r --arg iface $interface --arg toAS "$toAS" '.border_routers[].interfaces[$iface] | select(.isd_as==$toAS) | .underlay.public ' "$as_topofile"
}

function underlay_remote() {
    local fromAS=$1
    local toAS=$2
    local interface=$3
    local as_topofile=$(as_topo_file "$fromAS")
    jq -r --arg iface "$interface" --arg toAS "$toAS" '.border_routers[].interfaces[$iface] | select(.isd_as==$toAS) | .underlay.remote ' "$as_topofile"
}

echo "Resetting tc config for $net_interface"
sudo tcdel $net_interface --all

tc_links_from_topo_file $src_topo_file | while read link ; do
  from=$(echo "$link" | yq '.a' | cut -d "#" -f 1)
  fromif=$(echo "$link" | yq '.a' | cut -d "#" -f 2)
  to=$(echo "$link" | yq '.b' | cut -d "#" -f 1)
  toif=$(echo "$link" | yq '.b' | cut -d "#" -f 2)
  bw=$(echo "$link" | yq '.tcBw')

  remote=$(underlay_remote $from $to $toif)
  remote_ip=$(echo $remote | cut -d ':' -f1)
  remote_port=$(echo $remote | cut -d ':' -f2)
  public=$(underlay_public $from $to $fromif)
  public_ip=$(echo $public | cut -d ':' -f1)
  public_port=$(echo $public | cut -d ':' -f2)

  echo "link: $link"
  echo "from: $from"
  echo "to: $to"
  echo "toif: $toif"
  echo "remote: $remote"
  echo "remote ip: $remote_ip"
  echo


  echo "Throttling $from ($public) <-> $to ($remote) to a maximum of $bw."
  # tcset lo --rate 10Mbps --src-port 5000 --port 5000 --src-network 127.0.0.5 --network 127.0.0.6 --add --direction outgoing
  # one tc rule for each direction
  sudo tcset $net_interface --rate $bw \
    --src-port $public_port --port $remote_port \
    --src-network $public_ip --network $remote_ip \
    --add --direction outgoing
  sudo tcset $net_interface --rate $bw \
      --src-port $remote_port --port $public_port \
      --src-network $remote_ip --network $public_ip \
      --add --direction outgoing

  #sudo tcset $net_interface --rate $bw --port $public_port --network $public_ip --add --direction outgoing
  #sudo tcset $net_interface --rate $bw --port $remote_port --network $remote_ip --add --direction outgoing
done

echo "Current tc config for $net_interface"
tcshow $net_interface
