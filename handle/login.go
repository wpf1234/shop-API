package handle

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"testAPI/confInit"
	"testAPI/models"
	"testAPI/utils"
	"time"
)

func (g *Gin) Login(c *gin.Context)  {
	var login models.Login
	var res models.LoginRes
	var id int
	var password string

	err:=c.BindJSON(&login)
	// d41d8cd98f00b204e9800998ecf8427e
	fmt.Println("Login: ",login)
	if err!=nil{
		log.Error("获取请求事变: ",err)
		c.JSON( http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": err,
			"message": "请求失败!",
		})
		return
	}

	// 获取到登录信息之后，用用户名和密码去查询数据库，并获取 ID
	db := confInit.DB.Raw("select id,password from user where username=? ",
		login.Username)
	err = db.Row().Scan(&id,&password)
	if err!=nil{
		log.Error("查询用户 ID 失败: ",err)
		c.JSON( http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": err,
			"message": "用户信息有误!",
		})
		return
	}
	if id == 0{
		log.Info("用户信息不存在!")
		c.JSON( http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": nil,
			"message": "用户信息不存在!",
		})
		return
	}
	// d41d8cd98f00b204e9800998ecf8427e
	fmt.Println("密码: ",password)
	// 查到该用户，比较登录密码
	pwd := utils.StrMd5(login.Password)
	fmt.Println("MD5: ",pwd)
	if password != pwd {
		log.Info("用户密码错误!")
		c.JSON( http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"data": nil,
			"message": "用户密码错误!",
		})
		return
	}

	// 生成 Token
	claims := &utils.MyClaims{
		Id: id,
		Username: login.Username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                                             // 签名生效时间
			ExpiresAt: time.Now().Add(time.Duration(7*utils.ExpireTime) * time.Hour).Unix(), // 过期时间
		},
	}

	token, err := utils.GetToken(claims)
	if err != nil {
		log.Error("生成Token失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "生成Token失败!",
		})
		return
	}

	res.Id = id
	res.Username = login.Username
	res.Token = token

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "登录成功!",
	})

	t := time.Now().UnixNano() / 10e5
	db = confInit.DB.Exec("update user set login_tm=? where id=?",t,id)

	fmt.Println("数据更新: ",db.RowsAffected)
}
