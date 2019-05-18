<?php
// Blaster Nation
// This one has ended, so we can hard-code a few things
// We'll work from start to end, unlike most of the other archives

// I chose to use the main comic page URLs instead of the comic filenames here
// since they're actually sequential that way. Makes it much easier to read.

$start = 'http://www.blasternation.com/comic/1-one-day-in-the-life-of-matthew-palmer';
$end = 'http://www.blasternation.com/comic/516-end';

if (!is_dir('blasternation')) {
    mkdir('blasternation');
}

$url = $start;
while ($url != $end) {
    $html = file_get_contents($url);
    preg_match('@src="http://www.blasternation.com/comics/([0-9a-zA-Z-]+\\.[a-z]{3,4})" id="cc-comic"@', $html, $matches);
    if (!empty($matches[1])) {
        $name = trim(substr($url, 35), '/');
        if (glob("blasternation/{$name}*")) {
            return;
        }

        echo "Downloading {$name}\n";
        $data = @file_get_contents('http://www.blasternation.com/comics/' . $matches[1]);
        if ($data) {
            $ext = pathinfo(parse_url($matches[1])['path'], PATHINFO_EXTENSION);
            file_put_contents('blasternation/' . $name . '.' . $ext, $data);
        }
    }

    preg_match('@rel="next" href="(http://www.blasternation.com/comic/[0-9a-zA-Z-]+/?)"@', $html, $matches);
    $url = $matches[1];

    usleep(500000);
}
