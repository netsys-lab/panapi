#!/bin/bash
while :; do 
	netcat 127.0.0.1 5555
	sleep 1;
	echo "Closed."
done
