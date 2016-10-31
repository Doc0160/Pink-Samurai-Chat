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
    "strings"
)

type BroadCastSender struct {
    conn *net.UDPConn
    input chan *BroadCastMessage
}
func NewBroadCastSender()*BroadCastSender{
    bcsender := BroadCastSender{}
    bcsender.input = make(chan *BroadCastMessage)
    bcsender.Setup()
    return &bcsender
}

func (bcsender *BroadCastSender)Setup() error{
    /*addr, err := GetCastAddr()
    CheckError(err)
    ip := net.ParseIP(addr)*/
    conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
        IP: net.ParseIP("255.255.255.255"),
        Port: *Port,
    })
    if err != nil {
        return err
    }
    go func(){
        for msg := range bcsender.input {
            _,err := conn.Write(msg.ToBytes())
            if err != nil {
                fmt.Println(msg, err)
            }
        }
    }()
    bcsender.conn = conn
    return nil
}

func (bcsender *BroadCastSender)Send(bcmsg BroadCastMessage) {
    bcsender.input <- &bcmsg
}

func GetCastAddr() (string,error) {
	addrs,err := net.InterfaceAddrs()
	if err!=nil {
		return "",err
	}
	result := ""
	for _,addr := range addrs {
		str := addr.String()
		parts1 := strings.Split(str,"/")
		parts2 := strings.Split(parts1[0],".")
		if parts2[0]!="127" {
			switch(parts1[1]){
				case "24": {
					result = fmt.Sprintf("%v.%v.%v.255",parts2[0],parts2[1],parts2[2])
				}
				case "16": {
					result = fmt.Sprintf("%v.%v.255.255",parts2[0],parts2[1])
				}
				case "8": {
					result = fmt.Sprintf("%v.255.255.255",parts2[0])
				}
			}
		}
		if result!=""{
			break
		}
	}
	return result,nil
}
