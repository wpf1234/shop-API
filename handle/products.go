package handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"testAPI/confInit"
	"testAPI/models"
)

const defaultRepertory = 100

func (g *Gin) GetList(c *gin.Context) {
	// 获取 claims ，直接解析 claims 获取登录用户的信息
	//_, ok := c.Get("claims")
	//if !ok {
	//	log.Error("Claims字段不存在!")
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    http.StatusInternalServerError,
	//		"data":    nil,
	//		"message": "没有获取到信息!",
	//	})
	//	return
	//}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		log.Error("获取页码失败: ", err)
		return
	}

	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		log.Error("获取每页条数失败: ", err)
		return
	}
	var (
		product models.Product
		list    []models.Product
		total   int
		res models.ProdRes
	)

	db := confInit.DB.Raw("select count(*) from list where is_del=?",0)
	err = db.Row().Scan(&total)
	if err!=nil{
		log.Error("查询总数失败: ",err)
		return
	}
	db = confInit.DB.Raw(fmt.Sprintf(`select id,name,price,repertory from list where is_del=%d order by id limit %d,%d`,
		0,(page-1)*size, size))
	rows, err := db.Rows()
	if err != nil {
		log.Error("查询失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "查询商品列表失败!",
		})
		return
	}

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Repertory)
		if err != nil {
			log.Error("Error: ", err)
			return
		}
		list = append(list, product)
	}
	_ = rows.Close()
	fmt.Println("商品列表：", list)

	res.Total = total
	res.List = list

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "查询成功!",
	})
}

func (g *Gin) GetInfoByID( c *gin.Context){
	id,err:= strconv.Atoi(c.Query("id"))
	if err!=nil{
		log.Error("获取 ID 失败: ",err)
		return
	}

	var product models.Product
	product.Id = id
	db:= confInit.DB.Raw("select name,price,repertory from list where id=?",id)
	err = db.Row().Scan(&product.Name, &product.Price, &product.Repertory)
	if err!=nil{
		log.Error("Error: ",err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    product,
		"message": "查询成功!",
	})
}

func (g *Gin) ModifyProd(c *gin.Context)  {
	var product models.Product
	err :=c.BindJSON(&product)
	if err!=nil{
		log.Error("获取请求数据失败: ",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "获取请求数据失败!",
		})
		return
	}
	fmt.Println("Response data: ",product)

	db:=confInit.DB.Exec("update list set name=?,price=?,repertory=? where id=?",
		product.Name,product.Price,product.Repertory,product.Id)
	fmt.Println("Update: ",db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    nil,
		"message": "商品已更新!",
	})
}

func (g *Gin) AddNewProd(c *gin.Context)  {
	var product models.Product
	var isDel int
	err := c.BindJSON(&product)
	if err!=nil{
		log.Error("获取请求数据失败: ",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "获取请求数据失败!",
		})
		return
	}

	// 新增之前需要先判断是否有该商品了
	db := confInit.DB.Raw("select id,name,price,repertory,is_del from list where name=?",product.Name)
	_ = db.Row().Scan(&product.Id,&product.Name,&product.Price,&product.Repertory,&isDel)
	fmt.Println("ID: ",product.Id)
	if product.Id != 0{
		// 商品已经存在
		if isDel == 0{
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    "The data already exists!",
				"message": "商品已存在!",
			})
			return
		}else{
			db = confInit.DB.Exec("update list set is_del=? where id=?",0,product.Id)
			fmt.Println("Is delete: ",0," & rows affected: ",db.RowsAffected)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"data":    nil,
				"message": "新增商品成功!",
			})
			return
		}

	}

	if product.Repertory == 0{
		product.Repertory = defaultRepertory
	}

	db = confInit.DB.Exec("insert into list set name=?,price=?,repertory=?",
		product.Name,product.Price,product.Repertory)

	fmt.Println("Insert: ",db.RowsAffected)

	if db.RowsAffected == 1{
		// 插入数据完成
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"data":    nil,
			"message": "新增商品成功!",
		})
	}
}

func (g *Gin) DeleteProd(c *gin.Context)  {
	id,err:=strconv.Atoi(c.Query("id"))
	if err!=nil{
		log.Error("ID转换失败: ",err)
		return
	}

	fmt.Println("ID=====>",id)
	// 将 ID 为传入值的数据的 is_del 设置为 true
	db:=confInit.DB.Exec("update list set is_del=? where id=?",1,id)
	fmt.Println("Update delete: ",db.RowsAffected)

	if db.RowsAffected == 1{
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"data":    nil,
			"message": "删除商品成功!",
		})
	}
}