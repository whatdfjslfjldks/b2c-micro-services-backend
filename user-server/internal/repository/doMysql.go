package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/config"
	"time"
)

func IsEmailExist(email string) bool {
	// 假设你使用的是数据库连接对象 db
	query := "SELECT EXISTS(SELECT 1 FROM b2c_user.users WHERE email = ?)"

	var exists bool
	err := config.MySqlClient.QueryRow(query, email).Scan(&exists)
	if err != nil {
		// 处理错误
		fmt.Println("Error checking if email exists:", err)
		return false
	}

	return exists
}

func GetUserInfoByEmail(email string) (*int64, string, string, error) {
	var id *int64
	var name, role string
	query := "SELECT user_id, username, role FROM b2c_user.users WHERE email = ?"
	row := config.MySqlClient.QueryRow(query, email) // 使用 QueryRow 获取单行结果
	err := row.Scan(&id, &name, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", errors.New("用户不存在")
		} else {
			return nil, "", "", errors.New("数据库错误")
		}
	} else {
		return id, name, role, nil
	}
}

func GetAvatarUrlById(id int64) (string, error) {
	var avatarUrl string
	query := "SELECT avatar_url FROM b2c_user.user_profiles WHERE user_id=?"
	row := config.MySqlClient.QueryRow(query, id)
	err := row.Scan(&avatarUrl)
	if err != nil {
		return "", err
	}
	return avatarUrl, nil
}

// 初始化用户信息表
func SaveUserInfo(name string, email string, role string) (
	userId int64, err error) {
	currentTime := time.Now()
	// 初始化 users 表
	query := "INSERT INTO b2c_user.users (username, email, role,create_at,update_at) VALUES (?, ?, ?, ?, ?)"
	result, err := config.MySqlClient.Exec(query, name, email, role, currentTime, currentTime)
	if err != nil {
		return 0, fmt.Errorf("could not insert user info: %v", err)
	}

	userId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("could not get last insert ID: %v", err)
	}

	// 初始化 user_profiles 表
	query = "INSERT INTO b2c_user.user_profiles (user_id,avatar_url,create_at,update_at) VALUES (?,?,?,?)"
	_, err = config.MySqlClient.Exec(query, userId, "default", currentTime, currentTime)
	if err != nil {
		return 0, fmt.Errorf("could not insert user profile: %v", err)
	}
	return userId, nil
}

func CheckNameAndPwd(name string, pwd string) (
	int64, string, string, string, error) {
	if pwd == "" {
		return 0, "", "", "", errors.New("密码不能为空")
	}
	var id int64
	var role string
	// TODO 还没有定义密码的加密解密方法，这里把密码加密后和数据库的对比

	query := "SELECT user_id,role FROM b2c_user.users WHERE username=? AND password_hash=?"
	row := config.MySqlClient.QueryRow(query, name, pwd)
	err := row.Scan(&id, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", "", "", errors.New("用户名或密码错误")
		} else {
			return 0, "", "", "", errors.New("数据库错误")
		}
	}
	query = "SELECT avatar_url FROM b2c_user.user_profiles WHERE user_id=?"
	row = config.MySqlClient.QueryRow(query, id)
	var avatarUrl string
	err = row.Scan(&avatarUrl)
	if err != nil {
		return 0, "", "", "", errors.New("数据库错误")
	}
	return id, name, role, avatarUrl, nil
}
func IsUsernameExist(username string) bool {
	query := "SELECT EXISTS(SELECT 1 FROM b2c_user.users WHERE username = ?)"
	var exists bool
	err := config.MySqlClient.QueryRow(query, username).Scan(&exists)
	if err != nil {
		// 处理错误
		fmt.Println("Error checking if username exists:", err)
		return false
	}
	return exists
}

func ChangeUsername(id int64, username string) error {
	currentTime := time.Now()
	query := "UPDATE b2c_user.users SET username=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, username, currentTime, id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ChangeEmail(id int64, email string) error {
	currentTime := time.Now()
	query := "UPDATE b2c_user.users SET email=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, email, currentTime, id)
	if err != nil {
		return err
	}
	return nil
}

func CheckOldPassword(id int64, oldPwd string) error {
	query := "SELECT password_hash FROM b2c_user.users WHERE user_id=?"
	row := config.MySqlClient.QueryRow(query, id)
	var oldPwdHash string
	err := row.Scan(&oldPwdHash)
	if err != nil {
		return err
	}
	// 为了避免被 <时间攻击>，采用bcrypt内置函数进行比较
	// 而且存的密码有“加盐”计算，即时间戳，每次加密的结果都
	// 不一样，不能直接比较
	err = bcrypt.CompareHashAndPassword([]byte(oldPwdHash), []byte(oldPwd))
	if err != nil {
		return errors.New("旧密码错误")
	}
	return nil
}

func SaveNewPassword(id int64, newPwd string) error {
	currentTime := time.Now()
	// 对密码进行加密
	hs, err := userPkg.HashPassword(newPwd)
	if err != nil {
		return err
	}
	query := "UPDATE b2c_user.users SET password_hash=?,update_at=? WHERE user_id=?"
	_, err = config.MySqlClient.Exec(query, hs, currentTime, id)
	if err != nil {
		fmt.Println("密码存入出错：", err)
		return err
	}
	return nil
}
