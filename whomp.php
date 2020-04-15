<?php
// Whomp

// Another Hiveworks comic, pretty straightforward.

$html = file_get_contents('https://www.whompcomic.com/');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('whomp')) {
    mkdir('whomp');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('whomp/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "https://www.whompcomic.com/comics/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("whomp/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@rel="prev" href="(https://www.whompcomic.com/comic/[0-9a-zA-Z-]+)"@';
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
