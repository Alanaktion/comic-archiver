<?php
// Camp Weedonwantcha

// A fairly simple ComicPress strip, currently on extended hiatus, but
// potentially updating again. We'll start with the current strip and work back
// to the beginning.

$html = file_get_contents('http://campcomic.com/comic');
$path = 'http://hw1.pa-cdn.com/camp/assets/img/katie/comics/';
preg_match('@/katie/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('campcomic')) {
    mkdir('campcomic');
}

while (true) {
    if (is_file('campcomic/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $data = @file_get_contents($path . $matches[1]);
    if ($data) {
        file_put_contents("campcomic/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@class="btn btnPrev" href="(http://campcomic.com/comic/[0-9a-zA-Z-]+)"@';
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

    preg_match('@/katie/comics/(.+\\.(jpg|png|gif))@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
