package archivers

import (
	"regexp"
)

// Comic config struct
type Comic struct {
	Archiver      string
	StartURL      string
	Dir           string
	FileMatch     *regexp.Regexp
	FilePrefix    string
	PrevLinkMatch *regexp.Regexp
}

// Comics supported by the archiver
var Comics = map[string]Comic{
	"doa": Comic{
		Archiver:      "ComicPress",
		StartURL:      "https://www.dumbingofage.com/",
		Dir:           "doa",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.png)"),
		FilePrefix:    "https://www.dumbingofage.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(https://www.dumbingofage.com/[0-9a-zA-Z/-]+)\" class=\"navi navi-prev\""),
	},
	"campcomic": Comic{
		Archiver:      "ComicPress",
		StartURL:      "http://campcomic.com/comic",
		Dir:           "campcomic",
		FileMatch:     regexp.MustCompile("/katie/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://hw1.pa-cdn.com/camp/assets/img/katie/comics/",
		PrevLinkMatch: regexp.MustCompile("class=\"btn btnPrev\" href=\"(http://campcomic.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"gaia": Comic{
		Archiver:      "ComicPress",
		StartURL:      "http://www.sandraandwoo.com/gaia/",
		Dir:           "gaia",
		FileMatch:     regexp.MustCompile("/gaia/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://www.sandraandwoo.com/gaia/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://www.sandraandwoo.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)\" rel=\"prev\""),
	},
	"sandraandwoo": Comic{
		Archiver:      "ComicPress",
		StartURL:      "http://www.sandraandwoo.com/",
		Dir:           "sandraandwoo",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://www.sandraandwoo.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://www.sandraandwoo.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)\" rel=\"prev\""),
	},
	"gogetaroomie": Comic{
		Archiver:      "ComicPress",
		StartURL:      "https://www.gogetaroomie.com/",
		Dir:           "gogetaroomie",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.gogetaroomie.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.gogetaroomie.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"kiwiblitz": Comic{
		Archiver:      "ComicPress",
		StartURL:      "https://www.kiwiblitz.com/",
		Dir:           "kiwiblitz",
		FileMatch:     regexp.MustCompile("/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.kiwiblitz.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.kiwiblitz.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"sleeplessdomain": Comic{
		Archiver:      "ComicPress",
		StartURL:      "https://www.sleeplessdomain.com/",
		Dir:           "sleeplessdomain",
		FileMatch:     regexp.MustCompile("/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.sleeplessdomain.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.sleeplessdomain.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"superredundant": Comic{
		Archiver:      "ComicPress",
		StartURL:      "http://superredundant.com/",
		Dir:           "superredundant",
		FileMatch:     regexp.MustCompile("<img src=\"http://superredundant.com/wp-content/uploads/([^\"]+\\.(jpg|png|gif))"),
		FilePrefix:    "http://superredundant.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://superredundant.com/[\\?0-9a-zA-Z/=-]+)\" class=\"navi comic-nav-previous navi-prev\""),
	},
	"xkcd": Comic{
		Archiver:      "Xkcd",
		StartURL:      "https://xkcd.com/",
		Dir:           "xkcd",
		FileMatch:     regexp.MustCompile("//imgs.xkcd.com/comics/([^\"]+\\.png)"),
		FilePrefix:    "http://imgs.xkcd.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"/([0-9]+/)\""),
	},
}
