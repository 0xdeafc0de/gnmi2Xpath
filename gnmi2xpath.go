package main

import (
	"fmt"
	"strings"

	gnmi "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/protobuf/encoding/prototext"

)

func printGNMIPath(path *gnmi.Path) string {
	return prototext.Format(path)
}

// XpathToGNMIPath converts an OpenConfig-style XPATH string to a gNMI Path.
func XpathToGNMIPath(xpath string) *gnmi.Path {
	if xpath == "" {
		return &gnmi.Path{}
	}
	segments := strings.Split(xpath, "/")
	elems := []*gnmi.PathElem{}
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		name := seg
		keys := map[string]string{}
		if strings.Contains(seg, "[") {
			parts := strings.SplitN(seg, "[", 2)
			name = parts[0]
			keystr := strings.TrimSuffix(parts[1], "]")
			for _, kv := range strings.Split(keystr, "][") {
				kvparts := strings.SplitN(kv, "=", 2)
				if len(kvparts) == 2 {
					keys[kvparts[0]] = kvparts[1]
				}
			}
		}
		elems = append(elems, &gnmi.PathElem{Name: name, Key: keys})
	}
	return &gnmi.Path{Elem: elems}
}

// GNMIPathToXpath converts a gNMI Path back to XPATH-style OC path.
func GNMIPathToXpath(path *gnmi.Path) string {
	parts := []string{}
	for _, elem := range path.Elem {
		seg := elem.Name
		if len(elem.Key) > 0 {
			keyparts := []string{}
			for k, v := range elem.Key {
				keyparts = append(keyparts, fmt.Sprintf("%s=%s", k, v))
			}
			seg += "[" + strings.Join(keyparts, "][") + "]"
		}
		parts = append(parts, seg)
	}
	return "/" + strings.Join(parts, "/")
}

func main() {
	xpath := "/tunnel/serialnos/serialno[sn=A99Z99999454]/ssids/ssid[id=Nile-Secure]/radios/radio[id=2]/state/counters"
	fmt.Println("Input XPATH:", xpath)

	gnmiPath := XpathToGNMIPath(xpath)
	fmt.Println("Converted to gNMI Path:", printGNMIPath(gnmiPath))

	xpathBack := GNMIPathToXpath(gnmiPath)
	fmt.Println("Converted back to XPATH:", xpathBack)
}

