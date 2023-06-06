package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/sys/windows/registry"
)

var paths = []string{
	"Software\\miHoYo\\崩坏：星穹铁道",
	"Software\\Cognosphere\\Star Rail",
}

func init() {
	fmt.Print("Honkai: Star Rail FPS Unlocker\nSource code: https://github.com/gizmo-ds/star-rail-fps-unlocker\n\n")
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
		exit(0, "can't not find game registry key")
	}
	defer starRail.Close()

	names, err := starRail.ReadValueNames(0)
	if err != nil {
		exit(1, err.Error())
	}

	modelValue := lo.Filter(names, func(v string, i int) bool { return strings.Index(v, "GraphicsSettings_Model_h") == 0 })
	if len(modelValue) < 1 {
		exit(1, "no model value found")
	}

	value, _, err := starRail.GetBinaryValue(modelValue[0])
	if err != nil {
		exit(1, err.Error())
	}
	if value[len(value)-1] == 0x00 {
		value = value[:len(value)-1]
	}
	var m map[string]any
	if err = json.Unmarshal(value, &m); err != nil {
		exit(1, err.Error())
	}
	m["FPS"] = 120
	data, err := json.Marshal(m)
	if err != nil {
		exit(1, err.Error())
	}
	data = append(data, 0x00)
	if starRail.SetBinaryValue(modelValue[0], data) != nil {
		exit(1, err.Error())
	}
	exit(0, "Done")
}

func exit(code int, args ...any) {
	fmt.Println(args...)
	fmt.Println("Press enter to exit...")
	_, _ = fmt.Scanln()
	os.Exit(code)
}
