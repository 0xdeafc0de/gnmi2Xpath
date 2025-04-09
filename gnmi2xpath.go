package main

import (
	"fmt"
	"strings"

	gnmi "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

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

// Prints gNMI Path
func printGNMIPath(path *gnmi.Path) string {
	return prototext.Format(path)
}

// Compares 2 gNMI paths using proto.Equal
func compareGnmi(p1 *gnmi.Path, p2 *gnmi.Path) bool {
	return proto.Equal(p1, p2)
}

func main() {
	xpath := "/switch/clusters/cluster[name=cluster-1]/ports/port[id=eth0]/interfaces/interface[name=eth0.100]/state/counters"

	fmt.Println("Input XPATH:", xpath)

	gnmiPath1, err := ygot.StringToPath(xpath, ygot.StructuredPath)
	if err != nil {
		fmt.Printf("Error in converting xpath to gNMI using ygot. err = %v\n", err)
		return
	}

	// Our own implementation
	gnmiPath2 := XpathToGNMIPath(xpath)

	// Compare to see if both are same
	if !compareGnmiPaths(gnmiPath1, gnmiPath2) {
		fmt.Printf("Error in gNMI conversion. Ours - %v, ygot - %v\n", gnmiPath2, gnmiPath1)
		return
	}
	fmt.Println("Converted to gNMI Path:", printGNMIPath(gnmiPath2))

	// If we are here, we have successfully converted XPATH to gNMI path and compared it matching with ygot generated path
	// Now we will convert the gNMI path back to XPATH.

	xpathBack1, err := ygot.PathToString(gnmiPath1)
	if err != nil {
		fmt.Printf("Error in getting XPATH using ygot. error = %v\n", err)
		return
	}

	// Our implementation
	xpathBack2 := GNMIPathToXpath(gnmiPath2)
	if xpathBack1 != xpathBack2 {
		fmt.Printf("ERROR in XPATH conversion. Ours - %s, Ygot - %s\n", xpathBack2, xpathBack1)
	}
	fmt.Println("Converted back to XPATH: ", xpathBack1)
}
