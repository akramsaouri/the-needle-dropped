package main

// VideoLink ...
type VideoLink struct {
	HREF string `xml:"href,attr" json:"href"`
}

// VideoAuthor ...
type VideoAuthor struct {
	Name string `xml:"name" json:"name"`
	URI  string `xml:"uri" json:"uri"`
}

// Video represents xml feed entry
type Video struct {
	Title  string      `xml:"title" json:"title"`
	Link   VideoLink   `xml:"link" json:"link"`
	Author VideoAuthor `xml:"author" json:"author"`
}

// Feed represents feed body
type Feed struct {
	Entry Video `xml:"entry"`
}
