package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fakeCurl()
}

func fakeCurl() {

	// http://placehold.it:87/120x120&text=image3
	var url string = "127.0.0.1:8080"
	conn, err := net.Dial("tcp", url)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	//works with WriteString too
	data, err := rw.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
	fcheck(err)

	err2 := rw.Flush()
	fcheck(err2)

	response, err := rw.ReadString('\n')
	fcheck(err)
	fmt.Println(response)

	fcheck(err)
	fmt.Println(data)

}

func fcheck(err error) {
	if err != nil {
		panic(err)
	}
}
