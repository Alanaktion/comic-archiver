<?php
// Dumbing of Age
// A bit tricky, but not too bad

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://www.dumbingofage.com/');
preg_match('@/comics/(.+\\.png)@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('doa')) {
    mkdir('doa');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('doa/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "http://www.dumbingofage.com/comics/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("doa/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://www.dumbingofage.com/[0-9a-zA-Z/-]+)" class="navi navi-prev"@';
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

    preg_match('@/comics/(.+\\.png)@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
