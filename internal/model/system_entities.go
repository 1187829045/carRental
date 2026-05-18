package model

import "time"

type Role struct {
	RoleID        int    `json:"roleid"`
	RoleName      string `json:"rolename"`
	RoleDesc      string `json:"roledesc"`
	Available     int    `json:"available"`
	AssignedUsers string `json:"assignedusers,omitempty"`
}

type MenuManagerNode struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parentId"`
	Title    string `json:"title"`
	Spread   bool   `json:"spread"`
	CheckArr string `json:"checkArr,omitempty"`
}

type LoginLog struct {
	ID        int       `json:"id"`
	LoginName string    `json:"loginname"`
	LoginIP   string    `json:"loginip"`
	LoginTime time.Time `json:"logintime"`
}

type News struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createtime"`
	OperName   string    `json:"opername"`
}

type Message struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"createtime"`
	OperName   string    `json:"opername"`
}
