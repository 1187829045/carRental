package model

import "time"

type Car struct {
	CarNumber   string    `json:"carnumber"`
	CarType     string    `json:"cartype"`
	Color       string    `json:"color"`
	Price       float64   `json:"price"`
	RentPrice   float64   `json:"rentprice"`
	Deposit     float64   `json:"deposit"`
	IsRenting   int       `json:"isrenting"`
	Description string    `json:"description"`
	CarImg      string    `json:"carimg"`
	CreateTime  time.Time `json:"createtime"`
}

type Customer struct {
	Identity   string    `json:"identity"`
	CustName   string    `json:"custname"`
	Sex        int       `json:"sex"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Career     string    `json:"career"`
	CreateTime time.Time `json:"createtime"`
}

type Rent struct {
	RentID     string     `json:"rentid"`
	Price      float64    `json:"price"`
	BeginDate  time.Time  `json:"begindate"`
	ReturnDate *time.Time `json:"returndate"`
	RentFlag   int        `json:"rentflag"`
	Identity   string     `json:"identity"`
	CarNumber  string     `json:"carnumber"`
	OperName   string     `json:"opername"`
	CreateTime time.Time  `json:"createtime"`
}

type Check struct {
	CheckID    string    `json:"checkid"`
	CheckDate  time.Time `json:"checkdate"`
	CheckDesc  string    `json:"checkdesc"`
	Problem    string    `json:"problem"`
	PayMoney   float64   `json:"paymoney"`
	OperName   string    `json:"opername"`
	RentID     string    `json:"rentid"`
	CreateTime time.Time `json:"createtime"`
}

type Franchisee struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

