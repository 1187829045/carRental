package model

// User mirrors the subset of fields used by JSP pages and session checks.
// Table: sys_user
type User struct {
	UserID    int    `json:"userid"`
	LoginName string `json:"loginname"`
	RealName  string `json:"realname"`
	Type      int    `json:"type"` // 1 super admin, 2 normal
}

// SysUser is full user entity for user management pages.
type SysUser struct {
	UserID    int     `json:"userid"`
	LoginName string  `json:"loginname"`
	Identity  *string `json:"identity,omitempty"`
	RealName  *string `json:"realname,omitempty"`
	Sex       *int    `json:"sex,omitempty"`
	Address   *string `json:"address,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Pwd       *string `json:"pwd,omitempty"`
	Position  *string `json:"position,omitempty"`
	Type      *int    `json:"type,omitempty"`
	Available *int    `json:"available,omitempty"`
}
