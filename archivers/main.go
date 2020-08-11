package archivers

import (
	"regexp"
)

var protocolMatch = regexp.MustCompile("^https?:")
var basenameMatch = regexp.MustCompile("\\/?([^\\/]+\\.[A-Za-z]{3,4})$")

// Archive a comic
func Archive(dir string, comic Comic) {
	if comic.Archiver == "Generic" {
		Generic(comic.StartURL, dir, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch)
	}
	if comic.Archiver == "MultiImageGeneric" {
		MultiImageGeneric(comic.StartURL, dir, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch)
	}
	if comic.Archiver == "Sequential" {
		Sequential(dir, comic.FilePrefix, comic.SeqPattern, comic.SeqStart, comic.SeqEnd)
	}
}

// Comic base type for all archivers
type Comic struct {
	Archiver      string
	StartURL      string
	FileMatch     *regexp.Regexp
	FilePrefix    string
	PrevLinkMatch *regexp.Regexp
	SeqPattern    string
	SeqStart      int
	SeqEnd        int
}

// Comics supported by the archiver
var Comics = map[string]Comic{
	"doa": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.dumbingofage.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.png)"),
		FilePrefix:    "https://www.dumbingofage.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(https://www.dumbingofage.com/[0-9a-zA-Z/-]+)\" class=\"navi navi-prev\""),
	},
	"campcomic": Comic{
		Archiver:      "Generic",
		StartURL:      "http://campcomic.com/comic",
		FileMatch:     regexp.MustCompile("/katie/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://hw1.pa-cdn.com/camp/assets/img/katie/comics/",
		PrevLinkMatch: regexp.MustCompile("class=\"btn btnPrev\" href=\"(http://campcomic.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"gaia": Comic{
		Archiver:      "Generic",
		StartURL:      "http://www.sandraandwoo.com/gaia/",
		FileMatch:     regexp.MustCompile("/gaia/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://www.sandraandwoo.com/gaia/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://www.sandraandwoo.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)\" rel=\"prev\""),
	},
	"sandraandwoo": Comic{
		Archiver:      "Generic",
		StartURL:      "http://www.sandraandwoo.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://www.sandraandwoo.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://www.sandraandwoo.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)\" rel=\"prev\""),
	},
	"gogetaroomie": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.gogetaroomie.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.gogetaroomie.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.gogetaroomie.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"kiwiblitz": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.kiwiblitz.com/",
		FileMatch:     regexp.MustCompile("/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.kiwiblitz.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.kiwiblitz.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"sleeplessdomain": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.sleeplessdomain.com/",
		FileMatch:     regexp.MustCompile("/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.sleeplessdomain.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.sleeplessdomain.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"superredundant": Comic{
		Archiver:      "Generic",
		StartURL:      "http://superredundant.com/",
		FileMatch:     regexp.MustCompile("<img src=\"http://superredundant.com/wp-content/uploads/([^\"]+\\.(jpg|png|gif))"),
		FilePrefix:    "http://superredundant.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://superredundant.com/[\\?0-9a-zA-Z/=-]+)\" class=\"navi comic-nav-previous navi-prev\""),
	},
	"gunnerkrigg": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.gunnerkrigg.com/",
		FileMatch:     regexp.MustCompile("class=\"comic_image\" src=\"/comics/([0-9]+\\.(jpg|png|gif))\""),
		FilePrefix:    "https://www.gunnerkrigg.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(\\?p\\=[0-9]+)\"><img src=\"/images/prev"),
	},
	"channelate": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.channelate.com/",
		FileMatch:     regexp.MustCompile("img src=\"https://www.channelate.com/wp-content/uploads/(.+\\.png)"),
		FilePrefix:    "https://www.channelate.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile("href=\"(https://www.channelate.com/comic/[0-9a-zA-Z/-]+)\" class=\"navi comic-nav-previous navi-prev\""),
	},
	"iamarg": Comic{
		Archiver:      "Generic",
		StartURL:      "https://iamarg.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "https://iamarg.com/comics/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://iamarg.com/[0-9a-zA-Z/-]+)\" class=\"navi navi-prev\""),
	},
	"itswalky": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.itswalky.com/",
		FileMatch:     regexp.MustCompile("img src=\"https://www.itswalky.com/wp-content/uploads/(.+\\.png)"),
		FilePrefix:    "https://www.itswalky.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile("href=\"(https://www.itswalky.com/comic/[0-9a-zA-Z/-]+)\" class=\"comic-nav-base comic-nav-previous\""),
	},
	"letsspeakenglish": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.marycagle.com/letsspeakenglish/134-slow-motion",
		FileMatch:     regexp.MustCompile("/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.marycagle.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.marycagle.com/letsspeakenglish/[0-9a-zA-Z-]+)\""),
	},
	"loadingartist": Comic{
		Archiver:      "Generic",
		StartURL:      "https://loadingartist.com/latest",
		FileMatch:     regexp.MustCompile("/uploads/([0-9]+/[0-9]+/[0-9a-zA-Z-]+\\.[a-z]{3,4})"),
		FilePrefix:    "https://loadingartist.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile("class=\"normal highlight prev comic-thumb\" href=\"(https://loadingartist.com/comic/[0-9a-zA-Z-]+/?)\""),
	},
	"octopuspie": Comic{
		Archiver:      "Generic",
		StartURL:      "http://www.octopuspie.com/2017-06-05/1023-1026-thats-it/",
		FileMatch:     regexp.MustCompile("/strippy/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://www.octopuspie.com/strippy/",
		PrevLinkMatch: regexp.MustCompile("href=\"(http://www.octopuspie.com/[0-9a-zA-Z/_-]+)\" rel=\"prev\""),
	},
	"twogag": Comic{
		Archiver:      "Generic",
		StartURL:      "http://twogag.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "http://octopuspie.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(http://twogag.com/archives/[0-9a-zA-Z-]+)\""),
	},
	"whomp": Comic{
		Archiver:      "Generic",
		StartURL:      "https://www.whompcomic.com/",
		FileMatch:     regexp.MustCompile("/comics/(.+\\.(jpg|png|gif))"),
		FilePrefix:    "https://www.whompcomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"(https://www.whompcomic.com/comic/[0-9a-zA-Z-]+)\""),
	},
	"xkcd": Comic{
		Archiver:      "MultiImageGeneric",
		StartURL:      "https://xkcd.com/",
		FileMatch:     regexp.MustCompile("//imgs.xkcd.com/comics/([^\"]+\\.png)"),
		FilePrefix:    "http://imgs.xkcd.com/comics/",
		PrevLinkMatch: regexp.MustCompile("rel=\"prev\" href=\"/([0-9]+/)\""),
	},
	"wigu-adventures": Comic{
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/adventures/comics/",
		SeqPattern: "WADV%04d.png",
		SeqStart:   1,
		SeqEnd:     1179,
	},
	"wigu-havin-fun": Comic{
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/fun/comics/",
		SeqPattern: "WOO%04d.png",
		SeqStart:   1,
		SeqEnd:     61,
	},
	"wigu-when-i-grow-up": Comic{
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/whenigrowup/comics/",
		SeqPattern: "WIGU%04d.jpg",
		SeqStart:   1,
		SeqEnd:     679,
	},
	"overcompensating": Comic{
		Archiver:   "Sequential",
		FilePrefix: "http://www.overcompensating.com/oc/comics/",
		SeqPattern: "OC%04d.png",
		SeqStart:   1,
		SeqEnd:     1543,
	},
	// This one hasn't updated in a while, but isn't "finished" yet, so we may
	// need to update the max comic ID over time.
	"iverly": Comic{
		Archiver:   "Sequential",
		FilePrefix: "http://www.iverly.com/iverly/comics/",
		SeqPattern: "IVE%04d.png",
		SeqStart:   1,
		SeqEnd:     86,
	},
}
