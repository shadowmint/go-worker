package resources 

import "bufio"
import "os"
import "fmt"
import "strings"
import "bytes"

/*================================================================================*
 * {{{ Public api
 *================================================================================*/

func newFileConfig() Config {
  var rtn = fileConfig {
    configBase : configBase {
      values : make(map[string] string),
      requires : make(map[string] string),
      defaults : make(map[string] string),
    },
  }
  return &rtn
}

/*================================================================================*
 * }}}
 *================================================================================*/

/*================================================================================*
 * {{{ Config implementation
 *================================================================================*/

type fileConfig struct {
  configBase
}

// Read a string and parse the lines into a config
//
// We assume every line is in the format:
// blag df sadf dsf adsf = adsfdsaf adsf dasf adsf adf afd
//
// Lines starting with '#' are ignored.
func (self *fileConfig) Parse(raw string) error {
  var b = bytes.NewBufferString(raw)
  var r = bufio.NewReader(b)
  self.loadData(r)
  return nil
}

// Load a file and parse the lines into a config
//
// We assume every line is in the format:
// blag df sadf dsf adsf = adsfdsaf adsf dasf adsf adf afd
//
// Lines starting with '#' are ignored.
func (self *fileConfig) Load(path string) error {
  var rtn = true
  f, err := os.Open(path)
  if err != nil {
      fmt.Println("error opening file= ",err)
      rtn = false
  }
  if rtn {
    var r = bufio.NewReader(f)
    self.loadData(r)
  }
  return nil
}

/*================================================================================*
 * }}}
 *================================================================================*/

/*================================================================================*
 * {{{ Internal
 *================================================================================*/

// Load config from stream
func (self *fileConfig) loadData(reader *bufio.Reader) {
  self.values = make(map[string] string)
  s, e := self.read(reader)
  for e == nil {
    var k, v = self.parseLine(s)
    if k != "" {
      self.values[k] = v
    }
    s, e = self.read(reader)
  }
}

// Parse a line
func (self *fileConfig) parseLine(s string) (string, string) {
  s = strings.TrimSpace(s)
  var k, v string
  if !strings.HasPrefix(s, "#") {
    if strings.Contains(s, "=") {
      var set = strings.Split(s, "=")
      k = strings.TrimSpace(set[0])
      v = strings.TrimSpace(set[1])
    }
  }
  return k, v
}

// Read a single line from the buffer
func (self *fileConfig) read(r *bufio.Reader) (string, error) {
  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}

/*================================================================================*
 * }}}
 *================================================================================*/
