package archivers

import (
	"fmt"
	"log"
	"regexp"
	"sync"
)

var protocolMatch = regexp.MustCompile("^https?:")
var basenameMatch = regexp.MustCompile(`\/?([^\/]+\.[A-Za-z]{3,4})$`)

// Archive a comic
func Archive(dir string, comic Comic, skipExisting bool, wg *sync.WaitGroup) {
	prefix := fmt.Sprintf("[%s] ", dir)
	logger := log.New(log.Writer(), prefix, log.Flags())
	logger.Println("Starting archive")

	if comic.Archiver == "Generic" {
		Generic(comic.StartURL, dir, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch, skipExisting, logger)
	}
	if comic.Archiver == "GenericCustomStart" {
		GenericCustomStart(comic.StartURL, comic.StartMatch, dir, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch, skipExisting, logger)
	}
	if comic.Archiver == "MultiImageGeneric" {
		MultiImageGeneric(comic.StartURL, dir, comic.FileMatch, comic.FilePrefix, comic.PrevLinkMatch, skipExisting, logger)
	}
	if comic.Archiver == "Sequential" {
		Sequential(dir, comic.FilePrefix, comic.SeqPattern, comic.SeqStart, comic.SeqEnd, skipExisting, logger)
	}

	// Custom archivers
	if comic.Archiver == "AliceGrove" {
		AliceGrove(dir, comic.FilePrefix, comic.SeqEnd, skipExisting, logger)
	}
	if comic.Archiver == "Floraverse" {
		Floraverse(comic.StartURL, dir, skipExisting, logger)
	}

	wg.Done()
}

// Comic base type for all archivers
type Comic struct {
	Archiver      string
	StartURL      string
	StartMatch    *regexp.Regexp
	FileMatch     *regexp.Regexp
	FilePrefix    string
	PrevLinkMatch *regexp.Regexp
	SeqPattern    string
	SeqStart      int
	SeqEnd        int
}

// Comics supported by the archiver
var Comics = map[string]Comic{
	"doa": {
		Archiver:      "Generic",
		StartURL:      "https://www.dumbingofage.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.png)`),
		FilePrefix:    "https://www.dumbingofage.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://www.dumbingofage.com/[0-9a-zA-Z/-]+)" class="navi navi-prev"`),
	},
	"stuckat32": {
		Archiver:      "Generic",
		StartURL:      "http://stuckat32.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "http://stuckat32.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(http://stuckat32.com/[0-9a-zA-Z/-]+)"`),
	},
	"campcomic": {
		Archiver:      "Generic",
		StartURL:      "http://campcomic.com/comic",
		FileMatch:     regexp.MustCompile(`/katie/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "http://hw1.pa-cdn.com/camp/assets/img/katie/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="btn btnPrev" href="(http://campcomic.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"gaia": {
		Archiver:      "Generic",
		StartURL:      "http://www.sandraandwoo.com/gaia/",
		FileMatch:     regexp.MustCompile(`/gaia/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.sandraandwoo.com/gaia/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.sandraandwoo.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"`),
	},
	"sandraandwoo": {
		Archiver:      "Generic",
		StartURL:      "https://www.sandraandwoo.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.sandraandwoo.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.sandraandwoo.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"`),
},
	"scarlet": {
		Archiver:      "Generic",
		StartURL:      "https://www.sandraandwoo.com/scarlet/",
		FileMatch:     regexp.MustCompile(`size-full wp-image-[\d]+" src="https://www.sandraandwoo.com/scarlet/wp-content/uploads/([^"]+\.(jpg|png|gif))"`),
		FilePrefix:    "https://www.sandraandwoo.com/scarlet/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`previous-comic oli_hover" href="(https?://www.sandraandwoo.com/scarlet/comic/[0-9a-zA-Z-]+/?)"`),
	},
	"gogetaroomie": {
		Archiver:      "Generic",
		StartURL:      "https://www.gogetaroomie.com/comic/outro4",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.gogetaroomie.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" title="Previous" href="(https://www.gogetaroomie.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"kiwiblitz": {
		Archiver:      "Generic",
		StartURL:      "https://www.kiwiblitz.com/",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.kiwiblitz.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.kiwiblitz.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"sleeplessdomain": {
		Archiver:      "Generic",
		StartURL:      "https://www.sleeplessdomain.com/",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.sleeplessdomain.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.sleeplessdomain.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"superredundant": {
		Archiver:      "Generic",
		StartURL:      "http://superredundant.com/",
		FileMatch:     regexp.MustCompile(`<img src="http://superredundant.com/wp-content/uploads/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "http://superredundant.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(http://superredundant.com/[\?0-9a-zA-Z/=-]+)" class="navi comic-nav-previous navi-prev"`),
	},
	"gunnerkrigg": {
		Archiver:      "Generic",
		StartURL:      "https://www.gunnerkrigg.com/",
		FileMatch:     regexp.MustCompile(`class="comic_image" src="/comics/([0-9]+\.(jpg|png|gif))"`),
		FilePrefix:    "https://www.gunnerkrigg.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(\?p=[0-9]+)"><img src="/images/prev`),
	},
	"channelate": {
		Archiver:      "Generic",
		StartURL:      "https://www.channelate.com/",
		FileMatch:     regexp.MustCompile(`img src="https://www.channelate.com/wp-content/uploads/(.+\.png)`),
		FilePrefix:    "https://www.channelate.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://www.channelate.com/comic/[0-9a-zA-Z/-]+)" class="navi comic-nav-previous navi-prev"`),
	},
	"questionablecontent": {
		Archiver:      "Generic",
		StartURL:      "https://www.questionablecontent.net/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.questionablecontent.net/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(view.php\?comic=[0-9]+)">Previous`),
	},
	"iamarg": {
		Archiver:      "Generic",
		StartURL:      "http://iamarg.com/",
		FileMatch:     regexp.MustCompile(`/comics/([^'"]+\.(jpg|png|gif))`),
		FilePrefix:    "http://iamarg.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://iamarg.com/[0-9a-zA-Z/-]+)" class="navi navi-prev"`),
	},
	"itswalky": {
		Archiver:      "Generic",
		StartURL:      "https://www.itswalky.com/",
		FileMatch:     regexp.MustCompile(`img src="https://www.itswalky.com/wp-content/uploads/(.+\.png)`),
		FilePrefix:    "https://www.itswalky.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://www.itswalky.com/comic/[0-9a-zA-Z/-]+)" class="comic-nav-base comic-nav-previous"`),
	},
	"letsspeakenglish": {
		Archiver:      "Generic",
		StartURL:      "https://www.marycagle.com/letsspeakenglish/134-slow-motion",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.marycagle.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.marycagle.com/letsspeakenglish/[0-9a-zA-Z-]+)"`),
	},
	"loadingartist": {
		Archiver:      "Generic",
		StartURL:      "https://loadingartist.com/latest",
		FileMatch:     regexp.MustCompile(`/uploads/([0-9]+/[0-9]+/[0-9a-zA-Z-]+\.[a-z]{3,4})`),
		FilePrefix:    "https://loadingartist.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="normal highlight prev comic-thumb" href="(https://loadingartist.com/comic/[0-9a-zA-Z-]+/?)"`),
	},
	"octopuspie": {
		Archiver:      "Generic",
		StartURL:      "http://www.octopuspie.com/2017-06-05/1023-1026-thats-it/",
		FileMatch:     regexp.MustCompile(`src="https://test.octopuspie.com/wp-content/uploads/([^"]+\.(jpg|png|gif))" class="attachment-full size-full`),
		FilePrefix:    "https://test.octopuspie.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="previous-comic" href="(http://www.octopuspie.com/[0-9a-zA-Z/_-]+)"`),
	},
	"twogag": {
		Archiver:      "Generic",
		StartURL:      "http://twogag.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "http://twogag.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(http://twogag.com/archives/[0-9a-zA-Z-]+)"`),
	},
	"whomp": {
		Archiver:      "Generic",
		StartURL:      "https://www.whompcomic.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.whompcomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.whompcomic.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"xkcd": {
		Archiver:      "MultiImageGeneric",
		StartURL:      "https://xkcd.com/",
		FileMatch:     regexp.MustCompile(`//imgs.xkcd.com/comics/([^"]+\.png)`),
		FilePrefix:    "http://imgs.xkcd.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="/([0-9]+/)"`),
	},
	"wigu-adventures": {
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/adventures/comics/",
		SeqPattern: "WADV%04d.png",
		SeqStart:   1,
		SeqEnd:     1179,
	},
	"wigu-havin-fun": {
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/fun/comics/",
		SeqPattern: "WOO%04d.png",
		SeqStart:   1,
		SeqEnd:     61,
	},
	"wigu-when-i-grow-up": {
		Archiver:   "Sequential",
		FilePrefix: "https://www.wigucomics.com/whenigrowup/comics/",
		SeqPattern: "WIGU%04d.jpg",
		SeqStart:   1,
		SeqEnd:     679,
	},
	"overcompensating": {
		Archiver:   "Sequential",
		FilePrefix: "http://www.wigucomics.com/oc/comics/",
		SeqPattern: "OC%04d.png",
		SeqStart:   1,
		SeqEnd:     1543,
	},
		"iverly": {
		Archiver:   "Sequential",
		FilePrefix: "http://www.wigucomics.com/iverly/comics/",
		SeqPattern: "IVE%04d.png",
		SeqStart:   1,
		SeqEnd:     86,
	},
	"alicegrove": {
		Archiver:   "AliceGrove",
		FilePrefix: "https://www.questionablecontent.net/images/alice/",
		SeqEnd:     205,
	},
	"beefpaper": {
		Archiver:      "Generic",
		StartURL:      "http://beefpaper.com/",
		FileMatch:     regexp.MustCompile(`<img src="http://beefpaper.com/wp-content/uploads/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "http://beefpaper.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(http://beefpaper.com/comic/[0-9a-zA-Z/_-]+)" class="navi comic-nav-previous`),
	},
	"cucumberquest": {
		Archiver:      "Generic",
		StartURL:      "https://cucumber.gigidigi.com/cq/page-931/",
		FileMatch:     regexp.MustCompile(`src="https://cucumber.gigidigi.com/wp-content/uploads/([^"]+\.(jpg|png|gif))" class="attachment-full size-full`),
		FilePrefix:    "https://cucumber.gigidigi.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href=['"](https://cucumber.gigidigi.com/cq/[0-9a-zA-Z/_-]+)['"] class=['"]webcomic-link webcomic1-link previous`),
	},
	"treadingground": {
		Archiver:      "Generic",
		StartURL:      "https://www.treadingground.com/",
		FileMatch:     regexp.MustCompile(`<img src="(https://www.treadingground.com/comics/[^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.treadingground.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://www.treadingground.com/\?p=[0-9]+)" title="[^"]+" class="previous-comic`),
	},
	"pbf": {
		Archiver:      "Generic",
		StartURL:      "https://pbfcomics.com/",
		FileMatch:     regexp.MustCompile(`<img src='https://pbfcomics.com/wp-content/uploads/([^']+\.(jpg|png|gif))`),
		FilePrefix:    "https://pbfcomics.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://pbfcomics.com/comics/[0-9a-zA-Z/_-]+)" rel="prev"`),
	},
	"bunny": {
		Archiver:      "Generic",
		StartURL:      "http://www.bunny-comic.com/",
		FileMatch:     regexp.MustCompile(`src='strips/([^']+\.(jpg|png|gif))'`),
		FilePrefix:    "http://www.bunny-comic.com/strips/",
		PrevLinkMatch: regexp.MustCompile(`id="strip">\s+<a href="([0-9]+\.html)`),
	},
	"licd": {
		Archiver:      "GenericCustomStart",
		StartURL:      "https://leasticoulddo.com/",
		StartMatch:    regexp.MustCompile(`href="(https://leasticoulddo.com/comic/[0-9]+)" id="latest-comic"`),
		FileMatch:     regexp.MustCompile(`class="comic" src="https://leasticoulddo.com/wp-content/uploads/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://leasticoulddo.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://leasticoulddo.com/comic/[0-9]+)" rel="prev"`),
	},
	"asp": {
		Archiver:      "Generic",
		StartURL:      "https://www.amazingsuperpowers.com/",
		FileMatch:     regexp.MustCompile(`img src="https?://www.amazingsuperpowers.com/comics/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.amazingsuperpowers.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.amazingsuperpowers.com/[0-9a-zA-Z/_-]+)" class="navi navi-prev`),
	},
	"littletinythings": {
		Archiver:      "Generic",
		StartURL:      "https://littletinythings.com/comic/",
		FileMatch:     regexp.MustCompile(`src="https://littletinythings.com/comics/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://littletinythings.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" title="Previous" href="(https://littletinythings.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"samandfuzzy": {
		Archiver:      "Generic",
		StartURL:      "https://www.samandfuzzy.com/",
		FileMatch:     regexp.MustCompile(`src="https://www.samandfuzzy.com/img/comics/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.samandfuzzy.com/img/comics/",
		PrevLinkMatch: regexp.MustCompile(`prev-page"><a href="(https://www.samandfuzzy.com/[0-9]+)"`),
	},
	"nerfnow": {
		Archiver:      "Generic",
		StartURL:      "https://www.nerfnow.com/",
		FileMatch:     regexp.MustCompile(`og:image" content="https?://www.nerfnow.com/img/(\d+/\d+\.(jpg|png|gif))"`),
		FilePrefix:    "https://www.nerfnow.com/img/",
		PrevLinkMatch: regexp.MustCompile(`nav_previous"><a class="nav-link" href="/(comic/\d+)"`),
	},
	"devilscandy": {
		Archiver:      "Generic",
		StartURL:      "https://www.devilscandycomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www.devilscandycomic.com/comics/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://www.devilscandycomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.devilscandycomic.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"floraverse": {
		Archiver: "Floraverse",
		StartURL: "https://floraverse.com/",
	},
	"elephanttown": {
		Archiver:      "Generic",
		StartURL:      "https://elephant.town/comic/",
		FileMatch:     regexp.MustCompile(`src="https://elephant.town/comics/([^"]+\.(jpg|png|gif))`),
		FilePrefix:    "https://elephant.town/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://elephant.town/comic/[0-9a-zA-Z/_-]+)"`),
	},
}
