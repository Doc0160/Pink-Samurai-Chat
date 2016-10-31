/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "net"
    "encoding/gob"
    "bytes"
)

type BroadCastMessage struct {
    Nick string
    Channel string
    Text string
    IP net.IP
}

func (bcmsg *BroadCastMessage)ToBytes()[]byte{
    var b bytes.Buffer
    gob.NewEncoder(&b).Encode(bcmsg)
    return b.Bytes()
}

func (bcmsg *BroadCastMessage)FromBytes(buf []byte){
    gob.NewDecoder(bytes.NewReader(buf)).Decode(&bcmsg)
}
