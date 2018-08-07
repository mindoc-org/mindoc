package commands

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/json"

	"github.com/astaxie/beego"
	beegoCache "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/gocaptcha"
	"github.com/lifei6671/mindoc/cache"
	"github.com/lifei6671/mindoc/conf"
	"github.com/lifei6671/mindoc/models"
	"github.com/lifei6671/mindoc/utils/filetil"
	"github.com/astaxie/beego/cache/redis"
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
		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			beego.Error("加载时区配置信息失败,请检查是否存在ZONEINFO环境变量->", err)
		}

		port := beego.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

		if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
			beego.Error("注册默认数据库失败->", err)
			os.Exit(1)
		}
	} else if adapter == "sqlite3" {
		orm.DefaultTimeLoc = time.UTC
		database := beego.AppConfig.String("db_database")
		if strings.HasPrefix(database, "./") {
			database = filepath.Join(conf.WorkingDirectory, string(database[1:]))
		}

		dbPath := filepath.Dir(database)
		os.MkdirAll(dbPath, 0777)

		err := orm.RegisterDataBase("default", "sqlite3", database)

		if err != nil {
			beego.Error("注册默认数据库失败->", err)
		}
	} else {
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
		new(models.Blog),
	)
	gob.Register(models.Blog{})
	gob.Register(models.Document{})
	//migrate.RegisterMigration()
}

// RegisterLogger 注册日志
func RegisterLogger(log string) {

	logs.SetLogFuncCall(true)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)

	if beego.AppConfig.DefaultBool("log_is_async", true) {
		logs.Async(1e3)
	}
	if log == "" {
		log = conf.WorkingDir("runtime","logs")
	}

	logPath := filepath.Join(log, "log.log")

	if _, err := os.Stat(log); os.IsNotExist(err) {
		os.MkdirAll(log, 0777)
	}

	config := make(map[string]interface{}, 1)

	config["filename"] = logPath
	config["perm"] = "0755"
	config["rotate"] = true

	if maxLines := beego.AppConfig.DefaultInt("log_maxlines", 1000000); maxLines > 0 {
		config["maxLines"] = maxLines
	}
	if maxSize := beego.AppConfig.DefaultInt("log_maxsize", 1<<28); maxSize > 0 {
		config["maxsize"] = maxSize
	}
	if !beego.AppConfig.DefaultBool("log_daily", true) {
		config["daily"] = false
	}
	if maxDays := beego.AppConfig.DefaultInt("log_maxdays", 7); maxDays > 0 {
		config["maxdays"] = maxDays
	}
	if level := beego.AppConfig.DefaultString("log_level", "Trace"); level != "" {
		switch level {
		case "Emergency":
			config["level"] = beego.LevelEmergency;break
		case "Alert":
			config["level"] = beego.LevelAlert;break
		case "Critical":
			config["level"] = beego.LevelCritical;break
		case "Error":
			config["level"] = beego.LevelError; break
		case "Warning":
			config["level"] = beego.LevelWarning; break
		case "Notice":
			config["level"] = beego.LevelNotice; break
		case "Informational":
			config["level"] = beego.LevelInformational;break
		case "Debug":
			config["level"] = beego.LevelDebug;break
		}
	}
	b, err := json.Marshal(config);
	if  err != nil {
		beego.Error("初始化文件日志时出错 ->",err)
		beego.SetLogger("file", `{"filename":"`+ logPath + `"}`)
	}else{
		beego.SetLogger(logs.AdapterFile, string(b))
	}



	beego.SetLogFuncCall(true)
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if len(os.Args) >= 2 && os.Args[1] == "install" {
		ResolveCommand(os.Args[2:])
		Install()
	} else if len(os.Args) >= 2 && os.Args[1] == "version" {
		CheckUpdate()
		os.Exit(0)
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
		//如果没有设置cdn，则使用baseURL拼接
		if cdn == "" {
			baseUrl := beego.AppConfig.DefaultString("baseurl", "")

			if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
				return baseUrl + p[1:]
			}
			if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
				return baseUrl + "/" + p
			}
			return baseUrl + p
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})

	beego.AddFuncMap("cdnjs", conf.URLForWithCdnJs)
	beego.AddFuncMap("cdncss", conf.URLForWithCdnCss)
	beego.AddFuncMap("cdnimg", conf.URLForWithCdnImage)
	//重写url生成，支持配置域名以及域名前缀
	beego.AddFuncMap("urlfor", conf.URLFor)
	beego.AddFuncMap("date_format", func(t time.Time, format string) string {
		return t.Local().Format(format)
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
		conf.LogFile = conf.WorkingDir("runtime","logs")
	}
	if conf.ConfigurationFile == "" {
		conf.ConfigurationFile = conf.WorkingDir( "conf", "app.conf")
		config := conf.WorkingDir("conf", "app.conf.example")
		if !filetil.FileExists(conf.ConfigurationFile) && filetil.FileExists(config) {
			filetil.CopyFile(conf.ConfigurationFile, config)
		}
	}
	if err := gocaptcha.ReadFonts(conf.WorkingDir( "static", "fonts"), ".ttf");err != nil {
		log.Fatal("读取字体文件时出错 -> ",err)
	}

	if err := beego.LoadAppConfig("ini", conf.ConfigurationFile);err != nil {
		log.Fatal("An error occurred:", err)
	}
	conf.AutoLoadDelay = beego.AppConfig.DefaultInt("config_auto_delay",0)
	uploads := conf.WorkingDir("uploads")

	os.MkdirAll(uploads, 0666)

	beego.BConfig.WebConfig.StaticDir["/static"] = filepath.Join(conf.WorkingDirectory, "static")
	beego.BConfig.WebConfig.StaticDir["/uploads"] = uploads
	beego.BConfig.WebConfig.ViewsPath = conf.WorkingDir("views")

	fonts := conf.WorkingDir("static", "fonts")

	if !filetil.FileExists(fonts) {
		log.Fatal("Font path not exist.")
	}
	gocaptcha.ReadFonts(filepath.Join(conf.WorkingDirectory, "static", "fonts"), ".ttf")

	RegisterDataBase()
	RegisterCache()
	RegisterModel()
	RegisterLogger(conf.LogFile)

	ModifyPassword()

}

//注册缓存管道
func RegisterCache() {
	isOpenCache := beego.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		cache.Init(&cache.NullCache{})
	}
	beego.Info("正常初始化缓存配置.")
	cacheProvider := beego.AppConfig.String("cache_provider")
	if cacheProvider == "file" {
		cacheFilePath := beego.AppConfig.DefaultString("cache_file_path", "./runtime/cache/")
		if strings.HasPrefix(cacheFilePath, "./") {
			cacheFilePath = filepath.Join(conf.WorkingDirectory, string(cacheFilePath[1:]))
		}
		fileCache := beegoCache.NewFileCache()

		fileConfig := make(map[string]string, 0)

		fileConfig["CachePath"] = cacheFilePath
		fileConfig["DirectoryLevel"] = beego.AppConfig.DefaultString("cache_file_dir_level", "2")
		fileConfig["EmbedExpiry"] = beego.AppConfig.DefaultString("cache_file_expiry", "120")
		fileConfig["FileSuffix"] = beego.AppConfig.DefaultString("cache_file_suffix", ".bin")

		bc, err := json.Marshal(&fileConfig)
		if err != nil {
			beego.Error("初始化file缓存失败:", err)
			os.Exit(1)
		}

		fileCache.StartAndGC(string(bc))

		cache.Init(fileCache)

	} else if cacheProvider == "memory" {
		cacheInterval := beego.AppConfig.DefaultInt("cache_memory_interval", 60)
		memory := beegoCache.NewMemoryCache()
		beegoCache.DefaultEvery = cacheInterval
		cache.Init(memory)
	} else if cacheProvider == "redis" {
		//设置Redis前缀
		if key := beego.AppConfig.DefaultString("cache_redis_prefix",""); key != "" {
			redis.DefaultKey = key
		}
		var redisConfig struct {
			Conn     string `json:"conn"`
			Password string `json:"password"`
			DbNum    int    `json:"dbNum"`
		}
		redisConfig.DbNum = 0
		redisConfig.Conn = beego.AppConfig.DefaultString("cache_redis_host", "")
		if pwd := beego.AppConfig.DefaultString("cache_redis_password", ""); pwd != "" {
			redisConfig.Password = pwd
		}
		if dbNum := beego.AppConfig.DefaultInt("cache_redis_db", 0); dbNum > 0 {
			redisConfig.DbNum = dbNum
		}

		bc, err := json.Marshal(&redisConfig)
		if err != nil {
			beego.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}
		redisCache, err := beegoCache.NewCache("redis", string(bc))

		if err != nil {
			beego.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}

		cache.Init(redisCache)
	} else if cacheProvider == "memcache" {

		var memcacheConfig struct {
			Conn string `json:"conn"`
		}
		memcacheConfig.Conn = beego.AppConfig.DefaultString("cache_memcache_host", "")

		bc, err := json.Marshal(&memcacheConfig)
		if err != nil {
			beego.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}
		memcache, err := beegoCache.NewCache("memcache", string(bc))

		if err != nil {
			beego.Error("初始化Memcache缓存失败:", err)
			os.Exit(1)
		}

		cache.Init(memcache)

	} else {
		cache.Init(&cache.NullCache{})
		beego.Warn("不支持的缓存管道,缓存将禁用 ->" ,cacheProvider)
		return
	}
	beego.Info("缓存初始化完成.")
}

//自动加载配置文件.修改了监听端口号和数据库配置无法自动生效.
func RegisterAutoLoadConfig()  {
	if conf.AutoLoadDelay > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(conf.AutoLoadDelay))

		go func() {
			f,err := os.Stat(conf.ConfigurationFile)
			if err != nil {
				beego.Error("读取配置文件时出错 ->",err)
				return
			}
			modTime := f.ModTime()
			for {
				select {
				case <-ticker.C:
					f,err := os.Stat(conf.ConfigurationFile)
					if err != nil {
						beego.Error("读取配置文件时出错 ->",err)
						break
					}
					if modTime != f.ModTime() {
						if err := beego.LoadAppConfig("ini", conf.ConfigurationFile); err != nil {
							beego.Error("An error occurred:", err)
							break
						}
						modTime = f.ModTime()
						beego.Info("配置文件已加载")
					}
				}
			}
		}()
	}
}

func init() {

	if configPath, err := filepath.Abs(conf.ConfigurationFile); err == nil {
		conf.ConfigurationFile = configPath
	}
	gocaptcha.ReadFonts("./static/fonts", ".ttf")
	gob.Register(models.Member{})

	if p, err := filepath.Abs(os.Args[0]); err == nil {
		conf.WorkingDirectory = filepath.Dir(p)
	}
}
