package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"micro-services/pkg/utils"
	userPkg "micro-services/user-server/pkg"
	"micro-services/user-server/pkg/config"
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

func GetUserInfoByEmail(email string) (*int64, string, string, string, error) {
	var id *int64
	var name, role, createAt string
	query := "SELECT user_id, username, role,create_at FROM b2c_user.users WHERE email = ?"
	row := config.MySqlClient.QueryRow(query, email) // 使用 QueryRow 获取单行结果
	err := row.Scan(&id, &name, &role, &createAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", "", "", errors.New("用户不存在")
		} else {
			return nil, "", "", "", errors.New("数据库错误")
		}
	} else {
		return id, name, role, createAt, nil
	}
}

func GetAvatarUrlAndBioById(id int64) (string, string, error) {
	var avatarUrl, bio string
	query := "SELECT avatar_url,bio FROM b2c_user.user_profiles WHERE user_id=?"
	row := config.MySqlClient.QueryRow(query, id)
	err := row.Scan(&avatarUrl, &bio)
	if err != nil {
		return "", "", err
	}
	return avatarUrl, bio, nil
}

// 初始化用户信息表
func SaveUserInfo(name string, email string, role string) (
	userId int64, createAt string, err error) {
	currentTime := utils.GetTime()
	// 初始化 users 表
	query := "INSERT INTO b2c_user.users (username, email, role,create_at,update_at) VALUES (?, ?, ?, ?, ?)"
	result, err := config.MySqlClient.Exec(query, name, email, role, currentTime, currentTime)
	if err != nil {
		return 0, "", fmt.Errorf("could not insert user info: %v", err)
	}

	userId, err = result.LastInsertId()
	if err != nil {
		return 0, "", fmt.Errorf("could not get last insert ID: %v", err)
	}

	// 初始化 user_profiles 表
	query = "INSERT INTO b2c_user.user_profiles (user_id,avatar_url,create_at,update_at) VALUES (?,?,?,?)"
	_, err = config.MySqlClient.Exec(query, userId, "b2c/default.jpg", currentTime, currentTime)
	if err != nil {
		return 0, "", fmt.Errorf("could not insert user profile: %v", err)
	}
	return userId, currentTime, nil
}

func CheckNameAndPwd(name string, pwd string) (
	int64, string, string, string, string, string, string, error) {
	if pwd == "" {
		return 0, "", "", "", "", "", "", errors.New("GLB-001")
	}
	query := "SELECT user_id,role,password_hash,create_at,email FROM b2c_user.users WHERE username=?"
	row := config.MySqlClient.QueryRow(query, name)
	var id int64
	var role string
	var oldPwdHash string
	var createAt string
	var email string
	err := row.Scan(&id, &role, &oldPwdHash, &createAt, &email)
	if err != nil {
		log.Printf("Error checking if name exists: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", "", "", "", "", "", errors.New("GLB-001")
		} else {
			return 0, "", "", "", "", "", "", errors.New("GLB-003")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(oldPwdHash), []byte(pwd))

	if err != nil {
		log.Printf("Error checking if password is correct: %v", err)
		return 0, "", "", "", "", "", "", errors.New("GLB-001")
	}

	query = "SELECT avatar_url,bio FROM b2c_user.user_profiles WHERE user_id=?"
	row = config.MySqlClient.QueryRow(query, id)
	var avatarUrl string
	var bio sql.NullString
	err = row.Scan(&avatarUrl, &bio)
	if err != nil {
		log.Printf("Error checking if user profile exists: %v", err)
		return 0, "", "", "", "", "", "", errors.New("GBL-003")
	}
	if !bio.Valid {
		bio.String = ""
	}
	return id, name, role, avatarUrl, createAt, bio.String, email, nil
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
	currentTime := utils.GetTime()
	name := utils.Filter(username)
	query := "UPDATE b2c_user.users SET username=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, name, currentTime, id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ChangeEmail(id int64, email string) error {
	currentTime := utils.GetTime()
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
	currentTime := utils.GetTime()
	// 对密码进行加密
	hs, err := userPkg.HashPassword(newPwd)
	if err != nil {
		return err
	}
	query := "UPDATE b2c_user.users SET password_hash=?,update_at=? WHERE user_id=?"
	_, err = config.MySqlClient.Exec(query, hs, currentTime, id)
	if err != nil {
		//fmt.Println("密码存入出错：", err)
		return err
	}
	return nil
}

func GetUserIdByEmail(email string) (int64, error) {
	query := "SELECT user_id FROM b2c_user.users WHERE email=?"
	row := config.MySqlClient.QueryRow(query, email)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetEmailByUserId(id int64) (string, error) {
	query := "SELECT email FROM b2c_user.users WHERE user_id=?"
	row := config.MySqlClient.QueryRow(query, id)
	var email string
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func IsUserIdExist(userId int64) bool {
	var exists bool
	err := config.MySqlClient.QueryRow("SELECT EXISTS(SELECT 1 FROM b2c_user.users WHERE user_id = ?)", userId).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetUserInfoByUserId(userId int64) (avatarUrl, name, email, bio, createAt string, err error) {
	// 使用 sql.NullString 来处理可能为 NULL 的字段
	var avatarUrlNull, nameNull, emailNull, bioNull, createAtNull sql.NullString

	// 查询用户基本信息和头像URL
	query := `
		SELECT u.username, u.email, u.create_at, p.avatar_url, p.bio 
		FROM b2c_user.users u 
		LEFT JOIN b2c_user.user_profiles p 
		ON u.user_id = p.user_id 
		WHERE u.user_id = ?`

	row := config.MySqlClient.QueryRow(query, userId)
	err = row.Scan(&nameNull, &emailNull, &createAtNull, &avatarUrlNull, &bioNull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", "", "", "", errors.New("用户不存在")
		}
		return "", "", "", "", "", fmt.Errorf("数据库错误: %v", err)
	}

	// 检查 NullString 是否有效
	if nameNull.Valid {
		name = nameNull.String
	} else {
		name = ""
	}

	if emailNull.Valid {
		email = emailNull.String
	} else {
		email = ""
	}

	if createAtNull.Valid {
		createAt = createAtNull.String
	} else {
		createAt = ""
	}

	if avatarUrlNull.Valid {
		avatarUrl = avatarUrlNull.String
	} else {
		avatarUrl = "" // 默认头像 URL 或空字符串
	}

	if bioNull.Valid {
		bio = bioNull.String
	} else {
		bio = "" // 默认空的个人简介
	}

	return avatarUrl, name, email, bio, createAt, nil
}

func ChangeAvatar(userId int64, avatarUrl string) error {
	currentTime := utils.GetTime()
	query := "UPDATE b2c_user.user_profiles SET avatar_url=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, avatarUrl, currentTime, userId)
	if err != nil {
		return err
	}
	return nil
}

func ChangeBio(userId int64, bio string) error {
	currentTime := utils.GetTime()
	query := "UPDATE b2c_user.user_profiles SET bio=?,update_at=? WHERE user_id=?"
	_, err := config.MySqlClient.Exec(query, bio, currentTime, userId)
	if err != nil {
		return err
	}
	return nil
}
