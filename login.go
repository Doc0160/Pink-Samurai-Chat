/* ========================================================================
   $File: $
   $Date: $
   $Revision: $
   $Creator: Tristan Magniez $
   ======================================================================== */

package main

import (
    "crypto/sha1"
    "encoding/base64"
    "encoding/json"
    "compress/gzip"
    "io/ioutil"
    "bytes"
    "log"
)

type HashThingU struct {
    Username string
    Channel string
}
type HashThing map[string]HashThingU
type Member map[string][20]byte

type Members struct{
    members Member
    hashs HashThing
}
var members = NewMembers()

func (m*Members) Add(u, p string) string{
    m.members[u] = sha1.Sum([]byte(p))
    hasher := sha1.New()
    hasher.Write([]byte(u+":"+p))
    sum := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
    m.hashs[sum] = HashThingU{
        Username: u,
    }
    return sum
}

func (m Members) Save() {
    b, err := json.Marshal(m)
	if err != nil {
		println(err.Error())
	}
	var bg bytes.Buffer
	w, _ := gzip.NewWriterLevel(&bg, gzip.BestCompression)
	w.Write(b)
	w.Close()
	err = ioutil.WriteFile("bdd.gjson", bg.Bytes(), 0777)
    if err != nil {
        log.Panic(err.Error())
    }
}

func NewMembers() Members{
    m := Members{}
    m.members = make(Member)
    m.hashs = make(HashThing)
    var gzr *gzip.Reader
    var bs []byte
    var err error
    bs, err = ioutil.ReadFile("bdd.gjson")
    if err != nil {
        m.Save()
        log.Panic(err.Error())
    }
    gzr, err = gzip.NewReader(bytes.NewReader(bs))
	if err != nil {
        m.Save()
		panic(err.Error())
	}
	bs, err = ioutil.ReadAll(gzr)
    gzr.Close()
	if err != nil {
        m.Save()
		panic(err.Error())
	}
	err = json.Unmarshal(bs, &m)
	if err != nil {
        m.Save()
		panic(err.Error())
	}
    m.Add("Doc0160", "badwolf")
   //     m["kaze"] = sha1.Sum([]byte("test"))
    return m
}
