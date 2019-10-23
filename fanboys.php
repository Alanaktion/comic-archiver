<?php
// Fanboys Online
// A bit tricky, but not too bad

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://www.fanboys-online.com/');
preg_match('@/comics/([^"]+\\.[a-z]{3,4})@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('fanboys')) {
    mkdir('fanboys');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('fanboys/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "http://www.fanboys-online.com/comics/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("fanboys/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@href="(/index\\.php\\?id=[0-9]+)" class="prev"@';
    preg_match($regex, $html, $prevMatch);

    if (empty($prevMatch[1])) {
        echo "No previous URL found!\n";
        return;
    }

    $html = @file_get_contents('http://www.fanboys-online.com' . $prevMatch[1]);
    if (!$html) {
        echo "Failed to load previous page!\n";
        return;
    }

    preg_match('@/comics/([^"]+\\.[a-z]{3,4})@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found! Exiting.\n";
        return;
    }

    usleep(5e5);
}
