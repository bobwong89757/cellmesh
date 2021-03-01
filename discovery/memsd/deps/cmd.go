package deps

import (
	"fmt"
	"os"
	"sort"
)

func ViewSvc(flagAddr *string) {

	sd := InitSD(flagAddr)

	list := sd.QueryAll()

	sort.Slice(list, func(i, j int) bool {

		a := list[i]
		b := list[j]

		if a.GetMeta("SvcGroup") != b.GetMeta("SvcGroup") {
			return a.GetMeta("SvcGroup") < b.GetMeta("SvcGroup")
		}

		if a.Port != b.Port {
			return a.Port < b.Port
		}

		if a.Host != b.Host {
			return a.Host < b.Host
		}

		return a.ID < b.ID
	})

	for _, desc := range list {

		fmt.Println(desc.FormatString())
	}
}

func ViewKey(flagAddr *string) {
	sd := InitSD(flagAddr)
	list := sd.GetRawValueList("")
	sort.Slice(list, func(i, j int) bool {

		a := list[i]
		b := list[j]

		return a.Key < b.Key
	})

	for _, meta := range list {
		fmt.Printf("  %s = (size %d)\n", meta.Key, len(meta.Value))
	}
}

func GetValue(flagAddr *string,key string) {
	sd := InitSD(flagAddr)
	var value string
	err := sd.GetValue(key, &value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(value)
}

func SetValue(flagAddr *string,key, value string) {
	sd := InitSD(flagAddr)
	err := sd.SetValue(key, value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ClearSvc(flagAddr *string) {

	sd := InitSD(flagAddr)
	sd.ClearService()
}

func ClearValue(flagAddr *string) {

	sd := InitSD(flagAddr)
	sd.ClearKey()
}

func DeleteValue(flagAddr *string,key string) {

	sd := InitSD(flagAddr)
	sd.DeleteValue(key)
}
