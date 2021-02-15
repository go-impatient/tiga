package schema

import (
	"time"

	"moocss.com/tiga/pkg/bcrypt"
	"moocss.com/tiga/pkg/database"
)

func init() {
	database.RegisterModel(&User{})
}

// User, 用户表
type User struct {
	Model

	Username              string    `gorm:"column:username;type:varchar(50);not null;default:''" json:"username"`          // 用户名
	Password              string    `gorm:"column:password;type:varchar(50);not null;default:''" json:"password"`          // 密码
	Salt                  string    `gorm:"column:salt;type:varchar(20);not null;default:''" json:"salt"`                  // 盐
	Realname              string    `gorm:"column:realname;type:varchar(50);default:''" json:"realname"`                   // 真实姓名
	Nickname              string    `gorm:"column:nickname;type:varchar(50);default:''" json:"nickname"`                   // 昵称
	Email                 string    `gorm:"column:email;type:varchar(80);default:''" json:"email"`                         // 邮箱
	Phone                 string    `gorm:"column:phone;type:varchar(20);not null;default:''" json:"phone"`                // 手机号码
	Sex                   int       `gorm:"column:sex;type:int(2);default:2" json:"sex"`                                   // 性别: 0 男性, 1 女性, 2未知
	Status                int       `gorm:"index:idx_user_status;column:status;type:int(2);default:1" json:"status"`       // 状态: 1启用, 0禁用
	IsDeleted             int       `gorm:"column:is_deleted;type:int(2);default:0" json:"is_deleted"`                     // 是否已删除 : 1删除, 0未删除
	CreateUser            string    `gorm:"column:create_user;type:varchar(64);default:''" json:"create_user"`             // 创建人
	UpdateUser            string    `gorm:"column:update_user;type:varchar(64);default:''" json:"update_user"`             // 修改人
	PasswordErrorLastTime time.Time `gorm:"column:password_error_last_time;type:datetime" json:"password_error_last_time"` // 最后一次输错密码时间
	PasswordErrorNum      int       `gorm:"column:password_error_num;type:int(11);default:0" json:"password_error_num"`    // 密码错误次数
	PasswordExpireTime    time.Time `gorm:"column:password_expire_time;type:datetime" json:"password_expire_time"`         // 密码过期时间
}

// TableName, 获取用户表名称
func (u *User) TableName() string {
	return "sys_user"
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = bcrypt.Compare(u.Password, pwd)
	return
}

// Encrypt the user password.
func (u *User) Encrypt() (err error) {
	password, err := bcrypt.Encrypt(u.Password)
	if err != nil {
		return err
	}
	u.Password = password
	return
}
