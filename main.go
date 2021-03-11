package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"omo-msa-group/config"
	"omo-msa-group/handler"
	"omo-msa-group/model"
	"omo-msa-group/publisher"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-group/proto/group"
)

func main() {
	config.Setup()
	model.Setup()
    defer model.Cancel()
	model.AutoMigrateDatabase()

	// New Service
	service := micro.NewService(
		micro.Name(config.Schema.Service.Name),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)

	// Initialise service
	service.Init()

	// Register publisher
	publisher.DefaultPublisher = micro.NewPublisher(config.Schema.Service.Name + ".notification", service.Client())

	// Register Handler
	proto.RegisterHealthyHandler(service.Server(), new(handler.Healthy))
	proto.RegisterCollectionHandler(service.Server(), new(handler.Collection))
	proto.RegisterMemberHandler(service.Server(), new(handler.Member))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Error(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
