package main

import (
    "fmt"
    "flag"
)

var Nickname *string = flag.String("nick","Anonymous","Your Nickname")
var Channel *string = flag.String("channel","foyer","Channel")
var IP *string = flag.String("ip","","The broadcast IP to use")
var Port *int = flag.Int("port",1234,"The port used to broadcast")

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
	}
}

func main(){
    flag.Parse()

    NewPinkSamuraiChat()
    
    select{}
}

type PinkSamuraiChat struct {
    bcrecv *BroadCastReceiver
    bcsender *BroadCastSender
}
func NewPinkSamuraiChat()*PinkSamuraiChat{
    psc := PinkSamuraiChat{}
    psc.bcrecv = NewBroadCastReceiver()
    psc.bcsender = NewBroadCastSender()
    go func(){
        for {
            m := psc.bcrecv.GetMessage()
            println(m.Nick+"@"+m.Channel+"["+m.IP.String()+"]: " + m.Text)
        }
    }()
    go func(){
        for {
            var input string
            fmt.Scanln(&input)
            m := BroadCastMessage{
                Text: input,
                Nick: *Nickname,
                Channel: *Channel,
            }
            psc.bcsender.Send(m)
        }
    }()
    return &psc
}

func (psc *PinkSamuraiChat)f(){

}
