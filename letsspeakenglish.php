<?php
// Kiwi Blitz
// Mary Cagle's finished comic, hosted on her main personal domain.
// We'll start from the last actual strip, but you should *TOTALLY* buy the
// book as long as it's still available, it's great!
// https://hivemill.com/products/lets-speak-english

$html = file_get_contents('https://www.marycagle.com/letsspeakenglish/134-slow-motion');
preg_match('@/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('letsspeakenglish')) {
    mkdir('letsspeakenglish');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('letsspeakenglish/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "https://www.marycagle.com/comics/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("letsspeakenglish/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@rel="prev" href="(https://www.marycagle.com/letsspeakenglish/[0-9a-zA-Z-]+)"@';
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

    preg_match('@/comics/([0-9a-zA-Z_-]+\\.(jpg|png|gif))@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
