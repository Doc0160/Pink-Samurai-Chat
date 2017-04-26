/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "net/http"
    "log"
    "strings"
    "compress/flate"
    "compress/gzip"
)

func serve404(w http.ResponseWriter, r *http.Request){
    log.Println("Serving 404")
    res, err := Asset("404.html")
    if err != nil {
        Log("404 Borken", err) 
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Not found\nCongrats, you've borken it all.\n"+err.Error()))
    } else {
        Log("404", r)
        w.Header().Set("Content-Type", "text/html; charset=utf-8") 
        w.WriteHeader(http.StatusNotFound)
        w.Write(res)
    }
}

func Write(w http.ResponseWriter, r *http.Request, b []byte) {

    // TODO(doc): consider using gzip before deflate
    // even if deflate quicker than gzip to compress
    // Another problem found while deploying HTTP compression
    // on large scale is due to the deflate encoding definition:
    // while HTTP 1.1 defines the deflate encoding as data compressed
    // with deflate (RFC 1951) inside a zlib formatted stream (RFC 1950),
    // Microsoft server and client products historically implemented it as
    // a "raw" deflated stream, making its deployment unreliable.
    // https://en.wikipedia.org/wiki/HTTP_compression#Problems_preventing_the_use_of_HTTP_compression

    /*
    // NOTE(doc): if less than 1KO, don't bother compressing it ?
    if len(b) < 1000 { 
       w.Write(b)
       return
    } else 
    //*/

    if strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
        w.Header().Set("Content-Encoding", "deflate")
        g, _ := flate.NewWriter(w, -1)
        g.Write(b)
        g.Close()
        
    } else if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Header().Set("Content-Encoding", "gzip")
        g := gzip.NewWriter(w)
        g.Write(b)
        g.Close()

        //TODO(doc): add br / brotli (firefox), compress (unix), SDCH (chrome)
        // exi / Efficient XML Interchange, 

    } else {
        w.Write(b)
    }
}
