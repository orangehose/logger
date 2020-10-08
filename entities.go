package main

// Record ...
type Record struct {
	Time    string `gorm:"primaryKey"`
	Topic   string
	Message string
}
