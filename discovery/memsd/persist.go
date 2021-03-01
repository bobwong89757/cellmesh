package main

import (
	"github.com/bobwong89757/cellmesh/discovery/memsd/model"
	"github.com/bobwong89757/cellnet/log"
	"os"
	"time"
)

func LoadPersistFile(fileName string) {

	fileHandle, err := os.OpenFile(fileName, os.O_RDONLY, 0666)

	// 可能文件不存在，忽略
	if err != nil {
		return
	}

	log.GetLog().Info("Load values...")

	err = model.LoadValue(fileHandle)
	if err != nil {
		log.GetLog().Error("load values failed: %s %s", fileName, err.Error())
		return
	}

	log.GetLog().Info("Load %d values", model.ValueCount())
}

func StartPersistCheck(fileName string) {

	ticker := time.NewTicker(time.Minute)

	for {

		<-ticker.C

		// 与收发在一个队列中，保证无锁
		model.Queue.Post(func() {

			if model.ValueDirty {

				log.GetLog().Info("Save values...")

				fileHandle, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					log.GetLog().Error("save persist file failed: %s %s", fileName, err.Error())
					return
				}

				err = model.SaveValue(fileHandle)

				if err != nil {
					log.GetLog().Error("save values failed: %s %s", fileName, err.Error())
					return
				}

				log.GetLog().Info("Save %d values", model.ValueCount())

				model.ValueDirty = false

			}

		})

	}

}
