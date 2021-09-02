package models

import "time"

type T_commits struct {
	Commiter   string    `json:"commiter"`
	Addlines   int       `json:"addlines"`
	Dellines   int       `json:"dellines"`
	Totallines int       `json:"totallines"`
	Cdate      time.Time `json:"cdate"`
	Project    string    `json:"project"`
	Pid        string    `json:"pid"`
	Cid        string    `json:"cid"`
	Branch     string    `json:"branch"`
}
