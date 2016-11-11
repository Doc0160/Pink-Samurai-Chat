/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (

)
const (
    _Stub string = ""
    _Message string = "message"
    _UsernameChange string = "username_change"
)


type Stub struct {
    Type string `json:"type"`
}

type Message struct {
    Type string `json:"type"`
	Username string `json:"username"`
    Channel string `json:"channel"`
	Text string `json:"text"`
}

type UsernameChange struct {
    Type string `json:"type"`
	Username string `json:"username"`
    OldUsername string `json:"old_username"`
}
