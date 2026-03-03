package commands

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"bytes"
	"encoding/json"
	"net/http"

	beegoCache "github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/memcache"
	"github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/howeyc/fsnotify"
	_ "github.com/lib/pq"
	"github.com/lifei6671/gocaptcha"
	"github.com/mindoc-org/mindoc/cache"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils/filetil"
)

// RegisterDataBase 注册数据库
func RegisterDataBase() {
	logs.Info("正在初始化数据库配置.")
	dbadapter, _ := web.AppConfig.String("db_adapter")
	orm.DefaultTimeLoc = time.Local
	orm.DefaultRowsLimit = -1

	if strings.EqualFold(dbadapter, "mysql") {
		host, _ := web.AppConfig.String("db_host")
		database, _ := web.AppConfig.String("db_database")
		username, _ := web.AppConfig.String("db_username")
		password, _ := web.AppConfig.String("db_password")

		timezone, _ := web.AppConfig.String("timezone")
		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			logs.Error("加载时区配置信息失败,请检查是否存在 ZONEINFO 环境变量->", err)
		}

		port, _ := web.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

		if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
			logs.Error("注册默认数据库失败->", err)
			os.Exit(1)
		}

	} else if strings.EqualFold(dbadapter, "sqlite3") {

		database, _ := web.AppConfig.String("db_database")
		if strings.HasPrefix(database, "./") {
			database = filepath.Join(conf.WorkingDirectory, string(database[1:]))
		}
		if p, err := filepath.Abs(database); err == nil {
			database = p
		}

		dbPath := filepath.Dir(database)

		if _, err := os.Stat(dbPath); err != nil && os.IsNotExist(err) {
			_ = os.MkdirAll(dbPath, 0777)
		}

		err := orm.RegisterDataBase("default", "sqlite3", database)

		if err != nil {
			logs.Error("注册默认数据库失败->", err)
		}
	} else if strings.EqualFold(dbadapter, "postgres") {
		host, _ := web.AppConfig.String("db_host")
		database, _ := web.AppConfig.String("db_database")
		username, _ := web.AppConfig.String("db_username")
		password, _ := web.AppConfig.String("db_password")
		sslmode, _ := web.AppConfig.String("db_sslmode")

		timezone, _ := web.AppConfig.String("timezone")
		location, err := time.LoadLocation(timezone)
		if err == nil {
			orm.DefaultTimeLoc = location
		} else {
			logs.Error("加载时区配置信息失败,请检查是否存在 ZONEINFO 环境变量->", err)
		}

		port, _ := web.AppConfig.String("db_port")

		dataSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, database, sslmode)

		if err := orm.RegisterDataBase("default", "postgres", dataSource); err != nil {
			logs.Error("注册默认数据库失败->", err)
			os.Exit(1)
		}

	} else {
		logs.Error("不支持的数据库类型.")
		os.Exit(1)
	}
	logs.Info("数据库初始化完成.")
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
		new(models.Template),
		new(models.Team),
		new(models.TeamMember),
		new(models.TeamRelationship),
		new(models.Itemsets),
		new(models.Comment),
		new(models.WorkWeixinAccount),
		new(models.DingTalkAccount),
	)
	gob.Register(models.Blog{})
	gob.Register(models.Document{})
	gob.Register(models.Template{})
	//migrate.RegisterMigration()
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		logs.Error("注册Model失败 ->", err)
		os.Exit(1)
	}
}

// RegisterLogger 注册日志
func RegisterLogger(log string) {

	logs.Reset()
	logs.SetLogFuncCall(true)
	_ = logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)

	if web.AppConfig.DefaultBool("log_is_async", true) {
		logs.Async(1e3)
	}
	if log == "" {
		logPath, err := filepath.Abs(web.AppConfig.DefaultString("log_path", conf.WorkingDir("runtime", "logs")))
		if err == nil {
			log = logPath
		} else {
			log = conf.WorkingDir("runtime", "logs")
		}
	}

	logPath := filepath.Join(log, "log.log")

	if _, err := os.Stat(log); os.IsNotExist(err) {
		_ = os.MkdirAll(log, 0755)
	}

	config := make(map[string]interface{}, 1)

	config["filename"] = logPath
	config["perm"] = "0755"
	config["rotate"] = true

	if maxLines := web.AppConfig.DefaultInt("log_maxlines", 1000000); maxLines > 0 {
		config["maxLines"] = maxLines
	}
	if maxSize := web.AppConfig.DefaultInt("log_maxsize", 1<<28); maxSize > 0 {
		config["maxsize"] = maxSize
	}
	if !web.AppConfig.DefaultBool("log_daily", true) {
		config["daily"] = false
	}
	if maxDays := web.AppConfig.DefaultInt("log_maxdays", 7); maxDays > 0 {
		config["maxdays"] = maxDays
	}
	if level := web.AppConfig.DefaultString("log_level", "Trace"); level != "" {
		switch level {
		case "Emergency":
			config["level"] = logs.LevelEmergency
		case "Alert":
			config["level"] = logs.LevelAlert
		case "Critical":
			config["level"] = logs.LevelCritical
		case "Error":
			config["level"] = logs.LevelError
		case "Warning":
			config["level"] = logs.LevelWarning
		case "Notice":
			config["level"] = logs.LevelNotice
		case "Informational":
			config["level"] = logs.LevelInformational
		case "Debug":
			config["level"] = logs.LevelDebug
		}
	}
	b, err := json.Marshal(config)
	if err != nil {
		logs.Error("初始化文件日志时出错 ->", err)
		_ = logs.SetLogger("file", `{"filename":"`+logPath+`"}`)
	} else {
		_ = logs.SetLogger(logs.AdapterFile, string(b))
	}

	logs.SetLogFuncCall(true)
}

// RunCommand 注册orm命令行工具
func RegisterCommand() {

	if len(os.Args) >= 2 && os.Args[1] == "install" {
		ResolveCommand(os.Args[2:])
		Install()
	} else if len(os.Args) >= 2 && os.Args[1] == "version" {
		CheckUpdate()
		os.Exit(0)
	} else if len(os.Args) >= 2 && os.Args[1] == "update" {
		Update()
		os.Exit(0)
	}

}

// 注册模板函数
func RegisterFunction() {
	err := web.AddFuncMap("config", models.GetOptionValue)

	if err != nil {
		logs.Error("注册函数 config 出错 ->", err)
		os.Exit(-1)
	}
	err = web.AddFuncMap("cdn", func(p string) string {
		cdn := web.AppConfig.DefaultString("cdn", "")
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			return p
		}
		//如果没有设置cdn，则使用baseURL拼接
		if cdn == "" {
			baseUrl := web.AppConfig.DefaultString("baseurl", "")

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
	if err != nil {
		logs.Error("注册函数 cdn 出错 ->", err)
		os.Exit(-1)
	}

	err = web.AddFuncMap("cdnjs", conf.URLForWithCdnJs)
	if err != nil {
		logs.Error("注册函数 cdnjs 出错 ->", err)
		os.Exit(-1)
	}
	err = web.AddFuncMap("cdncss", conf.URLForWithCdnCss)
	if err != nil {
		logs.Error("注册函数 cdncss 出错 ->", err)
		os.Exit(-1)
	}
	err = web.AddFuncMap("cdnimg", conf.URLForWithCdnImage)
	if err != nil {
		logs.Error("注册函数 cdnimg 出错 ->", err)
		os.Exit(-1)
	}
	//重写url生成，支持配置域名以及域名前缀
	err = web.AddFuncMap("urlfor", conf.URLFor)
	if err != nil {
		logs.Error("注册函数 urlfor 出错 ->", err)
		os.Exit(-1)
	}
	//读取配置值(未作任何转换)
	err = web.AddFuncMap("conf", conf.CONF)
	if err != nil {
		logs.Error("注册函数 conf 出错 ->", err)
		os.Exit(-1)
	}
	err = web.AddFuncMap("date_format", func(t time.Time, format string) string {
		return t.Local().Format(format)
	})
	if err != nil {
		logs.Error("注册函数 date_format 出错 ->", err)
		os.Exit(-1)
	}

	err = web.AddFuncMap("i18n", i18n.Tr)
	if err != nil {
		logs.Error("注册函数 i18n 出错 ->", err)
		os.Exit(-1)
	}

	i18nList, err := web.AppConfig.String("i18n_list")
	if err != nil {
		logs.Error("error : failed to read i18n_list config ->", err)
		i18nList = ""
	}
	if i18nList == "" { // 之所以分开判断是因为读取出的配置也可能是空串
		logs.Error("error : config `i18n_list` is empty, please add config item like format: `i18n_list=zh-CN:简体中文|en-US:English`")
		i18nList = "zh-cn:简体中文|en-us:English|ru-ru:Русский" // 没有配置时给个默认配置，避免啥语言都没有
	}

	langs := strings.Split(i18nList, "|")
	i18nMap := make(map[string]string)
	for _, langItem := range langs {
		langItemSplit := strings.Split(langItem, ":")
		if len(langItemSplit) < 2 {
			logs.Error("error: language config value `" + langItem + "` for `i18n_list` format error")
			continue
		}
		lang := langItemSplit[0]
		i18nMap[lang] = langItemSplit[1]
		if err := i18n.SetMessage(lang, "conf/lang/"+lang+".ini"); err != nil {
			logs.Error("Fail to set message file: " + err.Error())
			return
		}
	}

	i18nMapBytes, err := json.Marshal(i18nMap)
	if err != nil {
		logs.Error("error: Fail to marshal i18n map, " + err.Error())
		i18nMapBytes = []byte("{}")
	}

	err = web.AppConfig.Set("i18n_map", string(i18nMapBytes))
	if err != nil {
		logs.Error("error: Fail to set i18n_map, " + err.Error())
	}
}

// 解析命令
func ResolveCommand(args []string) {
	flagSet := flag.NewFlagSet("MinDoc command: ", flag.ExitOnError)
	flagSet.StringVar(&conf.ConfigurationFile, "config", "", "MinDoc configuration file.")
	flagSet.StringVar(&conf.WorkingDirectory, "dir", "", "MinDoc working directory.")
	flagSet.StringVar(&conf.LogFile, "log", "", "MinDoc log file path.")

	if err := flagSet.Parse(args); err != nil {
		log.Fatal("解析命令失败 ->", err)
	}

	if conf.WorkingDirectory == "" {
		if p, err := filepath.Abs(os.Args[0]); err == nil {
			conf.WorkingDirectory = filepath.Dir(p)
		}
	}

	if conf.ConfigurationFile == "" {
		conf.ConfigurationFile = conf.WorkingDir("conf", "app.conf")
		config := conf.WorkingDir("conf", "app.conf.example")
		if !filetil.FileExists(conf.ConfigurationFile) && filetil.FileExists(config) {
			_ = filetil.CopyFile(conf.ConfigurationFile, config)
		}
	}
	if err := gocaptcha.ReadFonts(conf.WorkingDir("static", "fonts"), ".ttf"); err != nil {
		log.Fatal("读取字体文件时出错 -> ", err)
	}

	if err := web.LoadAppConfig("ini", conf.ConfigurationFile); err != nil {
		log.Fatal("An error occurred:", err)
	}
	if conf.LogFile == "" {
		logPath, err := filepath.Abs(web.AppConfig.DefaultString("log_path", conf.WorkingDir("runtime", "logs")))
		if err == nil {
			conf.LogFile = logPath
		} else {
			conf.LogFile = conf.WorkingDir("runtime", "logs")
		}
	}

	conf.AutoLoadDelay = web.AppConfig.DefaultInt("config_auto_delay", 0)
	uploads := conf.WorkingDir("uploads")

	_ = os.MkdirAll(uploads, 0666)

	web.BConfig.WebConfig.StaticDir["/static"] = filepath.Join(conf.WorkingDirectory, "static")
	web.BConfig.WebConfig.StaticDir["/uploads"] = uploads
	web.BConfig.WebConfig.ViewsPath = conf.WorkingDir("views")
	web.BConfig.WebConfig.Session.SessionCookieSameSite = http.SameSiteDefaultMode
	var upload_file_size = conf.GetUploadFileSize()
	if upload_file_size > web.BConfig.MaxUploadSize {
		web.BConfig.MaxUploadSize = upload_file_size
	}

	fonts := conf.WorkingDir("static", "fonts")

	if !filetil.FileExists(fonts) {
		log.Fatal("Font path not exist.")
	}
	if err := gocaptcha.ReadFonts(filepath.Join(conf.WorkingDirectory, "static", "fonts"), ".ttf"); err != nil {
		log.Fatal("读取字体失败 ->", err)
	}

	RegisterDataBase()
	RegisterCache()
	RegisterModel()
	RegisterLogger(conf.LogFile)

	ModifyPassword()

}

// 注册缓存管道
func RegisterCache() {
	isOpenCache := web.AppConfig.DefaultBool("cache", false)
	if !isOpenCache {
		cache.Init(&cache.NullCache{})
		return
	}
	logs.Info("正常初始化缓存配置.")
	cacheProvider, _ := web.AppConfig.String("cache_provider")
	if cacheProvider == "file" {
		cacheFilePath := web.AppConfig.DefaultString("cache_file_path", "./runtime/cache/")
		if strings.HasPrefix(cacheFilePath, "./") {
			cacheFilePath = filepath.Join(conf.WorkingDirectory, string(cacheFilePath[1:]))
		}
		fileCache := beegoCache.NewFileCache()

		fileConfig := make(map[string]string, 0)

		fileConfig["CachePath"] = cacheFilePath
		fileConfig["DirectoryLevel"] = web.AppConfig.DefaultString("cache_file_dir_level", "2")
		fileConfig["EmbedExpiry"] = web.AppConfig.DefaultString("cache_file_expiry", "120")
		fileConfig["FileSuffix"] = web.AppConfig.DefaultString("cache_file_suffix", ".bin")

		bc, err := json.Marshal(&fileConfig)
		if err != nil {
			logs.Error("初始化file缓存失败:", err)
			os.Exit(1)
		}

		_ = fileCache.StartAndGC(string(bc))

		cache.Init(fileCache)

	} else if cacheProvider == "memory" {
		cacheInterval := web.AppConfig.DefaultInt("cache_memory_interval", 60)
		memory := beegoCache.NewMemoryCache()
		beegoCache.DefaultEvery = cacheInterval
		cache.Init(memory)
	} else if cacheProvider == "redis" {
		//设置Redis前缀
		if key := web.AppConfig.DefaultString("cache_redis_prefix", ""); key != "" {
			redis.DefaultKey = key
		}
		var redisConfig struct {
			Conn     string `json:"conn"`
			Password string `json:"password"`
			DbNum    string `json:"dbNum"`
		}
		redisConfig.DbNum = "0"
		redisConfig.Conn = web.AppConfig.DefaultString("cache_redis_host", "")
		if pwd := web.AppConfig.DefaultString("cache_redis_password", ""); pwd != "" {
			redisConfig.Password = pwd
		}
		if dbNum := web.AppConfig.DefaultInt("cache_redis_db", 0); dbNum > 0 {
			redisConfig.DbNum = strconv.Itoa(dbNum)
		}

		bc, err := json.Marshal(&redisConfig)
		if err != nil {
			logs.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}
		redisCache, err := beegoCache.NewCache("redis", string(bc))

		if err != nil {
			logs.Error("初始化Redis缓存失败:", err)
			os.Exit(1)
		}

		cache.Init(redisCache)
	} else if cacheProvider == "memcache" {

		var memcacheConfig struct {
			Conn string `json:"conn"`
		}
		memcacheConfig.Conn = web.AppConfig.DefaultString("cache_memcache_host", "")

		bc, err := json.Marshal(&memcacheConfig)
		if err != nil {
			logs.Error("初始化 Memcache 缓存失败 ->", err)
			os.Exit(1)
		}
		memcache, err := beegoCache.NewCache("memcache", string(bc))

		if err != nil {
			logs.Error("初始化 Memcache 缓存失败 ->", err)
			os.Exit(1)
		}

		cache.Init(memcache)

	} else {
		cache.Init(&cache.NullCache{})
		logs.Warn("不支持的缓存管道,缓存将禁用 ->", cacheProvider)
		return
	}
	logs.Info("缓存初始化完成.")
}

// 自动加载配置文件.修改了监听端口号和数据库配置无法自动生效.
func RegisterAutoLoadConfig() {
	if conf.AutoLoadDelay > 0 {

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			logs.Error("创建配置文件监控器失败 ->", err)
		}
		go func() {
			for {
				select {
				case ev := <-watcher.Event:
					//如果是修改了配置文件
					if ev.IsModify() {
						if err := web.LoadAppConfig("ini", conf.ConfigurationFile); err != nil {
							logs.Error("An error occurred ->", err)
							continue
						}
						RegisterCache()
						RegisterLogger("")
						logs.Info("配置文件已加载 ->", conf.ConfigurationFile)
					} else if ev.IsRename() {
						_ = watcher.WatchFlags(conf.ConfigurationFile, fsnotify.FSN_MODIFY|fsnotify.FSN_RENAME)
					}
					logs.Info(ev.String())
				case err := <-watcher.Error:
					logs.Error("配置文件监控器错误 ->", err)

				}
			}
		}()

		err = watcher.WatchFlags(conf.ConfigurationFile, fsnotify.FSN_MODIFY|fsnotify.FSN_RENAME)

		if err != nil {
			logs.Error("监控配置文件失败 ->", err)
		}
	}
}

// 注册错误处理方法.
func RegisterError() {
	web.ErrorHandler("404", func(writer http.ResponseWriter, request *http.Request) {
		var buf bytes.Buffer

		data := make(map[string]interface{})
		data["ErrorCode"] = 404
		data["ErrorMessage"] = "页面未找到或已删除"

		if err := web.ExecuteViewPathTemplate(&buf, "errors/error.tpl", web.BConfig.WebConfig.ViewsPath, data); err == nil {
			_, _ = fmt.Fprint(writer, buf.String())
		} else {
			_, _ = fmt.Fprint(writer, data["ErrorMessage"])
		}
	})
	web.ErrorHandler("401", func(writer http.ResponseWriter, request *http.Request) {
		var buf bytes.Buffer

		data := make(map[string]interface{})
		data["ErrorCode"] = 401
		data["ErrorMessage"] = "请与 Web 服务器的管理员联系，以确认您是否具有访问所请求资源的权限。"

		if err := web.ExecuteViewPathTemplate(&buf, "errors/error.tpl", web.BConfig.WebConfig.ViewsPath, data); err == nil {
			_, _ = fmt.Fprint(writer, buf.String())
		} else {
			_, _ = fmt.Fprint(writer, data["ErrorMessage"])
		}
	})
}

func init() {

	if configPath, err := filepath.Abs(conf.ConfigurationFile); err == nil {
		conf.ConfigurationFile = configPath
	}
	if err := gocaptcha.ReadFonts(conf.WorkingDir("static", "fonts"), ".ttf"); err != nil {
		log.Fatal("读取字体文件失败 ->", err)
	}
	gob.Register(models.Member{})

	if p, err := filepath.Abs(os.Args[0]); err == nil {
		conf.WorkingDirectory = filepath.Dir(p)
	}
}
