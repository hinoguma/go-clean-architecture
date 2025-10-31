package config

func GetDBHost() string {
	// switch by environment variables
	return "aurora-cluster.cluster- xxxxxxxx.ap-northeast-1.rds.amazonaws.com"
}

func GetDBPort() int {
	// switch by environment variables
	return 5432
}

func GetDBUser() string {
	// switch by environment variables
	return "username"
}

func GetDBPassword() string {
	// switch by environment variables
	return "password"
}

func GetDBName() string {
	// switch by environment variables
	return "your db name"
}

func GetSSLMode() string {
	// switch by environment variables
	return "required"
}
