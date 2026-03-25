use regex::Regex;
use std::collections::HashMap;
use std::sync::LazyLock;

/// Defines the archiving strategy and parameters for a webcomic.
#[derive(Clone)]
pub struct Comic {
    pub archiver: &'static str,
    pub start_url: &'static str,
    /// Regex to find the actual start URL (for GenericCustomStart)
    pub start_match: Option<Regex>,
    /// Regex to extract the image filename from page HTML
    pub file_match: Option<Regex>,
    /// URL prefix prepended to the matched image filename
    pub file_prefix: &'static str,
    /// Regex to extract the previous-page URL from page HTML
    pub prev_link_match: Option<Regex>,
    /// printf-style format pattern for sequential filenames (e.g. "IMG%04d.png")
    pub seq_pattern: &'static str,
    pub seq_start: i32,
    pub seq_end: i32,
}

pub static COMICS: LazyLock<HashMap<&'static str, Comic>> = LazyLock::new(|| {
    let mut m = HashMap::new();

    m.insert("doa", Comic {
        archiver: "Generic",
        start_url: "https://www.dumbingofage.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"id="spliced-comic".+ src="https://www\.dumbingofage\.com/wp-content/uploads/([^"]+\.png)"#).unwrap()),
        file_prefix: "https://www.dumbingofage.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"class="previous-comic" href="(https://www\.dumbingofage\.com/[0-9a-zA-Z/-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("stuckat32", Comic {
        archiver: "Generic",
        start_url: "http://stuckat32.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "http://stuckat32.com/comics/",
        prev_link_match: Some(Regex::new(r#"class="cc-prev" rel="prev" href="(http://stuckat32\.com/[0-9a-zA-Z/-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("campcomic", Comic {
        archiver: "Generic",
        start_url: "http://campcomic.com/comic",
        start_match: None,
        file_match: Some(Regex::new(r"/katie/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "http://hw1.pa-cdn.com/camp/assets/img/katie/comics/",
        prev_link_match: Some(Regex::new(r#"class="btn btnPrev" href="(http://campcomic\.com/comic/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("gaia", Comic {
        archiver: "Generic",
        start_url: "http://www.sandraandwoo.com/gaia/",
        start_match: None,
        file_match: Some(Regex::new(r"/gaia/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.sandraandwoo.com/gaia/comics/",
        prev_link_match: Some(Regex::new(r#"href="(https?://www\.sandraandwoo\.com/gaia/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("sandraandwoo", Comic {
        archiver: "Generic",
        start_url: "https://www.sandraandwoo.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.sandraandwoo.com/comics/",
        prev_link_match: Some(Regex::new(r#"href="(https?://www\.sandraandwoo\.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("scarlet", Comic {
        archiver: "Generic",
        start_url: "https://www.sandraandwoo.com/scarlet/",
        start_match: None,
        file_match: Some(Regex::new(r#"size-full wp-image-[\d]+" src="https://www\.sandraandwoo\.com/scarlet/wp-content/uploads/([^"]+\.(jpg|png|gif))""#).unwrap()),
        file_prefix: "https://www.sandraandwoo.com/scarlet/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"previous-comic oli_hover" href="(https?://www\.sandraandwoo\.com/scarlet/comic/[0-9a-zA-Z-]+/?)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("gogetaroomie", Comic {
        archiver: "Generic",
        start_url: "https://www.gogetaroomie.com/comic/outro4",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.gogetaroomie.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" title="Previous" href="(https://www\.gogetaroomie\.com/comic/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("kiwiblitz", Comic {
        archiver: "Generic",
        start_url: "https://www.kiwiblitz.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.kiwiblitz.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://www\.kiwiblitz\.com/comic/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("sleeplessdomain", Comic {
        archiver: "Generic",
        start_url: "https://www.sleeplessdomain.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.sleeplessdomain.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://www\.sleeplessdomain\.com/comic/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("superredundant", Comic {
        archiver: "Generic",
        start_url: "http://superredundant.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"<img src="http://superredundant\.com/wp-content/uploads/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "http://superredundant.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(http://superredundant\.com/[\?0-9a-zA-Z/=-]+)" class="navi comic-nav-previous navi-prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("gunnerkrigg", Comic {
        archiver: "Generic",
        start_url: "https://www.gunnerkrigg.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"class="comic_image" src="/comics/([0-9]+\.(jpg|png|gif))""#).unwrap()),
        file_prefix: "https://www.gunnerkrigg.com/comics/",
        prev_link_match: Some(Regex::new(r#"href="(\?p=[0-9]+)"><img src="/images/prev"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("channelate", Comic {
        archiver: "Generic",
        start_url: "https://www.channelate.com/",
        start_match: None,
        file_match: Some(Regex::new(r"img src=.https://www\.channelate\.com/wp-content/uploads/(.+\.png)").unwrap()),
        file_prefix: "https://www.channelate.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(https://www\.channelate\.com/comic/[0-9a-zA-Z/-]+)" class="navi comic-nav-previous navi-prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("questionablecontent", Comic {
        archiver: "Generic",
        start_url: "https://www.questionablecontent.net/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.questionablecontent.net/comics/",
        prev_link_match: Some(Regex::new(r"href=.(view\.php\?comic=[0-9]+).>Previous").unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("iamarg", Comic {
        archiver: "Generic",
        start_url: "http://iamarg.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"/comics/([^'"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "http://iamarg.com/comics/",
        prev_link_match: Some(Regex::new(r#"href="(https?://iamarg\.com/[0-9a-zA-Z/-]+)" class="navi navi-prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("itswalky", Comic {
        archiver: "Generic",
        start_url: "https://www.itswalky.com/",
        start_match: None,
        file_match: Some(Regex::new(r"img src=.https://www\.itswalky\.com/wp-content/uploads/(.+\.png)").unwrap()),
        file_prefix: "https://www.itswalky.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(https://www\.itswalky\.com/comic/[0-9a-zA-Z/-]+)" class="comic-nav-base comic-nav-previous""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("letsspeakenglish", Comic {
        archiver: "Generic",
        start_url: "https://www.marycagle.com/letsspeakenglish/134-slow-motion",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/([0-9a-zA-Z_-]+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.marycagle.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://www\.marycagle\.com/letsspeakenglish/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("loadingartist", Comic {
        archiver: "Generic",
        start_url: "https://loadingartist.com/latest",
        start_match: None,
        file_match: Some(Regex::new(r"/uploads/([0-9]+/[0-9]+/[0-9a-zA-Z-]+\.[a-z]{3,4})").unwrap()),
        file_prefix: "https://loadingartist.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"class="normal highlight prev comic-thumb" href="(https://loadingartist\.com/comic/[0-9a-zA-Z-]+/?)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("octopuspie", Comic {
        archiver: "Generic",
        start_url: "http://www.octopuspie.com/2017-06-05/1023-1026-thats-it/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://test\.octopuspie\.com/wp-content/uploads/([^"]+\.(jpg|png|gif))" class="attachment-full size-full"#).unwrap()),
        file_prefix: "https://test.octopuspie.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"class="previous-comic" href="(http://www\.octopuspie\.com/[0-9a-zA-Z/_-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("twogag", Comic {
        archiver: "Generic",
        start_url: "http://twogag.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "http://twogag.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(http://twogag\.com/archives/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("whomp", Comic {
        archiver: "Generic",
        start_url: "https://www.whompcomic.com/",
        start_match: None,
        file_match: Some(Regex::new(r"/comics/(.+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://www.whompcomic.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://www\.whompcomic\.com/comic/[0-9a-zA-Z-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("xkcd", Comic {
        archiver: "MultiImageGeneric",
        start_url: "https://xkcd.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"//imgs\.xkcd\.com/comics/([^"]+\.png)"#).unwrap()),
        file_prefix: "http://imgs.xkcd.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="/([0-9]+/)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("wigu-adventures", Comic {
        archiver: "Sequential",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "https://www.wigucomics.com/adventures/comics/",
        prev_link_match: None,
        seq_pattern: "WADV%04d.png",
        seq_start: 1,
        seq_end: 1179,
    });

    m.insert("wigu-havin-fun", Comic {
        archiver: "Sequential",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "https://www.wigucomics.com/fun/comics/",
        prev_link_match: None,
        seq_pattern: "WOO%04d.png",
        seq_start: 1,
        seq_end: 61,
    });

    m.insert("wigu-when-i-grow-up", Comic {
        archiver: "Sequential",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "https://www.wigucomics.com/whenigrowup/comics/",
        prev_link_match: None,
        seq_pattern: "WIGU%04d.jpg",
        seq_start: 1,
        seq_end: 679,
    });

    m.insert("overcompensating", Comic {
        archiver: "Sequential",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "http://www.wigucomics.com/oc/comics/",
        prev_link_match: None,
        seq_pattern: "OC%04d.png",
        seq_start: 1,
        seq_end: 1543,
    });

    m.insert("iverly", Comic {
        archiver: "Sequential",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "http://www.wigucomics.com/iverly/comics/",
        prev_link_match: None,
        seq_pattern: "IVE%04d.png",
        seq_start: 1,
        seq_end: 86,
    });

    m.insert("alicegrove", Comic {
        archiver: "AliceGrove",
        start_url: "",
        start_match: None,
        file_match: None,
        file_prefix: "https://www.questionablecontent.net/images/alice/",
        prev_link_match: None,
        seq_pattern: "",
        seq_start: 0,
        seq_end: 205,
    });

    m.insert("beefpaper", Comic {
        archiver: "Generic",
        start_url: "http://beefpaper.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"<img src="http://beefpaper\.com/wp-content/uploads/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "http://beefpaper.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(http://beefpaper\.com/comic/[0-9a-zA-Z/_-]+)" class="navi comic-nav-previous"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("cucumberquest", Comic {
        archiver: "Generic",
        start_url: "https://cucumber.gigidigi.com/cq/page-931/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://cucumber\.gigidigi\.com/wp-content/uploads/([^"]+\.(jpg|png|gif))" class="attachment-full size-full"#).unwrap()),
        file_prefix: "https://cucumber.gigidigi.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href=['\"](https://cucumber\.gigidigi\.com/cq/[0-9a-zA-Z/_-]+)['\"] class=['"]webcomic-link webcomic1-link previous"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("treadingground", Comic {
        archiver: "Generic",
        start_url: "https://www.treadingground.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"<img src="(https://www\.treadingground\.com/comics/[^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://www.treadingground.com/comics/",
        prev_link_match: Some(Regex::new(r#"href="(https://www\.treadingground\.com/\?p=[0-9]+)" title="[^"]+" class="previous-comic"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("pbf", Comic {
        archiver: "Generic",
        start_url: "https://pbfcomics.com/",
        start_match: None,
        file_match: Some(Regex::new(r"<img src='https://pbfcomics\.com/wp-content/uploads/([^']+\.(jpg|png|gif))").unwrap()),
        file_prefix: "https://pbfcomics.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(https://pbfcomics\.com/comics/[0-9a-zA-Z/_-]+)" rel="prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("bunny", Comic {
        archiver: "Generic",
        start_url: "http://www.bunny-comic.com/",
        start_match: None,
        file_match: Some(Regex::new(r"src='strips/([^']+\.(jpg|png|gif))'").unwrap()),
        file_prefix: "http://www.bunny-comic.com/strips/",
        prev_link_match: Some(Regex::new(r#"id="strip">\s+<a href="([0-9]+\.html)"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("licd", Comic {
        archiver: "GenericCustomStart",
        start_url: "https://leasticoulddo.com/",
        start_match: Some(Regex::new(r#"href="(https://leasticoulddo\.com/comic/[0-9]+)" id="latest-comic""#).unwrap()),
        file_match: Some(Regex::new(r#"class="comic" src="https://leasticoulddo\.com/wp-content/uploads/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://leasticoulddo.com/wp-content/uploads/",
        prev_link_match: Some(Regex::new(r#"href="(https://leasticoulddo\.com/comic/[0-9]+)" rel="prev""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("asp", Comic {
        archiver: "Generic",
        start_url: "https://www.amazingsuperpowers.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"img src=.https?://www\.amazingsuperpowers\.com/comics/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://www.amazingsuperpowers.com/comics/",
        prev_link_match: Some(Regex::new(r#"href="(https?://www\.amazingsuperpowers\.com/[0-9a-zA-Z/_-]+)" class="navi navi-prev"#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("littletinythings", Comic {
        archiver: "Generic",
        start_url: "https://littletinythings.com/comic/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://littletinythings\.com/comics/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://littletinythings.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" title="Previous" href="(https://littletinythings\.com/comic/[0-9a-zA-Z/_-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("samandfuzzy", Comic {
        archiver: "Generic",
        start_url: "https://www.samandfuzzy.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://www\.samandfuzzy\.com/img/comics/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://www.samandfuzzy.com/img/comics/",
        prev_link_match: Some(Regex::new(r#"prev-page"><a href="(https://www\.samandfuzzy\.com/[0-9]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("nerfnow", Comic {
        archiver: "Generic",
        start_url: "https://www.nerfnow.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"og:image" content="https?://www\.nerfnow\.com/img/(\d+/\d+\.(jpg|png|gif))""#).unwrap()),
        file_prefix: "https://www.nerfnow.com/img/",
        prev_link_match: Some(Regex::new(r#"nav_previous"><a class="nav-link" href="/(comic/\d+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("devilscandy", Comic {
        archiver: "Generic",
        start_url: "https://www.devilscandycomic.com/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://www\.devilscandycomic\.com/comics/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://www.devilscandycomic.com/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://www\.devilscandycomic\.com/comic/[0-9a-zA-Z/_-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("floraverse", Comic {
        archiver: "Floraverse",
        start_url: "https://floraverse.com/",
        start_match: None,
        file_match: None,
        file_prefix: "",
        prev_link_match: None,
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m.insert("elephanttown", Comic {
        archiver: "Generic",
        start_url: "https://elephant.town/comic/",
        start_match: None,
        file_match: Some(Regex::new(r#"src="https://elephant\.town/comics/([^"]+\.(jpg|png|gif))"#).unwrap()),
        file_prefix: "https://elephant.town/comics/",
        prev_link_match: Some(Regex::new(r#"rel="prev" href="(https://elephant\.town/comic/[0-9a-zA-Z/_-]+)""#).unwrap()),
        seq_pattern: "", seq_start: 0, seq_end: 0,
    });

    m
});
