<?php
// Alice Grove

// This one is finished, and Jeph re-hosted it as a single sequential list of
// files, so we're no longer reliant on scraping the Tumblr blog!

// This is basically the perfect scenario for super simple archiving, thanks so
// much Jeph! Not that you're reading this, but can we please get new QC books?

if (!is_dir('alicegrove')) {
    mkdir('alicegrove');
}

for ($i = 1; $i <= 205; $i++) {
    $path = "alicegrove/$i.png";
    if (!is_file($path)) {
        $url = "https://www.questionablecontent.net/images/alice/$i.png";
        echo "Downloading $i.png\n";
        file_put_contents($path, file_get_contents($url));
        usleep(5e5);
    }
}
