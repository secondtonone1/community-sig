package config

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/zookeeper/v2"
)

//文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func LoadConfigAndSetDefault() error {
	var configfile string
	var rootdir, vardir string

	flag.StringVar(&configfile, "conf", "", "conf file path")
	flag.StringVar(&rootdir, "root_dir", "", "root dir")
	flag.StringVar(&vardir, "var_dir", "", "var dir")

	flag.Parse()

	//community
	if configfile == "" {
		configfile = "/etc/communitysig.toml"
	}

	//分环境
	env := os.Getenv("profiles.active")
	if env != "" {
		dir, fileName := filepath.Split(configfile)
		fileNameArr := strings.Split(fileName, ".")
		if len(fileNameArr) == 2 {
			tem := dir + fileNameArr[0] + "_" + env + "." + fileNameArr[1]
			if FileExist(tem) {
				configfile = tem
			}
		}
	}
	var conf Config
	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		return errors.New("load toml conf file fail:" + err.Error())
	}

	if len(conf.RegisterCenter.Address) > 0 {
		r := zookeeper.NewRegistry(func(op *registry.Options) {
			op.Addrs = conf.RegisterCenter.Address
			op.Context = context.Background()
			if conf.RegisterCenter.Timeout > 0 {
				op.Timeout = time.Second * time.Duration(conf.RegisterCenter.Timeout)
			} else {
				op.Timeout = time.Second * 5
			}
		})
		conf.RegisterCenter.register = r
	}

	if conf.Base.GRPCAddr == "" {
		conf.Base.GRPCAddr = "0.0.0.0:8080"
	}

	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr != "" {
		conf.Base.GRPCAddr = grpcAddr
	}

	if conf.Base.ServiceName == "" {
		conf.Base.ServiceName = conf.LogConf.Project
	}

	if conf.Base.PeerServiceName == "" {
		conf.Base.PeerServiceName = "status-service"
	}

	if conf.Base.WebAddr == "" {
		conf.Base.WebAddr = "0.0.0.0:8080"
	}

	webAddr := os.Getenv("WEB_ADDR")
	if webAddr != "" {
		conf.Base.WebAddr = webAddr
	}

	if rootdir != "" {
		conf.Base.RootDir = rootdir
	}
	if vardir != "" {
		conf.Base.VarDir = vardir
	}
	if conf.Base.VarDir == "" {
		conf.Base.VarDir = conf.Base.RootDir
	}

	if conf.Base.HeartMax == 0 {
		conf.Base.HeartMax = 1800
	}

	conf.LogConf.LogDir = conf.Base.VarDir + "/" + conf.LogConf.LogDir

	fmt.Printf("conf is :%+v\n", conf)

	config = &conf

	return nil
}
