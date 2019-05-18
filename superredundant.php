<?php
// League of Super Redundant Heroes
// This one is also running ComicPress, but has a slightly different setup.
// Images are stored directly in /wp-content/uploads/, but so are other things!

// This one can take ~1 hour to run as of early 2019.

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://superredundant.com/');
preg_match('@<img src="http://superredundant.com/wp-content/uploads/([^"]+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
    echo "No comic found on home page! :(\n";
    exit(1);
}

if (!is_dir('superredundant')) {
    mkdir('superredundant');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    if (is_file('superredundant/' . $matches[1])) {
        return;
    }

    echo "Downloading {$matches[1]}\n";
    $url = "http://superredundant.com/wp-content/uploads/{$matches[1]}";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("superredundant/{$matches[1]}", $data);
    }

    // Find previous page link
    $regex = '@href="(http://superredundant.com/[\\?0-9a-zA-Z/=-]+)" class="navi comic-nav-previous navi-prev"@';
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

    preg_match('@<img src="http://superredundant.com/wp-content/uploads/([^"]+\\.(jpg|png|gif))@', $html, $matches);
    if (empty($matches[1])) {
        echo "No image found on page!\n";
        return;
    }

    usleep(5e5);
}
