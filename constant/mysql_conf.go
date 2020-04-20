package constant

var (
	MYSQL_CONF = &MysqlConf{
		UserName: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "fantim",
	}
)

type MysqlConf struct {
	UserName string
	Password string
	Host     string
	Port     int32
	DbName   string
}
