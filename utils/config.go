package utils

import (
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

func Parser(globConf *models.Config, conf string) {

	cDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath(cDir)
	if conf == "" {
		viper.SetConfigName("config")
	} else {
		viper.AddConfigPath("./conf/")
		viper.SetConfigName(conf)
	}
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(errors.WithStack(err))
	} else {
		// API
		globConf.Api.Username = viper.GetString("API.username")
		globConf.Api.Password = viper.GetString("API.password")
		globConf.Api.GetPerPage = viper.GetInt("API.get_per_page")

		// RT
		globConf.RackTables.Production = viper.GetString("RT.pro")
		globConf.RackTables.Stage = viper.GetString("RT.stage")

		// DB
		globConf.Database.Host = viper.GetString("DB.db_host")
		globConf.Database.Provider = viper.GetString("DB.db_provider")
		globConf.Database.Username = viper.GetString("DB.db_user")
		globConf.Database.Password = viper.GetString("DB.db_password")
		globConf.Database.DBName = viper.GetString("DB.db_schema")

		// WEB
		globConf.Web.Port = viper.GetInt("WEB.port")
		globConf.Web.JWTSecret = viper.GetString("WEB.jwt_secret")

		// LOGGING
		globConf.Logging.ErrorLog = viper.GetString("LOGGING.err_log")
		globConf.Logging.AccessLog = viper.GetString("LOGGING.acc_log")
		globConf.Logging.TraceLog = viper.GetString("LOGGING.trace_log")

		// LDAP
		globConf.LDAP.BindUser = viper.GetString("LDAP.bin_user")
		globConf.LDAP.BindPassword = viper.GetString("LDAP.bin_pass")
		globConf.LDAP.LdapServer = viper.GetStringSlice("LDAP.ldap_server")
		globConf.LDAP.LdapServerPort = viper.GetInt("LDAP.ldap_server_port")
		globConf.LDAP.BaseDn = viper.GetString("LDAP.base_dn")
		globConf.LDAP.MatchStr = viper.GetString("LDAP.match_string")

		// Other
		globConf.MasterHost = viper.GetString("master_host")
	}
}
