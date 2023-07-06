package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tealeg/xlsx"
)

type ep struct {
	Id  int    `db:"id"`
	S1  string `db:"s1"`
	S2  string `db:"s2"`
	S3  string `db:"s3"`
	S4  string `db:"s4"`
	S5  string `db:"s5"`
	S6  string `db:"s6"`
	S7  string `db:"s7"`
	S8  string `db:"s8"`
	S9  string `db:"s9"`
	S10 string `db:"s10"`
	S11 string `db:"s11"`
	S12 string `db:"s12"`
	S13 string `db:"s13"`
	S14 string `db:"s14"`
	S15 string `db:"s15"`
	S16 string `db:"s16"`
	S17 string `db:"s17"`
	S18 string `db:"s18"`
	S19 string `db:"s19"`
	S20 string `db:"s20"`
	S21 string `db:"s21"`
	S22 string `db:"s22"`
	S23 string `db:"s23"`
	S24 string `db:"s24"`
	S25 string `db:"s25"`
	S26 string `db:"s26"`
	S27 string `db:"s27"`
	S28 string `db:"s28"`
	S29 string `db:"s29"`
	S30 string `db:"s30"`
	S31 string `db:"s31"`
	S32 string `db:"s32"`
	S33 string `db:"s33"`
	S34 string `db:"s34"`
	S35 string `db:"s35"`
	S36 string `db:"s36"`
	S37 string `db:"s37"`
	S38 string `db:"s38"`
	S39 string `db:"s39"`
	S40 string `db:"s40"`
	S41 string `db:"s41"`
	S42 string `db:"s42"`
	S43 string `db:"s43"`
	S44 string `db:"s44"`
	S45 string `db:"s45"`
	S46 string `db:"s46"`
	S47 string `db:"s47"`
	S48 string `db:"s48"`
	S49 string `db:"s49"`
}

// 用户结构体
type Users struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"4`
}

// 数据库指针
var db *sqlx.DB

// 初始化数据库连接，init()方法系统会在动在main方法之前执行。
func init() {
	database, err := sqlx.Open("mysql", "root:qy20021003@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Println("open mysql failed,", err)
	}
	db = database
}
func mselect(userchan chan []ep, downchan chan int) {
	for i := 0; i <= 1000000000; {
		var users []ep

		sql := "select * from ep where " + strconv.Itoa(i) + "<=id limit 10000"
		//sql := "select user_id,username,sex,email from user where 10000<=user_id limit 10000"
		err := db.Select(&users, sql)
		if err != nil {
			fmt.Println("exec failed, ", err)
			return
		}
		if len(users) == 0 {
			//fmt.Println(0)
			downchan <- 0
			return
		} else {
			//fmt.Println(1)
			userchan <- users
			//fmt.Println(users)
		}
		i = users[len(users)-1].Id + 1
	}
	//userchan := make(chan []Users)
	//downchan := make(chan int)

	//fmt.Println("select succ:", users)
	return
}

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	ser := &Server{
		// 用的gin.Default()引擎，自带Logger and Recovery两个中间件，也可以用gin.New()，不带中间件
		engine: gin.New(),
	}
	return ser
}

//	type Users struct {
//		UserId   int    `db:"user_id"`
//		Username string `db:"username"`
//		Sex      string `db:"sex"`
//		Email    string `db:"email"`
//	}
//
// 造数据的方法
func (s *Server) getData(userchan chan []ep, downchan chan int) {
	mselect(userchan, downchan)
}

// 这是对应的4个接口，下面会具体说明一下不同接口的作用。
func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	s.engine.GET("/csv", s.csvApi)   //1,请求后直接下载csv文件
	s.engine.GET("/xlsx", s.xlsxApi) //请求后直接下载xlsx文件
	s.engine.Run(":9999")
}

// 启动进程
func main() {
	NewServer().Start()
}

// 接口一
// 请求接口后，会直接下载csv格式文件，使用 "encoding/csv" 包实现，代码里直接创建了文件，最后还删除，不然也会给服务器压力，或定期删除。
func (s *Server) csvApi(c *gin.Context) {
	downchan := make(chan int)
	userchan := make(chan []ep)
	filename, err := s.toCsv(userchan, downchan)
	if err != nil {
		fmt.Println("t.toCsv() failed == ", err)
	}
	if filename == "" {
		fmt.Println("export excel file failed == ", filename)
	}
	defer func() {
		err := os.Remove("./" + filename) //下载后，删除文件
		if err != nil {
			fmt.Println("remove  excel file failed", err)
		}
	}()
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream") //设置下载文件格式，流式下载
	c.File("./" + filename)                                           //直接返回文件

}

// 接口一 function
func (t *Server) toCsv(userchan chan []ep, downchan chan int) (string, error) {
	//获取数据
	go t.getData(userchan, downchan)
	strTime := time.Now().Format("20060102150405")
	//创建csv文件
	filename := fmt.Sprintf("学生信息-%s.csv", strTime)
	xlsFile, fErr := os.OpenFile("./"+filename, os.O_RDWR|os.O_CREATE, 0766)
	if fErr != nil {
		fmt.Println("Export:created excel file failed ==", fErr)
		return "", fErr
	}
	defer xlsFile.Close()
	//开始写入内容
	//写入UTF-8 BOM,此处如果不写入就会导致写入的汉字乱码
	xlsFile.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(xlsFile)
	wStr.Write([]string{"学号", "姓名", "性别", "邮箱地址"})
	var data []ep
	var check int
	for {
		select {
		case data = <-userchan:
			for _, s := range data {
				wStr.Write([]string{strconv.Itoa(s.Id), s.S1, s.S2, s.S3, s.S4, s.S5, s.S6, s.S7, s.S8, s.S9, s.S10, s.S11, s.S12, s.S13, s.S14, s.S15, s.S16, s.S17, s.S18, s.S19, s.S20, s.S21, s.S22, s.S23, s.S24, s.S25, s.S26, s.S27, s.S28, s.S29, s.S30, s.S31, s.S32, s.S33, s.S34, s.S35, s.S36, s.S37, s.S38, s.S39, s.S40, s.S41, s.S42, s.S43, s.S44, s.S45, s.S46, s.S47, s.S48, s.S49})
			}
		case check = <-downchan:
			wStr.Flush() //写入文件
			fmt.Println(check)
			return filename, nil
		}
	}
	return "", nil
}

// xlsx接口

func (s *Server) xlsxApi(c *gin.Context) {
	downchan := make(chan int)
	userchan := make(chan []ep)
	filename, err := s.toxlsx(userchan, downchan)
	if err != nil {
		fmt.Println("t.toxlsx() failed == ", err)
	}
	if filename == "" {
		fmt.Println("export excel file failed == ", filename)
	}
	defer func() {
		err := os.Remove("./" + filename) //下载后，删除文件
		if err != nil {
			fmt.Println("remove  excel file failed", err)
		}
	}()
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream") //设置下载文件格式，流式下载
	c.File("./" + filename)                                           //直接返回文件
}

// 接口一 function
func (t *Server) toxlsx(userchan chan []ep, downchan chan int) (string, error) {
	//获取数据
	go t.getData(userchan, downchan)
	strTime := time.Now().Format("20060102150405")
	//创建csv文件
	filename := fmt.Sprintf("学生信息-%s.xlsx", strTime)
	//开始写入内容
	//写入UTF-8 BOM,此处如果不写入就会导致写入的汉字乱码
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Println("init xlsx file failed, err == ", err.Error())
		return "", err
	}
	row := sheet.AddRow()
	row.SetHeightCM(1)
	cell := row.AddCell()
	cell.Value = "id"
	cell = row.AddCell()
	cell.Value = "s1"
	cell = row.AddCell()
	cell.Value = "s2"
	cell = row.AddCell()
	cell.Value = "s3"
	cell = row.AddCell()
	cell.Value = "s4"
	cell = row.AddCell()
	cell.Value = "s5"
	cell = row.AddCell()
	cell.Value = "s6"
	cell = row.AddCell()
	cell.Value = "s7"
	cell = row.AddCell()
	cell.Value = "s8"
	cell = row.AddCell()
	cell.Value = "s9"
	cell = row.AddCell()
	cell.Value = "s10"
	cell = row.AddCell()
	cell.Value = "s11"
	cell = row.AddCell()
	cell.Value = "s12"
	cell = row.AddCell()
	cell.Value = "s13"
	cell = row.AddCell()
	cell.Value = "s14"
	cell = row.AddCell()
	cell.Value = "s15"
	cell = row.AddCell()
	cell.Value = "s16"
	cell = row.AddCell()
	cell.Value = "s17"
	cell = row.AddCell()
	cell.Value = "s18"
	cell = row.AddCell()
	cell.Value = "s19"
	cell = row.AddCell()
	cell.Value = "s20"
	cell = row.AddCell()
	cell.Value = "s21"
	cell = row.AddCell()
	cell.Value = "s22"
	cell = row.AddCell()
	cell.Value = "s23"
	cell = row.AddCell()
	cell.Value = "s24"
	cell = row.AddCell()
	cell.Value = "s25"
	cell = row.AddCell()
	cell.Value = "s26"
	cell = row.AddCell()
	cell.Value = "s27"
	cell = row.AddCell()
	cell.Value = "s28"
	cell = row.AddCell()
	cell.Value = "s29"
	cell = row.AddCell()
	cell.Value = "s30"
	cell = row.AddCell()
	cell.Value = "s31"
	cell = row.AddCell()
	cell.Value = "s32"
	cell = row.AddCell()
	cell.Value = "s33"
	cell = row.AddCell()
	cell.Value = "s34"
	cell = row.AddCell()
	cell.Value = "s35"
	cell = row.AddCell()
	cell.Value = "s36"
	cell = row.AddCell()
	cell.Value = "s37"
	cell = row.AddCell()
	cell.Value = "s38"
	cell = row.AddCell()
	cell.Value = "s39"
	cell = row.AddCell()
	cell.Value = "s40"
	cell = row.AddCell()
	cell.Value = "s41"
	cell = row.AddCell()
	cell.Value = "s42"
	cell = row.AddCell()
	cell.Value = "s43"
	cell = row.AddCell()
	cell.Value = "s44"
	cell = row.AddCell()
	cell.Value = "s45"
	cell = row.AddCell()
	cell.Value = "s46"
	cell = row.AddCell()
	cell.Value = "s47"
	cell = row.AddCell()
	cell.Value = "s48"
	cell = row.AddCell()
	cell.Value = "s49"
	cell = row.AddCell()
	cell.Value = "s50"
	var data []ep
	for {
		select {
		case <-userchan:
			data = <-userchan
			for _, v := range data {
				row1 := sheet.AddRow()
				cell = row1.AddCell()
				cell.Value = strconv.Itoa(v.Id)
				cell = row1.AddCell()
				cell.Value = v.S1
				cell = row1.AddCell()
				cell.Value = v.S2
				cell = row1.AddCell()
				cell.Value = v.S3
				cell = row1.AddCell()
				cell.Value = v.S4
				cell = row1.AddCell()
				cell.Value = v.S5
				cell = row1.AddCell()
				cell.Value = v.S6
				cell = row1.AddCell()
				cell.Value = v.S7
				cell = row1.AddCell()
				cell.Value = v.S8
				cell = row1.AddCell()
				cell.Value = v.S9
				cell = row1.AddCell()
				cell.Value = v.S10
				cell = row1.AddCell()
				cell.Value = v.S11
				cell = row1.AddCell()
				cell.Value = v.S12
				cell = row1.AddCell()
				cell.Value = v.S13
				cell = row1.AddCell()
				cell.Value = v.S14
				cell = row1.AddCell()
				cell.Value = v.S15
				cell = row1.AddCell()
				cell.Value = v.S16
				cell = row1.AddCell()
				cell.Value = v.S17
				cell = row1.AddCell()
				cell.Value = v.S18
				cell = row1.AddCell()
				cell.Value = v.S19
				cell = row1.AddCell()
				cell.Value = v.S20
				cell = row1.AddCell()
				cell.Value = v.S21
				cell = row1.AddCell()
				cell.Value = v.S22
				cell = row1.AddCell()
				cell.Value = v.S23
				cell = row1.AddCell()
				cell.Value = v.S24
				cell = row1.AddCell()
				cell.Value = v.S25
				cell = row1.AddCell()
				cell.Value = v.S26
				cell = row1.AddCell()
				cell.Value = v.S27
				cell = row1.AddCell()
				cell.Value = v.S28
				cell = row1.AddCell()
				cell.Value = v.S29
				cell = row1.AddCell()
				cell.Value = v.S30
				cell = row1.AddCell()
				cell.Value = v.S31
				cell = row1.AddCell()
				cell.Value = v.S32
				cell = row1.AddCell()
				cell.Value = v.S33
				cell = row1.AddCell()
				cell.Value = v.S34
				cell = row1.AddCell()
				cell.Value = v.S35
				cell = row1.AddCell()
				cell.Value = v.S36
				cell = row1.AddCell()
				cell.Value = v.S37
				cell = row1.AddCell()
				cell.Value = v.S38
				cell = row1.AddCell()
				cell.Value = v.S39
				cell = row1.AddCell()
				cell.Value = v.S40
				cell = row1.AddCell()
				cell.Value = v.S41
				cell = row1.AddCell()
				cell.Value = v.S42
				cell = row1.AddCell()
				cell.Value = v.S43
				cell = row1.AddCell()
				cell.Value = v.S44
				cell = row1.AddCell()
				cell.Value = v.S45
				cell = row1.AddCell()
				cell.Value = v.S46
				cell = row1.AddCell()
				cell.Value = v.S47
				cell = row1.AddCell()
				cell.Value = v.S48
				cell = row1.AddCell()
				cell.Value = v.S49
			}
		case <-downchan:
			err = file.Save("./" + filename)
			if err != nil {
				panic(err)
			}
			return filename, nil
		}

	}

	return "", err
}
