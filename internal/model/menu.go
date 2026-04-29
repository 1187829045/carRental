package model

// Menu mirrors sys_menu.
type Menu struct {
	ID        int     `json:"id"`
	Pid       int     `json:"pid"`
	Title     string  `json:"title"`
	Href      *string `json:"href"`
	Spread    int     `json:"spread"`
	Target    *string `json:"target"`
	Icon      *string `json:"icon"`
	Available int     `json:"available"`
}

// TreeNode matches the JSON fields consumed by LayUI on the index left menu.
type TreeNode struct {
	ID     int        `json:"id"`
	Pid    int        `json:"pid"`
	Title  string     `json:"title"`
	Icon   *string    `json:"icon,omitempty"`
	Href   *string    `json:"href,omitempty"`
	Spread bool       `json:"spread"`
	Target *string    `json:"target,omitempty"`
	Child  []TreeNode `json:"children,omitempty"`
}

