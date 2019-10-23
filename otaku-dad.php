<?php
// Otaku Dad
// This one has ended, so we can hard-code a few things
// We'll work from start to end, unlike most of the other archives

// We force the file names to be sequential. Makes it much easier to read.

$start = 'http://www.otaku-dad.com/comic/i39m-embarrassed-by-my-otaku-dad';
$end = 'http://www.otaku-dad.com/comic/the-end';

if (!is_dir('otaku-dad')) {
    mkdir('otaku-dad');
}

$url = $start;
$i = 0;
while ($url != $end) {
    $html = file_get_contents($url);
    preg_match('@src="http://www.otaku-dad.com/comics/([0-9a-zA-Z-]+\\.[a-z]{3,4})" id="cc-comic"@', $html, $matches);
    if (!empty($matches[1])) {
        $name = trim(substr($url, 31), '/');
        if (glob("otaku-dad/{$name}*")) {
            return;
        }

        echo "Downloading {$name}\n";
        $data = @file_get_contents('http://www.otaku-dad.com/comics/' . $matches[1]);
        if ($data) {
            $ext = pathinfo(parse_url($matches[1])['path'], PATHINFO_EXTENSION);
            file_put_contents("otaku-dad/" . str_pad($i, 2, '0', STR_PAD_LEFT) . "-$name.$ext", $data);
        }
    }

    preg_match('@href="(http://www.otaku-dad.com/comic/[0-9a-zA-Z-]+/?)" class="next"@', $html, $matches);
    $url = $matches[1];
    $i++;

    usleep(5e5);
}
