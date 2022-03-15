package db

import (
	"database/sql" // 标准库，用于操作数据库
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 导入mysql的驱动
	"strings"
	"time"
)

var db *sql.DB

func InitMysql(userName, password, ip, port, dbName string) (err error) {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8&parseTime=true&loc=Local"}, "")
	db, err = sql.Open("mysql", path)
	if err != nil {
		return err
	}
	fmt.Println("connect success")
	var adminTableSql = `CREATE TABLE IF NOT EXISTS admin(
		id int NOT NULL AUTO_INCREMENT,
		maxEther  int DEFAULT '200',
		TransEther  int DEFAULT '0',
		AddrAmount int DEFAULT '0',
		Dates varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
        PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`
	var userTableSql = `CREATE TABLE IF NOT EXISTS user(
		id int NOT NULL AUTO_INCREMENT,
		address varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
		applyCount  int DEFAULT '0',
		ApplyTimes int DEFAULT '0',
		latestTrans int DEFAULT NULL,
        Dates varchar(255) DEFAULT NULL,
		PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;`

	// 执行建表语句
	_, err = db.Exec(adminTableSql)
	_, err = db.Exec(userTableSql)
	if err != nil {
		return err
	}

	return
}
func ShowUser(_address string) (err error, strErr string) {
	nowTime := time.Now().Format("2006-01-02")
	var sqlStr = `SELECT address,applyCount,ApplyTimes,LatestTrans,Dates FROM user where address = ? AND Dates = ?`
	rows, err := db.Query(sqlStr, _address, nowTime)
	defer rows.Close()
	if err != nil {
		return err, ""
	}
	var adminStr = `SELECT maxEther,TransEther,AddrAmount,Dates FROM admin where Dates = ?`
	adRows, err := db.Query(adminStr, nowTime)
	defer adRows.Close()
	if err != nil {
		return err, ""
	}
	var ad Admin
	for adRows.Next() {
		err = adRows.Scan(&ad.MaxEther, &ad.TransEther, &ad.AddrAmount, &ad.Date)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return err, "数据库出错，请稍后重试"
		}
	}
	var u User
	for rows.Next() {
		err = rows.Scan(&u.Address, &u.ApplyCount, &u.ApplyTimes, &u.LatestTrans, &u.Dates)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return err, "数据库出错，请稍后重试"
		}
	}
	var sum = `SELECT count(address) FROM user WHERE Dates = ?`
	addrCount, err := db.Query(sum, nowTime)
	for addrCount.Next() {
		err = addrCount.Scan(&ad.AddrAmount)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return err, "数据库出错，请稍后重试"
		}
	}
	var insSum = `UPDATE admin SET AddrAmount = ? WHERE Dates = ?`
	_, err = db.Exec(insSum, ad.AddrAmount, nowTime)
	if ad.Date == "" {
		var update = `INSERT INTO admin (TransEther,Dates) VALUES(?,?)`
		var insSql = `INSERT INTO user (address,applyCount,ApplyTimes,LatestTrans,Dates) VALUES (?,?,?,?,?)`
		_, err = db.Exec(update, 1, nowTime)
		_, err = db.Exec(insSql, _address, 1, 1, time.Now().Unix(), nowTime)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return err, "数据库出错，请稍后重试"
		}
	} else if ad.TransEther < ad.MaxEther {

		if u.Address == "" {
			var newSql = `INSERT INTO user (address,applyCount,ApplyTimes,LatestTrans,Dates) VALUES (?,?,?,?,?)`
			var tranStr = `UPDATE admin SET TransEther = ?,AddrAmount = ? where Dates = ? `
			_, err = db.Exec(newSql, _address, 1, 1, time.Now().Unix(), nowTime)
			_, err = db.Exec(tranStr, ad.TransEther+1, ad.AddrAmount+1, nowTime)
			//_,err = db.Exec(adAmountSql,ad.AddrAmount+1,nowTime)
			if err != nil {
				return err, "数据库出错，请稍后重试"
			}
		} else if u.Dates != nowTime {

			var insertData = `INSERT INTO user (address,applyCount,ApplyTimes,LatestTrans,Dates) VALUES(?,?,?,?,?)`
			var tranStr = `UPDATE admin SET TransEther = ? where Dates = ? `
			_, err = db.Exec(insertData, u.Address, 1, 1, time.Now().Unix(), nowTime)
			_, err = db.Exec(tranStr, ad.TransEther+1, nowTime)
			if err != nil {
				return err, "数据库出错，请稍后重试"
			}
		} else if time.Now().Unix()-u.LatestTrans > 60 && u.ApplyTimes < 10 && u.ApplyCount < 10 {
			var updateSql = `UPDATE user SET applyCount = ?,ApplyTimes = ?,latestTrans = ? where address = ? AND Dates = ? `
			var tranStr = `UPDATE admin SET TransEther = ? where Dates = ? `
			_, err = db.Exec(tranStr, ad.TransEther+1, nowTime)
			_, err = db.Exec(updateSql, u.ApplyCount+1, u.ApplyTimes+1, time.Now().Unix(), u.Address, nowTime)
			if err != nil {
				return err, "申请出错"
			}
		} else {
			return errors.New("not allow to receive"), "申 请 间 隔 需 大 于 一 分 钟  或 今 日 申 请 次 数 用 尽！"
		}

	} else {
		return errors.New("today is out of maxBalance"), "今 日 试 币 已 发 放 完 毕 ！"
	}

	return
}
