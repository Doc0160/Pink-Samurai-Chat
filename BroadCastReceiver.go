/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "net"
    "fmt"
)

type BroadCastReceiver struct {
    conn *net.UDPConn
    msgs chan *BroadCastMessage
}

func NewBroadCastReceiver() *BroadCastReceiver{
    bcrecv := BroadCastReceiver{}
    bcrecv.msgs = make(chan *BroadCastMessage, 64)
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{
        IP:   net.IPv4(0, 0, 0, 0),
        Port: *Port,
    })
    CheckError(err)
    bcrecv.conn = conn

    go func(){
	    data := make([]byte, (1<<16)-1)  //maximal udp packet size
        var msg BroadCastMessage
		for {
		    read, addr , err := bcrecv.conn.ReadFromUDP(data[0:])
		    if err!=nil {
		    	fmt.Print("BCastReceiver:: ",err,addr)
		    	continue
		    }
		    msg.FromBytes(data[:read])
		    msg.IP = addr.IP
		    bcrecv.msgs <- &msg
		}
    }()
    return &bcrecv
}

func (r *BroadCastReceiver) GetMessage() *BroadCastMessage {
	return <-r.msgs
}
