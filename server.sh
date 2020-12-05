#!/bin/bash
while :; do 
	netcat -l 5555; 
	sleep 1;
	echo "Closed."
done
