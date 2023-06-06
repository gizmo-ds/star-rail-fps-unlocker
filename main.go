package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/sys/windows/registry"
)

var paths = []string{
	"Software\\miHoYo\\崩坏：星穹铁道",
	"Software\\Cognosphere\\Star Rail",
}

func main() {
	var starRail *registry.Key
	var err error
	for _, p := range paths {
		sr, err := registry.OpenKey(registry.CURRENT_USER, p, registry.READ|registry.WRITE)
		if err == nil {
			starRail = &sr
			break
		}
	}
	if starRail == nil {
		return
	}
	defer starRail.Close()

	names, err := starRail.ReadValueNames(0)
	if err != nil {
		panic(err)
	}

	modelValue := lo.Filter(names, func(v string, i int) bool { return strings.Index(v, "GraphicsSettings_Model_h") == 0 })
	if len(modelValue) < 1 {
		panic("not found")
	}

	value, _, err := starRail.GetBinaryValue(modelValue[0])
	if err != nil {
		panic(err)
	}
	if value[len(value)-1] == 0x00 {
		value = value[:len(value)-1]
	}
	var m map[string]any
	if err = json.Unmarshal(value, &m); err != nil {
		panic(err)
	}
	m["FPS"] = 120
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	data = append(data, 0x00)
	if starRail.SetBinaryValue(modelValue[0], data) != nil {
		panic(err)
	}
	fmt.Println("done")
}
