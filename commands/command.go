package commands

import (
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
	"time"
	"log"
	"flag"
	"path/filepath"
	"strings"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/gocaptcha"
	"github.com/lifei6671/mindoc/commands/migrate"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils"
	"github.com/lifei6671/mindoc/cache"
	beegoCache "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
)

// RegisterDataBase 注册数据库
func RegisterDataBase() {
	beego.Info("正在初始化数据库配置.")
	adapter := beego.AppConfig.String("db_adapter")

	if adapter == "mysql" {
		host := beego.AppConfig.String("db_host")
		database := beego.AppConfig.String("db_database")
		username := beego.AppConfig.String("db_username")
		password := beego.AppConfig.String("db_password")
		timezone := beego.AppConfig.String("timezone")

		port := beego.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

		err := orm.RegisterDataBase("default", "mysql", dataSource)
		if err != nil {
			beego.Error("注册默认数据库失败:",err)
			os.Exit(1)
		}

		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			beego.Error("加载时区配置信息失败,请检查是否存在ZONEINFO环境变量:",err)
		}
	} else if adapter == "sqlite3" {
		database := beego.AppConfig.String("db_database")
		if strings.HasPrefix(database, "./") {
			database = filepath.Join(conf.WorkingDirectory, string(database[1:]))
		}

		dbPath := filepath.Dir(database)
		os.MkdirAll(dbPath, 0777)

		err := orm.RegisterDataBase("default", "sqlite3", database)

		if err != nil {
			beego.Error("注册默认数据库失败:",err)
		}
	}else{
		beego.Error("不支持的数据库类型.")
		os.Exit(1)
	}
	beego.Info("数据库初始化完成.")
}

// RegisterModel 注册Model
func RegisterModel() {
	orm.RegisterModelWithPrefix(conf.GetDatabasePrefix(),
		new(models.Member),
		new(models.Book),
		new(models.Relationship),
		new(models.Option),
		new(models.Document),
		new(models.Attachment),
		new(models.Logger),
		new(models.MemberToken),
		new(models.DocumentHistory),
		new(models.Migration),
		new(models.Label),
	)
	//migrate.RegisterMigration()
}

// RegisterLogger 注册日志
func RegisterLogger(log string) {

	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.Async()

	logPath := filepath.Join(log, "log.log")

	if _, err := os.Stat(logPath); os.IsNotExist(err) {

		os.MkdirAll(log, 0777)

		if f, err := os.Create(logPath); err == nil {
			f.Close()
			config := make(map[string]interface{}, 1)

			config["filename"] = logPath

			b, _ := json.Marshal(config)

			beego.SetLogger("file", string(b))
		}
	}

	beego.SetLogFuncCall(true)
	beego.BeeLogger.Async()
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if len(os.Args) >= 2 && os.Args[1] == "install" {
		ResolveCommand(os.Args[2:])
		Install()
	} else if len(os.Args) >= 2 && os.Args[1] == "version" {
		CheckUpdate()
		os.Exit(0)
	} else if len(os.Args) >= 2 && os.Args[1] == "migrate" {
		ResolveCommand(os.Args[2:])
		migrate.RunMigration()
	}
}
//注册模板函数
func RegisterFunction() {
	beego.AddFuncMap("config", models.GetOptionValue)

	beego.AddFuncMap("cdn", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdn", "")
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})

	beego.AddFuncMap("cdnjs", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdnjs", "")
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
	beego.AddFuncMap("cdncss", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdncss", "")
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
	beego.AddFuncMap("cdnimg", func(p string) string {
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		cdn := beego.AppConfig.DefaultString("cdnimg", "")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
}

//解析命令
func ResolveCommand(args []string) {
	flagSet := flag.NewFlagSet("MinDoc command: ", flag.ExitOnError)
	flagSet.StringVar(&conf.ConfigurationFile, "config", "", "MinDoc configuration file.")
	flagSet.StringVar(&conf.WorkingDirectory, "dir", "", "MinDoc working directory.")
	flagSet.StringVar(&conf.LogFile, "log", "", "MinDoc log file path.")

	flagSet.Parse(args)

	if conf.WorkingDirectory == "" {
		if p, err := filepath.Abs(os.Args[0]); err == nil {
			conf.WorkingDirectory = filepath.Dir(p)
		}
	}
	if conf.LogFile == "" {
		conf.LogFile = filepath.Join(conf.WorkingDirectory, "logs")
	}
	if conf.ConfigurationFile == "" {
		conf.ConfigurationFile = filepath.Join(conf.WorkingDirectory, "conf", "app.conf")
		config := filepath.Join(conf.WorkingDirectory, "conf", "app.conf.example")
		if !utils.FileExists(conf.ConfigurationFile) && utils.FileExists(config) {
			utils.CopyFile(conf.ConfigurationFile, config)
		}
	}
	gocaptcha.ReadFonts(filepath.Join(conf.WorkingDirectory, "static", "fonts"), ".ttf")

	err := beego.LoadAppConfig("ini", conf.ConfigurationFile)

	if err != nil {
		log.Println("An error occurred:", err)
		os.Exit(1)
	}
	uploads := filepath.Join(conf.WorkingDirectory, "uploads")

	os.MkdirAll(uploads, 0666)

	beego.BConfig.WebConfig.StaticDir["/static"] = filepath.Join(conf.WorkingDirectory, "static")
	beego.BConfig.WebConfig.StaticDir["/uploads"] = uploads
	beego.BConfig.WebConfig.ViewsPath = filepath.Join(conf.WorkingDirectory, "views")

	fonts := filepath.Join(conf.WorkingDirectory, "static", "fonts")

	if !utils.FileExists(fonts) {
		log.Fatal("Font path not exist.")
	}
	gocaptcha.ReadFonts(filepath.Join(conf.WorkingDirectory, "static", "fonts"), ".ttf")

	RegisterDataBase()
	RegisterCache()
	RegisterModel()
	RegisterLogger(conf.LogFile)
}

//注册缓存管道
func RegisterCache()  {
	isOpenCache := beego.AppConfig.DefaultBool("cache",false)
	if !isOpenCache {
		cache.Init(&cache.NullCache{})
	}
	beego.Info("正常初始化缓存配置.")
	cacheProvider := beego.AppConfig.String("cache_provider")
	if cacheProvider == "file" {
		cacheFilePath := beego.AppConfig.DefaultString("cache_file_path","./runtime/cache/")
		if strings.HasPrefix(cacheFilePath, "./") {
			cacheFilePath = filepath.Join(conf.WorkingDirectory, string(cacheFilePath[1:]))
		}
		fileCache := beegoCache.NewFileCache()


		fileConfig := make(map[string]string,0)

		fileConfig["CachePath"] =  cacheFilePath
		fileConfig["DirectoryLevel"] = beego.AppConfig.DefaultString("cache_file_dir_level","2")
		fileConfig["EmbedExpiry"] = beego.AppConfig.DefaultString("cache_file_expiry","120")
		fileConfig["FileSuffix"] = beego.AppConfig.DefaultString("cache_file_suffix",".bin")

		bc,err := json.Marshal(&fileConfig)
		if err != nil {
			beego.Error("初始化Redis缓存失败:",err)
			os.Exit(1)
		}

		fileCache.StartAndGC(string(bc))

		cache.Init(fileCache)

	}else if cacheProvider == "memory" {
		cacheInterval := beego.AppConfig.DefaultInt("cache_memory_interval",60)
		memory := beegoCache.NewMemoryCache()
		beegoCache.DefaultEvery = cacheInterval
		cache.Init(memory)
	}else if cacheProvider == "redis"{
		var redisConfig struct{
			Conn string `json:"conn"`
			Password string `json:"password"`
			DbNum int `json:"dbNum"`
		}
		redisConfig.DbNum = 0
		redisConfig.Conn = beego.AppConfig.DefaultString("cache_redis_host","")
		if pwd := beego.AppConfig.DefaultString("cache_redis_password","");pwd != "" {
			redisConfig.Password = pwd
		}
		if dbNum := beego.AppConfig.DefaultInt("cache_redis_db",0); dbNum > 0 {
			redisConfig.DbNum = dbNum
		}

		bc,err := json.Marshal(&redisConfig)
		if err != nil {
			beego.Error("初始化Redis缓存失败:",err)
			os.Exit(1)
		}
		redisCache,err := beegoCache.NewCache("redis",string(bc))

		if err != nil {
			beego.Error("初始化Redis缓存失败:",err)
			os.Exit(1)
		}

		cache.Init(redisCache)
	}else if cacheProvider == "memcache" {

		var memcacheConfig struct{
			Conn string `json:"conn"`
		}
		memcacheConfig.Conn = beego.AppConfig.DefaultString("cache_memcache_host","")

		bc,err := json.Marshal(&memcacheConfig)
		if err != nil {
			beego.Error("初始化Redis缓存失败:",err)
			os.Exit(1)
		}
		memcache,err := beegoCache.NewCache("memcache",string(bc))

		if err != nil {
			beego.Error("初始化Memcache缓存失败:",err)
			os.Exit(1)
		}

		cache.Init(memcache)

	}else {
		cache.Init(&cache.NullCache{})
		beego.Warn("不支持的缓存管道,缓存将禁用.")
		return
	}
	beego.Info("缓存初始化完成.")
}

func init() {

	if configPath ,err := filepath.Abs(conf.ConfigurationFile); err == nil {
		conf.ConfigurationFile = configPath
	}
	gocaptcha.ReadFonts("./static/fonts", ".ttf")
	gob.Register(models.Member{})

	if p, err := filepath.Abs(os.Args[0]); err == nil {
		conf.WorkingDirectory = filepath.Dir(p)
	}
}
