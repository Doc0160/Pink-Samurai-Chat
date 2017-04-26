/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import(
    "time"
)

type Stub struct {
    Type string `json:"type"`
}
const _Hello = "hello"
func NewHello() Stub {
    return Stub{Type: _Hello}
}

const _Message string = "message"
type Message struct {
    Type string     `json:"type"`
    Username string `json:"username"`
    Channel string  `json:"channel"`
    Text string     `json:"text"`
    Time int64      `json:"time"`
}
func NewMessage(username, channel, text string) Message {
    return Message{
        Type: _Message,
        Username: username,
        Channel: channel,
        Text: text,
        Time: time.Now().UnixNano(),
    }
}

const _ChannelJoin string = "channel_join"
type ChannelJoin struct {
    Type string `json:"type"`
    Username string `json:"username"`
    Channel string `json:"channel"`
    Time int64      `json:"time"`
}
func NewChannelJoin(u, c string) ChannelJoin {
    return ChannelJoin{
        Type: _ChannelJoin,
        Username: u,
        Channel: c,
        Time: time.Now().UnixNano(),
    }
}

const _ChannelLeave string = "channel_leave"
type ChannelLeave struct {
    Type string      `json:"type"`
    Username string `json:"username"`
    Channel string  `json:"channel"`
    Time int64      `json:"time"`
}
func NewChannelLeave(username, channel string) ChannelLeave {
    return ChannelLeave{
        Type: _ChannelLeave,
        Username: username,
        Channel: channel,
        Time: time.Now().UnixNano(),
    }
}

const _Disconnect string = "disconnect"
type Disconnect struct {
    Type string     `json:"type"`
    Username string `json:"username"`
    Time int64      `json:"time"`
}
func NewDisconnect(username string) Disconnect {
    return Disconnect{
        Type: _Disconnect,
        Username: username,
        Time: time.Now().UnixNano(),
    }
}

const _Command string = "command"
type Command struct {
    Type string     `json:"type"`
    Username string `json:"username"`
    Channel string  `json:"channel"`
    Command string  `json:"command"`
    Time    int64   `json:"time"`
}
