/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "net/http"
    _ "net/http/pprof"
)

type MemInfoDumper struct {
}

func NewMemInfoDumper() *MemInfoDumper {
    m := MemInfoDumper{}

	go http.ListenAndServe("localhost:8080", nil)
    
    return &m
}
