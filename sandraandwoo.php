<?php
// Sandra and Woo
// Another ComicPress one, still updating regularly

// Note that you can get official downloads of this comic, including
// high-resolution and draft versions by supporting the creators on Patreon!
// https://www.patreon.com/sandraandwoo

$html = file_get_contents('http://www.sandraandwoo.com/');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('sandraandwoo')) {
    mkdir('sandraandwoo');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('sandraandwoo/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "http://www.sandraandwoo.com/comics/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("sandraandwoo/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://www.sandraandwoo.com/[0-9]{4}/[0-9]+/[0-9]+/[0-9a-zA-Z-]+/?)" rel="prev"@';
    preg_match($regex, $html, $prevMatch);

    if (empty($prevMatch[1])) {
        echo "No previous URL found!\n";
        return;
    }

    $html = @file_get_contents($prevMatch[1]);
    if (!$html) {
        echo "Failed to load previous page!\n";
        return;
    }

    preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
