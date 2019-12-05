<?php
// It's Walky
// Based on the Dumbing of Age script

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://www.itswalky.com/');
preg_match('@img src="http://www.itswalky.com/wp-content/uploads/(.+\\.png)@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('itswalky')) {
    mkdir('itswalky');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    $filename = basename($matches[1]);
    if (is_file('itswalky/' . $filename)) {
        return;
    }

    echo "Downloading {$filename}\n";
    $url = "http://www.itswalky.com/wp-content/uploads/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("itswalky/{$filename}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://www.itswalky.com/comic/[0-9a-zA-Z/-]+)" class="comic-nav-base comic-nav-previous"@';
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

    preg_match('@img src="http://www.itswalky.com/wp-content/uploads/(.+\\.png)@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
