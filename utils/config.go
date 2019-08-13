package utils

import (
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

func Parser(cfg *models.Config, conf string) {

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
		cfg.Api.Username = viper.GetString("API.username")
		cfg.Api.Password = viper.GetString("API.password")
		cfg.Api.GetPerPage = viper.GetInt("API.get_per_page")

		// GIT
		cfg.Git.Repo = viper.GetString("GIT.repo")
		cfg.Git.Directory = viper.GetString("GIT.directory")
		cfg.Git.Token = viper.GetString("GIT.token")

		// RT
		cfg.RackTables.Production = viper.GetString("RT.pro")
		cfg.RackTables.Stage = viper.GetString("RT.stage")

		// DB
		cfg.Database.Host = viper.GetString("DB.db_host")
		cfg.Database.Provider = viper.GetString("DB.db_provider")
		cfg.Database.Username = viper.GetString("DB.db_user")
		cfg.Database.Password = viper.GetString("DB.db_password")
		cfg.Database.DBName = viper.GetString("DB.db_schema")

		// WEB
		cfg.Web.Port = viper.GetInt("WEB.port")
		cfg.Web.JWTSecret = viper.GetString("WEB.jwt_secret")

		// LOGGING
		cfg.Logging.ErrorLog = viper.GetString("LOGGING.err_log")
		cfg.Logging.AccessLog = viper.GetString("LOGGING.acc_log")
		cfg.Logging.TraceLog = viper.GetString("LOGGING.trace_log")

		// LDAP
		cfg.LDAP.BindUser = viper.GetString("LDAP.bin_user")
		cfg.LDAP.BindPassword = viper.GetString("LDAP.bin_pass")
		cfg.LDAP.LdapServer = viper.GetStringSlice("LDAP.ldap_server")
		cfg.LDAP.LdapServerPort = viper.GetInt("LDAP.ldap_server_port")
		cfg.LDAP.BaseDn = viper.GetString("LDAP.base_dn")
		cfg.LDAP.MatchStr = viper.GetString("LDAP.match_string")

		// Other
		cfg.MasterHost = viper.GetString("master_host")
	}
}
