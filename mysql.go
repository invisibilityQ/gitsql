package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func mwrite() {

}

//s1 := gosql{
//	user: Users{
//		UserId:   0,
//		Username: "",
//		Sex:      "",
//		Email:    "",
//	},
//}
//s2 := gosql{
//	user: Users{
//		UserId:   0,
//		Username: "",
//		Sex:      "",
//		Email:    "",
//	},
//}
//sql1 := "insert into user(username,sex, email)values (?,?,?)"
//sql2 := "insert into user(username,sex, email)values (?,?,?)"
//s1.user.Username = "user01"
//s1.user.Sex = "man"
//s1.user.Email = "user01@163.com"
//s1.user.Username = "user02"
//s1.user.Sex = "woman"
//s1.user.Email = "user02@163.com"
////执行SQL语句
//r1, err := db.Exec(sql1, s1.user.Username, s1.user.Sex, s1.user.Email)
//if err != nil {
//	fmt.Println("exec failed,", err)
//	return
//}
//r2, err := db.Exec(sql2, s2.user.Username, s2.user.Sex, s2.user.Email)
//if err != nil {
//	fmt.Println("exec failed,", err)
//	return
//}
//
////查询最后一天用户ID，判断是否插入成功
//id1, err := r1.LastInsertId()
//if err != nil {
//	fmt.Println("exec failed,", err)
//	return
//}
//fmt.Println("insert succ", id1)
//id2, err := r2.LastInsertId()
//if err != nil {
//	fmt.Println("exec failed,", err)
//	return
//}
//fmt.Println("insert succ", id2)
