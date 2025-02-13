package internal

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"micro-services/pkg/utils"
	"micro-services/support-server/pkg/config"
	userPkg "micro-services/user-server/pkg"
	"net/http"
)

func SupportRegister(c *gin.Context) {
	type Register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var register Register
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "注册失败！",
		})
		return
	}
	uid, _ := uuid2.NewUUID()
	supportID := "sppt" + uid.String()

	hs, _ := userPkg.HashPassword(register.Password)

	_, err := config.MySqlClient.Exec("INSERT INTO b2c_support.agents ( name, password, role, create_at, support_id) VALUES (?,?,?,?,?)", register.Username, hs, 0, utils.GetTime(), supportID)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"supportID": supportID,
		},
	})
}
func SupportLogin(c *gin.Context) {
	type Login struct {
		SupportID string `json:"support_id"`
		Password  string `json:"password"`
	}
	var login Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "登录失败！",
		})
		return
	}
	var name, pwd string
	err := config.MySqlClient.QueryRow("SELECT name,password FROM b2c_support.agents WHERE support_id=?", login.SupportID).Scan(&name, &pwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg":  "登录失败，用户不存在！",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg":  "登录失败,数据库挂了！",
			})
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 400,
			"msg":  "登录失败，密码错误！",
		})
		return
	}
	// 把客服状态存入数据库 1 在线
	//config.RdClient.HMSet(config.Ctx, login.SupportID, "status", "1")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"name":      name,
			"supportID": login.SupportID,
		},
	})

}
