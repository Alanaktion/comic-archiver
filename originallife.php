<?php
// Original Life
// By Jay Naylor, probably complete, or at least abandoned. Custom site.

$html = file_get_contents('http://jaynaylor.com/originallife/');
preg_match('@/originallife/comic/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('originallife')) {
    mkdir('originallife');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    echo "Downloading {$matches[1]}\n";
    $url = "http://jaynaylor.com/originallife/comic/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("originallife/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://jaynaylor.com/originallife/archives/[0-9A-Za-z/_-]+\.html)">&laquo; Previous@';
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

    preg_match('@/originallife/comic/(.+\\.(jpg|png|gif))@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
