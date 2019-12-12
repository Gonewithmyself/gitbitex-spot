// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"fmt"
	"sync"

	"github.com/gitbitex/gitbitex-spot/conf"
	"github.com/gitbitex/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

var dbs sync.Map
var loc sync.Mutex

type Store struct {
	db *gorm.DB
}

func get(dbname string) models.Store {
	val, ok := dbs.Load(dbname)
	if ok {
		return val.(models.Store)
	}

	loc.Lock()
	defer loc.Unlock()
	val, ok = dbs.Load(dbname)
	if ok {
		return val.(models.Store)
	}

	db, err := initDb(dbname)
	if err != nil {
		panic(err)
	}

	store := NewStore(db)
	dbs.Store(dbname, store)
	return store
}

func SharedStore(dbname string) models.Store {
	return get(dbname)
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func initDb(dbname string) (db *gorm.DB, err error) {
	cfg := conf.GetConfig()
	if dbname == "" {
		dbname = cfg.DataSource.Database
	}

	url := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		cfg.DataSource.User, cfg.DataSource.Password, cfg.DataSource.Addr, dbname)
	db, err = gorm.Open(cfg.DataSource.DriverName, url)
	if err != nil {
		return
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "g_" + defaultTableName
	}

	// if cfg.DataSource.EnableAutoMigrate {
	// 	var tables = []interface{}{
	// 		&models.Account{},
	// 		&models.Order{},
	// 		&models.Product{},
	// 		&models.Trade{},
	// 		&models.Fill{},
	// 		&models.User{},
	// 		&models.Bill{},
	// 		&models.Tick{},
	// 		&models.Config{},
	// 	}
	// 	for _, table := range tables {
	// 		log.Infof("migrating database, table: %v", reflect.TypeOf(table))
	// 		if err = gdb.AutoMigrate(table).Error; err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return
}

func (s *Store) BeginTx() (models.Store, error) {
	db := s.db.Begin()
	if db.Error != nil {
		return nil, db.Error
	}
	return NewStore(db), nil
}

func (s *Store) Rollback() error {
	return s.db.Rollback().Error
}

func (s *Store) CommitTx() error {
	return s.db.Commit().Error
}
