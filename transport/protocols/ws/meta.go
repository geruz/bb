package ws

import (
	"fmt"
	"strings"
	"strconv"
)

type Meta struct {
	Major int
	Resource string
	Action string
}
func (this Meta) ToBytes() []byte {
	return []byte(fmt.Sprintf("%v/%v/%v\n", this.Major, this.Resource, this.Action))
}
func ParseMeta(data []byte) (Meta, error) {
	s := string(data)
	parts := strings.Split(s, "/");
	if len(parts) != 3 {
		return Meta{}, fmt.Errorf("Bad parts count")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil{
		return Meta{}, err
	}
	return Meta{
		Major: major,
		Resource: parts[1],
		Action: parts[2],
	}, nil
}