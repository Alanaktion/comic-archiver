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
// TODO: send error return values through channel without blocking
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
	// TODO: add Shortpacked! comic
	"doa": {
		Archiver:      "Generic",
		StartURL:      "https://www.dumbingofage.com/",
		FileMatch:     regexp.MustCompile(`id="spliced-comic".+ src="https://www.dumbingofage.com/wp-content/uploads/([^"]+\.png)`),
		FilePrefix:    "https://www.dumbingofage.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="previous-comic" href="(https://www.dumbingofage.com/[0-9a-zA-Z/-]+)"`),
	},
	"stuckat32": {
		Archiver:      "Generic",
		StartURL:      "http://stuckat32.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "http://stuckat32.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(http://stuckat32.com/[0-9a-zA-Z/-]+)"`),
	},
	"campcomic": {
		// TODO: is this site gone?
		Archiver:      "Generic",
		StartURL:      "http://campcomic.com/comic",
		FileMatch:     regexp.MustCompile(`/katie/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "http://hw1.pa-cdn.com/camp/assets/img/katie/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="btn btnPrev" href="(http://campcomic.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"fragile": {
		Archiver: "Generic",
		// StartURL:      "https://www.fragilestory.com/fragile/",
		// Original source seems to be gone, this mirror works for now.
		StartURL:      "https://fragile.webcomic.ws/",
		FileMatch:     regexp.MustCompile(`https://img.comicfury.com/comics/245/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://img.comicfury.com/comics/245/",
		PrevLinkMatch: regexp.MustCompile(`href="(/comics/[0-9]+)" rel="prev"`),
	},
	"gaia": {
		// TODO: find correct end-page URL for original run
		Archiver:      "Generic",
		StartURL:      "https://www.sandraandwoo.com/gaia/2023/03/25/some-more-great-youtube-channels/",
		FileMatch:     regexp.MustCompile(`/gaia/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.sandraandwoo.com/gaia/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.sandraandwoo.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"`),
	},
	"sandraandwoo": {
		Archiver:      "Generic",
		StartURL:      "https://www.sandraandwoo.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.sandraandwoo.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.sandraandwoo.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"`),
	},
	"scarlet": {
		Archiver:      "Generic",
		StartURL:      "https://www.sandraandwoo.com/scarlet/",
		FileMatch:     regexp.MustCompile(`size-full wp-image-[\d]+" src="https://www.sandraandwoo.com/scarlet/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://www.sandraandwoo.com/scarlet/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`previous-comic oli_hover" href="(https?://www.sandraandwoo.com/scarlet/comic/[0-9a-zA-Z-]+/?)"`),
	},
	"gogetaroomie": {
		Archiver:      "Generic",
		StartURL:      "https://www.gogetaroomie.com/comic/outro4",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.gogetaroomie.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" title="Previous" href="(https://www.gogetaroomie.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"kiwiblitz": {
		Archiver:      "Generic",
		StartURL:      "https://www.kiwiblitz.com/",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.kiwiblitz.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.kiwiblitz.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"sleeplessdomain": {
		Archiver:      "Generic",
		StartURL:      "https://www.sleeplessdomain.com/",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.sleeplessdomain.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.sleeplessdomain.com/comic/[0-9a-zA-Z-]+)"`),
	},
	"superredundant": {
		Archiver:      "Generic",
		StartURL:      "http://superredundant.com/",
		FileMatch:     regexp.MustCompile(`<img src="http://superredundant.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "http://superredundant.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(http://superredundant.com/[\?0-9a-zA-Z/=-]+)" class="navi comic-nav-previous navi-prev"`),
	},
	"gunnerkrigg": {
		Archiver:      "Generic",
		StartURL:      "https://www.gunnerkrigg.com/",
		FileMatch:     regexp.MustCompile(`class="comic_image" src="/comics/([0-9]+\.(jpg|jpeg|png|gif))"`),
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
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.questionablecontent.net/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(view.php\?comic=[0-9]+)">Previous`),
	},
	"iamarg": {
		Archiver:      "Generic",
		StartURL:      "http://iamarg.com/",
		FileMatch:     regexp.MustCompile(`/comics/([^'"]+\.(jpg|jpeg|png|gif))`),
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
	// TODO: verify all of this:
	"shortpacked": {
		Archiver:      "Generic",
		StartURL:      "https://www.shortpacked.com/comic",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.png)`),
		FilePrefix:    "https://www.shortpacked.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://www.shortpacked.com/[0-9a-zA-Z/-]+)" class="navi navi-prev"`),
	},
	"girlswithslingshots": {
		Archiver: "Generic",
		// TODO: determine ideal starting URL for this since it's complete and chaser includes direct duplicates of originals for a certain range.
		// StartURL: "https://www.girlswithslingshots.com/",
		// TODO: determine if this is the ideal chaser start position. There are direct duplicate strip images after the originals started being colorized, and I'm doing this blindly with a single dumped page I found in my filesystem.
		StartURL:      "https://www.girlswithslingshots.com/comic/gws-chaser-260",
		FileMatch:     regexp.MustCompile(`src="https://www.girlswithslingshots.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.girlswithslingshots.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.girlswithslingshots.com/comic/[0-9a-zA-Z/-]+)"`),
	},
	"letsspeakenglish": {
		Archiver:      "Generic",
		StartURL:      "https://www.marycagle.com/letsspeakenglish/134-slow-motion",
		FileMatch:     regexp.MustCompile(`/comics/([0-9a-zA-Z_-]+\.(jpg|jpeg|png|gif))`),
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
		FileMatch:     regexp.MustCompile(`src="https://test.octopuspie.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))" class="attachment-full size-full`),
		FilePrefix:    "https://test.octopuspie.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="previous-comic" href="(http://www.octopuspie.com/[0-9a-zA-Z/_-]+)"`),
	},
	"twogag": {
		Archiver:      "Generic",
		StartURL:      "http://twogag.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "http://twogag.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(http://twogag.com/archives/[0-9a-zA-Z-]+)"`),
	},
	"whomp": {
		Archiver:      "Generic",
		StartURL:      "https://www.whompcomic.com/",
		FileMatch:     regexp.MustCompile(`/comics/(.+\.(jpg|jpeg|png|gif))`),
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
		FileMatch:     regexp.MustCompile(`<img src="http://beefpaper.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "http://beefpaper.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(http://beefpaper.com/comic/[0-9a-zA-Z/_-]+)" class="navi comic-nav-previous`),
	},
	"cucumberquest": {
		Archiver:      "Generic",
		StartURL:      "https://cucumber.gigidigi.com/cq/page-931/",
		FileMatch:     regexp.MustCompile(`src="https://cucumber.gigidigi.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))" class="attachment-full size-full`),
		FilePrefix:    "https://cucumber.gigidigi.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href=['"](https://cucumber.gigidigi.com/cq/[0-9a-zA-Z/_-]+)['"] class=['"]webcomic-link webcomic1-link previous`),
	},
	// TODO: Fix pattern matching errors on some pages
	"treadingground": {
		Archiver:      "Generic",
		StartURL:      "https://og.treadingground.com/comic/251/",
		FileMatch:     regexp.MustCompile(`size-full wp-image-\d+" src="https://og.treadingground.com/wp-content/uploads/([^'"]+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://og.treadingground.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="previous-comic" href="(https?://og.treadingground.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"treadingground2": {
		Archiver:      "Generic",
		StartURL:      "https://www.treadingground.com/",
		FileMatch:     regexp.MustCompile(`(min-width: 1200px)" srcset="https://www.treadingground.com/wp-content/uploads/([^'"]+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://www.treadingground.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`class="previous-comic" href="(https?://www.treadingground.com/comic/[0-9a-zA-Z/_?=-]+)"`),
	},
	"pbf": {
		Archiver:      "Generic",
		StartURL:      "https://pbfcomics.com/",
		FileMatch:     regexp.MustCompile(`<img src='https://pbfcomics.com/wp-content/uploads/([^']+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://pbfcomics.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://pbfcomics.com/comics/[0-9a-zA-Z/_-]+)" rel="prev"`),
	},
	"bunny": {
		Archiver:      "Generic",
		StartURL:      "http://www.bunny-comic.com/",
		FileMatch:     regexp.MustCompile(`src='strips/([^']+\.(jpg|jpeg|png|gif))'`),
		FilePrefix:    "http://www.bunny-comic.com/strips/",
		PrevLinkMatch: regexp.MustCompile(`id="strip">\s+<a href="([0-9]+\.html)`),
	},
	"licd": {
		Archiver:      "GenericCustomStart",
		StartURL:      "https://leasticoulddo.com/",
		StartMatch:    regexp.MustCompile(`href="(https://leasticoulddo.com/comic/[0-9]+)" id="latest-comic"`),
		FileMatch:     regexp.MustCompile(`class="comic" src="https://leasticoulddo.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://leasticoulddo.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://leasticoulddo.com/comic/[0-9]+)" rel="prev"`),
	},
	"asp": {
		Archiver:      "Generic",
		StartURL:      "https://www.amazingsuperpowers.com/",
		FileMatch:     regexp.MustCompile(`img src="https?://www.amazingsuperpowers.com/comics/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.amazingsuperpowers.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`href="(https?://www.amazingsuperpowers.com/[0-9a-zA-Z/_-]+)" class="navi navi-prev`),
	},
	"littletinythings": {
		Archiver:      "Generic",
		StartURL:      "https://littletinythings.com/comic/",
		FileMatch:     regexp.MustCompile(`src="https://littletinythings.com/comics/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://littletinythings.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" title="Previous" href="(https://littletinythings.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"samandfuzzy": {
		Archiver:      "Generic",
		StartURL:      "https://www.samandfuzzy.com/",
		FileMatch:     regexp.MustCompile(`src="https://www.samandfuzzy.com/img/comics/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.samandfuzzy.com/img/comics/",
		PrevLinkMatch: regexp.MustCompile(`prev-page"><a href="(https://www.samandfuzzy.com/[0-9]+)"`),
	},
	"nerfnow": {
		Archiver:      "Generic",
		StartURL:      "https://www.nerfnow.com/",
		FileMatch:     regexp.MustCompile(`og:image" content="https?://www.nerfnow.com/img/(\d+/\d+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://www.nerfnow.com/img/",
		PrevLinkMatch: regexp.MustCompile(`nav_previous"><a class="nav-link" href="/(comic/\d+)"`),
	},
	"devilscandy": {
		Archiver:      "Generic",
		StartURL:      "https://www.devilscandycomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www.devilscandycomic.com/comics/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://www.devilscandycomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://www.devilscandycomic.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	// TODO: needs support for filenames derived from page content, otherwise each chapter will overwrite previous pages.
	"drugsandwires": {
		Archiver:      "Generic",
		StartURL:      "https://www.drugsandwires.fail/",
		FileMatch:     regexp.MustCompile(`src="/assets/img/pages/([^"]+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://www.drugsandwires.fail/assets/img/pages/",
		PrevLinkMatch: regexp.MustCompile(`href="(/comic/[0-9a-zA-Z/_-]+/)" class="comic__navigation__option" id="prev"`),
	},
	"floraverse": {
		Archiver: "Floraverse",
		StartURL: "https://floraverse.com/",
	},
	"elephanttown": {
		Archiver:      "Generic",
		StartURL:      "https://elephant.town/comic/",
		FileMatch:     regexp.MustCompile(`src="https://elephant.town/comics/([^"]+\.(jpg|jpeg|png|gif))`),
		FilePrefix:    "https://elephant.town/comics/",
		PrevLinkMatch: regexp.MustCompile(`rel="prev" href="(https://elephant.town/comic/[0-9a-zA-Z/_-]+)"`),
	},

	// Hiveworks hiatus comics
	"anacrine-complex": {
		Archiver:      "Generic",
		StartURL:      "https://www.pigeoncomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.pigeoncomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.pigeoncomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.pigeoncomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"anarchy-dreamers": {
		Archiver:      "Generic",
		StartURL:      "https://www.anarchydreamers.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.anarchydreamers\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.anarchydreamers.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.anarchydreamers\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"aquapunk": {
		Archiver:      "Generic",
		StartURL:      "https://www.aquapunk.co/",
		FileMatch:     regexp.MustCompile(`src="https://www\.aquapunk\.co/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.aquapunk.co/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.aquapunk\.co/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"astral-aves": {
		Archiver:      "Generic",
		StartURL:      "https://www.astralaves.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.astralaves\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.astralaves.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.astralaves\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"blindsprings": {
		Archiver:      "Generic",
		StartURL:      "https://www.blindsprings.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.blindsprings\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.blindsprings.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.blindsprings\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"clockwork": {
		Archiver:      "Generic",
		StartURL:      "https://www.clockwork-comic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.clockwork-comic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.clockwork-comic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.clockwork-comic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"demo-street": {
		Archiver:      "Generic",
		StartURL:      "https://www.demonstreet.co/",
		FileMatch:     regexp.MustCompile(`src="https://www\.demonstreet\.co/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.demonstreet.co/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.demonstreet\.co/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"edison-rex": {
		Archiver:      "Generic",
		StartURL:      "https://www.edisonrex.net/",
		FileMatch:     regexp.MustCompile(`src="https://www\.edisonrex\.net/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.edisonrex.net/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.edisonrex\.net/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"fireweed-moors": {
		Archiver:      "Generic",
		StartURL:      "https://www.fireweedmoors.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.fireweedmoors\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.fireweedmoors.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.fireweedmoors\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"goodbye-to-halos": {
		Archiver:      "Generic",
		StartURL:      "https://www.goodbyetohalos.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.goodbyetohalos\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.goodbyetohalos.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.goodbyetohalos\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"heart-of-gold": {
		// Home page shows archive; start from the final comic page.
		// Site template generates malformed "https:/..." URLs (missing //domain).
		Archiver:      "Generic",
		StartURL:      "https://www.heartofgoldcomic.com/comic/ii-196",
		FileMatch:     regexp.MustCompile(`src="https:/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.heartofgoldcomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" title="Previous" href="(https:/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"helvetica": {
		Archiver:      "Generic",
		StartURL:      "https://helvetica.jnwiedle.com/",
		FileMatch:     regexp.MustCompile(`id="comic">[^<]*(?:<a[^>]*>)?[^<]*<img src="https://helvetica\.jnwiedle\.com/wp-content/uploads/([^"]+\.(png|jpg|gif))"`),
		FilePrefix:    "https://helvetica.jnwiedle.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://helvetica\.jnwiedle\.com/comic/[^"]+)" class="comic-nav-base comic-nav-previous"`),
	},
	"hemlock": {
		Archiver:      "Generic",
		StartURL:      "https://www.hemlockcomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.hemlockcomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.hemlockcomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.hemlockcomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"killjoys": {
		Archiver:      "Generic",
		StartURL:      "https://www.killjoyscomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.killjoyscomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.killjoyscomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.killjoyscomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"knights-errant": {
		// Hosted at theyoungdoyler.com (no www)
		Archiver:      "Generic",
		StartURL:      "https://theyoungdoyler.com/",
		FileMatch:     regexp.MustCompile(`src="https://theyoungdoyler\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://theyoungdoyler.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://theyoungdoyler\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"last-diplomat": {
		Archiver:      "Generic",
		StartURL:      "https://www.thelastdiplomat.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.thelastdiplomat\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.thelastdiplomat.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.thelastdiplomat\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"lonely-vincent": {
		Archiver:      "Generic",
		StartURL:      "https://www.lonelyvincent.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.lonelyvincent\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.lonelyvincent.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.lonelyvincent\.com/lonelyvincent/[0-9a-zA-Z/_-]+)"`),
	},
	"manly-guys": {
		// Manly Guys Doing Manly Things (WordPress/ComicEasel)
		Archiver:      "Generic",
		StartURL:      "https://thepunchlineismachismo.com/",
		FileMatch:     regexp.MustCompile(`id="comic">[^<]*(?:<a[^>]*>)?[^<]*<img src="https://thepunchlineismachismo\.com/wp-content/uploads/([^"]+\.(jpg|jpeg|png|gif))"`),
		FilePrefix:    "https://thepunchlineismachismo.com/wp-content/uploads/",
		PrevLinkMatch: regexp.MustCompile(`href="(https://thepunchlineismachismo\.com/archives/comic/[^"]+)" class="comic-nav-base comic-nav-previous"`),
	},
	"monsterkind": {
		// Hosted at monsterkind.enenkay.com (http only, no www)
		Archiver:      "Generic",
		StartURL:      "http://monsterkind.enenkay.com/",
		FileMatch:     regexp.MustCompile(`src="http://monsterkind\.enenkay\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "http://monsterkind.enenkay.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(http://monsterkind\.enenkay\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"monsters-garden": {
		Archiver:      "Generic",
		StartURL:      "https://www.monstersgarden.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.monstersgarden\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.monstersgarden.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.monstersgarden\.com/monsters-garden/[0-9a-zA-Z/_-]+)"`),
	},
	"never-satisfied": {
		Archiver:      "Generic",
		StartURL:      "https://www.neversatisfiedcomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.neversatisfiedcomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.neversatisfiedcomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.neversatisfiedcomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"not-drunk-enough": {
		Archiver:      "Generic",
		StartURL:      "https://www.ndecomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.ndecomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.ndecomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.ndecomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"patrik-the-vampire": {
		Archiver:      "Generic",
		StartURL:      "https://www.patrikthevampire.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.patrikthevampire\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.patrikthevampire.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.patrikthevampire\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"peritale": {
		Archiver:      "Generic",
		StartURL:      "https://www.peritale.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.peritale\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.peritale.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.peritale\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"quicksilver": {
		Archiver:      "Generic",
		StartURL:      "https://www.quicksilvercomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.quicksilvercomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.quicksilvercomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.quicksilvercomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"real-science-adventures": {
		// Hosted on atomic-robo.com under /rsa/ path
		Archiver:      "Generic",
		StartURL:      "https://www.atomic-robo.com/rsa",
		FileMatch:     regexp.MustCompile(`src="https://www\.atomic-robo\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.atomic-robo.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.atomic-robo\.com/rsa/[0-9a-zA-Z/_-]+)"`),
	},
	"saint-forrent": {
		Archiver:      "Generic",
		StartURL:      "https://www.saintforrent.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.saintforrent\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.saintforrent.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.saintforrent\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"sakana": {
		Archiver:      "Generic",
		StartURL:      "https://www.sakana-comic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.sakana-comic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.sakana-comic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev"[^>]* href="(https://www\.sakana-comic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"scape": {
		Archiver:      "Generic",
		StartURL:      "https://www.scapecomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.scapecomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.scapecomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.scapecomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"stutterhug": {
		Archiver:      "Generic",
		StartURL:      "https://www.stutterhug.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.stutterhug\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.stutterhug.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.stutterhug\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"the-substitutes": {
		Archiver:      "Generic",
		StartURL:      "https://www.thesubstitutescomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.thesubstitutescomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.thesubstitutescomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.thesubstitutescomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"tove": {
		Archiver:      "Generic",
		StartURL:      "https://www.tovecomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.tovecomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.tovecomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.tovecomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"undivine": {
		Archiver:      "Generic",
		StartURL:      "https://www.undivinecomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.undivinecomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.undivinecomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.undivinecomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"within": {
		Archiver:      "Generic",
		StartURL:      "https://www.withincomic.net/",
		FileMatch:     regexp.MustCompile(`src="https://www\.withincomic\.net/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.withincomic.net/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.withincomic\.net/comic/[0-9a-zA-Z/_-]+)"`),
	},
	"witchy": {
		Archiver:      "Generic",
		StartURL:      "https://www.witchycomic.com/",
		FileMatch:     regexp.MustCompile(`src="https://www\.witchycomic\.com/comics/([^"]+\.(jpg|jpeg|png|gif))" id="cc-comic"`),
		FilePrefix:    "https://www.witchycomic.com/comics/",
		PrevLinkMatch: regexp.MustCompile(`class="cc-prev" rel="prev" href="(https://www\.witchycomic\.com/comic/[0-9a-zA-Z/_-]+)"`),
	},
}
