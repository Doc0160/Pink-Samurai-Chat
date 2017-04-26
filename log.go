/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "log"
    "runtime"
    "path"
    "strconv"
)

func Log(typ string, info interface{}){
    _, file, line, _ := runtime.Caller(1)
    log.Println("["+typ+"]",
        path.Base(file)+":"+strconv.Itoa(line),
        info)
}
