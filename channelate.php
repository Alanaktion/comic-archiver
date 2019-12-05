<?php
// Channelate
// A fairly simple WordPress site

// Start by navigating to the latest comic
$html = file_get_contents('https://www.channelate.com/');
preg_match('@img src="https://www.channelate.com/wp-content/uploads/(.+\\.png)@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('channelate')) {
    mkdir('channelate');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    $filename = basename($matches[1]);
    if (is_file('channelate/' . $filename)) {
        return;
    }

    echo "Downloading {$filename}\n";
    $url = "https://www.channelate.com/wp-content/uploads/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("channelate/{$filename}", $data);
    }

    // Find previous page link
    $regex = '@href="(https://www.channelate.com/comic/[0-9a-zA-Z/-]+)" class="navi comic-nav-previous navi-prev"@';
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

    preg_match('@img src="https://www.channelate.com/wp-content/uploads/(.+\\.png)@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
