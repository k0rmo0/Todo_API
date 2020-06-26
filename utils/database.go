package utils

import (
	"fmt"
	"io/ioutil"

	//Needs to be imported, or it will cause an error .
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

//DBAccess ...
type DBAccess struct {
	SQLDB *sqlx.DB
}

// SQLAcc ...
var SQLAcc DBAccess

//DBConf ...
type DBConf struct {
	Name    string `yaml:"name"`
	User    string `yaml:"user"`
	Pass    string `yaml:"pass"`
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

//Configs ...
type Configs map[string]DBConf

//GetSQLDB ...
func GetSQLDB(env, path string) error {

	data, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error opening configuration", err.Error())
	}

	var cs Configs

	err = yaml.Unmarshal(data, &cs)
	dbconf := cs[env]

	dbDriver := "mysql"
	dbUser := dbconf.User
	dbPass := dbconf.Pass
	dbName := dbconf.Name
	dbAddress := dbconf.Address
	dbPort := dbconf.Port
	db, err := sqlx.Connect(dbDriver, dbUser+":"+dbPass+"@"+"tcp("+dbAddress+":"+dbPort+")/"+dbName)

	if err != nil {
		return err
	}

	SQLAcc.SQLDB = db

	return err
}

//GetSQLDB method ...
func (a DBAccess) GetSQLDB() *sqlx.DB {
	return a.SQLDB
}
