<?php
// Least I Could Do
// A bit tricky, but not too bad

// Start with home page and navigate to latest comic
$home = file_get_contents('http://www.leasticoulddo.com/');
$regex = '@href="(http://www.leasticoulddo.com/comic/[0-9]+/?)" id="feature-comic"@';
preg_match($regex, $home, $homeMatch);
if (!empty($homeMatch[1])) {
    $html = file_get_contents($homeMatch[1]);
    preg_match('@/uploads/([0-9]+/[0-9]+)/([0-9]+\\.[a-z]{3,4})@', $html, $matches);
} else {
    echo "Unable to find link to latest comic! :(\n";
    return;
}

if (empty($matches[2])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('licd')) {
    mkdir('licd');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('licd/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[2]}\n";
    $url = "http://www.leasticoulddo.com/wp-content/uploads/{$matches[1]}/{$matches[2]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("licd/{$matches[2]}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://www.leasticoulddo.com/comic/[0-9]+/?)" id="nav-large-prev"@';
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

    preg_match('@/uploads/([0-9]+/[0-9]+)/([0-9]+\\.[a-z]{3,4})@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(500000);
}
